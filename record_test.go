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

	"github.com/hsanjuan/go-ndef/types/generic"
)

func TestRecordMarshalUnmarshal(t *testing.T) {
	t.Log("Testing a Record created with a provided chunk")
	r := NewRecord(
		NFCForumExternalType,
		"test",
		"#ab",
		&generic.Payload{[]byte("abc")},
	)

	rBytes, err := r.Marshal()
	if err != nil {
		t.Error(err)
	}
	r2 := new(Record)
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
		t.Error("We cannot produce the same bytes after re-parsing a Record")
	}
}

func TestRecordString(t *testing.T) {
	// Just test we are not crashing
	type args struct {
		TNF     byte
		Type    string
		ID      string
		Payload RecordPayload
	}

	tcs := []args{
		args{NFCForumWellKnownType, "X", "#ab", &generic.Payload{[]byte("abc")}},
		args{Empty, "", "", nil},
		args{MediaType, "image/jpeg", "", &generic.Payload{[]byte("\x03abc")}},
		args{AbsoluteURI, "http://resource", "#ab", &generic.Payload{[]byte("")}},
		args{NFCForumExternalType, "T", "#ab", &generic.Payload{[]byte("abc")}},
		args{Unknown, "", "", &generic.Payload{[]byte("abc")}},
		args{Unchanged, "", "", &generic.Payload{[]byte("abc")}},
	}

	m := NewTextRecord("abc", "en")
	t.Log(m)

	m = NewURIRecord("http://abc")
	t.Log(m)

	for _, tc := range tcs {
		r := NewRecord(tc.TNF, tc.Type, tc.ID, tc.Payload)
		t.Log(r)
	}
}

func TestRecordBadChunksTest(t *testing.T) {
	cases := []struct{ Expected string }{
		{eNOCHUNKS},
		{eFIRSTCHUNKED},
		{eLASTCHUNKED},
		{eCFMISSING},
		{eBADIL},
		{eBADTYPELENGTH},
		{eBADTNF},
	}

	errs := []error{}
	r := &Record{}
	r.chunks = []*recordChunk{}
	errs = append(errs, r.check())

	// First and only record is chuncked
	r.chunks = []*recordChunk{
		&recordChunk{
			MB: true,
			ME: true,
			CF: true,
		},
	}
	errs = append(errs, r.check())

	// Last record is Chunked
	r.chunks = []*recordChunk{
		&recordChunk{
			MB: true,
			ME: false,
			CF: true,
		},
		&recordChunk{
			MB: false,
			ME: true,
			CF: true,
		},
	}
	errs = append(errs, r.check())

	// recordChunk missing CF
	r.chunks = []*recordChunk{
		&recordChunk{
			MB: true,
			ME: false,
			CF: false,
		},
		&recordChunk{
			MB: false,
			ME: true,
			CF: false,
		},
	}
	errs = append(errs, r.check())

	// Non-first record with IL
	r.chunks = []*recordChunk{
		&recordChunk{
			MB:       true,
			ME:       false,
			CF:       true,
			IL:       true,
			IDLength: 1,
			ID:       "a",
		},
		&recordChunk{
			MB:       false,
			ME:       true,
			CF:       false,
			IL:       true,
			IDLength: 1,
			ID:       "a",
		},
	}
	errs = append(errs, r.check())

	// Non-first record with TypeLength
	r.chunks = []*recordChunk{
		&recordChunk{
			MB:         true,
			ME:         false,
			CF:         true,
			TypeLength: 1,
			Type:       "U",
		},
		&recordChunk{
			MB:         false,
			ME:         true,
			CF:         false,
			TypeLength: 1,
			Type:       "U",
		},
	}
	errs = append(errs, r.check())

	// Non-first record with BAD TNF
	r.chunks = []*recordChunk{
		&recordChunk{
			MB:         true,
			ME:         false,
			CF:         true,
			TypeLength: 1,
			Type:       "U",
			TNF:        Empty,
		},
		&recordChunk{
			MB:         false,
			ME:         true,
			CF:         false,
			TypeLength: 0,
			TNF:        Unknown,
		},
	}
	errs = append(errs, r.check())

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

func TestNDEFGoodrecordChunkTest(t *testing.T) {
	// Non-first record with TypeLength
	chunks := []*recordChunk{
		&recordChunk{
			MB:            true,
			ME:            false,
			CF:            true,
			SR:            true,
			TNF:           NFCForumWellKnownType,
			TypeLength:    1,
			Type:          "U",
			PayloadLength: 2,
			Payload:       []byte("\x00a"),
		},
		&recordChunk{
			MB:            false,
			ME:            false,
			CF:            true,
			SR:            true,
			TypeLength:    0,
			TNF:           Unchanged,
			PayloadLength: 2,
			Payload:       []byte("bc"),
		},
		&recordChunk{
			MB:            false,
			ME:            true,
			CF:            false,
			SR:            true,
			TypeLength:    0,
			TNF:           Unchanged,
			PayloadLength: 1,
			Payload:       []byte("d"),
		},
	}
	r := &Record{chunks: chunks}
	err := r.check()
	if err != nil {
		t.Log("recordChunk was good but failed because:", err)
		t.FailNow()
	}

	// Since we are here, test that we can reparse correctly
	var buf bytes.Buffer
	for _, c := range chunks {
		cBytes, err := c.Marshal()
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		t.Logf("% 02x", cBytes)
		buf.Write(cBytes)
	}

	r = &Record{}
	_, err = r.Unmarshal(buf.Bytes())
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	pl, err := r.Payload()
	if err != nil {
		t.Fatal(err)
	}
	if pl.String() != "abcd" {
		t.Error("Payload is not what we would expect!")
	}
}
