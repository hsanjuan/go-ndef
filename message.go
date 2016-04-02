/***
    Copyright (c) 2016, Hector Sanjuan

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Lesser General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Lesser General Public License for more details.

    You should have received a copy of the GNU Lesser General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
***/

package ndef

import (
	"bytes"
	"errors"
	"fmt"
)

// An NDEF Message is a collection of one or more NDEF records
// When there are several records in an NDEF message, it means
// the contents are chunked and the payload needs to be formed by
// putting together all the chunks.
type Message struct {
	// NDEF records are set when parsing bytes, and re-used
	// when generating bytes. They are not meant to be set
	// directly, but there is a setter useful for testing and
	// for someone that really needs to created a chunked Message
	records []*Record

	TNF     byte // See possible values in ndef.go
	Type    []byte
	ID      []byte // Optional
	Payload []byte // Payload built from the records
}

// Return a string with some information about the message, and if it's easy
// enough to be placed in a string, the payload.
func (m *Message) String() string {
	var str string
	switch m.TNF {
	case EMPTY:
		str += fmt.Sprintf("Payload is EMPTY.")
	case NFC_FORUM_WELL_KNOWN_TYPE:
		str += fmt.Sprintf("urn:nfc:wkt:%s:", string(m.Type))
		switch string(m.Type) {
		case "T": // Plain text
			str += fmt.Sprintln(string(m.Payload))
		case "U": // URI
			str += fmt.Sprintf("%s%s\n",
				URIProtocols(m.Payload[0]),
				string(m.Payload[1:]))
		default:
			str += fmt.Sprintln("Payload is a NFC Forum Well" +
				"Known Type but don't know how to print it.")
		}
	case MEDIA_TYPE: // as defined at https://www.ietf.org/rfc/rfc2046.txt
		str += fmt.Sprintf("Payload is a media type: %s. ",
			string(m.Type))
		str += fmt.Sprintln("Payload will not be printed")
	case ABSOLUTE_URI: // as defined https://www.ietf.org/rfc/rfc3986.txt
		str += fmt.Sprintln(string(m.Payload))
	case NFC_FORUM_EXTERNAL_TYPE:
		str += fmt.Sprintf("Payload is of type EXTERNAL.")
	case UNKNOWN:
		str += fmt.Sprintf("Payload is of type UNKNOWN.")
	case UNCHANGED:
		str += fmt.Sprintf("Payload is of type UNCHANGED.")
	case RESERVED:
		str += fmt.Sprintf("Payload is of type RESERVED.")
	}
	return str
}

// Parses a byte slice into an Message. It will parse each record found until
// a Message End. Then it will assemble the Payload and fill in the fields
// with information about the message.
//
// Returns the number of bytes processed (message length), or an error if something looks wrong,
// with the message or something happened parsing the Records.
func (m *Message) ParseBytes(byte_array []byte) (int, error) {
	i := 0
	for i < len(byte_array) {
		r := new(Record)
		r_len, err := r.ParseBytes(byte_array[i:])
		if err != nil {
			return 0, err
		}
		i += r_len
		m.records = append(m.records, r)
		// In case our byte
		if r.ME {
			break
		}
	}

	err := m.TestRecords()
	if err != nil {
		return 0, err
	}

	first_record := m.records[0]
	m.TNF = first_record.TNF
	m.Type = first_record.Type
	m.ID = first_record.ID

	var buffer bytes.Buffer
	for _, r := range m.records {
		buffer.Write(r.Payload)
	}
	m.Payload = buffer.Bytes()
	return i, nil
}

