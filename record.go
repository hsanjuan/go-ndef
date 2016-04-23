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
	"runtime"
)

// Record represents a NDEF Record, which is part of an NDEF Message.
// Records follow some strict rules when they go together in a Message
// (see the Message TestRecords()). Records can have two forms:
// a ShortRecord (SR) only uses 1 byte for the Payload Length, but a regular
// record uses 4 bytes instead.
type Record struct {
	// First byte
	MB            bool    // Message begin
	ME            bool    // Message end
	CF            bool    // Chunk Flag
	SR            bool    // Short record
	IL            bool    // ID length field present
	TNF           byte    // Type name format (3 bits)
	TypeLength    byte    // Type Length
	IDLength      byte    // Length of the ID field
	PayloadLength [4]byte // Length of the Payload. For SR: only first byte.
	Type          []byte  // Type of the payload. Must follow TNF
	ID            []byte  // Unique ID, only in MB record
	Payload       []byte  // Payload
}

// Reset clears up all the fields of the Record and sets them to their
// default values.
func (r *Record) Reset() {
	r.MB = false
	r.ME = false
	r.CF = false
	r.SR = false
	r.IL = false
	r.TNF = 0
	r.TypeLength = 0
	r.IDLength = 0
	r.PayloadLength = [4]byte{0, 0, 0, 0}
	r.Type = []byte{}
	r.ID = []byte{}
	r.Payload = []byte{}
}

// Provide a string with information about this record.
// Records' payload do not make sense without having compiled a whole Message
// so they are not dealed with here.
func (r *Record) String() string {
	var str string
	str += fmt.Sprintf("MB: %t | ME: %t | CF: %t | SR: %t | IL: %t | TNF: %d\n",
		r.MB, r.ME, r.CF, r.SR, r.IL, r.TNF)
	str += fmt.Sprintf("TypeLength: %d", r.TypeLength)
	str += fmt.Sprintf(" | Type: %s\n", string(r.Type))
	if r.SR {
		str += fmt.Sprintf("Record Payload Length: %d",
			r.PayloadLength[0])
	} else {
		str += fmt.Sprintf("Record Payload Length: %d",
			bytesToUint64(r.PayloadLength[:]))
	}
	if r.IL {
		str += fmt.Sprintf(" | IDLength: %d", r.IDLength)
		str += fmt.Sprintf(" | ID: %x", bytesToUint64(r.ID))
	}
	str += fmt.Sprintf("\n")
	return str
}

// Unmarshal parses a byte slice into a single Record struct (the slice can
// have extra bytes which are ignored). The Record is always reset before
// parsing.
//
// Returns how many bytes were parsed from the slice (record length) or
// an error if something went wrong.
func (r *Record) Unmarshal(buf []byte) (rlen int, err error) {
	// Handle errors that are produced by getByte() and getBytes()
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			err = errors.New("Record.Unmarshal: " + err.Error())
		}
	}()
	r.Reset()
	bytesBuf := bytes.NewBuffer(buf)

	firstByte := getByte(bytesBuf)
	r.MB = (firstByte >> 7 & 0x1) == 1
	r.ME = (firstByte >> 6 & 0x1) == 1
	r.CF = (firstByte >> 5 & 0x1) == 1
	r.SR = (firstByte >> 4 & 0x1) == 1
	r.IL = (firstByte >> 3 & 0x1) == 1
	r.TNF = firstByte & 0x7

	r.TypeLength = getByte(bytesBuf)

	var payloadLen int
	if r.SR { //This is a short record
		r.PayloadLength[0] = getByte(bytesBuf)
		payloadLen = int(r.PayloadLength[0])
	} else { // Regular record
		var pl [4]byte
		copy(pl[:], getBytes(bytesBuf, 4))
		r.PayloadLength = pl
		payloadLen = int(bytesToUint64(r.PayloadLength[:]))
	}
	if r.IL {
		r.IDLength = getByte(bytesBuf)
	}
	r.Type = getBytes(bytesBuf, int(r.TypeLength))
	if r.IL {
		r.ID = getBytes(bytesBuf, int(r.IDLength))
	}
	r.Payload = getBytes(bytesBuf, payloadLen)
	// Return the records length:
	// length of original buffer - length of unread portion.
	return len(buf) - bytesBuf.Len(), nil
}

// Marshal returns the byte representation of a Record, or an error
// if something went wrong
func (r *Record) Marshal() ([]byte, error) {
	var buffer bytes.Buffer
	firstByte := byte(0)
	if r.MB {
		firstByte |= 0x1 << 7
	}
	if r.ME {
		firstByte |= 0x1 << 6
	}
	if r.CF {
		firstByte |= 0x1 << 5
	}
	if r.SR {
		firstByte |= 0x1 << 4
	}
	if r.IL {
		firstByte |= 0x1 << 3
	}
	firstByte |= (r.TNF & 0x7) //Last 3 bits are from TNF
	buffer.WriteByte(firstByte)
	// TypeLength byte
	buffer.WriteByte(r.TypeLength)

	// Payload Length byte (for SR) or 4 bytes for the regular case
	if r.SR {
		buffer.WriteByte(r.PayloadLength[0])
	} else {
		buffer.Write(r.PayloadLength[:])
	}

	// ID Length byte if we are meant to have it
	if r.IL {
		buffer.WriteByte(r.IDLength)
	}

	// Write the type bytes if we have something
	if r.TypeLength > 0 {
		buffer.Write(r.Type)
	}

	// Write the ID bytes if we have something
	if r.IL && r.IDLength > 0 {
		buffer.Write(r.ID)
	}

	buffer.Write(r.Payload)
	return buffer.Bytes(), nil
}
