// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factoid

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

/************************
 * RCD 2
 ************************/

// Type 2 RCD implement multisig
// m of n
// Must have m addresses from which to choose, no fewer, no more
// Must have n RCD, no fewer no more.
// NOTE: This does mean you can have a multisig nested in a
// multisig.  It just works.

type RCD_2 struct {
	M           int                   // Number signatures required
	N           int                   // Total sigatures possible
	N_Addresses []interfaces.IAddress // n addresses
}

var _ interfaces.IRCD = (*RCD_2)(nil)

/*************************************
 *       Stubs
 *************************************/

func (b RCD_2) GetAddress() (interfaces.IAddress, error) {
	return nil, nil
}

func (b RCD_2) GetHash() interfaces.IHash {
	return nil
}

func (b RCD_2) NumberOfSignatures() int {
	return 1
}

/***************************************
 *       Methods
 ***************************************/

func (b RCD_2) UnmarshalBinary(data []byte) error {
	_, err := b.UnmarshalBinaryData(data)
	return err
}

func (b RCD_2) CheckSig(trans interfaces.ITransaction, sigblk interfaces.ISignatureBlock) bool {
	return false
}

func (e *RCD_2) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *RCD_2) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (e *RCD_2) JSONBuffer(b *bytes.Buffer) error {
	return primitives.EncodeJSONToBuffer(e, b)
}

func (b RCD_2) String() string {
	txt, err := b.CustomMarshalText()
	if err != nil {
		return "<error>"
	}
	return string(txt)
}

func (w RCD_2) Clone() interfaces.IRCD {
	c := new(RCD_2)
	c.M = w.M
	c.N = w.N
	c.N_Addresses = make([]interfaces.IAddress, len(w.N_Addresses))
	for i, address := range w.N_Addresses {
		c.N_Addresses[i] = CreateAddress(address)
	}
	return c
}

func (a1 *RCD_2) IsEqual(addr interfaces.IBlock) []interfaces.IBlock {
	a2, ok := addr.(*RCD_2)
	if !ok || // Not the right kind of interfaces.IBlock
		a1.N != a2.N || // Size of sig has to match
		a1.M != a2.M || // Size of sig has to match
		len(a1.N_Addresses) != len(a2.N_Addresses) { // Size of arrays has to match
		r := make([]interfaces.IBlock, 0, 5)
		return append(r, a1)
	}

	for i, addr := range a1.N_Addresses {
		r := addr.IsEqual(a2.N_Addresses[i])
		if r != nil {
			return append(r, a1)
		}
	}

	return nil
}

func (t *RCD_2) UnmarshalBinaryData(data []byte) (newData []byte, err error) {

	typ := int8(data[0])
	data = data[1:]
	if typ != 2 {
		return nil, fmt.Errorf("Bad data fed to RCD_2 UnmarshalBinaryData()")
	}

	t.N, data = int(binary.BigEndian.Uint16(data[0:2])), data[2:]
	t.M, data = int(binary.BigEndian.Uint16(data[0:2])), data[2:]

	t.N_Addresses = make([]interfaces.IAddress, t.M, t.M)

	for i, _ := range t.N_Addresses {
		t.N_Addresses[i] = new(Address)
		data, err = t.N_Addresses[i].UnmarshalBinaryData(data)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (a RCD_2) MarshalBinary() ([]byte, error) {
	var out bytes.Buffer

	binary.Write(&out, binary.BigEndian, uint8(2))
	binary.Write(&out, binary.BigEndian, uint16(a.N))
	binary.Write(&out, binary.BigEndian, uint16(a.M))
	for i := 0; i < a.M; i++ {
		data, err := a.N_Addresses[i].MarshalBinary()
		if err != nil {
			return nil, err
		}
		out.Write(data)
	}

	return out.Bytes(), nil
}

func (a RCD_2) CustomMarshalText() ([]byte, error) {
	var out bytes.Buffer

	primitives.WriteNumber8(&out, uint8(2)) // Type 2 Authorization
	out.WriteString("\n n: ")
	primitives.WriteNumber16(&out, uint16(a.N))
	out.WriteString(" m: ")
	primitives.WriteNumber16(&out, uint16(a.M))
	out.WriteString("\n")
	for i := 0; i < a.M; i++ {
		out.WriteString("  m: ")
		out.WriteString(hex.EncodeToString(a.N_Addresses[i].Bytes()))
		out.WriteString("\n")
	}

	return out.Bytes(), nil
}
