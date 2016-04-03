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

func TestMessageBytesAndParsing(t *testing.T) {
	t.Log("Testing a Message created with a provided NDEF Record")
	r := &Record{
		MB:            true,
		ME:            true,
		CF:            false,
		SR:            true,
		IL:            true,
		TNF:           Unknown,
		TypeLength:    4,
		IDLength:      3,
		PayloadLength: [4]byte{3, 0, 0, 0},
		Type:          []byte("test"),
		ID:            []byte{1, 2, 3},
		Payload:       []byte("abc"),
	}

	m := new(Message)
	m.SetRecords([]*Record{r})
	mBytes, err := m.Bytes()
	if err != nil {
		t.Error(err)
	}
	m2 := new(Message)
	_, err = m2.ParseBytes(mBytes)
	if err != nil {
		t.Error(err)
	}
	m2Bytes, err := m2.Bytes()
	if err != nil {
		t.Error(err)
	}
	t.Log("M1:", FmtBytes(mBytes, len(mBytes)))
	t.Log("M2:", FmtBytes(m2Bytes, len(m2Bytes)))
	if !bytes.Equal(mBytes, m2Bytes) {
		t.Error("We cannot produce the same bytes after re-parsing a Message")
	}

	t.Log("Testing with a Message created by defining its fields")
	m = &Message{
		TNF:     Unknown,
		ID:      []byte{1, 2, 3},
		Type:    []byte("test"),
		Payload: []byte("abc"),
	}
	mBytes, _ = m.Bytes()
	m2 = new(Message)
	_, err = m2.ParseBytes(mBytes)
	if err != nil {
		t.Error(err)
	}
	m2Bytes, err = m2.Bytes()
	if err != nil {
		t.Error(err)
	}
	t.Log("M1:", FmtBytes(mBytes, len(mBytes)))
	t.Log("M2:", FmtBytes(m2Bytes, len(m2Bytes)))
	if !bytes.Equal(mBytes, m2Bytes) {
		t.Error("We cannot produce the same bytes after re-parsing a Message")
	}
}

func TestMessageString(t *testing.T) {
	// Just test we are not crashing

	m := &Message{
		TNF:     NFCForumWellKnownType,
		ID:      []byte{1, 2, 3},
		Type:    []byte("T"),
		Payload: []byte("abc"),
	}
	t.Log(m)

	m = &Message{
		TNF:     NFCForumWellKnownType,
		ID:      []byte{1, 2, 3},
		Type:    []byte("U"),
		Payload: []byte("\x03abc"),
	}
	t.Log(m)

	m = &Message{
		TNF:     NFCForumWellKnownType,
		ID:      []byte{1, 2, 3},
		Type:    []byte("X"),
		Payload: []byte("abc"),
	}
	t.Log(m)

	m = &Message{
		TNF: Empty,
	}
	t.Log(m)

	m = &Message{
		TNF:     MediaType,
		Type:    []byte("image/jpeg"),
		Payload: []byte("abc"),
	}
	t.Log(m)

	m = &Message{
		TNF: AbsoluteURI,
		ID:  []byte{1, 2, 3},
		//Type:    []byte("T"),
		Payload: []byte("http://abc.de"),
	}
	t.Log(m)

	m = &Message{
		TNF:     NFCForumExternalType,
		ID:      []byte{1, 2, 3},
		Type:    []byte("T"),
		Payload: []byte("abc"),
	}
	t.Log(m)

	m = &Message{
		TNF:     Unknown,
		ID:      []byte{1, 2, 3},
		Type:    []byte("T"),
		Payload: []byte("abc"),
	}
	t.Log(m)

	m = &Message{
		TNF:     Unchanged,
		ID:      []byte{1, 2, 3},
		Type:    []byte("T"),
		Payload: []byte("abc"),
	}
	t.Log(m)

	m = &Message{
		TNF:     Reserved,
		ID:      []byte{1, 2, 3},
		Type:    []byte("T"),
		Payload: []byte("abc"),
	}
	t.Log(m)
}

