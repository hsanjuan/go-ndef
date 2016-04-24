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

// Record represents a consolidated NDEF Record (assembled, non-chunked),
// which is a part of an NDEF Message.
type Record struct {
	TNF     byte          // Type name format (3 bits)
	Type    string        // Type of the payload. Must follow TNF
	ID      string        // An URI (per RFC 3986)
	Payload RecordPayload // Payload
}

// Reset clears up all the fields of the Record and sets them to their
// default values.
func (r *Record) Reset() {
	r.TNF = 0
	r.Type = ""
	r.ID = ""
	r.Payload = nil
}

// String a string representation of the payload of the record, prefixed
// by the URN of the resource.
//
// Note that not all NDEF Payloads are supported, and that custom types/payloads
// are considered not printable. In those cases, a generic RecordPayload is
// used and an explanatory message is returned instead.
// See submodules under "types/" for a list of supported types.
func (r *Record) String() string {
	return r.Payload.URN() + ":" + r.Payload.String()
}

// Inspect provides a string with information about this record.
// For a String representation of the contents use String().
func (r *Record) Inspect() string {
	var str string
	str += fmt.Sprintf("TNF: %d\n", r.TNF)
	str += fmt.Sprintf("Type: %s\n", r.Type)
	str += fmt.Sprintf("ID: %s\n", r.ID)
	str += fmt.Sprintf("Payload Length: %d", r.Payload.Len())
	return str
}

// Unmarshal parses a byte slice into a Record struct (the slice can
// have extra bytes which are ignored). The Record is always reset before
// parsing.
//
// It does this by parsing every record chunk until a MessageEnd chunk
// is read. Then it consolidates the chunks into a single Record and sets
// the TNF, Type and ID fields.
//
// Returns how many bytes were parsed from the slice (record length) or
// an error if something went wrong.
func (r *Record) Unmarshal(buf []byte) (rLen int, err error) {
	r.Reset()
	rLen = 0
	var chunks []*recordChunk
	for rLen < len(buf) {
		chunk := new(recordChunk)
		chunkSize, err := chunk.Unmarshal(buf[rLen:])
		rLen += chunkSize
		if err != nil {
			return rLen, err
		}
		chunks = append(chunks, chunk)
		if chunk.ME {
			break
		}
	}

	err = checkChunks(chunks)
	if err != nil {
		return rLen, err
	}

	r.TNF = chunks[0].TNF
	r.Type = chunks[0].Type
	r.ID = chunks[0].ID

	var buffer bytes.Buffer
	for _, c := range chunks {
		buffer.Write(c.Payload)
	}
	payloadBytes := buffer.Bytes()
	r.Payload = makeRecordType(r.TNF, r.Type, payloadBytes)

	r.Payload.Unmarshal(payloadBytes)
	err = r.check()
	return rLen, err
}

// Marshal returns the byte representation of a Record. It does this
// by producing a single record chunk.
//
// Note that if the original Record was unmarshaled from many chunks,
// the recovery is not possible anymore.
func (r *Record) Marshal() ([]byte, error) {
	err := r.check()
	if err != nil {
		return nil, err
	}
	tempChunk := new(recordChunk)
	tempChunk.MB = true
	tempChunk.ME = true
	tempChunk.CF = false
	tempChunk.IL = len(r.ID) > 0
	tempChunk.TNF = r.TNF
	tempChunk.TypeLength = byte(len([]byte(r.Type)))
	tempChunk.Type = r.Type
	tempChunk.IDLength = byte(len([]byte(r.ID)))
	tempChunk.ID = r.ID

	rPayload := r.Payload.Marshal()
	payloadLen := len(rPayload)

	if payloadLen > 4294967295 { //2^32-1. 4GB message max.
		payloadLen = 4294967295
	}
	tempChunk.SR = payloadLen < 256 // Short record vs. Long
	tempChunk.PayloadLength = uint64(payloadLen)

	// FIXME: If payload is greater than 2^32 - 1
	// we'll truncate without warning with this
	tempChunk.Payload = rPayload[:payloadLen]

	rBytes, err := tempChunk.Marshal()
	return rBytes, err
}

func (r *Record) check() error {
	return nil
}

// Set some short-hands for the errors that can happen on checkChunks().
const (
	eNORECORDS     = "checkChunks: No records"
	eNOMB          = "checkChunks: First record has not the MessageBegin flag set"
	eFIRSTCHUNKED  = "checkChunks: A single record cannot have the Chunk flag set"
	eNOME          = "checkChunks: Last record has not the MessageEnd flag set"
	eLASTCHUNKED   = "checkChunks: Last record cannot have the Chunk flag set"
	eCFMISSING     = "checkChunks: Chunk Flag missing from some records"
	eBADIL         = "checkChunks: IL flag is set on a middle or final chunk"
	eBADTYPELENGTH = "checkChunks: A middle or last chunk has TypeLength != 0"
	eBADTNF        = "checkChunks: A middle or last chunk TNF is not UNCHANGED"
)

func checkChunks(chunks []*recordChunk) error {
	chunksLen := len(chunks)
	last := chunksLen - 1
	if chunksLen == 0 {
		return errors.New(eNORECORDS)
	}
	if !chunks[0].MB {
		return errors.New(eNOMB)
	}
	if chunksLen == 1 && chunks[0].CF {
		return errors.New(eFIRSTCHUNKED)
	}
	if !chunks[last].ME {
		return errors.New(eNOME)
	}
	if chunks[0].CF && chunks[last].CF {
		return errors.New(eLASTCHUNKED)
	}

	if chunksLen > 1 {
		chunksWithoutCF := 0
		chunksWithIL := 0
		chunksWithTypeLength := 0
		chunksWithoutUnchangedType := 0
		for i, r := range chunks {
			// Check CF in all but the last
			if !r.CF && i != last {
				chunksWithoutCF++
			}
			// Check IL in all but first
			if r.IL && i != 0 {
				chunksWithIL++
			}
			// TypeLength should be zero except in the first
			if r.TypeLength > 0 && i != 0 {
				chunksWithTypeLength++
			}
			// All but first chunk should have TNF to 0x06
			if r.TNF != Unchanged && i != 0 {
				chunksWithoutUnchangedType++
			}
		}
		if chunksWithoutCF > 0 {
			return errors.New(eCFMISSING)
		}
		if chunksWithIL > 0 {
			return errors.New(eBADIL)
		}
		if chunksWithTypeLength > 0 {
			return errors.New(eBADTYPELENGTH)
		}
		if chunksWithoutUnchangedType > 0 {
			return errors.New(eBADTNF)
		}
	}
	return nil
}