// Get the byte representation of a Message
//
// There are two ways this can happen. If there are any Records, we just
// concat the Bytes() for each record. Otherwise, we produce a single
// record from the information we have and return its Bytes(). This
// allows the possibility of creating an NDEF Message by either setting the
// Payload, TNF and Type in the struct, or by manually providing the
// NDEF Record(s) with SetRecords().
//
// Returns an error if somehow things don't look OK.
func (m *Message) Bytes() ([]byte, error) {
	if len(m.records) > 0 {
		// We have records. Just concat their Bytes. But test first
		if err := m.TestRecords(); err != nil {
			return nil, err
		}
		var buffer bytes.Buffer
		for _, r := range m.records {
			r_bytes, err := r.Bytes()
			if err != nil {
				return nil, err
			}
			buffer.Write(r_bytes)
		}
		return buffer.Bytes(), nil
	} else {
		// We have no records.
		// FIXME: Truncates when data > 4GB
		temp_record := new(Record)
		temp_record.MB = true
		temp_record.ME = true
		temp_record.CF = false
		temp_record.IL = len(m.ID) > 0
		temp_record.TypeLength = byte(len(m.Type))
		temp_record.Type = m.Type
		temp_record.IDLength = byte(len(m.ID))
		temp_record.ID = m.ID
		pl_length := len(m.Payload)
		if pl_length > 4294967295 { //2^32-1. 4GB message max.
			pl_length = 2 ^ 32 - 1
		}
		if pl_length < 256 { // Short Record
			temp_record.SR = true
			temp_record.PayloadLength = [4]byte{
				byte(pl_length), 0, 0, 0}
		} else { // Long record
			temp_record.SR = false
			copy(temp_record.PayloadLength[:],
				Uint64ToBytes(uint64(pl_length), 4))
		}
		// FIXME: If payload is greater than 2^32 - 1
		// we'll truncate without warning with this
		temp_record.Payload = m.Payload[:pl_length]
		temp_message := new(Message)
		temp_message.SetRecords([]*Record{temp_record})
		return temp_message.Bytes() // A message with 1 record
	}
}

const (
	ERROR_NO_RECORDS     = "No records"
	ERROR_NO_MB          = "First record has not the MessageBegin flag set"
	ERROR_FIRST_CHUNKED  = "A single record cannot have the Chunk flag set"
	ERROR_NO_ME          = "Last record has not the MessageEnd flag set"
	ERROR_LAST_CHUNKED   = "Last record cannot have the Chunk flag set"
	ERROR_CF_MISSING     = "Chunk Flag missing from some records"
	ERROR_BAD_IL         = "IL flag is set on a middle or final chunk"
	ERROR_BAD_TYPELENGTH = "A middle or last chunk has TypeLength != 0"
	ERROR_BAD_TNF        = "A middle or last chunk TNF is not UNCHANGED"
)

// This function performs checks which are inspired in 2.5 NDEF Mechanisms
// Test Requirements section of the standard,
// and return an error if the NDEF Message Records don't look good.
func (m *Message) TestRecords() error {
	records := m.records
	len_records := len(records)
	last := len_records - 1
	if len_records == 0 {
		return errors.New(ERROR_NO_RECORDS)
	}
	if !records[0].MB {
		return errors.New(ERROR_NO_MB)
	}
	if len_records == 1 && records[0].CF {
		return errors.New(ERROR_FIRST_CHUNKED)
	}
	if !records[last].ME {
		return errors.New(ERROR_NO_ME)
	}
	if records[0].CF && records[last].CF {
		return errors.New(ERROR_LAST_CHUNKED)
	}

	if len_records > 1 {
		records_without_CF := 0
		records_with_IL := 0
		records_with_TypeLength := 0
		records_without_unchanged_type := 0
		for i, r := range records {
			// Check CF in all but the last
			if !r.CF && i != last {
				records_without_CF++
			}
			// Check IL in all but first
			if r.IL && i != 0 {
				records_with_IL++
			}
			// TypeLength should be zero except in the first
			if r.TypeLength > 0 && i != 0 {
				records_with_TypeLength++
			}
			// All but first chunk should have TNF to 0x06
			if r.TNF != UNCHANGED && i != 0 {
				records_without_unchanged_type++
			}
		}
		if records_without_CF > 0 {
			return errors.New(ERROR_CF_MISSING)
		}
		if records_with_IL > 0 {
			return errors.New(ERROR_BAD_IL)
		}
		if records_with_TypeLength > 0 {
			return errors.New(ERROR_BAD_TYPELENGTH)
		}
		if records_without_unchanged_type > 0 {
			return errors.New(ERROR_BAD_TNF)
		}
	}
	return nil
}

// To manually set the Records. Useful for testing and for advanced usage.
func (m *Message) SetRecords(records []*Record) {
	m.records = records
}
