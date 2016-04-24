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

// Package types provides a generic implementation for NDEF Payloads which are
// either custom or not supported yet.
package types

// Generic is a wrapper to store a Payload
type Generic struct {
	Payload []byte
}

// String returns a string explaining that we are not sure how to print
// this type
func (g *Generic) String() string {
	return "<Non standard type: contents not printable>"
}

// URN returns the Uniform Resource Name for Generics
func (g *Generic) URN() string {
	return "urn:nfc:ext:go-ndef:generic"
}

// Marshal returns the bytes representing the payload of a Record of Generic type
func (g *Generic) Marshal() []byte {
	return g.Payload
}

// Unmarshal parses the payload of a Generic type record
func (g *Generic) Unmarshal(buf []byte) {
	g.Payload = buf
}

// Len is the length of the byte slice resulting of Marshaling
func (g *Generic) Len() int {
	return len(g.Marshal())
}
