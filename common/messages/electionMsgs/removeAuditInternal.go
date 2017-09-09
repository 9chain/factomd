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
	log "github.com/FactomProject/logrus"
	"github.com/FactomProject/factomd/elections"
)

//General acknowledge message
type RemoveAuditInternal struct {
	msgbase.MessageBase
	NName       string
	ServerID    interfaces.IHash // Hash of message acknowledged
	DBHeight    uint32           // Directory Block Height that owns this ack
	Height      uint32           // Height of this ack in this process list
	MessageHash interfaces.IHash
}

var _ interfaces.IMsg = (*RemoveAuditInternal)(nil)

func (m *RemoveAuditInternal) ElectionProcess(state interfaces.IState, elections interfaces.IElections) {
	e, ok := elections.(*elections.Elections)
	if !ok {
		panic("Invalid elections object")
	}
}

func (m *RemoveAuditInternal) GetServerID() interfaces.IHash {
	return m.ServerID
}

func (m *RemoveAuditInternal) LogFields() log.Fields {
	return log.Fields{"category": "message", "messagetype": "RemoveAuditInternal", "dbheight": m.DBHeight, "newleader": m.ServerID.String()[4:12]}
}

func (m *RemoveAuditInternal) GetRepeatHash() interfaces.IHash {
	return m.GetMsgHash()
}

// We have to return the haswh of the underlying message.
func (m *RemoveAuditInternal) GetHash() interfaces.IHash {
	return m.MessageHash
}

func (m *RemoveAuditInternal) GetTimestamp() interfaces.Timestamp {
	return primitives.NewTimestampNow()
}

func (m *RemoveAuditInternal) GetMsgHash() interfaces.IHash {
	if m.MsgHash == nil {
	}
	return m.MsgHash
}

func (m *RemoveAuditInternal) Type() byte {
	return constants.INTERNALREMOVEAUDIT
}

func (m *RemoveAuditInternal) Validate(state interfaces.IState) int {
	return 1
}

// Returns true if this is a message for this server to execute as
// a leader.
func (m *RemoveAuditInternal) ComputeVMIndex(state interfaces.IState) {
}

// Execute the leader functions of the given message
// Leader, follower, do the same thing.
func (m *RemoveAuditInternal) LeaderExecute(state interfaces.IState) {
	m.FollowerExecute(state)
}

func (m *RemoveAuditInternal) FollowerExecute(state interfaces.IState) {

}

// Acknowledgements do not go into the process list.
func (e *RemoveAuditInternal) Process(dbheight uint32, state interfaces.IState) bool {
	panic("Ack object should never have its Process() method called")
}

func (e *RemoveAuditInternal) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *RemoveAuditInternal) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (m *RemoveAuditInternal) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling: %v", r)
		}
	}()
	return
}

func (m *RemoveAuditInternal) UnmarshalBinary(data []byte) error {
	_, err := m.UnmarshalBinaryData(data)
	return err
}

func (m *RemoveAuditInternal) MarshalBinary() (data []byte, err error) {
	return
}

func (m *RemoveAuditInternal) String() string {
	if m.LeaderChainID == nil {
		m.LeaderChainID = primitives.NewZeroHash()
	}
	return fmt.Sprintf("%20s %x %10s dbheight %d", "Remove Audit Internal", m.ServerID.Bytes(), m.NName, m.DBHeight)
}

func (a *RemoveAuditInternal) IsSameAs(b *RemoveAuditInternal) bool {
	return true
}
