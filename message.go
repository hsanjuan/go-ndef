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

// Message represents an NDEF Message, which is a collection of one or
// more NDEF records.
//
// Each message has a Type Name Field (TNF), a Type and an optional ID,
// which are indicated by its first record. It also contains a Payload,
// which is the concatenation of the payloads of all the chunks.
type Message struct {
	// NDEF records are set when parsing bytes, and re-used
	// when generating bytes. They are not meant to be set
	// directly, but there is a setter useful for testing and
	// for someone that really needs to created a chunked Message
	records []*Record

	TNF     byte   // See possible values in the constants.
	Type    []byte // Message type. Typically
	ID      []byte // Message ID is optional
	Payload []byte // Message payload
}

// Return a string with some information about the message, and if it's easy
// enough to be placed in a string, the payload.
func (m *Message) String() string {
	var str string
	switch m.TNF {
	case Empty:
		str += fmt.Sprintf("Payload is EMPTY.")
	case NFCForumWellKnownType:
		str += fmt.Sprintf("urn:nfc:wkt:%s:", string(m.Type))
		switch string(m.Type) {
		case "T": // Plain text
			str += fmt.Sprintln(string(m.Payload))
		case "U": // URI
			str += fmt.Sprintf("%s%s\n",
				URIProtocols[m.Payload[0]],
				string(m.Payload[1:]))
		default:
			str += fmt.Sprintln("Payload is a NFC Forum Well" +
				"Known Type but don't know how to print it.")
		}
	case MediaType: // as defined at https://www.ietf.org/rfc/rfc2046.txt
		str += fmt.Sprintf("Payload is a media type: %s. ",
			string(m.Type))
		str += fmt.Sprintln("Payload will not be printed")
	case AbsoluteURI: // as defined https://www.ietf.org/rfc/rfc3986.txt
		str += fmt.Sprintln(string(m.Payload))
	case NFCForumExternalType:
		str += fmt.Sprintf("Payload is of type EXTERNAL.")
	case Unknown:
		str += fmt.Sprintf("Payload is of type UNKNOWN.")
	case Unchanged:
		str += fmt.Sprintf("Payload is of type UNCHANGED.")
	case Reserved:
		str += fmt.Sprintf("Payload is of type RESERVED.")
	}
	return str
}

