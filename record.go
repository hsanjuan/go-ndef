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
	"fmt"
	//	"errors"
)

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
			BytesToUint64(r.PayloadLength[:]))
	}
	if r.IL {
		str += fmt.Sprintf(" | IDLength: %d", r.IDLength)
		str += fmt.Sprintf(" | ID: %x", BytesToUint64(r.ID))
	}
	str += fmt.Sprintf("\n")
	return str
}

// Parse a byte slice into a Record struct. The byte slice should start
// with the NDEF record even though it could contain more records after it.
//
// Returns how many bytes were parsed from the slice (record length) or
// an error if something went wrong.
// FIXME: It will panic badly if there are not enough bytes to read.
func (r *Record) ParseBytes(bytes []byte) (int, error) {
	i := 0
	first_byte := bytes[i]
	r.MB = (first_byte >> 7 & 0x1) == 1
	r.ME = (first_byte >> 6 & 0x1) == 1
	r.CF = (first_byte >> 5 & 0x1) == 1
	r.SR = (first_byte >> 4 & 0x1) == 1
	r.IL = (first_byte >> 3 & 0x1) == 1
	r.TNF = first_byte & 0x7
	i++

	r.TypeLength = bytes[i]
	i++

	var pl_length int
	if r.SR { //This is a short record
		r.PayloadLength[0] = bytes[i]
		i++
		pl_length = int(r.PayloadLength[0])
	} else { // Regular record
		var pl [4]byte
		copy(pl[:], bytes[i:i+4])
		r.PayloadLength = pl
		i += 4
		pl_length = int(BytesToUint64(r.PayloadLength[:]))
	}
	if r.IL {
		r.IDLength = bytes[i]
		i++
	}
	r.Type = bytes[i : i+int(r.TypeLength)]
	i += int(r.TypeLength)
	if r.IL {
		r.ID = bytes[i : i+int(r.IDLength)]
		i += int(r.IDLength)
	}
	r.Payload = bytes[i : i+pl_length]
	i += pl_length
	// Return the records length
	return i, nil
}

// Returns the byte representation of the Record
func (r *Record) Bytes() ([]byte, error) {
	var buffer bytes.Buffer
	first_byte := byte(0)
	if r.MB {
		first_byte |= 0x1 << 7
	}
	if r.ME {
		first_byte |= 0x1 << 6
	}
	if r.CF {
		first_byte |= 0x1 << 5
	}
	if r.SR {
		first_byte |= 0x1 << 4
	}
	if r.IL {
		first_byte |= 0x1 << 3
	}
	first_byte |= (r.TNF & 0x7) //Last 3 bits are from TNF
	buffer.WriteByte(first_byte)
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
