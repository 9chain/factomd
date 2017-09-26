// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package msgsupport

//https://docs.google.com/spreadsheets/d/1wy9JDEqyM2uRYhZ6Y1e9C3hIDm2prIILebztQ5BGlr8/edit#gid=1997221100

import (
	"fmt"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"

	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/common/messages/electionMsgs"
	log "github.com/FactomProject/logrus"
)

// packageLogger is the general logger for all message related logs. You can add additional fields,
// or create more context loggers off of this
var packageLogger = log.WithFields(log.Fields{"package": "messages"})

func UnmarshalMessage(data []byte) (interfaces.IMsg, error) {
	_, msg, err := UnmarshalMessageData(data)
	return msg, err
}

func CreateMsg(messageType byte) interfaces.IMsg {
	switch messageType {
	case constants.EOM_MSG:
		return new(messages.EOM)
	case constants.ACK_MSG:
		return new(messages.Ack)
	case constants.AUDIT_SERVER_FAULT_MSG:
		return new(messages.AuditServerFault)
	case constants.FED_SERVER_FAULT_MSG:
		return new(messages.ServerFault)
	case constants.FULL_SERVER_FAULT_MSG:
		return new(messages.FullServerFault)
	case constants.COMMIT_CHAIN_MSG:
		return new(messages.CommitChainMsg)
	case constants.COMMIT_ENTRY_MSG:
		return new(messages.CommitEntryMsg)
	case constants.DIRECTORY_BLOCK_SIGNATURE_MSG:
		return new(messages.DirectoryBlockSignature)
	case constants.FACTOID_TRANSACTION_MSG:
		return new(messages.FactoidTransaction)
	case constants.HEARTBEAT_MSG:
		return new(messages.Heartbeat)
	case constants.MISSING_MSG:
		return new(messages.MissingMsg)
	case constants.MISSING_MSG_RESPONSE:
		return new(messages.MissingMsgResponse)
	case constants.MISSING_DATA:
		return new(messages.MissingData)
	case constants.DATA_RESPONSE:
		return new(messages.DataResponse)
	case constants.REVEAL_ENTRY_MSG:
		return new(messages.RevealEntryMsg)
	case constants.REQUEST_BLOCK_MSG:
		return new(messages.RequestBlock)
	case constants.DBSTATE_MISSING_MSG:
		return new(messages.DBStateMissing)
	case constants.DBSTATE_MSG:
		return new(messages.DBStateMsg)
	case constants.ADDSERVER_MSG:
		return new(messages.AddServerMsg)
	case constants.CHANGESERVER_KEY_MSG:
		return new(messages.ChangeServerKeyMsg)
	case constants.REMOVESERVER_MSG:
		return new(messages.RemoveServerMsg)
	case constants.BOUNCE_MSG:
		return new(messages.Bounce)
	case constants.BOUNCEREPLY_MSG:
		return new(messages.BounceReply)
	case constants.SYNC_MSG:
		return new(electionMsgs.SyncMsg)
	case constants.VOLUNTEERAUDIT:
		return new(electionMsgs.VolunteerAudit)
	case constants.LEADER_ACK_VOLUNTEER:
		return new(electionMsgs.LeaderAck)
	default:
		return nil
	}
}

func UnmarshalMessageData(data []byte) (newdata []byte, msg interfaces.IMsg, err error) {
	if data == nil {
		return nil, nil, fmt.Errorf("No data provided")
	}
	if len(data) == 0 {
		return nil, nil, fmt.Errorf("No data provided")
	}
	messageType := data[0]
fmt.Println("messagetype",messageType,MessageName(messageType))
	msg = CreateMsg(messageType)
	if msg == nil {
		fmt.Sprintf("Transaction Failed to Validate %x", data[0])
		return data, nil, fmt.Errorf("Unknown message type %d %x", messageType, data[0])
	}

	newdata, err = msg.UnmarshalBinaryData(data[:])
	if err != nil {
		fmt.Sprintf("Transaction Failed to Unmarshal %x", data[0])
		return data, nil, err
	}

	return newdata, msg, nil

}

func MessageName(Type byte) string {
	switch Type {
	case constants.EOM_MSG:
		return "EOM"
	case constants.ACK_MSG:
		return "Ack"
	case constants.AUDIT_SERVER_FAULT_MSG:
		return "Audit Server Fault"
	case constants.FED_SERVER_FAULT_MSG:
		return "Fed Server Fault"
	case constants.FULL_SERVER_FAULT_MSG:
		return "Full Server Fault"
	case constants.COMMIT_CHAIN_MSG:
		return "Commit Chain"
	case constants.COMMIT_ENTRY_MSG:
		return "Commit Entry"
	case constants.DIRECTORY_BLOCK_SIGNATURE_MSG:
		return "Directory Block Signature"
	case constants.EOM_TIMEOUT_MSG:
		return "EOM Timeout"
	case constants.FACTOID_TRANSACTION_MSG:
		return "Factoid Transaction"
	case constants.HEARTBEAT_MSG:
		return "HeartBeat"
	case constants.INVALID_ACK_MSG:
		return "Invalid Ack"
	case constants.INVALID_DIRECTORY_BLOCK_MSG:
		return "Invalid Directory Block"
	case constants.MISSING_MSG:
		return "Missing Msg"
	case constants.MISSING_MSG_RESPONSE:
		return "Missing Msg Response"
	case constants.MISSING_DATA:
		return "Missing Data"
	case constants.DATA_RESPONSE:
		return "Data Response"
	case constants.REVEAL_ENTRY_MSG:
		return "Reveal Entry"
	case constants.REQUEST_BLOCK_MSG:
		return "Request Block"
	case constants.SIGNATURE_TIMEOUT_MSG:
		return "Signature Timeout"
	case constants.DBSTATE_MISSING_MSG:
		return "DBState Missing"
	case constants.DBSTATE_MSG:
		return "DBState"
	case constants.BOUNCE_MSG:
		return "Bounce Message"
	case constants.BOUNCEREPLY_MSG:
		return "Bounce Reply Message"
	case constants.SYNC_MSG:
		return "Sync Msg"
	case constants.VOLUNTEERAUDIT:
		return "Volunteer Audit"
	case constants.LEADER_ACK_VOLUNTEER:
		return "Leader Ack Volunteer"
	default:
		return "Unknown:" + fmt.Sprintf(" %d", Type)
	}
}

// GeneralFactory is used to get around package import loops.
type GeneralFactory struct {
}

var _ interfaces.IGeneralMsg = (*GeneralFactory)(nil)

func (GeneralFactory) CreateMsg(messageType byte) interfaces.IMsg {
	return CreateMsg(messageType)
}

func (GeneralFactory) MessageName(Type byte) string {
	return MessageName(Type)
}

func (GeneralFactory) UnmarshalMessageData(data []byte) (newdata []byte, msg interfaces.IMsg, err error) {
	return UnmarshalMessageData(data)
}

func (GeneralFactory) UnmarshalMessage(data []byte) (interfaces.IMsg, error) {
	return UnmarshalMessage(data)
}
