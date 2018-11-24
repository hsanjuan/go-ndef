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

package absoluteuri

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	absu := New("http://a.b", []byte{0x00})
	if !bytes.Equal(absu.Payload, []byte{0x00}) {
		t.Error("The type should hold the given payload")
	}
	if absu.Type() != "http://a.b" {
		t.Error("Unexpected Type name")
	}
}

func TestString(t *testing.T) {
	absu := New("http://a.b", []byte{0x00})
	if absu.String() != "<The message contains a binary payload>" {
		t.Error("Bad string generation")
	}

	absu = New("http://a.b", []byte{})
	if absu.String() != "" {
		t.Error("Expected an empty string")
	}
}

func TestMarshal(t *testing.T) {
	absu := New("http://a.b", []byte{0x04})
	pl := absu.Marshal()
	if !bytes.Equal(pl, []byte{0x04}) {
		t.Error("Bad payload generation")
	}
}

func TestUnmarshal(t *testing.T) {
	bts := []byte{0x79}
	absu := new(Payload)
	absu.Unmarshal(bts)
	if !bytes.Equal(absu.Payload, []byte{0x79}) {
		t.Error("Bad unmarshaling")
	}
}

func TestLen(t *testing.T) {
	absu := New("http://a.b", []byte{1, 2, 3})
	if absu.Len() != 3 {
		t.Error("Unexpected length")
	}
}
