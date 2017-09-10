// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package electionMsgs

import (
	"fmt"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/messages/msgbase"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/elections"
	"github.com/FactomProject/factomd/state"
	log "github.com/FactomProject/logrus"
)

//General acknowledge message
type AddLeaderInternal struct {
	msgbase.MessageBase
	NName       string
	ServerID    interfaces.IHash // Hash of message acknowledged
	DBHeight    uint32           // Directory Block Height that owns this ack
	Height      uint32           // Height of this ack in this process list
	MessageHash interfaces.IHash
}

var _ interfaces.IMsg = (*AddLeaderInternal)(nil)

func (m *AddLeaderInternal) ElectionProcess(s interfaces.IState, elect interfaces.IElectionMsg) {
	e, ok := elect.(*elections.Elections)
	if !ok {
		panic("Invalid elections object")
	}
	e.Federated = append(e.Federated, &state.Server{ChainID: m.ServerID, Online: true})
	Sort(e.Federated)
}

func (m *AddLeaderInternal) GetServerID() interfaces.IHash {
	return m.ServerID
}

func (m *AddLeaderInternal) LogFields() log.Fields {
	return log.Fields{"category": "message", "messagetype": "addleaderinternal", "dbheight": m.DBHeight, "newleader": m.ServerID.String()[4:12]}
}

func (m *AddLeaderInternal) GetRepeatHash() interfaces.IHash {
	return m.GetMsgHash()
}

// We have to return the haswh of the underlying message.
func (m *AddLeaderInternal) GetHash() interfaces.IHash {
	return m.MessageHash
}

func (m *AddLeaderInternal) GetTimestamp() interfaces.Timestamp {
	return primitives.NewTimestampNow()
}

func (m *AddLeaderInternal) GetMsgHash() interfaces.IHash {
	if m.MsgHash == nil {
	}
	return m.MsgHash
}

func (m *AddLeaderInternal) Type() byte {
	return constants.INTERNALADDLEADER
}

func (m *AddLeaderInternal) Validate(state interfaces.IState) int {
	return 1
}

// Returns true if this is a message for this server to execute as
// a leader.
func (m *AddLeaderInternal) ComputeVMIndex(state interfaces.IState) {
}

// Execute the leader functions of the given message
// Leader, follower, do the same thing.
func (m *AddLeaderInternal) LeaderExecute(state interfaces.IState) {
	m.FollowerExecute(state)
}

func (m *AddLeaderInternal) FollowerExecute(state interfaces.IState) {

}

// Acknowledgements do not go into the process list.
func (e *AddLeaderInternal) Process(dbheight uint32, state interfaces.IState) bool {
	panic("Ack object should never have its Process() method called")
}

func (e *AddLeaderInternal) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *AddLeaderInternal) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (m *AddLeaderInternal) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling: %v", r)
		}
	}()
	return
}

func (m *AddLeaderInternal) UnmarshalBinary(data []byte) error {
	_, err := m.UnmarshalBinaryData(data)
	return err
}

func (m *AddLeaderInternal) MarshalBinary() (data []byte, err error) {
	return
}

func (m *AddLeaderInternal) String() string {
	if m.LeaderChainID == nil {
		m.LeaderChainID = primitives.NewZeroHash()
	}
	return fmt.Sprintf("%20s %x %10s dbheight %d", "Add Leader Internal", m.ServerID.Bytes(), m.NName, m.DBHeight)
}

func (a *AddLeaderInternal) IsSameAs(b *AddLeaderInternal) bool {
	return true
}
