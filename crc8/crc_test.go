package crc8

import "testing"

func TestCRC8(t *testing.T) {
	if err := Checksum(Default, []byte{92, 93}, 224); err != nil {
		t.Fatal(err)
	}
}
