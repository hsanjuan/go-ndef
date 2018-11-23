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

package generic

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	gen := New([]byte{0x00})
	if !bytes.Equal(gen.Payload, []byte{0x00}) {
		t.Error("The type should hold the given payload")
	}
	if gen.Type() != "go-ndef-generic" {
		t.Error("Unexpected type name")
	}
}

func TestString(t *testing.T) {
	gen := New([]byte{0x00})
	if gen.String() != "<The message contains a binary payload>" {
		t.Error("Bad string generation")
	}

	gen = New([]byte{})
	if gen.String() != "" {
		t.Error("Expected an empty string")
	}
}

func TestMarshal(t *testing.T) {
	gen := New([]byte{0x04})
	pl := gen.Marshal()
	if !bytes.Equal(pl, []byte{0x04}) {
		t.Error("Bad payload generation")
	}
}

func TestUnmarshal(t *testing.T) {
	bts := []byte{0x79}
	gen := new(Payload)
	gen.Unmarshal(bts)
	if !bytes.Equal(gen.Payload, []byte{0x79}) {
		t.Error("Bad unmarshaling")
	}
}

func TestLen(t *testing.T) {
	gen := New([]byte{1, 2, 3})
	if gen.Len() != 3 {
		t.Error("Unexpected length")
	}
}
