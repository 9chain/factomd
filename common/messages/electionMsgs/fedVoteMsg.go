// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package electionMsgs

import (
	"bytes"
	"fmt"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/messages/msgbase"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/elections"
	"github.com/FactomProject/goleveldb/leveldb/errors"
	log "github.com/sirupsen/logrus"
	//"github.com/FactomProject/factomd/state"
	"time"
)

var _ = fmt.Print

// FedVoteMsg
// We vote on the Audit Server to replace a Federated Server that fails
// We vote to move to the next round, if the audit server fails.
// Could make these two messages, but for now we will do it in one.
type FedVoteMsg struct {
	msgbase.MessageBase
	TS       interfaces.Timestamp // Message Timestamp
	TypeMsg  byte                 // Can be either a Volunteer from an Audit Server, or End of round
	DBHeight uint32               // Directory Block Height that owns this ack
	Minute   byte                 // Minute (-1 for dbsig)
	Round    int                  // Voting Round
	Sigs     []interfaces.IHash   // Federated Server signatures.
	// End of Round Fields
	NewRound int // The new Round of elections (requires a majority)
	// Volunteer fields
	EOM        bool             // True if an EOM, false if a DBSig
	Name       string           // Server name
	FedIdx     uint32           // Server faulting
	FedID      interfaces.IHash // Server faulting
	ServerIdx  uint32           // Index of Server replacing
	ServerID   interfaces.IHash // Volunteer Server ChainID
	ServerName string           // Volunteer Name
	Weight     interfaces.IHash // Computed Weight at this DBHeight, Minute, Round
	Missing    interfaces.IMsg  // The Missing DBSig or EOM
	Ack        interfaces.IMsg  // The acknowledgement for the missing message

	messageHash interfaces.IHash
}

func delayVol(is interfaces.IState, e *elections.Elections, m *FedVoteMsg) {
	time.Sleep(100 * time.Millisecond)
	is.ElectionsQueue().Enqueue(m)
}

func (m *FedVoteMsg) ElectionProcess(is interfaces.IState, elect interfaces.IElections) {
	//s := is.(*state.State)
	e := elect.(*elections.Elections)

	if e.DBHeight > int(m.DBHeight) || e.Minute > int(m.Minute) {
		return
	}

	// If we don't have a timeout ourselves, then wait on this for a bit and try again.
	if e.Electing < 0 {
		go delayVol(is, e, m)
		return
	}

	idx := e.LeaderIndex(is.GetIdentityChainID())
	aidx := e.AuditIndex(is.GetIdentityChainID())

	//if m.DBHeight < uint32(e.DBHeight) || m.Minute < byte(e.Minute) || m.Round < e.Round[m.ServerIdx] {
	//	return
	//}
	auditIdx := e.AuditPriority()
	if aidx >= 0 && auditIdx == aidx {
		e.FeedBackStr(fmt.Sprintf("V%d", m.ServerIdx), false, aidx)
	} else if idx >= 0 {
		e.FeedBackStr(fmt.Sprintf("V%d", m.ServerIdx), true, idx)
	} else if aidx >= 0 {
		e.FeedBackStr(fmt.Sprintf("*%d", m.ServerIdx), false, aidx)
	}
	e.Msg = m.Missing
	e.Ack = m.Ack
	e.VName = m.ServerName
}

var _ interfaces.IMsg = (*FedVoteMsg)(nil)

