/***
    Copyright (c) 2018, Hector Sanjuan

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
	"testing"
)

func TestNewSmartPosterPayload(t *testing.T) {
	msg := NewTextMessage("hi", "en")
	sp := NewSmartPosterPayload(msg)
	if sp.Type() != "urn:nfc:wkt:Sp" {
		t.Error("Expected URN")
	}
}

func TestString(t *testing.T) {
	msg := NewTextMessage("hi", "en")
	sp := NewSmartPosterPayload(msg)
	if sp.String() != "\nurn:nfc:wkt:T:hi" {
		t.Error("Bad string generation")
	}
}

func TestMarshal(t *testing.T) {
	msg := NewURIMessage("http://www.nfc-forum.org")
	sp := NewSmartPosterPayload(msg)
	pl := sp.Marshal()
	if !bytes.Equal(pl, []byte{0xD1, 0x01, 0x0E, 'U', 0x01, 'n', 'f', 'c', '-', 'f', 'o', 'r', 'u', 'm', '.', 'o', 'r', 'g'}) {
		t.Error("Bad payload generation")
	}
}

func TestUnmarshal(t *testing.T) {
	bts := []byte{0xD1, 0x01, 0x0E, 'U', 0x01, 'n', 'f', 'c', '-', 'f', 'o', 'r', 'u', 'm', '.', 'o', 'r', 'g'}
	sp := new(SmartPosterPayload)
	sp.Unmarshal(bts)
	if sp.Message.String() != "urn:nfc:wkt:U:http://www.nfc-forum.org" {
		t.Error("Bad unmarshaling: ", sp.Message.String())
	}
}

func TestLen(t *testing.T) {
	msg := NewURIMessage("http://www.nfc-forum.org")
	sp := NewSmartPosterPayload(msg)
	if l := sp.Len(); l != 18 {
		t.Error("Unexpected length: ", l)
	}
}
