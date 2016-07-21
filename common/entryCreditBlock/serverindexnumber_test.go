// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package entryCreditBlock_test

import (
	. "github.com/FactomProject/factomd/common/entryCreditBlock"
	"testing"
)

func TestServerIndexMarshalUnmarshal(t *testing.T) {
	si1 := NewServerIndexNumber()
	si1.Number = 3
	b, err := si1.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	if len(b) != 1 {
		t.Error("Invalid byte length")
	}
	if b[0] != 3 {
		t.Error("Invalid byte")
	}

	si2 := NewServerIndexNumber()
	err = si2.UnmarshalBinary(b)
	if err != nil {
		t.Error(err)
	}
	if si1.Number != si2.Number {
		t.Error("Invalid data unmarshalled")
	}
}

func TestServerIndexNumberMisc(t *testing.T) {
	si := NewServerIndexNumber()
	si.Number = 4
	if si.IsInterpretable() == false {
		t.Fail()
	}
	if si.Interpret() != "ServerIndexNumber 4" {
		t.Fail()
	}
	if si.Hash().String() != "e52d9c508c502347344d8c07ad91cbd6068afc75ff6292f062a09ca381c89e71" {
		t.Fail()
	}
}
