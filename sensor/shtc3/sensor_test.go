package shtc3

import "testing"

func TestIdent_SHTC3(t *testing.T) {
	id := Ident(2183)
	if !id.SHTC3() {
		t.Fatal(id)
	}

	id = Ident(1234)
	if id.SHTC3() {
		t.Fatal(id)
	}
}
