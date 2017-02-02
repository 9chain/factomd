// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package messages

import (
	"encoding/binary"
	"fmt"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/entryBlock"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

//Requests entry blocks from a range of DBlocks

type EntryBlockResponse struct {
	MessageBase
	Timestamp interfaces.Timestamp

	EBlockCount uint32
	EBlocks     []interfaces.IEntryBlock
	EntryCount  uint32
	Entries     []interfaces.IEBEntry

	//Not signed!
}

var _ interfaces.IMsg = (*EntryBlockResponse)(nil)

func (a *EntryBlockResponse) IsSameAs(b *EntryBlockResponse) bool {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseIsSameAs.Observe(float64(time.Now().UnixNano() - callTime))
	if b == nil {
		return false
	}
	if a.Timestamp.GetTimeMilli() != b.Timestamp.GetTimeMilli() {
		return false
	}
	if a.EBlockCount != b.EBlockCount {
		return false
	}
	if a.EntryCount != b.EntryCount {
		return false
	}

	if len(a.EBlocks) != len(b.EBlocks) {
		return false
	}

	if len(a.Entries) != len(b.Entries) {
		return false
	}

	//TODO: check blocks and entries

	return true
}

func (m *EntryBlockResponse) GetRepeatHash() interfaces.IHash {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseGetRepeatHash.Observe(float64(time.Now().UnixNano() - callTime))
	return m.GetMsgHash()
}

func (m *EntryBlockResponse) GetHash() interfaces.IHash {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseGetHash.Observe(float64(time.Now().UnixNano() - callTime))
	return m.GetMsgHash()
}

func (m *EntryBlockResponse) GetMsgHash() interfaces.IHash {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseGetMsgHash.Observe(float64(time.Now().UnixNano() - callTime))
	if m.MsgHash == nil {
		data, err := m.MarshalBinary()
		if err != nil {
			return nil
		}
		m.MsgHash = primitives.Sha(data)
	}
	return m.MsgHash
}

func (m *EntryBlockResponse) Type() byte {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseType.Observe(float64(time.Now().UnixNano() - callTime))
	return constants.ENTRY_BLOCK_RESPONSE
}

func (m *EntryBlockResponse) GetTimestamp() interfaces.Timestamp {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseGetTimestamp.Observe(float64(time.Now().UnixNano() - callTime))
	return m.Timestamp
}

// Validate the message, given the state.  Three possible results:
//  < 0 -- Message is invalid.  Discard
//  0   -- Cannot tell if message is Valid
//  1   -- Message is valid
func (m *EntryBlockResponse) Validate(state interfaces.IState) int {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseValidate.Observe(float64(time.Now().UnixNano() - callTime))
	if m.EBlockCount != uint32(len(m.EBlocks)) {
		return -1
	}
	if m.EntryCount != uint32(len(m.Entries)) {
		return -1
	}

	return 1
}

func (m *EntryBlockResponse) ComputeVMIndex(state interfaces.IState) {
}

// Execute the leader functions of the given message
func (m *EntryBlockResponse) LeaderExecute(state interfaces.IState) {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseLeaderExecute.Observe(float64(time.Now().UnixNano() - callTime))
	m.FollowerExecute(state)
}

func (m *EntryBlockResponse) FollowerExecute(state interfaces.IState) {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseFollowerExecute.Observe(float64(time.Now().UnixNano() - callTime))
	if len(state.NetworkOutMsgQueue()) > 1000 {
		return
	}

	db := state.GetAndLockDB()
	defer state.UnlockDB()

	for _, v := range m.EBlocks {
		db.ProcessEBlockBatchWithoutHead(v, true)
	}
	for _, v := range m.Entries {
		db.InsertEntry(v)
	}

	return
}

// Acknowledgements do not go into the process list.
func (e *EntryBlockResponse) Process(dbheight uint32, state interfaces.IState) bool {
	panic("Ack object should never have its Process() method called")
}

func (e *EntryBlockResponse) JSONByte() ([]byte, error) {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseJSONByte.Observe(float64(time.Now().UnixNano() - callTime))
	return primitives.EncodeJSON(e)
}

func (e *EntryBlockResponse) JSONString() (string, error) {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseJSONString.Observe(float64(time.Now().UnixNano() - callTime))
	return primitives.EncodeJSONString(e)
}

func (m *EntryBlockResponse) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseUnmarshalBinaryData.Observe(float64(time.Now().UnixNano() - callTime))
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling Directory Block State Missing Message: %v", r)
		}
	}()
	newData = data
	if newData[0] != m.Type() {
		return nil, fmt.Errorf("Invalid Message type")
	}
	newData = newData[1:]

	m.Peer2Peer = true // This is always a Peer2peer message

	m.Timestamp = new(primitives.Timestamp)
	newData, err = m.Timestamp.UnmarshalBinaryData(newData)
	if err != nil {
		return nil, err
	}

	m.EBlockCount, newData = binary.BigEndian.Uint32(newData[0:4]), newData[4:]

	for i := 0; i < int(m.EBlockCount); i++ {
		eBlock := entryBlock.NewEBlock()
		newData, err = eBlock.UnmarshalBinaryData(newData)
		if err != nil {
			return nil, err
		}
		m.EBlocks = append(m.EBlocks, eBlock)
	}

	m.EntryCount, newData = binary.BigEndian.Uint32(newData[0:4]), newData[4:]

	for i := 0; i < int(m.EntryCount); i++ {
		entry := entryBlock.NewEntry()
		newData, err = entry.UnmarshalBinaryData(newData)
		if err != nil {
			return nil, err
		}
		m.Entries = append(m.Entries, entry)
	}

	return
}

func (m *EntryBlockResponse) UnmarshalBinary(data []byte) error {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseUnmarshalBinary.Observe(float64(time.Now().UnixNano() - callTime))
	_, err := m.UnmarshalBinaryData(data)
	return err
}

func (m *EntryBlockResponse) MarshalForSignature() ([]byte, error) {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseMarshalForSignature.Observe(float64(time.Now().UnixNano() - callTime))
	var buf primitives.Buffer

	binary.Write(&buf, binary.BigEndian, m.Type())

	t := m.GetTimestamp()
	data, err := t.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	m.EBlockCount = uint32(len(m.EBlocks))
	binary.Write(&buf, binary.BigEndian, m.EBlockCount)
	for _, eb := range m.EBlocks {
		bin, err := eb.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(bin)
	}

	m.EntryCount = uint32(len(m.Entries))
	binary.Write(&buf, binary.BigEndian, m.EntryCount)
	for _, e := range m.Entries {
		bin, err := e.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(bin)
	}

	return buf.DeepCopyBytes(), nil
}

func (m *EntryBlockResponse) MarshalBinary() ([]byte, error) {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseMarshalBinary.Observe(float64(time.Now().UnixNano() - callTime))
	return m.MarshalForSignature()
}

func (m *EntryBlockResponse) String() string {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseString.Observe(float64(time.Now().UnixNano() - callTime))
	str, _ := m.JSONString()
	return str
}

func NewEntryBlockResponse(state interfaces.IState) interfaces.IMsg {
	callTime := time.Now().UnixNano()
	defer messagesEntryBlockResponseNewEntryBlockResponse.Observe(float64(time.Now().UnixNano() - callTime))
	msg := new(EntryBlockResponse)

	msg.Peer2Peer = true // Always a peer2peer request.
	msg.Timestamp = state.GetTimestamp()

	return msg
}
