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
	"github.com/hsanjuan/go-ndef/types"
	"github.com/hsanjuan/go-ndef/types/wkt/text"
	"github.com/hsanjuan/go-ndef/types/wkt/uri"
)

// The RecordPayload interface should be implemented by supported
// NDEF Record types. It ensures that we have a way to interpret payloads
// into printable information and to produce NDEF Record payloads for a given
// type.
type RecordPayload interface {
	// Returns a string representation of the Payload
	String() string
	// Provide serialization for the Payload
	Marshal() []byte
	// Provide de-serialization for the Payload
	Unmarshal(buf []byte)
	// Returns an Unified Resource Name for the Payload Type
	URN() string
	// Return the length of the Payload (serialized)
	Len() int
}

func makeRecordPayload(tnf byte, rtype string, payload []byte) RecordPayload {
	var r RecordPayload
	switch tnf {
	case NFCForumWellKnownType:
		switch rtype {
		case "U":
			r = new(uri.URI)
		case "T":
			r = new(text.Text)
		default:
			r = new(types.Generic)
		}
	default:
		r = new(types.Generic)
	}
	r.Unmarshal(payload)
	return r
}
