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

package uri

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	u := New("https://something.com")
	if u.IdentCode != 4 {
		t.Error("Bad parsing of the provided protocol in New")
	}
	if u.URIField != "something.com" {
		t.Error("Bad removal of the protocol in New")
	}
	if u.URN() != "urn:nfc:wkt:U" {
		t.Error("Unexpected URN")
	}
}

func TestString(t *testing.T) {
	u := &URI{
		IdentCode: 20,
		URIField:  "mail",
	}
	if u.String() != "pop:mail" {
		t.Error("Bad string generation")
	}
}

func TestMarshal(t *testing.T) {
	u := New("https://www.a.a")
	pl := u.Marshal()
	if !bytes.Equal(pl, []byte{0x02, 0x61, 0x2e, 0x61}) {
		t.Error("Bad payload generation")
	}
}

func TestUnmarshal(t *testing.T) {
	bytes := []byte{0x02, 0x61, 0x2e, 0x61}
	u := new(URI)
	u.Unmarshal(bytes)
	if u.IdentCode != 2 || u.URIField != "a.a" {
		t.Error("Bad unmarshaling")
	}
}

func TestLen(t *testing.T) {
	u := New("http://ab.com")
	if u.Len() != 7 {
		t.Error("Unexpected length")
	}
}
