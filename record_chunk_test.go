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
	"testing"
)

func TestRecordChunkMarshalUnmarshal(t *testing.T) {
	t.Log("Testing with SR recordChunk")
	r := &recordChunk{
		MB:            true,
		ME:            true,
		CF:            false,
		SR:            true,
		IL:            true,
		TNF:           NFCForumExternalType,
		TypeLength:    4,
		IDLength:      3,
		PayloadLength: 3,
		Type:          "test",
		ID:            "#ab",
		Payload:       []byte("abc"),
	}

	rBytes, err := r.Marshal()
	if err != nil {
		t.Error(err)
	}
	r2 := new(recordChunk)
	_, err = r2.Unmarshal(rBytes)
	if err != nil {
		t.Error(err)
	}
	r2Bytes, err := r2.Marshal()
	if err != nil {
		t.Error(err)
	}
	t.Log("R1:", fmtBytes(rBytes, len(rBytes)))
	t.Log("R2:", fmtBytes(r2Bytes, len(r2Bytes)))
	if !bytes.Equal(rBytes, r2Bytes) {
		t.Error("We cannot produce the same bytes after re-parsing a recordChunk")
	}

	t.Log("Testing with regular recordChunk")

	r = &recordChunk{
		MB:            true,
		ME:            true,
		CF:            false,
		SR:            false,
		IL:            false,
		TNF:           NFCForumExternalType,
		TypeLength:    4,
		IDLength:      3,
		PayloadLength: 3,
		Type:          "test",
		ID:            "#ab",
		Payload:       []byte("abc"),
	}

	rBytes, err = r.Marshal()
	if err != nil {
		t.Error(err)
	}
	r2 = new(recordChunk)
	_, err = r2.Unmarshal(rBytes)
	if err != nil {
		t.Error(err)
	}
	r2Bytes, err = r2.Marshal()
	if err != nil {
		t.Error(err)
	}
	t.Log("R1:", fmtBytes(rBytes, len(rBytes)))
	t.Log("R2:", fmtBytes(r2Bytes, len(r2Bytes)))
	if !bytes.Equal(rBytes, r2Bytes) {
		t.Error("We cannot produce the same bytes after re-parsing a recordChunk")
	}
}

func TestRecordChunkUnmarshalTooShort(t *testing.T) {
	t.Log("Testing with a truncated recordChunk bytes")
	r := &recordChunk{
		MB:            true,
		ME:            true,
		CF:            false,
		SR:            true,
		IL:            true,
		TNF:           NFCForumWellKnownType,
		TypeLength:    1,
		IDLength:      3,
		PayloadLength: 3,
		Type:          "T",
		ID:            "#ab",
		Payload:       []byte("abcdefg"),
	}

	rBytes, err := r.Marshal()
	if err != nil {
		t.Error(err)
	}
	r2 := new(recordChunk)
	_, err = r2.Unmarshal(rBytes[:len(rBytes)-5])
	if err == nil {
		t.Error("It should have errored")
	}
	t.Log("Throws error:", err)
}

func TestRecordChunkString(t *testing.T) {
	r := &recordChunk{
		MB:            true,
		ME:            true,
		CF:            false,
		SR:            false,
		IL:            false,
		TNF:           Unknown,
		TypeLength:    4,
		IDLength:      3,
		PayloadLength: 3,
		Type:          "test",
		ID:            "#ab",
		Payload:       []byte("abc"),
	}
	t.Log(r)
	r = &recordChunk{
		MB:            true,
		ME:            true,
		CF:            false,
		SR:            true,
		IL:            true,
		TNF:           Unknown,
		TypeLength:    4,
		IDLength:      3,
		PayloadLength: 3,
		Type:          "test",
		ID:            "#ab",
		Payload:       []byte("abc"),
	}
	t.Log(r)
}

func TestRecordChunkCheck(t *testing.T) {
	testcasesbad := map[string]*recordChunk{
		"non_ascii_type": &recordChunk{
			MB:            true,
			ME:            true,
			CF:            false,
			SR:            true,
			IL:            true,
			TNF:           NFCForumWellKnownType,
			TypeLength:    3,
			IDLength:      3,
			PayloadLength: 3,
			Type:          "⌘",
			ID:            "#ab",
			Payload:       []byte("abc"),
		},
		"bad_empty": &recordChunk{
			MB:            true,
			ME:            true,
			CF:            false,
			SR:            true,
			IL:            true,
			TNF:           Empty,
			TypeLength:    1,
			IDLength:      3,
			PayloadLength: 3,
			Type:          "a",
			ID:            "#ab",
			Payload:       []byte("abc"),
		},
		"type_on_unknown": &recordChunk{
			MB:            true,
			ME:            true,
			CF:            false,
			SR:            true,
			IL:            true,
			TNF:           Unknown,
			TypeLength:    1,
			IDLength:      3,
			PayloadLength: 3,
			Type:          "a",
			ID:            "#ab",
			Payload:       []byte("abc"),
		},
		"reserved_tnf": &recordChunk{
			MB:            true,
			ME:            true,
			CF:            false,
			SR:            true,
			IL:            true,
			TNF:           Reserved,
			TypeLength:    1,
			IDLength:      3,
			PayloadLength: 3,
			Type:          "a",
			ID:            "#ab",
			Payload:       []byte("abc"),
		},
	}

	t.Log("Testing with bad recordChunks")
	for k, r := range testcasesbad {
		_, err := r.Marshal()
		if err == nil {
			t.Error("Testcase", k, "should have failed")
		} else {
			// FIXME: are we getting the error we expect?
			t.Logf("%s: %s", k, err.Error())
		}
	}
}