func TestNDEFBadMessageTest(t *testing.T) {
	cases := []struct{ Expected string }{
		{ENORECORDS},
		{ENOMB},
		{EFIRSTCHUNKED},
		{ENOME},
		{ELASTCHUNKED},
		{ECFMISSING},
		{EBADIL},
		{EBADTYPELENGTH},
		{EBADTNF},
	}

	errs := []error{}

	m := &Message{} // 0 records
	errs = append(errs, m.TestRecords())

	// First record is not MB
	rs := []*Record{
		&Record{
			MB: false,
			ME: true,
			CF: true,
		},
	}
	m.records = rs
	errs = append(errs, m.TestRecords())

	// First and only record is chuncked
	rs = []*Record{
		&Record{
			MB: true,
			ME: true,
			CF: true,
		},
	}
	m.records = rs
	errs = append(errs, m.TestRecords())

	// Last record is not ME
	rs = []*Record{
		&Record{
			MB: true,
			ME: false,
		},
	}
	m.records = rs
	errs = append(errs, m.TestRecords())

	// Last record is Chunked
	rs = []*Record{
		&Record{
			MB: true,
			ME: false,
			CF: true,
		},
		&Record{
			MB: false,
			ME: true,
			CF: true,
		},
	}
	m.records = rs
	errs = append(errs, m.TestRecords())

	// Record missing CF
	rs = []*Record{
		&Record{
			MB: true,
			ME: false,
			CF: false,
		},
		&Record{
			MB: false,
			ME: true,
			CF: false,
		},
	}
	m.records = rs
	errs = append(errs, m.TestRecords())

	// Non-first record with IL
	rs = []*Record{
		&Record{
			MB:       true,
			ME:       false,
			CF:       true,
			IL:       true,
			IDLength: 1,
			ID:       []byte("a"),
		},
		&Record{
			MB:       false,
			ME:       true,
			CF:       false,
			IL:       true,
			IDLength: 1,
			ID:       []byte("a"),
		},
	}
	m.records = rs
	errs = append(errs, m.TestRecords())

	// Non-first record with TypeLength
	rs = []*Record{
		&Record{
			MB:         true,
			ME:         false,
			CF:         true,
			TypeLength: 1,
			Type:       []byte("U"),
		},
		&Record{
			MB:         false,
			ME:         true,
			CF:         false,
			TypeLength: 1,
			Type:       []byte("U"),
		},
	}
	m.records = rs
	errs = append(errs, m.TestRecords())

	// Non-first record with BAD TNF
	rs = []*Record{
		&Record{
			MB:         true,
			ME:         false,
			CF:         true,
			TypeLength: 1,
			Type:       []byte("U"),
			TNF:        Empty,
		},
		&Record{
			MB:         false,
			ME:         true,
			CF:         false,
			TypeLength: 0,
			TNF:        Unknown,
		},
	}
	m.records = rs
	errs = append(errs, m.TestRecords())

	for i, err := range errs {
		t.Logf("Expected: %s...", cases[i].Expected)
		if err == nil {
			t.Error("Test didn't fail as expected")
		} else if e := cases[i].Expected; err.Error() != e {
			t.Errorf("Test failed unexpectedly because: %s.", err)
		} else {
			t.Log("Ok!")
		}
	}
}

func TestNDEFGoodMessageTest(t *testing.T) {
	m := &Message{}
	// Non-first record with TypeLength
	rs := []*Record{
		&Record{
			MB:            true,
			ME:            false,
			CF:            true,
			SR:            true,
			TNF:           NFCForumWellKnownType,
			TypeLength:    1,
			Type:          []byte("U"),
			PayloadLength: [4]byte{1, 0, 0, 0},
			Payload:       []byte("a"),
		},
		&Record{
			MB:            false,
			ME:            false,
			CF:            true,
			SR:            true,
			TypeLength:    0,
			TNF:           Unchanged,
			PayloadLength: [4]byte{2, 0, 0, 0},
			Payload:       []byte("bc"),
		},
		&Record{
			MB:            false,
			ME:            true,
			CF:            false,
			SR:            true,
			TypeLength:    0,
			TNF:           Unchanged,
			PayloadLength: [4]byte{1, 0, 0, 0},
			Payload:       []byte("d"),
		},
	}
	m.records = rs
	err := m.TestRecords()
	if err != nil {
		t.Error("Message was good but failed because:", err)
	}

	// Since we are here, test that we can reparse correctly
	mBytes, err := m.Bytes()
	if err != nil {
		t.Error(err)
	}
	m2 := &Message{}
	m2.ParseBytes(mBytes)
	if string(m2.Payload) != "abcd" {
		t.Error("Payload is not what we would expect!")
	}
}
