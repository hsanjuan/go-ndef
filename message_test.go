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
	"testing"
)

func ExampleMessage() {
	// Here we create a Message of type "U" (URI).
	// Note that the first byte of the payload is dedicated to encode the
	// URI Protocol
	ndefMessage := NewMessage(NFCForumWellKnownType, "U", "", []byte("\x04github.com/hsanjuan/go-ndef"))
	fmt.Println(ndefMessage)
	// Output:
	// urn:nfc:wkt:U:https://github.com/hsanjuan/go-ndef
}

func ExampleMessage_Unmarshal() {
	ndefMessageBytes := []byte{0xd1, 0x01, 0x23, 0x54, 0x02, 0x65, 0x6e,
		0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20,
		0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x20, 0x6f, 0x66,
		0x20, 0x54, 0x5b, 0x65, 0x78, 0x74, 0x5d, 0x20, 0x74, 0x79,
		0x70, 0x65}
	ndefMessage := &Message{}                         // Create uninitialized message
	_, err := ndefMessage.Unmarshal(ndefMessageBytes) // Parse bytes into it
	if err != nil {                                   // Your bytes don't look good
		fmt.Println(err)
		return
	}
	fmt.Println(ndefMessage) // Print the contents of every record
	// Output:
	// urn:nfc:wkt:T:This is a message of T[ext] type
}

func TestInspect(t *testing.T) {
	ndefMessage := NewMediaMessage("text/plain", []byte("abc"))
	t.Log(ndefMessage.Inspect())
}

func TestTypes(t *testing.T) {
	ndefMessage := NewTextMessage("abc", "en_US")
	t.Log(ndefMessage)

	ndefMessage = NewURIMessage("http://a.b")
	t.Log(ndefMessage)

	ndefMessage = NewMediaMessage("text/json", []byte(`{ "a" : 3 }`))
	t.Log(ndefMessage)

	ndefMessage = NewAbsoluteURIMessage("http://a.b", []byte("payload"))
	t.Log(ndefMessage)

	ndefMessage = NewExternalMessage("exttype", []byte("payload"))
	t.Log(ndefMessage)
}

func TestMarhsal(t *testing.T) {
	m := NewURIMessage("http://s.com")
	mBytes, err := m.Marshal()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	expectedBytes := []byte{
		0xd1,                         //Flags+TNF
		0x01,                         // TypeLen
		0x06,                         // Payload Len
		0x55,                         // Type
		0x03,                         // URI code
		0x73, 0x2e, 0x63, 0x6f, 0x6d, // URI
	}
	if !bytes.Equal(mBytes, expectedBytes) {
		t.Error("Unexpected marshal result")
	}
}
