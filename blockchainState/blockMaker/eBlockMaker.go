// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package blockMaker

import (
	"sort"

	"github.com/FactomProject/factomd/common/entryBlock"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

func (bm *BlockMaker) BuildEBlocks() ([]interfaces.IEntryBlock, []interfaces.IEBEntry, error) {
	sortedEntries := map[string][]*EBlockEntry{}
	for _, v := range bm.ProcessedEBEntries {
		sortedEntries[v.Entry.GetChainID().String()] = append(sortedEntries[v.Entry.GetChainID().String()], v)
	}
	keys := []string{}
	for k := range sortedEntries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	eBlocks := []interfaces.IEntryBlock{}
	ebEntries := []interfaces.IEBEntry{}
	for _, k := range keys {
		entries := sortedEntries[k]
		if len(entries) == 0 {
			continue
		}
		eb := entryBlock.NewEBlock()
		eb.GetHeader().SetChainID(entries[0].Entry.GetChainID())
		head := bm.BState.GetEBlockHead(entries[0].Entry.GetChainID().String())
		eb.GetHeader().SetPrevKeyMR(head.KeyMR)
		eb.GetHeader().SetPrevFullHash(head.Hash)
		eb.GetHeader().SetDBHeight(bm.BState.DBlockHeight + 1)

		minute := entries[0].Minute
		for _, v := range entries {
			if v.Minute != minute {
				eb.AddEndOfMinuteMarker(uint8(minute + 1))
				minute = v.Minute
			}
			err := eb.AddEBEntry(v.Entry)
			if err != nil {
				return nil, nil, err
			}
			ebEntries = append(ebEntries, v.Entry)
		}
		eb.AddEndOfMinuteMarker(uint8(minute + 1))
		eBlocks = append(eBlocks, eb)
	}
	return eBlocks, ebEntries, nil
}

func (bm *BlockMaker) ProcessEBEntry(e interfaces.IEntry, minute int) error {
	ebe := new(EBlockEntry)
	ebe.Entry = e
	ebe.Minute = minute
	err := bm.BState.ProcessEntryHash(e.GetHash(), primitives.NewZeroHash())
	if err != nil {
		return err
	}
	bm.ProcessedEBEntries = append(bm.ProcessedEBEntries, ebe)
	return nil
}

type EBlockEntry struct {
	Entry  interfaces.IEntry
	Minute int
}
