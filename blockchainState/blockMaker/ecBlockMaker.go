// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package blockMaker

import (
	"github.com/FactomProject/factomd/common/entryCreditBlock"
	"github.com/FactomProject/factomd/common/interfaces"
)

func (bm *BlockMaker) BuildECBlock() (interfaces.IEntryCreditBlock, error) {
	ecBlock := entryCreditBlock.NewECBlock()
	ecBlock.GetHeader().SetPrevHeaderHash(bm.BState.ECBlockHead.KeyMR)
	ecBlock.GetHeader().SetPrevFullHash(bm.BState.ECBlockHead.Hash)
	ecBlock.GetHeader().SetDBHeight(bm.BState.ECBlockHead.Height + 1)

	minute := 0
	for _, v := range bm.ProcessedECBEntries {
		for ; minute < v.Minute; minute++ {
			e := entryCreditBlock.NewMinuteNumber(uint8(minute + 1))
			ecBlock.GetBody().AddEntry(e)
		}
		ecBlock.GetBody().AddEntry(v.Entry)
	}
	for ; minute < 9; minute++ {
		e := entryCreditBlock.NewMinuteNumber(uint8(minute + 1))
		ecBlock.GetBody().AddEntry(e)
	}

	err := ecBlock.BuildHeader()
	if err != nil {
		return nil, err
	}

	return ecBlock, nil
}

func (bm *BlockMaker) ProcessECEntry(e interfaces.IECBlockEntry) error {
	ebe := new(ECBlockEntry)
	ebe.Entry = e
	ebe.Minute = bm.CurrentMinute
	err := bm.BState.ProcessECEntry(e)
	if err != nil {
		return err
	}
	bm.ProcessedECBEntries = append(bm.ProcessedECBEntries, ebe)
	return nil
}

type ECBlockEntry struct {
	Entry  interfaces.IECBlockEntry
	Minute int
}