func (a *FedVoteMsg) IsSameAs(msg interfaces.IMsg) bool {
	b, ok := msg.(*FedVoteMsg)
	if !ok {
		return false
	}
	if a.TS.GetTimeMilli() != b.TS.GetTimeMilli() {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	if a.EOM != b.EOM {
		return false
	}
	if a.ServerIdx != b.ServerIdx {
		return false
	}
	if a.ServerID.Fixed() != b.ServerID.Fixed() {
		return false
	}
	if a.Weight.Fixed() != b.Weight.Fixed() {
		return false
	}
	if a.DBHeight != b.DBHeight {
		return false
	}
	if a.VMIndex != b.VMIndex {
		return false
	}
	if a.Round != b.Round {
		return false
	}
	if a.Minute != b.Minute {
		return false
	}
	binA, errA := a.MarshalBinary()
	binB, errB := a.MarshalBinary()
	if errA != nil || errB != nil || bytes.Compare(binA, binB) != 0 {
		return false
	}
	return true
}

func (m *FedVoteMsg) GetServerID() interfaces.IHash {
	return m.ServerID
}

func (m *FedVoteMsg) LogFields() log.Fields {
	return log.Fields{"category": "message", "messagetype": "FedVoteMsg", "dbheight": m.DBHeight, "newleader": m.ServerID.String()[4:12]}
}

func (m *FedVoteMsg) GetRepeatHash() interfaces.IHash {
	return m.GetMsgHash()
}

// We have to return the haswh of the underlying message.

func (m *FedVoteMsg) GetHash() interfaces.IHash {
	return m.GetMsgHash()
}

func (m *FedVoteMsg) GetTimestamp() interfaces.Timestamp {
	return m.TS
}

func (m *FedVoteMsg) GetMsgHash() interfaces.IHash {
	if m.MsgHash == nil {
		data, err := m.MarshalBinary()
		if err != nil {
			return nil
		}
		m.MsgHash = primitives.Sha(data)
	}
	return m.MsgHash
}

func (m *FedVoteMsg) Type() byte {
	return constants.VOLUNTEERAUDIT
}

func (m *FedVoteMsg) Validate(state interfaces.IState) int {
	return 1
}

// Returns true if this is a message for this server to execute as
// a leader.
func (m *FedVoteMsg) ComputeVMIndex(state interfaces.IState) {
}

// Execute the leader functions of the given message
// Leader, follower, do the same thing.
func (m *FedVoteMsg) LeaderExecute(state interfaces.IState) {
	m.FollowerExecute(state)
}

func (m *FedVoteMsg) FollowerExecute(state interfaces.IState) {
	state.ElectionsQueue().Enqueue(m)
}

// Acknowledgements do not go into the process list.
func (e *FedVoteMsg) Process(dbheight uint32, state interfaces.IState) bool {
	panic("Ack object should never have its Process() method called")
}

func (e *FedVoteMsg) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *FedVoteMsg) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (m *FedVoteMsg) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling: %v", r)
		}
	}()

	buf := primitives.NewBuffer(data)
	if t, e := buf.PopByte(); e != nil || t != constants.VOLUNTEERAUDIT {
		return nil, errors.New("Not a Volunteer Audit type")
	}
	if m.TS, err = buf.PopTimestamp(); err != nil {
		return nil, err
	}
	if m.Name, err = buf.PopString(); err != nil {
		return nil, err
	}
	if m.EOM, err = buf.PopBool(); err != nil {
		return nil, err
	}
	if m.ServerIdx, err = buf.PopUInt32(); err != nil {
		return nil, err
	}
	if m.ServerID, err = buf.PopIHash(); err != nil {
		return nil, err
	}
	if m.ServerName, err = buf.PopString(); err != nil {
		return nil, err
	}
	if m.FedIdx, err = buf.PopUInt32(); err != nil {
		return nil, err
	}
	if m.FedID, err = buf.PopIHash(); err != nil {
		return nil, err
	}
	if m.Weight, err = buf.PopIHash(); err != nil {
		return nil, err
	}
	if m.DBHeight, err = buf.PopUInt32(); err != nil {
		return nil, err
	}
	if m.VMIndex, err = buf.PopInt(); err != nil {
		return nil, err
	}
	if m.Round, err = buf.PopInt(); err != nil {
		return nil, err
	}
	if m.Minute, err = buf.PopByte(); err != nil {
		return nil, err
	}
	if m.Ack, err = buf.PopMsg(); err != nil {
		return nil, err
	}
	if m.Missing, err = buf.PopMsg(); err != nil {
		return nil, err
	}
	newData, err = buf.PopBytes()
	return
}

func (m *FedVoteMsg) UnmarshalBinary(data []byte) error {
	_, err := m.UnmarshalBinaryData(data)
	return err
}

func (m *FedVoteMsg) MarshalBinary() (data []byte, err error) {
	var buf primitives.Buffer

	if err = buf.PushByte(constants.VOLUNTEERAUDIT); err != nil {
		return nil, err
	}
	if e := buf.PushTimestamp(m.TS); e != nil {
		return nil, e
	}
	if e := buf.PushString(m.Name); e != nil {
		return nil, e
	}
	if e := buf.PushBool(m.EOM); e != nil {
		return nil, e
	}
	if e := buf.PushUInt32(m.ServerIdx); e != nil {
		return nil, e
	}
	if e := buf.PushIHash(m.ServerID); e != nil {
		return nil, e
	}
	if e := buf.PushString(m.ServerName); e != nil {
		return nil, e
	}
	if e := buf.PushUInt32(m.FedIdx); e != nil {
		return nil, e
	}
	if e := buf.PushIHash(m.FedID); e != nil {
		return nil, e
	}
	if e := buf.PushIHash(m.Weight); e != nil {
		return nil, e
	}
	if e := buf.PushUInt32(m.DBHeight); e != nil {
		return nil, e
	}
	if e := buf.PushInt(m.VMIndex); e != nil {
		return nil, e
	}
	if e := buf.PushInt(m.Round); e != nil {
		return nil, e
	}
	if e := buf.PushByte(m.Minute); e != nil {
		return nil, e
	}
	if e := buf.PushMsg(m.Ack); e != nil {
		return nil, e
	}
	if e := buf.PushMsg(m.Missing); e != nil {
		return nil, e
	}
	data = buf.DeepCopyBytes()
	return data, nil
}

func (m *FedVoteMsg) String() string {
	if m.LeaderChainID == nil {
		m.LeaderChainID = primitives.NewZeroHash()
	}
	return fmt.Sprintf("%19s %20s %20s ID: %x WT: %x serverIdx: %d vmIdx: %d round: %d dbheight: %d minute: %d ",
		m.Name,
		"Volunteer Audit",
		m.TS.String(),
		m.ServerID.Bytes()[2:5],
		m.Weight.Bytes()[2:5],
		m.ServerIdx,
		m.VMIndex,
		m.Round,
		m.DBHeight,
		m.Minute)
}
