// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package p2p

import (
	"time"
)

// maximum time we wait for a partial message to arrive, old entries are cleaned up only when new part arrives
const MaxTimeWaitingForReassembly time.Duration = time.Second * 60 * 10

type PartialMessage struct {
	parts                  []*Parcel // array of message parts we've received so far
	firstPartReceived      time.Time // a timestamp indicating when the first part was received
	mostRecentPartReceived time.Time // a timestamp indicating when the mostRecent part was received
}

// PartsAssembler is responsible for assembling message parts into full messages
type PartsAssembler struct {
	messages map[uint64]*PartialMessage // a map of full message IDs to
}

// Initializes the assembler
func (assembler *PartsAssembler) Init() *PartsAssembler {
	assembler.messages = make(map[uint64]*PartialMessage)
	return assembler
}

// Handles a single message part, returns either a fully assembled message or nil
func (assembler *PartsAssembler) handlePart(parcel Parcel) *Parcel {
	debug("PartsAssembler", "Handling message part %d %d/%d", parcel.Header.MessageID, parcel.Header.PartNo, parcel.Header.PartsTotal)
	partial, exists := assembler.messages[parcel.Header.MessageID]

	if !exists {
		partial = createNewPartialMessage(parcel)
		assembler.messages[parcel.Header.MessageID] = partial
	}

	partial.parts[parcel.Header.PartNo] = &parcel
	partial.mostRecentPartReceived = time.Now()

	fullParcel := tryReassemblingMessage(partial)
	if fullParcel != nil {
		delete(assembler.messages, parcel.Header.MessageID)
		debug("PartsAssembler", "Fully assembled %d", parcel.Header.MessageID)
	}

	assembler.cleanupOldPartialMessages()

	return fullParcel
}

// Checks existing partial messages and if there is anything older than MaxTimeWaitingForReassembly,
// drops the partial message
func (assembler *PartsAssembler) cleanupOldPartialMessages() {
	for messageID, partial := range assembler.messages {
		timeWaiting := time.Since(partial.mostRecentPartReceived)
		timeSinceFirst := time.Since(partial.firstPartReceived)
		if timeWaiting > MaxTimeWaitingForReassembly {
			delete(assembler.messages, messageID)
			note("PartsAssembler", "Dropping message %d after %s secs, time since first part: %s secs",
				messageID, timeWaiting/time.Second, timeSinceFirst/time.Second)
		}
	}
}

// Creates a new PartialMessage from a given parcel
func createNewPartialMessage(parcel Parcel) *PartialMessage {
	partial := new(PartialMessage)
	partial.parts = make([]*Parcel, parcel.Header.PartsTotal)
	partial.firstPartReceived = time.Now()
	return partial
}

// Tries reassembling a full Message from existing MessageParts, returns nil if
// we don't have all the necessary parts yet
func tryReassemblingMessage(partial *PartialMessage) *Parcel {
	for _, part := range partial.parts {
		if part == nil {
			return nil
		}
	}

	return ReassembleParcel(partial.parts)
}
