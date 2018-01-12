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
	log "github.com/sirupsen/logrus"
)

var _ = state.MakeMap

//General acknowledge message
type TimeoutInternal struct {
	msgbase.MessageBase
	Name        string
	DBHeight    int
	Minute      int
	Round       int
	MessageHash interfaces.IHash
}

var _ interfaces.IMsg = (*TimeoutInternal)(nil)

func (m *TimeoutInternal) ElectionProcess(is interfaces.IState, elect interfaces.IElections) {
	s := is.(*state.State)

	e, ok := elect.(*elections.Elections)
	if !ok {
		panic("Invalid elections object")
	}

	aidx := e.AuditIndex(e.FedID)

	// We have advanced, so do nothing.  We can't reset anything because there
	// can be a timeout process that started before we got here (with short minutes)
	if e.DBHeight > m.DBHeight || e.Minute > m.Minute {
		return
	}

	cnt := 0
	e.Electing = -1
	for i, b := range e.Sync {
		if !b {
			cnt++
			if e.Electing < 0 {
				e.Electing = i
			}
		}
	}
	// Hey, if all is well, then continue.
	if e.Electing < 0 {
		return
	}

	e.State.(*state.State).Election2 = e.FeedBackStr("E", true, e.Electing)

	for len(e.Round) <= e.Electing {
		e.Round = append(e.Round, 0)
	}

	// New timeout, new round of elections.
	e.Round[e.Electing]++

	// If we don't have all our sync messages, we will have to come back around and see if all is well.
	go Fault(e, int(m.DBHeight), int(m.Minute), e.Round[e.Electing])

	// Can we see a majority of the federated servers?
	if cnt >= (len(e.Federated)+1)/2 {
		// Reset the timeout and give up if we can't see a majority.
		return
	}

	auditIdx := e.AuditPriority()

	if aidx >= 0 {
		serverMap := state.MakeMap(len(e.Federated), uint32(e.DBHeight))
		vm := state.FedServerVM(serverMap, len(e.Federated), e.Minute, e.Electing)

		if aidx == auditIdx {
			Sync := new(SyncMsg)
			Sync.SetLocal(true)
			Sync.VMIndex = vm
			Sync.TS = primitives.NewTimestampNow()
			Sync.Name = e.Name
			Sync.FedIdx = uint32(e.Electing)
			Sync.FedID = e.FedID
			Sync.Weight = e.APriority[auditIdx]
			Sync.DBHeight = uint32(e.DBHeight)
			Sync.Minute = byte(e.Minute)
			Sync.Round = e.Round[e.Electing]
			s.InMsgQueue().Enqueue(Sync)
			s.Election2 = e.FeedBackStr(fmt.Sprintf("%d", e.Round[e.Electing]), false, auditIdx)
		}
	}

	if aidx != auditIdx {
		s.Election2 = e.FeedBackStr(fmt.Sprintf("%d-%d", e.Round[e.Electing], auditIdx), true, e.Electing)
	}

}

func (m *TimeoutInternal) GetServerID() interfaces.IHash {
	return nil
}

func (m *TimeoutInternal) LogFields() log.Fields {
	return log.Fields{"category": "message", "messagetype": "TimeoutInternal", "dbheight": m.DBHeight}
}

func (m *TimeoutInternal) GetRepeatHash() interfaces.IHash {
	return m.GetMsgHash()
}

// We have to return the haswh of the underlying message.
func (m *TimeoutInternal) GetHash() interfaces.IHash {
	return m.MessageHash
}

func (m *TimeoutInternal) GetTimestamp() interfaces.Timestamp {
	return primitives.NewTimestampNow()
}

func (m *TimeoutInternal) GetMsgHash() interfaces.IHash {
	if m.MsgHash == nil {
	}
	return m.MsgHash
}

func (m *TimeoutInternal) Type() byte {
	return constants.INTERNALTIMEOUT
}

func (m *TimeoutInternal) Validate(state interfaces.IState) int {
	return 1
}

// Returns true if this is a message for this server to execute as
// a leader.
func (m *TimeoutInternal) ComputeVMIndex(state interfaces.IState) {
}

// Execute the leader functions of the given message
// Leader, follower, do the same thing.
func (m *TimeoutInternal) LeaderExecute(state interfaces.IState) {
	m.FollowerExecute(state)
}

func (m *TimeoutInternal) FollowerExecute(state interfaces.IState) {

}

// Acknowledgements do not go into the process list.
func (e *TimeoutInternal) Process(dbheight uint32, state interfaces.IState) bool {
	panic("Ack object should never have its Process() method called")
}

func (e *TimeoutInternal) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *TimeoutInternal) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (m *TimeoutInternal) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling: %v", r)
		}
	}()
	return
}

func (m *TimeoutInternal) UnmarshalBinary(data []byte) error {
	_, err := m.UnmarshalBinaryData(data)
	return err
}

func (m *TimeoutInternal) MarshalBinary() (data []byte, err error) {
	return
}

func (m *TimeoutInternal) String() string {
	if m.LeaderChainID == nil {
		m.LeaderChainID = primitives.NewZeroHash()
	}
	return fmt.Sprintf(" %20s %10s dbheight %d minute %d",
		m.Name,
		"Time Out",
		m.DBHeight,
		m.Minute)
}

func (a *TimeoutInternal) IsSameAs(b *TimeoutInternal) bool {
	return true
}
