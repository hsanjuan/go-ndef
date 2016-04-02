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
	"encoding/binary"
	"fmt"
)

// Parse some bytes into an uint64. If the slice is longer that 8 bytes, it's truncated.
func BytesToUint64(b []byte) uint64 {
	// Make sure we are not parsing more than 8 bytes (uint64 size)
	byte8 := make([]byte, 8)
	if len(b) > 8 {
		copy(byte8, b[len(b)-8:]) // use the last 8 bytes
	} else {
		copy(byte8[8-len(b):], b) // copy shorter arrays in the last positions of byte8
	}
	return binary.BigEndian.Uint64(byte8)
}

// Given an Uint64, return a slice of bytes.
func Uint64ToBytes(n uint64, desired_len int) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, n)
	if desired_len >= 8 {
		slice := make([]byte, desired_len)
		copy(slice[desired_len-8:], buf.Bytes())
		return slice
	} else {
		return buf.Bytes()[8-desired_len:]
	}
}

// func PrintBytes(bytes []byte, length int) {
// 	for i := 0; i < length; i++ {
// 		fmt.Printf("%02x ", bytes[i])
// 	}
// 	fmt.Println()
// }

func FmtBytes(bytes []byte, length int) (str string) {
	for i := 0; i < length; i++ {
		str += fmt.Sprintf("%02x ", bytes[i])
	}
	return str
}
