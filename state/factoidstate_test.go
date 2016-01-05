// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package state_test

import (
	"github.com/FactomProject/factomd/common/entryCreditBlock"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	. "github.com/FactomProject/factomd/state"
	"testing"
)

var state interfaces.IFactoidState

func TestBalances(t *testing.T) {
	s := new(State)
	s.Init("")
	state = s.GetFactoidState()
	state.SetFactoshisPerEC(1)
	add1, err := primitives.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		t.Error(err)
	}
	add2, err := primitives.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		t.Error(err)
	}
	add3, err := primitives.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
	if err != nil {
		t.Error(err)
	}

	tx := new(factoid.Transaction)
	tx.AddOutput(add1, 1000000)

	err = state.UpdateTransaction(tx)
	if err != nil {
		t.Error(err)
	}

	if state.GetFactoidBalance(add1.Fixed()) != 1000000 {
		t.Errorf("Invalid address balance - %v", state.GetFactoidBalance(add1.Fixed()))
	}
	if state.GetECBalance(add1.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add1.Fixed()))
	}
	if state.GetFactoidBalance(add2.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetFactoidBalance(add2.Fixed()))
	}
	if state.GetECBalance(add2.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add2.Fixed()))
	}
	if state.GetFactoidBalance(add3.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetFactoidBalance(add3.Fixed()))
	}
	if state.GetECBalance(add3.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add3.Fixed()))
	}

	tx = new(factoid.Transaction)
	tx.AddInput(add1, 1000)
	tx.AddOutput(add2, 1000)

	err = state.UpdateTransaction(tx)
	if err != nil {
		t.Error(err)
	}

	if state.GetFactoidBalance(add1.Fixed()) != 999000 {
		t.Errorf("Invalid address balance - %v", state.GetFactoidBalance(add1.Fixed()))
	}
	if state.GetECBalance(add1.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add1.Fixed()))
	}
	if state.GetFactoidBalance(add2.Fixed()) != 1000 {
		t.Errorf("Invalid address balance - %v", state.GetFactoidBalance(add2.Fixed()))
	}
	if state.GetECBalance(add2.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add2.Fixed()))
	}
	if state.GetFactoidBalance(add3.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetFactoidBalance(add3.Fixed()))
	}
	if state.GetECBalance(add3.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add3.Fixed()))
	}

	tx = new(factoid.Transaction)
	tx.AddInput(add1, 1000)
	tx.AddECOutput(add3, 1000)

	err = state.UpdateTransaction(tx)
	if err != nil {
		t.Error(err)
	}

	if state.GetFactoidBalance(add1.Fixed()) != 998000 {
		t.Errorf("Invalid address balance - %v", state.GetFactoidBalance(add1.Fixed()))
	}
	if state.GetECBalance(add1.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add1.Fixed()))
	}
	if state.GetFactoidBalance(add2.Fixed()) != 1000 {
		t.Errorf("Invalid address balance - %v", state.GetFactoidBalance(add2.Fixed()))
	}
	if state.GetECBalance(add2.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add2.Fixed()))
	}
	if state.GetFactoidBalance(add3.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetFactoidBalance(add3.Fixed()))
	}
	if state.GetECBalance(add3.Fixed()) != 1000 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add3.Fixed()))
	}

	state.ResetBalances()

	if state.GetFactoidBalance(add1.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetFactoidBalance(add1.Fixed()))
	}
	if state.GetECBalance(add1.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add1.Fixed()))
	}
	if state.GetFactoidBalance(add2.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetFactoidBalance(add2.Fixed()))
	}
	if state.GetECBalance(add2.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add2.Fixed()))
	}
	if state.GetFactoidBalance(add3.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetFactoidBalance(add3.Fixed()))
	}
	if state.GetECBalance(add3.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add3.Fixed()))
	}
}

func TestUpdateECTransaction(t *testing.T) {
	state.SetFactoshisPerEC(1)
	add1, err := primitives.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		t.Error(err)
		return
	}
	add1bs := primitives.StringToByteSlice32("0000000000000000000000000000000000000000000000000000000000000001")

	if state.GetECBalance(add1.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add1.Fixed()))
		return
	}

	var tx interfaces.IECBlockEntry
	tx = new(entryCreditBlock.ServerIndexNumber)

	err = state.UpdateECTransaction(tx)
	if err != nil {
		t.Error(err)
		return
	}
	if state.GetECBalance(add1.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add1.Fixed()))
	}

	tx = new(entryCreditBlock.MinuteNumber)

	err = state.UpdateECTransaction(tx)
	if err != nil {
		t.Error(err)
		return
	}
	if state.GetECBalance(add1.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add1.Fixed()))
		return
	}

	//Proper processing
	cc := new(entryCreditBlock.CommitChain)
	cc.ECPubKey = add1bs
	cc.Credits = 100
	tx = cc

	err = state.UpdateECTransaction(tx)
	if err != nil {
		t.Error(err)
		return
	}
	if state.GetECBalance(add1.Fixed()) != -100 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add1.Fixed()))
		return
	}

	ib := new(entryCreditBlock.IncreaseBalance)
	ib.ECPubKey = add1bs
	ib.NumEC = 100
	tx = ib

	err = state.UpdateECTransaction(tx)
	if err != nil {
		t.Error(err)
		return
	}
	if state.GetECBalance(add1.Fixed()) != 0 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add1.Fixed()))
		return
	}

	ce := new(entryCreditBlock.CommitEntry)
	ce.ECPubKey = add1bs
	ce.Credits = 100
	tx = ce

	err = state.UpdateECTransaction(tx)
	if err != nil {
		t.Error(err)
		return
	}
	if state.GetECBalance(add1.Fixed()) != -100 {
		t.Errorf("Invalid address balance - %v", state.GetECBalance(add1.Fixed()))
		return
	}

}

/*
import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/FactomProject/ed25519"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	. "github.com/FactomProject/factomd/database/boltdb"
	"math/rand"
	"testing"
)

var _ = hex.EncodeToString
var _ = fmt.Printf
var _ = ed25519.Sign
var _ = rand.New
var _ = binary.Write
var _ = Prtln

func GetDatabase() interfaces. {
	var bucketList [][]byte

	bucketList = make([][]byte, 5, 5)

	bucketList[0] = []byte("factoidAddress_balances")
	bucketList[0] = []byte("factoidOrphans_balances")
	bucketList[0] = []byte("factomAddress_balances")

	db := new(BoltDB)

	db.Init(bucketList, "/tmp/fs_test.db")

	return db
}

func Test_updating_balances_FactoidState(test *testing.T) {
	fs := new(FactoidState)
	fs.database = GetDatabase()

}*/
