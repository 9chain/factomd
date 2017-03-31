// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package state

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/database/databaseOverlay"
)

func has(s *State, entry interfaces.IHash) bool {
	exists, _ := s.DB.DoesKeyExist(databaseOverlay.ENTRY, entry.Bytes())
	return exists
}

var _ = fmt.Print

func (s *State) fetchByTorrent(height uint32) {
	if height == 0 {
		return
	}
	err := s.GetMissingDBState(height)
	if err != nil {
		fmt.Println("DEBUG: Error in torrent retrieve: " + err.Error())
	}
}

// This go routine checks every so often to see if we have any missing entries or entry blocks.  It then requests
// them if it finds entries in the missing lists.
func (s *State) MakeMissingEntryRequests() {

	missing := 0
	found := 0

	MissingEntryMap := make(map[[32]byte]*MissingEntry)
	last := time.Now()

	for {

		now := time.Now()

		newrequest := 0

		cnt := 0
		sum := 0
		avg := 0
		var _ = avg

		// Look through our map, and remove any entries we now have in our database.
		for k := range MissingEntryMap {
			if has(s, MissingEntryMap[k].EntryHash) {
				found++
				delete(MissingEntryMap, k)
			} else {
				cnt++
				sum += MissingEntryMap[k].Cnt
			}
		}
		if cnt > 0 {
			avg = (1000 * sum) / cnt
		}

		ESAsking.Set(float64(len(MissingEntryMap)))
		ESAsking.Set(float64(cnt))
		ESFound.Set(float64(found))
		ESAvgRequests.Set(float64(avg) / 1000)

		// Keep our map of entries that we are asking for filled up.
	fillMap:
		for len(MissingEntryMap) < 3000 {
			select {
			case et := <-s.MissingEntries:
				missing++
				MissingEntryMap[et.EntryHash.Fixed()] = et
			default:
				break fillMap
			}
		}

		sent := 0

		if len(s.inMsgQueue) < 500 {
			// Make requests for entries we don't have.
			for k := range MissingEntryMap {

				et := MissingEntryMap[k]

				if et.Cnt == 0 {
					et.Cnt = 1
					et.LastTime = now.Add(time.Duration((rand.Int() % 5000)) * time.Millisecond)
					continue
				}
				if now.Unix()-et.LastTime.Unix() > 5 && sent < 100 {
					sent++
					entryRequest := messages.NewMissingData(s, et.EntryHash)
					entryRequest.SendOut(s, entryRequest)
					newrequest++
					et.LastTime = now.Add(time.Duration((rand.Int() % 5000)) * time.Millisecond)
					et.Cnt++
				}

			}
		} else {
			time.Sleep(20 * time.Second)
		}

		if s.UsingTorrent() { // Torrent solution to assist
			// Only fetch unique heights.
			heightsRequested := make(map[uint32]struct{})

			// Only do once per 5 seconds, anything more is wasted calls to torrent
			if !last.Before(now.Add(-5 * time.Second)) {
				goto skipTorrent
			}
			last = now

			for k := range MissingEntryMap {
				et := MissingEntryMap[k]
				if _, ok := heightsRequested[et.DBHeight]; ok {
					continue // Already requested this
				}

				// The entries in this height will be returned on a channel
				s.fetchByTorrent(et.DBHeight)              // For printout
				heightsRequested[et.DBHeight] = struct{}{} // Only request each height once per batch of requests
			}
		skipTorrent:
		}

		// Insert the entries we have found into the database.
	InsertLoop:
		for {

			select {

			case entry := <-s.WriteEntry:

				asked := MissingEntryMap[entry.GetHash().Fixed()] != nil

				if asked {
					s.DB.InsertEntry(entry)
				}

			default:
				break InsertLoop
			}
		}
		if sent == 0 {
			if s.GetHighestKnownBlock()-s.GetHighestSavedBlk() > 100 {
				time.Sleep(10 * time.Second)
			} else {
				time.Sleep(1 * time.Second)
			}
			if s.EntryDBHeightComplete == s.GetHighestSavedBlk() {
				time.Sleep(20 * time.Second)
			}
		}
	}
}

func (s *State) GoSyncEntries() {
	go s.MakeMissingEntryRequests()

	// Map to track what I know is missing
	missingMap := make(map[[32]byte]interfaces.IHash)

	// Once I have found all the entries, we quit searching so much for missing entries.
	start := uint32(1)
	entryMissing := 0

	// If I find no missing entries, then the firstMissing will be -1
	firstMissing := -1

	lastfirstmissing := 0

	for {

		ESMissingQueue.Set(float64(len(missingMap)))
		ESDBHTComplete.Set(float64(s.EntryDBHeightComplete))
		ESFirstMissing.Set(float64(lastfirstmissing))
		ESHighestMissing.Set(float64(s.GetHighestSavedBlk()))

		entryMissing = 0

		for k := range missingMap {
			if has(s, missingMap[k]) {
				delete(missingMap, k)
			}
		}

		// Scan all the directory blocks, from start to the highest saved.  Once we catch up,
		// start will be the last block saved.

		// First reset first Missing back to -1 every time.
		firstMissing = -1

	dirblkSearch:
		for scan := start; scan <= s.GetHighestSavedBlk(); scan++ {

			if firstMissing < 0 {
				if scan > 1 {
					s.EntryDBHeightComplete = scan - 1
					start = scan
				}
			}

			db := s.GetDirectoryBlockByHeight(scan)

			// Wait for the database if we have to
			for db == nil {
				time.Sleep(1 * time.Second)
				db = s.GetDirectoryBlockByHeight(scan)
			}

			for _, ebKeyMR := range db.GetEntryHashes()[3:] {
				// The first three entries (0,1,2) in every directory block are blocks we already have by
				// definition.  If we decide to not have Factoid blocks or Entry Credit blocks in some cases,
				// then this assumption might not hold.  But it does for now.

				eBlock, _ := s.DB.FetchEBlock(ebKeyMR)

				// Dont have an eBlock?  Huh. We can go on, but we can't advance.  We just wait until it
				// does show up.
				for eBlock == nil {
					time.Sleep(1 * time.Second)
					eBlock, _ = s.DB.FetchEBlock(ebKeyMR)
				}

				// Go through all the entry hashes.
				for _, entryhash := range eBlock.GetEntryHashes() {
					if entryhash.IsMinuteMarker() {
						continue
					}
					ueh := new(EntryUpdate)
					ueh.Hash = entryhash
					ueh.Timestamp = db.GetTimestamp()
					s.UpdateEntryHash <- ueh

					// If I have the entry, then remove it from the Missing Entries list.
					if has(s, entryhash) {
						delete(missingMap, entryhash.Fixed())
						continue
					}

					if firstMissing < 0 {
						firstMissing = int(scan)
					}

					eh := missingMap[entryhash.Fixed()]
					if eh == nil {

						// If we have a full queue, break so we don't stall.
						if len(s.MissingEntries) > 9000 {
							break dirblkSearch
						}

						var v MissingEntry

						v.DBHeight = eBlock.GetHeader().GetDBHeight()
						v.EntryHash = entryhash
						v.EBHash = ebKeyMR
						entryMissing++
						missingMap[entryhash.Fixed()] = entryhash
						s.MissingEntries <- &v
					}

				}
			}
		}
		lastfirstmissing = firstMissing
		if firstMissing < 0 {
			time.Sleep(60 * time.Second)
		}

		time.Sleep(1 * time.Second)

	}
}
