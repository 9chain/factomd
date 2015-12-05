// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factoid_test

import (
	"encoding/hex"
	"github.com/FactomProject/ed25519"
	. "github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/testHelper"
	"testing"
)

func TestSignatureBlock(t *testing.T) {
	priv := testHelper.NewPrivKey(0)
	testData, err := hex.DecodeString("00112233445566778899")
	if err != nil {
		t.Error(err)
	}

	sig := NewSingleSignatureBlock(priv, testData)

	rcd := testHelper.NewFactoidRCDAddress(0)
	pub := rcd.(*RCD_1).GetPublicKey()
	pub2 := [32]byte{}
	copy(pub2[:], pub)

	s := sig.Signatures[0].(*FactoidSignature).Signature
	valid := ed25519.VerifyCanonical(&pub2, testData, &s)
	t.Errorf("Valid - %v", valid)
}
