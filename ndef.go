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

// Package ndef provides an implementation of the NFC Data Exchange
// Format (NDEF) specification.
//
// It allows to parse byte slices into a structured Message type,
// as well as to turn an Message into a byte slice.
package ndef

// Possible values for the TNF Field as defined in the specification.
const (
	Empty = byte(iota)
	NFCForumWellKnownType
	MediaType
	AbsoluteURI
	NFCForumExternalType
	Unknown
	Unchanged
	Reserved
)

// URIProtocols provides a mapping between the first byte if a NDEF Payload of
// type "U" (URI) and the string value for the protocol
var URIProtocols = map[byte]string{
	0:  "",
	1:  "http://www.",
	2:  "https://www.",
	3:  "http://",
	4:  "https://",
	5:  "tel:",
	6:  "mailto:",
	7:  "ftp://anonymous:anonymous@",
	8:  "ftp://ftp.",
	9:  "ftps://",
	10: "sftp://",
	11: "smb://",
	12: "nfs://",
	13: "ftp://",
	14: "dev://",
	15: "news:",
	16: "telnet://",
	17: "imap:",
	18: "rtsp://",
	19: "urn:",
	20: "pop:",
	21: "sip:",
	22: "sips:",
	23: "tftp:",
	24: "btspp://",
	25: "btl2cap://",
	26: "btgoep://",
	27: "tcpobex://",
	28: "irdaobex://",
	29: "file://",
	30: "urn:epc:id:",
	31: "urn:epc:tag:",
	32: "urn:epc:pat:",
	33: "urn:epc:raw:",
	34: "urn:epc:",
	35: "urn:nfc:",
}