// ParseBytes parses s a byte slice into a Message. It will parse
// each record until and including the  Message End Record. Then
// it will assemble the Payload and set the TNF, Type, ID fields with
// the correct information.
//
// Returns the number of bytes processed (message length), or an error
// if something looks wrong with the message or its records.
func (m *Message) ParseBytes(byteSlice []byte) (int, error) {
	i := 0
	for i < len(byteSlice) {
		r := new(Record)
		rLen, err := r.ParseBytes(byteSlice[i:])
		if err != nil {
			return 0, err
		}
		i += rLen
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

	firstRecord := m.records[0]
	m.TNF = firstRecord.TNF
	m.Type = firstRecord.Type
	m.ID = firstRecord.ID

	var buffer bytes.Buffer
	for _, r := range m.records {
		buffer.Write(r.Payload)
	}
	m.Payload = buffer.Bytes()
	return i, nil
}

// Bytes provides the byte slice representation of a Message
//
// There are two ways this can happen. If there are any Records,
// the concatenation of the Bytes() for each record is provided.
// Otherwise, a single record is produced from the Message fields
// (TNF, Type, ID, Payload) and its Bytes() returned.
//
// This allows the possibility of creating an NDEF Message by either
// setting the fields of the Message struct, or by manually providing the
// NDEF Record(s) with SetRecords().
//
// Returns an error if something goes wrong.
func (m *Message) Bytes() ([]byte, error) {
	if len(m.records) > 0 {
		// We have records. Just concat their Bytes. But test first
		if err := m.TestRecords(); err != nil {
			return nil, err
		}
		var buffer bytes.Buffer
		for _, r := range m.records {
			rBytes, err := r.Bytes()
			if err != nil {
				return nil, err
			}
			buffer.Write(rBytes)
		}
		return buffer.Bytes(), nil
	}

	// We have no records.
	// FIXME: Truncates when data > 4GB
	tempRecord := new(Record)
	tempRecord.MB = true
	tempRecord.ME = true
	tempRecord.CF = false
	tempRecord.IL = len(m.ID) > 0
	tempRecord.TypeLength = byte(len(m.Type))
	tempRecord.Type = m.Type
	tempRecord.IDLength = byte(len(m.ID))
	tempRecord.ID = m.ID
	tempRecord.TNF = m.TNF
	payloadLen := len(m.Payload)
	if payloadLen > 4294967295 { //2^32-1. 4GB message max.
		payloadLen = 2 ^ 32 - 1
	}
	if payloadLen < 256 { // Short Record
		tempRecord.SR = true
		tempRecord.PayloadLength = [4]byte{
			byte(payloadLen), 0, 0, 0}
	} else { // Long record
		tempRecord.SR = false
		copy(tempRecord.PayloadLength[:],
			Uint64ToBytes(uint64(payloadLen), 4))
	}
	// FIXME: If payload is greater than 2^32 - 1
	// we'll truncate without warning with this
	tempRecord.Payload = m.Payload[:payloadLen]
	tempMessage := new(Message)
	tempMessage.SetRecords([]*Record{tempRecord})
	return tempMessage.Bytes() // A message with 1 record
}

// Set some short-hands for the errors that can happen on TestRecords().
const (
	ENORECORDS     = "No records"
	ENOMB          = "First record has not the MessageBegin flag set"
	EFIRSTCHUNKED  = "A single record cannot have the Chunk flag set"
	ENOME          = "Last record has not the MessageEnd flag set"
	ELASTCHUNKED   = "Last record cannot have the Chunk flag set"
	ECFMISSING     = "Chunk Flag missing from some records"
	EBADIL         = "IL flag is set on a middle or final chunk"
	EBADTYPELENGTH = "A middle or last chunk has TypeLength != 0"
	EBADTNF        = "A middle or last chunk TNF is not UNCHANGED"
)

// TestRecords performs checks which are inspired in the "2.5 NDEF Mechanisms
// Test Requirements" section of the specification.
//
// Returns an error if the NDEF Message Records don't look good.
func (m *Message) TestRecords() error {
	records := m.records
	recordsLen := len(records)
	last := recordsLen - 1
	if recordsLen == 0 {
		return errors.New(ENORECORDS)
	}
	if !records[0].MB {
		return errors.New(ENOMB)
	}
	if recordsLen == 1 && records[0].CF {
		return errors.New(EFIRSTCHUNKED)
	}
	if !records[last].ME {
		return errors.New(ENOME)
	}
	if records[0].CF && records[last].CF {
		return errors.New(ELASTCHUNKED)
	}

	if recordsLen > 1 {
		recordsWithoutCF := 0
		recordsWithIL := 0
		recordsWithTypeLength := 0
		recordsWithoutUnchangedType := 0
		for i, r := range records {
			// Check CF in all but the last
			if !r.CF && i != last {
				recordsWithoutCF++
			}
			// Check IL in all but first
			if r.IL && i != 0 {
				recordsWithIL++
			}
			// TypeLength should be zero except in the first
			if r.TypeLength > 0 && i != 0 {
				recordsWithTypeLength++
			}
			// All but first chunk should have TNF to 0x06
			if r.TNF != Unchanged && i != 0 {
				recordsWithoutUnchangedType++
			}
		}
		if recordsWithoutCF > 0 {
			return errors.New(ECFMISSING)
		}
		if recordsWithIL > 0 {
			return errors.New(EBADIL)
		}
		if recordsWithTypeLength > 0 {
			return errors.New(EBADTYPELENGTH)
		}
		if recordsWithoutUnchangedType > 0 {
			return errors.New(EBADTNF)
		}
	}
	return nil
}

// SetRecords allows to manually set the private m.records field of a Message.
//
// This is useful for testing and for those who require to
// produce a chunked NDEF Message. In this case, manual construction of
// every record is necessary, along with a good read of the specification.
func (m *Message) SetRecords(records []*Record) {
	m.records = records
}
