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

	"github.com/hsanjuan/go-ndef/types/generic"
)

func TestRecordMarshalUnmarshal(t *testing.T) {
	t.Log("Testing a Record created with a provided chunk")
	r := &Record{
		TNF:  NFCForumExternalType,
		Type: "test",
		ID:   "#ab",
		Payload: &types.Generic{
			Payload: []byte("abc"),
		},
	}

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

	m := NewTextRecord("abc", "en")
	t.Log(m)

	m = NewURIRecord("http://abc")
	t.Log(m)

	m = &Record{
		TNF:  NFCForumWellKnownType,
		ID:   "#ab",
		Type: "X",
		Payload: &types.Generic{
			Payload: []byte("abc"),
		},
	}
	t.Log(m)

	m = &Record{
		TNF: Empty,
	}
	t.Log(m)

	m = &Record{
		TNF:  MediaType,
		Type: "image/jpeg",
		Payload: &types.Generic{
			Payload: []byte("\x03abc"),
		},
	}
	t.Log(m)

	m = &Record{
		TNF:  AbsoluteURI,
		ID:   "#ab",
		Type: "http://resource",
		Payload: &types.Generic{
			Payload: []byte(""),
		},
	}
	t.Log(m)

	m = &Record{
		TNF:  NFCForumExternalType,
		ID:   "#ab",
		Type: "T",
		Payload: &types.Generic{
			Payload: []byte("abc"),
		},
	}
	t.Log(m)

	m = &Record{
		TNF: Unknown,
		Payload: &types.Generic{
			Payload: []byte("abc"),
		},
	}
	t.Log(m)

	m = &Record{
		TNF: Unchanged,
		Payload: &types.Generic{
			Payload: []byte("abc"),
		},
	}
	t.Log(m)
}

func TestRecordBadChunksTest(t *testing.T) {
	cases := []struct{ Expected string }{
		{eNORECORDS},
		{eNOMB},
		{eFIRSTCHUNKED},
		{eNOME},
		{eLASTCHUNKED},
		{eCFMISSING},
		{eBADIL},
		{eBADTYPELENGTH},
		{eBADTNF},
	}

	errs := []error{}
	chunks := []*recordChunk{}
	errs = append(errs, checkChunks(chunks))

	// First record is not MB
	chunks = []*recordChunk{
		&recordChunk{
			MB: false,
			ME: true,
			CF: true,
		},
	}
	errs = append(errs, checkChunks(chunks))

	// First and only record is chuncked
	chunks = []*recordChunk{
		&recordChunk{
			MB: true,
			ME: true,
			CF: true,
		},
	}
	errs = append(errs, checkChunks(chunks))

	// Last record is not ME
	chunks = []*recordChunk{
		&recordChunk{
			MB: true,
			ME: false,
		},
	}
	errs = append(errs, checkChunks(chunks))

	// Last record is Chunked
	chunks = []*recordChunk{
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
	errs = append(errs, checkChunks(chunks))

	// recordChunk missing CF
	chunks = []*recordChunk{
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
	errs = append(errs, checkChunks(chunks))

	// Non-first record with IL
	chunks = []*recordChunk{
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
	errs = append(errs, checkChunks(chunks))

	// Non-first record with TypeLength
	chunks = []*recordChunk{
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
	errs = append(errs, checkChunks(chunks))

	// Non-first record with BAD TNF
	chunks = []*recordChunk{
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
	errs = append(errs, checkChunks(chunks))

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
	err := checkChunks(chunks)
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

	r := new(Record)
	_, err = r.Unmarshal(buf.Bytes())
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if r.Payload.String() != "abcd" {
		t.Error("Payload is not what we would expect!")
	}
}
