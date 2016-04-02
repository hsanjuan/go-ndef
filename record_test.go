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

func TestRecordBytesAndParsing(t *testing.T) {
	t.Log("Testing with SR Record")
	r := &Record{
		MB:            true,
		ME:            true,
		CF:            false,
		SR:            true,
		IL:            true,
		TNF:           UNKNOWN,
		TypeLength:    4,
		IDLength:      3,
		PayloadLength: [4]byte{3, 0, 0, 0},
		Type:          []byte("test"),
		ID:            []byte{1, 2, 3},
		Payload:       []byte("abc"),
	}

	r_bytes, err := r.Bytes()
	if err != nil {
		t.Error(err)
	}
	r2 := new(Record)
	_, err = r2.ParseBytes(r_bytes)
	if err != nil {
		t.Error(err)
	}
	r2_bytes, err := r2.Bytes()
	if err != nil {
		t.Error(err)
	}
	t.Log("R1:", FmtBytes(r_bytes, len(r_bytes)))
	t.Log("R2:", FmtBytes(r2_bytes, len(r2_bytes)))
	if !bytes.Equal(r_bytes, r2_bytes) {
		t.Error("We cannot produce the same bytes after re-parsing a Record")
	}

	t.Log("Testing with regular Record")

	r = &Record{
		MB:            true,
		ME:            true,
		CF:            false,
		SR:            false,
		IL:            false,
		TNF:           UNKNOWN,
		TypeLength:    4,
		IDLength:      3,
		PayloadLength: [4]byte{0, 0, 0, 3},
		Type:          []byte("test"),
		ID:            []byte{1, 2, 3},
		Payload:       []byte("abc"),
	}

	r_bytes, err = r.Bytes()
	if err != nil {
		t.Error(err)
	}
	r2 = new(Record)
	_, err = r2.ParseBytes(r_bytes)
	if err != nil {
		t.Error(err)
	}
	r2_bytes, err = r2.Bytes()
	if err != nil {
		t.Error(err)
	}
	t.Log("R1:", FmtBytes(r_bytes, len(r_bytes)))
	t.Log("R2:", FmtBytes(r2_bytes, len(r2_bytes)))
	if !bytes.Equal(r_bytes, r2_bytes) {
		t.Error("We cannot produce the same bytes after re-parsing a Record")
	}
}

func TestRecordString(t *testing.T) {
	r := &Record{
		MB:            true,
		ME:            true,
		CF:            false,
		SR:            false,
		IL:            false,
		TNF:           UNKNOWN,
		TypeLength:    4,
		IDLength:      3,
		PayloadLength: [4]byte{0, 0, 0, 3},
		Type:          []byte("test"),
		ID:            []byte{1, 2, 3},
		Payload:       []byte("abc"),
	}
	t.Log(r)
}