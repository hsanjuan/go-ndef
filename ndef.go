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

/*

This is an implementation of the NFC Data Exchange Format (NDEF) standard.

It allows to parse byte slices into a structured Message type,
as well as to turn an Message into a byte slice.

See ndef_message.go and ndef_record.go for details.

*/

// Possible values for the TNF Field
const (
	EMPTY = byte(iota)
	NFC_FORUM_WELL_KNOWN_TYPE
	MEDIA_TYPE
	ABSOLUTE_URI
	NFC_FORUM_EXTERNAL_TYPE
	UNKNOWN
	UNCHANGED
	RESERVED
)

// Given an URI identifier code (the first byte of a NDEF Payload of type 'U'), return the meaning
func URIProtocols(uri_identifier_code byte) string {
	switch uri_identifier_code {
	case 0:
		return ""
	case 1:
		return "http://www."
	case 2:
		return "https://www."
	case 3:
		return "http://"
	case 4:
		return "https://"
	case 5:
		return "tel:"
	case 6:
		return "mailto:"
	case 7:
		return "ftp://anonymous:anonymous@"
	case 8:
		return "ftp://ftp."
	case 9:
		return "ftps://"
	case 10:
		return "sftp://"
	case 11:
		return "smb://"
	case 12:
		return "nfs://"
	case 13:
		return "ftp://"
	case 14:
		return "dev://"
	case 15:
		return "news:"
	case 16:
		return "telnet://"
	case 17:
		return "imap:"
	case 18:
		return "rtsp://"
	case 19:
		return "urn:"
	case 20:
		return "pop:"
	case 21:
		return "sip:"
	case 22:
		return "sips:"
	case 23:
		return "tftp:"
	case 24:
		return "btspp://"
	case 25:
		return "btl2cap://"
	case 26:
		return "btgoep://"
	case 27:
		return "tcpobex://"
	case 28:
		return "irdaobex://"
	case 29:
		return "file://"
	case 30:
		return "urn:epc:id:"
	case 31:
		return "urn:epc:tag:"
	case 32:
		return "urn:epc:pat:"
	case 33:
		return "urn:epc:raw:"
	case 34:
		return "urn:epc:"
	case 35:
		return "urn:nfc:"
	default:
		return "RFU"
	}

}
