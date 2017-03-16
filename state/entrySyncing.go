// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package state

import (
	"fmt"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/database/databaseOverlay"
	"time"
)

func has(s *State, entry interfaces.IHash) bool {
	exists, _ := s.DB.DoesKeyExist(databaseOverlay.ENTRY, entry.Bytes())
	return exists
}

var _ = fmt.Print

func (s *State) TorrentMissingEntries() {
	for {
		// Copy missing list
		var low uint32 = 99999999
		var high uint32 = 0
		var prev uint32 = 0
		amt := 0

		s.MissingEntryMutex.Lock()
		for k := range s.MissingEntryMap {
			et := s.MissingEntryMap[k]
			s.MissingEntryMutex.Unlock()
			if et.DBHeight != prev {
				if et.DBHeight < low {
					low = et.DBHeight
				}
				if et.DBHeight > high {
					high = et.DBHeight
				}
				s.fetchByTorrent(et.DBHeight)
				amt++
				prev = et.DBHeight
			}

			s.MissingEntryMutex.Lock()
		}
		s.MissingEntryMutex.Unlock()

		if high != 0 {
			fmt.Printf("{{ Torrenting heights: Low: %d, High %d, Total: %d }} \n", low, high, amt)
		}

		s.SetDBStateManagerCompletedHeight(s.EntryDBHeightComplete)
		time.Sleep(5 * time.Second)
	}
}

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

	pause := time.Now().Add(60 * time.Second)

	for {

		now := time.Now()
		// While we are waiting to boot or whatever, just idle here for a while.  Otherwise we can lock up.
		for now.Before(pause) {
			time.Sleep(1 * time.Second)
			now = time.Now()
		}

		newrequest := 0

		cnt := 0
		sum := 0
		avg := 0
		var _ = avg

		// Look through our map, and remove any entries we now have in our database.
		s.MissingEntryMutex.Lock()
		for k := range s.MissingEntryMap {
			if has(s, s.MissingEntryMap[k].EntryHash) {
				found++
				delete(s.MissingEntryMap, k)
			} else {
				cnt++
				sum += s.MissingEntryMap[k].Cnt
			}
		}
		if cnt > 0 {
			avg = sum / cnt
		}

		// fmt.Printf("***es %-10s Avg %6d Missing: %6d  Found: %6d Queue: %d\n", s.FactomNodeName, avg, missing, found, len(s.MissingEntries))

		s.MissingEntryMutex.Unlock()

		time.Sleep(20 * time.Millisecond)

		s.MissingEntryMutex.Lock()
	fillMap:
		for len(s.MissingEntryMap) < 1500 {
			select {
			case et := <-s.MissingEntries:
				missing++
				s.MissingEntryMap[et.EntryHash.Fixed()] = &et
			default:
				break fillMap
			}
		}
		s.MissingEntryMutex.Unlock()

		s.MissingEntryMutex.Lock()
		if !s.UsingTorrent() {
			for k := range s.MissingEntryMap {
				et := s.MissingEntryMap[k]
				s.MissingEntryMutex.Unlock()

				if et.Cnt == 0 || now.Unix()-et.LastTime.Unix() > 40 {
					entryRequest := messages.NewMissingData(s, et.EntryHash)
					entryRequest.SendOut(s, entryRequest)
					newrequest++
					et.LastTime = now
					et.Cnt++
					if et.Cnt%25 == 25 {
						fmt.Printf("***es Can't get Entry Block %x Entry %x in %v attempts.\n", et.EBHash.Bytes(), et.EntryHash.Bytes(), et.Cnt)
					}
				}
				s.MissingEntryMutex.Lock()
			}
		}
		s.MissingEntryMutex.Unlock()
		// slow down as the number of retries per message goes up
		time.Sleep(time.Duration((200)) * time.Millisecond)
		for len(s.inMsgQueue) > 100 {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *State) GoWriteEntries() {
	for {
		time.Sleep(300 * time.Millisecond)

	entryWrite:
		for {

			select {

			case entry := <-s.WriteEntry:

				s.MissingEntryMutex.Lock()
				asked := s.MissingEntryMap[entry.GetHash().Fixed()] != nil
				s.MissingEntryMutex.Unlock()

				if asked {
					s.DB.InsertEntry(entry)
				}

			default:
				break entryWrite
			}
		}
	}
}

func (s *State) GoSyncEntries() {
	go s.MakeMissingEntryRequests()
	go s.GoWriteEntries() // Start a go routine to write the Entries to the DB

	// Map to track what I know is missing
	missingMap := make(map[[32]byte]interfaces.IHash)

	// Once I have found all the entries, we quit searching so much for missing entries.
	start := uint32(0)
	entryMissing := 0

	// If I find no missing entries, then the firstMissing will be -1
	firstMissing := -1
	lastfirstmissing := 0
	for {
		fmt.Printf("***es %10s Missing: %6d MissingMap %6d FirstMissing %6d\n", s.FactomNodeName, entryMissing, len(missingMap), lastfirstmissing)
		entryMissing = 0

		for k := range missingMap {
			if has(s, missingMap[k]) {
				delete(missingMap, k)
			}
		}

		// Scan all the directory blocks, from start to the highest saved.  Once we catch up,
		// start will be the last block saved.
	dirblkSearch:
		for scan := start; scan <= s.GetHighestSavedBlk(); scan++ {

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
					if !entryhash.IsMinuteMarker() {

						// If I have the entry, then remove it from the Missing Entries list.
						if has(s, entryhash) {
							delete(missingMap, entryhash.Fixed())
						} else {

							eh := missingMap[entryhash.Fixed()]
							if eh == nil {

								if firstMissing < 0 {
									firstMissing = int(scan)
								}

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
								s.MissingEntries <- v
							}
						}
						ueh := new(EntryUpdate)
						ueh.Hash = entryhash
						ueh.Timestamp = db.GetTimestamp()
						s.UpdateEntryHash <- ueh
					}
				}
			}
			start = scan
		}
		lastfirstmissing = firstMissing
		// If we are caught up, we hardly need to do anything.
		for start >= s.GetHighestSavedBlk() {
			time.Sleep(1 * time.Second)
		}
		if firstMissing >= 0 {
			if firstMissing > 0 {
				s.EntryDBHeightComplete = uint32(firstMissing - 1)
			}
			firstMissing = -1
		} else {
			s.EntryDBHeightComplete = start
		}
		start = s.EntryDBHeightComplete
		// sleep some time no matter what.
		for len(s.MissingEntries) > 1000 {
			time.Sleep(100 * time.Millisecond)
		}
		for len(s.inMsgQueue) > 100 {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
