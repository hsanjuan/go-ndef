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

package text

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	tx := New("Text", "en")
	if tx.Language != "en" {
		t.Error("Should set language to 'en'")
	}
	if tx.Text != "Text" {
		t.Error("Should set Text to 'Text'")
	}
	if tx.Type() != "urn:nfc:wkt:T" {
		t.Error("Unexpected URN")
	}
}

func TestString(t *testing.T) {
	tx := &Payload{
		Language: "en_US",
		Text:     "hey",
	}
	if tx.String() != "hey" {
		t.Error("Bad string generation")
	}
}

func TestMarshal(t *testing.T) {
	tx := New("hey", "en")
	pl := tx.Marshal()
	if !bytes.Equal(pl, []byte{0x02, 0x65, 0x6e, 0x68, 0x65, 0x79}) {
		t.Error("Bad payload generation")
	}
}

func TestUnmarshal(t *testing.T) {
	bts := []byte{0x02, 0x65, 0x6e, 0x68, 0x65, 0x79}
	tx := new(Payload)
	tx.Unmarshal(bts)
	if tx.Language != "en" || tx.Text != "hey" {
		t.Error("Bad unmarshaling")
	}

	// Now with utf16
	bts = []byte{0x82, 0x65, 0x6e, 0x00, 0x68, 0x00, 0x65, 0x00, 0x79}
	tx = new(Payload)
	tx.Unmarshal(bts)
	if tx.Language != "en" || tx.Text != "hey" {
		t.Error("Bad unmarshaling in utf16")
	}

}

func TestLen(t *testing.T) {
	tx := New("ab", "en")
	if tx.Len() != 5 {
		t.Error("Unexpected length")
	}
}
