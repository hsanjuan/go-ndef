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

func TestBytesToUint64(t *testing.T) {
	cases := map[uint64][]byte{
		1: []byte{0, 1},
		2: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		3: []byte{10, 12, 30, 1, 0, 0, 0, 0, 0, 0, 0, 3},
	}
	for expected, test := range cases {
		result := BytesToUint64(test)
		if result != expected {
			t.Error("Result/Expected:", result, expected)
		}
	}
}

func TestUint64ToBytes(t *testing.T) {
	cases := []struct {
		Number   uint64
		L        int
		Expected []byte
	}{
		{1, 5, []byte{0, 0, 0, 0, 1}},
		{2, 1, []byte{2}},
		{10, 10, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 10}},
	}
	for _, tc := range cases {
		result := Uint64ToBytes(tc.Number, tc.L)
		if !bytes.Equal(result, tc.Expected) {
			t.Error("Result/Expected:", result, tc.Expected)
		}
	}
}
