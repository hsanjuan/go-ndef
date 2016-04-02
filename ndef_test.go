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

import "testing"

func TestURIProtocols(t *testing.T) {
	for i := 0; i <= 36; i++ {
		r := URIProtocols(byte(i))
		t.Logf("URIProtocols(%d) -> %s", i, r)
		if r == "RFU" && i != 36 {
			t.Errorf("Error: URIProtocols(%d) -> %s", i, r)
		}
	}
}
