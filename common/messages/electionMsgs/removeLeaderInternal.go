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
	log "github.com/FactomProject/logrus"
)

//General acknowledge message
type RemoveLeaderInternal struct {
	msgbase.MessageBase
	NName       string
	ServerID    interfaces.IHash // Hash of message acknowledged
	DBHeight    uint32           // Directory Block Height that owns this ack
	Height      uint32           // Height of this ack in this process list
	MessageHash interfaces.IHash
}

var _ interfaces.IMsg = (*RemoveLeaderInternal)(nil)

func (m *RemoveLeaderInternal) ElectionProcess(state interfaces.IState, elect interfaces.IElections) {
	e, ok := elect.(*elections.Elections)
	if !ok {
		panic("Invalid elections object")
	}
	idx := 0
	for i, s := range e.Federated {
		idx = i
		if s.GetChainID().IsSameAs(m.GetServerID()) {
			break
		}
	}
	if idx < len(e.Federated) {
		e.Federated = append(e.Federated[:idx], e.Federated[idx+1:]...)
	}
}

func (m *RemoveLeaderInternal) GetServerID() interfaces.IHash {
	return m.ServerID
}

func (m *RemoveLeaderInternal) LogFields() log.Fields {
	return log.Fields{"category": "message", "messagetype": "RemoveLeaderInternal", "dbheight": m.DBHeight, "newleader": m.ServerID.String()[4:12]}
}

func (m *RemoveLeaderInternal) GetRepeatHash() interfaces.IHash {
	return m.GetMsgHash()
}

// We have to return the haswh of the underlying message.
func (m *RemoveLeaderInternal) GetHash() interfaces.IHash {
	return m.MessageHash
}

func (m *RemoveLeaderInternal) GetTimestamp() interfaces.Timestamp {
	return primitives.NewTimestampNow()
}

func (m *RemoveLeaderInternal) GetMsgHash() interfaces.IHash {
	if m.MsgHash == nil {
	}
	return m.MsgHash
}

func (m *RemoveLeaderInternal) Type() byte {
	return constants.INTERNALREMOVELEADER
}

func (m *RemoveLeaderInternal) Validate(state interfaces.IState) int {
	return 1
}

// Returns true if this is a message for this server to execute as
// a leader.
func (m *RemoveLeaderInternal) ComputeVMIndex(state interfaces.IState) {
}

// Execute the leader functions of the given message
// Leader, follower, do the same thing.
func (m *RemoveLeaderInternal) LeaderExecute(state interfaces.IState) {
	m.FollowerExecute(state)
}

func (m *RemoveLeaderInternal) FollowerExecute(state interfaces.IState) {

}

// Acknowledgements do not go into the process list.
func (e *RemoveLeaderInternal) Process(dbheight uint32, state interfaces.IState) bool {
	panic("Ack object should never have its Process() method called")
}

func (e *RemoveLeaderInternal) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *RemoveLeaderInternal) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (m *RemoveLeaderInternal) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling: %v", r)
		}
	}()
	return
}

func (m *RemoveLeaderInternal) UnmarshalBinary(data []byte) error {
	_, err := m.UnmarshalBinaryData(data)
	return err
}

func (m *RemoveLeaderInternal) MarshalBinary() (data []byte, err error) {
	return
}

func (m *RemoveLeaderInternal) String() string {
	if m.LeaderChainID == nil {
		m.LeaderChainID = primitives.NewZeroHash()
	}
	return fmt.Sprintf(" %20s %x %10s dbheight %d", "Remove Leader Internal", m.ServerID.Bytes(), m.NName, m.DBHeight)
}

func (a *RemoveLeaderInternal) IsSameAs(b *RemoveLeaderInternal) bool {
	return true
}
