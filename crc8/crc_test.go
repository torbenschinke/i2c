package crc8

import "testing"

func TestCRC8(t *testing.T) {
	if err := Checksum(Default, []byte{92, 93}, 224); err != nil {
		t.Fatal(err)
	}

	x := [3]byte{92, 134, 14}
	if err := Checksum(Default, x[:2], x[2]); err != nil {
		t.Fatal(err)
	}
}
