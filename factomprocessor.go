// Copyright 2015 FactomProject Authors. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// factomlog is based on github.com/alexcesaro/log and
// github.com/alexcesaro/log/golog (MIT License)

package btcd

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/FactomProject/FactomCode/anchor"
	"github.com/FactomProject/FactomCode/common"
	"github.com/FactomProject/FactomCode/consensus"
	"github.com/FactomProject/FactomCode/database"
	"github.com/FactomProject/FactomCode/factomlog"
    "github.com/FactomProject/FactomCode/util"
    sc "github.com/FactomProject/simplecoin"
    "github.com/FactomProject/simplecoin/block"
    "github.com/FactomProject/btcd/wire"
	"github.com/FactomProject/btcutil"
	"github.com/davecgh/go-spew/spew"
)

var _ = (*sc.Transaction)(nil)
var _ = (*block.SCBlock)(nil)


var (
	currentAddr btcutil.Address
	db          database.Db        // database
	dchain      *common.DChain     //Directory Block Chain
	ecchain     *common.ECChain    //Entry Credit Chain
	achain      *common.AdminChain //Admin Chain
	scchain     *common.SCChain    // SimpleCoin Chain
	fchainID    *common.Hash

	creditsPerChain   int32  = 10
	creditsPerFactoid uint64 = 1000

	// To be moved to ftmMemPool??
	chainIDMap     map[string]*common.EChain // ChainIDMap with chainID string([32]byte) as key
	commitChainMap = make(map[string]*common.CommitChain, 0)
	commitEntryMap = make(map[string]*common.CommitEntry, 0)
	eCreditMap     map[*[32]byte]int32 // eCreditMap with public key string([32]byte) as key, credit balance as value

	chainIDMapBackup map[string]*common.EChain //previous block bakcup - ChainIDMap with chainID string([32]byte) as key
	eCreditMapBackup map[*[32]byte]int32       // backup from previous block - eCreditMap with public key string([32]byte) as key, credit balance as value

	//Diretory Block meta data map
	//dbInfoMap map[string]*common.DBInfo // dbInfoMap with dbHash string([32]byte) as key

	fMemPool *ftmMemPool
	plMgr    *consensus.ProcessListMgr

	//Server Private key and Public key for milestone 1
	serverPrivKey common.PrivateKey
	serverPubKey  common.PublicKey

	FactoshisPerCredit uint64 // .001 / .15 * 100000000 (assuming a Factoid is .15 cents, entry credit = .1 cents

	factomdUser string
	factomdPass string
)

var (
	directoryBlockInSeconds int
	dataStorePath           string
	ldbpath                 string
	nodeMode                string
	devNet                  bool
	serverPrivKeyHex        string
)

func LoadConfigurations(cfg *util.FactomdConfig) {
	util.Trace()

	//setting the variables by the valued form the config file
	logLevel = cfg.Log.LogLevel
	dataStorePath = cfg.App.DataStorePath
	ldbpath = cfg.App.LdbPath
	directoryBlockInSeconds = cfg.App.DirectoryBlockInSeconds
	nodeMode = cfg.App.NodeMode
	serverPrivKeyHex = cfg.App.ServerPrivKey

	factomdUser = cfg.Btc.RpcUser
	factomdPass = cfg.Btc.RpcPass
}

func watchError(err error) {
	panic(err)
}

func readError(err error) {
	fmt.Println("error: ", err)
}

// Initialize the entry chains in memory from db
func initEChainFromDB(chain *common.EChain) {

	eBlocks, _ := db.FetchAllEBlocksByChain(chain.ChainID)
	sort.Sort(util.ByEBlockIDAccending(*eBlocks))

	for i := 0; i < len(*eBlocks); i = i + 1 {
		if uint32(i) != (*eBlocks)[i].Header.EBHeight {
			panic(errors.New("BlockID does not equal index for chain:" + chain.ChainID.String() + " block:" + fmt.Sprintf("%v", (*eBlocks)[i].Header.EBHeight)))
		}
	}

	if len(*eBlocks) == 0 {
		chain.NextBlockHeight = 0
		chain.NextBlock, _ = common.CreateBlock(chain, nil, 10)
	} else {
		chain.NextBlockHeight = uint32(len(*eBlocks))
		chain.NextBlock, _ = common.CreateBlock(chain, &(*eBlocks)[len(*eBlocks)-1], 10)
	}

	// Initialize chain with the first entry (Name and rules) for non-server mode
	if nodeMode != common.SERVER_NODE && chain.FirstEntry == nil && len(*eBlocks) > 0 {
		chain.FirstEntry, _ = db.FetchEntryByHash((*eBlocks)[0].EBEntries[0].EntryHash)
		if chain.FirstEntry != nil {
			db.InsertChain(chain)
		}
	}

	if chain.NextBlock.IsSealed == true {
		panic("chain.NextBlock.IsSealed for chain:" + chain.ChainID.String())
	}
}

func initProcess() {

	wire.Init()

	util.Trace()

	// init server private key or pub key
	initServerKeys()

	// init mem pools
	fMemPool = new(ftmMemPool)
	fMemPool.init_ftmMemPool()

	// init wire.FChainID
	wire.FChainID = new(common.Hash)
	wire.FChainID.SetBytes(common.FACTOID_CHAINID)

	FactoshisPerCredit = 666667 // .001 / .15 * 100000000 (assuming a Factoid is .15 cents, entry credit = .1 cents

	// init Directory Block Chain
	initDChain()
	fmt.Println("Loaded", dchain.NextBlockHeight, "Directory blocks for chain: "+dchain.ChainID.String())

	// init Entry Credit Chain
	initECChain()
	fmt.Println("Loaded", ecchain.NextBlockHeight, "Entry Credit blocks for chain: "+ecchain.ChainID.String())

	// init Admin Chain
	initAChain()
	fmt.Println("Loaded", achain.NextBlockHeight, "Admin blocks for chain: "+achain.ChainID.String())

    initSCChain()
    fmt.Println("Loaded", scchain.NextBlockHeight, "Simplecoin blocks for chain: "+scchain.ChainID.String())
    
	anchor.InitAnchor(db)

	// build the Genesis blocks if the current height is 0
	if dchain.NextBlockHeight == 0 {
		buildGenesisBlocks()
	} else {
		/*
			// still send a message to the btcd-side to start up the database; such as a current block height
			eomMsg := &wire.MsgInt_EOM{
				EOM_Type:         wire.INFO_CURRENT_HEIGHT,
				NextDBlockHeight: dchain.NextBlockHeight,
			}
			outCtlMsgQueue <- eomMsg
		*/

		// To be improved in milestone 2
		SignDirectoryBlock()
	}

	// init process list manager
	initProcessListMgr()

	// init Entry Chains
	initEChains()
	for _, chain := range chainIDMap {
		initEChainFromDB(chain)

		fmt.Println("Loaded", chain.NextBlockHeight, "blocks for chain: "+chain.ChainID.String())
		//fmt.Printf("PROCESSOR: echain=%s\n", spew.Sdump(chain))
	}

	// Validate all dir blocks
	err := validateDChain(dchain)
	if err != nil {
		if nodeMode == common.SERVER_NODE {
			panic("Error found in validating directory blocks: " + err.Error())
		} else {
			dchain.IsValidated = false
		}
	}
}

func Start_Processor(
	ldb database.Db,
	inMsgQ chan wire.FtmInternalMsg,
	outMsgQ chan wire.FtmInternalMsg,
	inCtlMsgQ chan wire.FtmInternalMsg,
	outCtlMsgQ chan wire.FtmInternalMsg) {
	db = ldb

	inMsgQueue = inMsgQ
	outMsgQueue = outMsgQ

	inCtlMsgQueue = inCtlMsgQ
	outCtlMsgQueue = outCtlMsgQ

	initProcess()

	// Initialize timer for the open dblock before processing messages
	if nodeMode == common.SERVER_NODE {
		timer := &BlockTimer{
			nextDBlockHeight: dchain.NextBlockHeight,
			inCtlMsgQueue:    inCtlMsgQueue,
		}
		go timer.StartBlockTimer()
	}

    
	// Process msg from the incoming queue one by one
	for {
		select {
		case msg := <-inMsgQ:
			fmt.Printf("PROCESSOR: in inMsgQ, msg:%+v\n", msg)

			if err := serveMsgRequest(msg); err != nil {
				log.Println(err)
			}

		case ctlMsg := <-inCtlMsgQueue:
			fmt.Printf("PROCESSOR: in ctlMsg, msg:%+v\n", ctlMsg)

			if err := serveMsgRequest(ctlMsg); err != nil {
				log.Println(err)
			}
		}

	}

	util.Trace()

}

func fileNotExists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return true
	}
	return err != nil
}

// Serve the "fast lane" incoming control msg from inCtlMsgQueue
func serveCtlMsgRequest(msg wire.FtmInternalMsg) error {

	util.Trace()

	switch msg.Command() {
	case wire.CmdCommitChain:

	default:
		return errors.New("Message type unsupported:" + fmt.Sprintf("%+v", msg))
	}
	return nil

}

// Serve incoming msg from inMsgQueue
func serveMsgRequest(msg wire.FtmInternalMsg) error {

	util.Trace()

	switch msg.Command() {
	case wire.CmdCommitChain:
		msgCommitChain, ok := msg.(*wire.MsgCommitChain)
		if ok && msgCommitChain.IsValid() {
			err := processCommitChain(msgCommitChain)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Error in processing msg:" + fmt.Sprintf("%+v", msg))
		}
		// Broadcast the msg to the network if no errors
		outMsgQueue <- msg

	case wire.CmdCommitEntry:
		msgCommitEntry, ok := msg.(*wire.MsgCommitEntry)
		if ok && msgCommitEntry.IsValid() {
			err := processCommitEntry(msgCommitEntry)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Error in processing msg:" + fmt.Sprintf("%+v", msg))
		}
		// Broadcast the msg to the network if no errors
		outMsgQueue <- msg

	case wire.CmdRevealEntry:
		msgRevealEntry, ok := msg.(*wire.MsgRevealEntry)
		if ok {
			err := processRevealEntry(msgRevealEntry)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Error in processing msg:" + fmt.Sprintf("%+v", msg))
		}
		// Broadcast the msg to the network if no errors
		outMsgQueue <- msg

	case wire.CmdInt_EOM:
		util.Trace("CmdInt_EOM")

		if nodeMode == common.SERVER_NODE {
			msgEom, ok := msg.(*wire.MsgInt_EOM)
			if !ok {
				return errors.New("Error in build blocks:" + fmt.Sprintf("%+v", msg))
			}
			fmt.Printf("PROCESSOR: End of minute msg - wire.CmdInt_EOM:%+v\n", msg)

			if msgEom.EOM_Type == wire.END_MINUTE_10 {
				// Process from Orphan pool before the end of process list
				processFromOrphanPool()

				// Pass the Entry Credit Exchange Rate into the Factoid component
				msgEom.EC_Exchange_Rate = FactoshisPerCredit
				plMgr.AddMyProcessListItem(msgEom, nil, wire.END_MINUTE_10)

				err := buildBlocks()
				if err != nil {
					return err
				}

			} else if msgEom.EOM_Type >= wire.END_MINUTE_1 && msgEom.EOM_Type < wire.END_MINUTE_10 {
				plMgr.AddMyProcessListItem(msgEom, nil, msgEom.EOM_Type)
			}
		}

	case wire.CmdInt_FactoidBlock: // to be removed??
		factoidBlock, ok := msg.(*wire.MsgInt_FactoidBlock)
		util.Trace("Factoid Block (GENERATED??) -- detected in the processor")
		fmt.Println("factoidBlock= ", factoidBlock, " ok= ", ok)

	case wire.CmdDirBlock:
		if nodeMode == common.SERVER_NODE {
			break
		}

		dirBlock, ok := msg.(*wire.MsgDirBlock)
		if ok {
			err := processDirBlock(dirBlock)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Error in processing msg:" + fmt.Sprintf("%+v", msg))
		}

    case wire.CmdSCBlock:
        if nodeMode == common.SERVER_NODE {
            break
        }
        
        scblock, ok := msg.(*wire.MsgSCBlock)
        if ok {
            err := processSCBlock(scblock)
            if err != nil {
                return err
            }
        } else {
            return errors.New("Error in processing msg:" + fmt.Sprintf("%+v", msg))
        }
        
	case wire.CmdABlock:
		if nodeMode == common.SERVER_NODE {
			break
		}

		ablock, ok := msg.(*wire.MsgABlock)
		if ok {
			err := processABlock(ablock)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Error in processing msg:" + fmt.Sprintf("%+v", msg))
		}

	case wire.CmdECBlock:
		if nodeMode == common.SERVER_NODE {
			break
		}

		cblock, ok := msg.(*wire.MsgECBlock)
		if ok {
			err := processCBlock(cblock)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Error in processing msg:" + fmt.Sprintf("%+v", msg))
		}

	case wire.CmdEBlock:
		if nodeMode == common.SERVER_NODE {
			break
		}

		eblock, ok := msg.(*wire.MsgEBlock)
		if ok {
			err := processEBlock(eblock)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Error in processing msg:" + fmt.Sprintf("%+v", msg))
		}

	case wire.CmdTestCredit:
		cred, ok := msg.(*wire.MsgTestCredit)
		if !ok {
			return fmt.Errorf("Error adding test entry credits")
		}
		if err := processTestCredit(cred); err != nil {
			return err
		}

	case wire.CmdEntry:
		if nodeMode == common.SERVER_NODE {
			break
		}

		entry, ok := msg.(*wire.MsgEntry)
		if ok {
			err := processEntry(entry)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Error in processing msg:" + fmt.Sprintf("%+v", msg))
		}

	default:
		return errors.New("Message type unsupported:" + fmt.Sprintf("%+v", msg))
	}

	return nil
}

// processTestCedits assignes credits to a specified publick key for testing
// against the local node. This credit purchase should never propigate across
// the network.
// TODO remove this before production
func processTestCredit(msg *wire.MsgTestCredit) error {
	if _, exists := eCreditMap[msg.ECKey]; !exists {
		eCreditMap[msg.ECKey] = 0
	}
	eCreditMap[msg.ECKey] += msg.Amt
	return nil
}

// processDirBlock validates dir block and save it to factom db.
// similar to blockChain.BC_ProcessBlock
func processDirBlock(msg *wire.MsgDirBlock) error {
	util.Trace()

	// Error condiftion for Milestone 1
	if nodeMode == common.SERVER_NODE {
		return errors.New("Server received msg:" + msg.Command())
	}

	blk, _ := db.FetchDBlockByHeight(msg.DBlk.Header.BlockHeight)
	if blk != nil {
		fmt.Println("DBlock already existing for height:" + string(msg.DBlk.Header.BlockHeight))
		return nil
	}

	msg.DBlk.IsSealed = true
	dchain.AddDBlockToDChain(msg.DBlk)

	db.ProcessDBlockBatch(msg.DBlk) //?? to be removed later

	fmt.Printf("PROCESSOR: MsgDirBlock=%s\n", spew.Sdump(msg.DBlk))
	fmt.Printf("PROCESSOR: dchain=%s\n", spew.Sdump(dchain))

	exportDChain(dchain)

	return nil
}

// processSCBlock validates admin block and save it to factom db.
// similar to blockChain.BC_ProcessBlock
func processSCBlock(msg *wire.MsgSCBlock) error {
    
    //Need to validate against Dchain??
    
    db.ProcessSCBlockBatch(msg.SC)
    
    exportSCChain(scchain)
    
    return nil
}

// processABlock validates admin block and save it to factom db.
// similar to blockChain.BC_ProcessBlock
func processABlock(msg *wire.MsgABlock) error {
	util.Trace()

	// Error condiftion for Milestone 1
	if nodeMode == common.SERVER_NODE {
		return errors.New("Server received msg:" + msg.Command())
	}

	//Need to validate against Dchain??

	db.ProcessABlockBatch(msg.ABlk)

	fmt.Printf("PROCESSOR: MsgABlock=%s\n", spew.Sdump(msg.ABlk))

	exportAChain(achain)

	return nil
}

// processCBlock validates entry credit block and save it to factom db.
// similar to blockChain.BC_ProcessBlock
func processCBlock(msg *wire.MsgECBlock) error {
	util.Trace()

	// Error condiftion for Milestone 1
	if nodeMode == common.SERVER_NODE {
		return errors.New("Server received msg:" + msg.Command())
	}

	//Need to validate against Dchain??

	// check if the block already exists
	h, _ := common.CreateHash(msg.ECBlock)
	cblk, _ := db.FetchECBlockByHash(h)
	if cblk != nil {
		return nil
	}

	db.ProcessECBlockBatch(msg.ECBlock)

	initializeECreditMap(msg.ECBlock)

	// for debugging??
	fmt.Printf("PROCESSOR: MsgCBlock=%s\n", spew.Sdump(msg.ECBlock))
	printCreditMap()

	exportECChain(ecchain)

	return nil
}

// processEBlock validates entry block and save it to factom db.
// similar to blockChain.BC_ProcessBlock
func processEBlock(msg *wire.MsgEBlock) error {
	util.Trace()

	// Error condiftion for Milestone 1
	if nodeMode == common.SERVER_NODE {
		return errors.New("Server received msg:" + msg.Command())
	}

	if msg.EBlk.Header.DBHeight >= dchain.NextBlockHeight || msg.EBlk.Header.DBHeight < 0 {
		return errors.New("MsgEBlock has an invalid DBHeight:" + strconv.Itoa(int(msg.EBlk.Header.DBHeight)))
	}

	dblock := dchain.Blocks[msg.EBlk.Header.DBHeight]

	if dblock == nil {
		return errors.New("MsgEBlock has an invalid DBHeight:" + strconv.Itoa(int(msg.EBlk.Header.DBHeight)))
	}

	msg.EBlk.BuildMerkleRoot()

	validEblock := false
	for _, dbEntry := range dblock.DBEntries {
		if msg.EBlk.MerkleRoot.IsSameAs(dbEntry.MerkleRoot) && dbEntry.ChainID.IsSameAs(msg.EBlk.Header.ChainID) {
			validEblock = true
			break
		}
	}

	if !validEblock {
		return errors.New("Invalid MsgEBlock with height:" + strconv.Itoa(int(msg.EBlk.Header.EBHeight)))
	}

	// create a chain in db if it's not existing
	chain := chainIDMap[msg.EBlk.Header.ChainID.String()]
	if chain == nil {
		chain = new(common.EChain)
		chain.ChainID = msg.EBlk.Header.ChainID

		if msg.EBlk.Header.EBHeight == 0 {
			chain.FirstEntry, _ = db.FetchEntryByHash(msg.EBlk.EBEntries[0].EntryHash)
		}

		db.InsertChain(chain)
		chainIDMap[chain.ChainID.String()] = chain
	} else if chain.FirstEntry == nil && msg.EBlk.Header.EBHeight == 0 {
		chain.FirstEntry, _ = db.FetchEntryByHash(msg.EBlk.EBEntries[0].EntryHash)
		db.InsertChain(chain)
	}

	db.ProcessEBlockBatch(msg.EBlk)


	exportEChain(chain)

	return nil
}

// processEntry validates entry and save it to factom db.
// similar to blockChain.BC_ProcessBlock
func processEntry(msg *wire.MsgEntry) error {
	util.Trace()

	// Error condiftion for Milestone 1
	if nodeMode == common.SERVER_NODE {
		return errors.New("Server received msg:" + msg.Command())
	}

	// store the new entry in db
	entryBinary, _ := msg.Entry.MarshalBinary()
	entryHash := common.Sha(entryBinary)
    b := msg.Entry.ChainID.Bytes()
	db.InsertEntry(entryHash, &entryBinary, msg.Entry, &b)

	fmt.Printf("PROCESSOR: MsgEntry=%s\n", spew.Sdump(msg.Entry))

	return nil
}

/* this should be processed on btcd side
// processFactoidBlock validates factoid block and save it to factom db.
func processFactoidBlock(msg *wire.MsgBlock) error {
	util.Trace()
	fmt.Printf("PROCESSOR: MsgFactoidBlock=%s\n", spew.Sdump(msg))
	return nil
}
*/

/*
// Process a factoid obj message and put it in the process list
func processFactoidTx(msg *wire.MsgInt_FactoidObj) error {

	// Update the credit balance in memory for each EC output
	for k, v := range msg.EntryCredits {
		pubKey := new([32]byte)
		copy(pubKey[:], k.Bytes())
		//credits := int32(creditsPerFactoid * v / 100000000)
		// Update the credit balance in memory
		balance, _ := eCreditMap[pubKey]
		eCreditMap[pubKey] = balance + int32(v)
	}

	// Add to MyPL if Server Node
	if nodeMode == common.SERVER_NODE {
		err := plMgr.AddMyProcessListItem(msg, msg.TxSha, wire.ACK_FACTOID_TX)
		if err != nil {
			return err
		}

	}

	return nil
}
*/

func processRevealEntry(msg *wire.MsgRevealEntry) error {
	e := msg.Entry
	bin, _ := e.MarshalBinary()
	h, _ := wire.NewShaHash(e.Hash().Bytes())

	if c, ok := commitEntryMap[e.Hash().String()]; ok {
		if chainIDMap[e.ChainID.String()] == nil {
			fMemPool.addOrphanMsg(msg, h)
			return fmt.Errorf("This chain is not supported: %s",
				msg.Entry.ChainID.String())
		}

		cred := int32(binary.Size(bin)/1024 + 1)
		if int32(c.Credits) < cred {
			fMemPool.addOrphanMsg(msg, h)
			return fmt.Errorf("Credit needs to paid first before an entry is revealed: %s", e.Hash().String())
			// Add the msg to the Mem pool
			fMemPool.addMsg(msg, h)

			// Add to MyPL if Server Node
			if nodeMode == common.SERVER_NODE {
				if err := plMgr.AddMyProcessListItem(msg, h,
					wire.ACK_REVEAL_ENTRY); err != nil {
					return err
				}
			}
		}

		delete(commitEntryMap, e.Hash().String())
		return nil
	} else if c, ok := commitChainMap[e.Hash().String()]; ok {
		if chainIDMap[e.ChainID.String()] != nil {
			fMemPool.addOrphanMsg(msg, h)
			return fmt.Errorf("This chain is not supported: %s",
				msg.Entry.ChainID.String())
		}

		// add new chain to chainIDMap
		newChain := new(common.EChain)
		newChain.ChainID = e.ChainID
		newChain.FirstEntry = e
		chainIDMap[e.ChainID.String()] = newChain

		cred := int32(binary.Size(bin)/1024 + 1 + 10)
		if int32(c.Credits) < cred {
			fMemPool.addOrphanMsg(msg, h)
			return fmt.Errorf("Credit needs to paid first before an entry is revealed: %s", e.Hash().String())
			// Add the msg to the Mem pool
			fMemPool.addMsg(msg, h)

			// Add to MyPL if Server Node
			if nodeMode == common.SERVER_NODE {
				if err := plMgr.AddMyProcessListItem(msg, h,
					wire.ACK_REVEAL_ENTRY); err != nil {
					return err
				}
			}
		}

		delete(commitChainMap, e.Hash().String())
		return nil
	} else {
		return fmt.Errorf("No commit for entry")
	}

	return nil
}

func processCommitEntry(msg *wire.MsgCommitEntry) error {
	c := msg.CommitEntry

	// check that the CommitChain is fresh
	if !c.InTime() {
		return fmt.Errorf("Cannot commit chain, CommitChain must be timestamped within 24 hours of commit")
	}

	// check to see if the EntryHash has already been committed
	if _, exist := commitEntryMap[c.EntryHash.String()]; exist {
		return fmt.Errorf("Cannot commit entry, entry has already been commited")
	}

	// add to the commitEntryMap
	commitEntryMap[c.EntryHash.String()] = c

	// Server: add to MyPL
	if nodeMode == common.SERVER_NODE {
		h, _ := msg.Sha()
		if err := plMgr.AddMyProcessListItem(msg, &h, wire.ACK_COMMIT_ENTRY); err != nil {
			return err
		}
	}

	return nil
}

func processCommitChain(msg *wire.MsgCommitChain) error {
	c := msg.CommitChain

	// check that the CommitChain is fresh
	if !c.InTime() {
		return fmt.Errorf("Cannot commit chain, CommitChain must be timestamped within 24 hours of commit")
	}

	// check to see if the EntryHash has already been committed
	if _, exist := commitChainMap[c.EntryHash.String()]; exist {
		return fmt.Errorf("Cannot commit chain, first entry for chain already exists")
	}

	// deduct the entry credits from the eCreditMap
	if eCreditMap[c.ECPubKey] < int32(c.Credits) {
		return fmt.Errorf("Not enough credits for CommitChain")
	}
	eCreditMap[c.ECPubKey] -= int32(c.Credits)

	// add to the commitChainMap
	commitChainMap[c.EntryHash.String()] = c

	// Server: add to MyPL
	if nodeMode == common.SERVER_NODE {
		h, _ := msg.Sha()
		if err := plMgr.AddMyProcessListItem(msg, &h,
			wire.ACK_COMMIT_CHAIN); err != nil {
			return err
		}
	}

	return nil
}

func processBuyEntryCredit(pubKey *[32]byte, credits int32, factoidTxHash *common.Hash) error {

	// Update the credit balance in memory
	balance, _ := eCreditMap[pubKey]
	eCreditMap[pubKey] = balance + credits

	return nil
}

// Process Orphan pool before the end of 10 min
func processFromOrphanPool() error {
	for k, msg := range fMemPool.orphans {
		switch msg.Command() {
		case wire.CmdCommitChain:
			msgCommitChain, _ := msg.(*wire.MsgCommitChain)
			err := processCommitChain(msgCommitChain)
			if err != nil {
				return err
			}
			delete(fMemPool.orphans, k)

		case wire.CmdCommitEntry:
			msgCommitEntry, _ := msg.(*wire.MsgCommitEntry)
			err := processCommitEntry(msgCommitEntry)
			if err != nil {
				return err
			}
			delete(fMemPool.orphans, k)

		case wire.CmdRevealEntry:
			msgRevealEntry, _ := msg.(*wire.MsgRevealEntry)
			err := processRevealEntry(msgRevealEntry)
			if err != nil {
				return err
			}
			delete(fMemPool.orphans, k)
		}
	}
	return nil
}

func buildRevealEntry(msg *wire.MsgRevealEntry) {

	chain := chainIDMap[msg.Entry.ChainID.String()]

	// store the new entry in db
	entryBinary, _ := msg.Entry.MarshalBinary()
	entryHash := common.Sha(entryBinary)
    b :=chain.ChainID.Bytes()
	db.InsertEntry(entryHash, &entryBinary, msg.Entry, &b)

	err := chain.NextBlock.AddEBEntry(msg.Entry)

	if err != nil {
		panic("Error while adding Entity to Block:" + err.Error())
	}

}

func buildCommitEntry(msg *wire.MsgCommitEntry) {
	ecchain.NextBlock.AddEntry(msg.CommitEntry)
}

func buildCommitChain(msg *wire.MsgCommitChain) {
	ecchain.NextBlock.AddEntry(msg.CommitChain)
}

/*
func buildFactoidObj(msg *wire.MsgInt_FactoidObj) {
	factoidTxHash := new(common.Hash)
	factoidTxHash.SetBytes(msg.TxSha.Bytes())

	for k, v := range msg.EntryCredits {
		pubkey := new([32]byte)
		copy(pubkey[:], k.Bytes())
		cbEntry := common.NewIncreaseBalance(pubkey, factoidTxHash, int32(v))
		ecchain.NextBlock.AddEntry(cbEntry)
	}
}
*/

func buildRevealChain(msg *wire.MsgRevealChain) {

	newChain := chainIDMap[msg.FirstEntry.ChainID.String()]

	// Store the new chain in db
	db.InsertChain(newChain)

	// Chain initialization
	initEChainFromDB(newChain)

	// store the new entry in db
	entryBinary, _ := newChain.FirstEntry.MarshalBinary()
	entryHash := common.Sha(entryBinary)
    b := newChain.ChainID.Bytes()
	db.InsertEntry(entryHash, &entryBinary, newChain.FirstEntry, &b)

	err := newChain.NextBlock.AddEBEntry(newChain.FirstEntry)

	if err != nil {
		panic(fmt.Sprintf(`Error while adding the First Entry to Block: %s`, err.Error()))
	}
}

// Loop through the Process List items and get the touched chains
// Put End-Of-Minute marker in the entry chains
func buildEndOfMinute(pl *consensus.ProcessList, pli *consensus.ProcessListItem) {
	tempChainMap := make(map[string]*common.EChain)
	items := pl.GetPLItems()
	for i := pli.Ack.Index; i >= 0; i-- {
		if wire.END_MINUTE_1 <= items[i].Ack.Type && items[i].Ack.Type <= wire.END_MINUTE_10 {
			break
		} else if items[i].Ack.Type == wire.ACK_REVEAL_ENTRY && tempChainMap[items[i].Ack.ChainID.String()] == nil {

			chain := chainIDMap[items[i].Ack.ChainID.String()]
			chain.NextBlock.AddEndOfMinuteMarker(pli.Ack.Type)
			// Add the new chain in the tempChainMap
			tempChainMap[chain.ChainID.String()] = chain
		}
	}

	// Add it to the entry credit chain
	entries := ecchain.NextBlock.Body.Entries
	if len(entries) > 0 && entries[len(entries)-1].ECID() != common.ECIDMinuteNumber {
		cbEntry := common.NewMinuteNumber()
		cbEntry.Number = pli.Ack.Type
		ecchain.NextBlock.AddEntry(cbEntry)
	}

	// Add it to the admin chain
	abEntries := achain.NextBlock.ABEntries
	if len(abEntries) > 0 && abEntries[len(abEntries)-1].Type() != common.TYPE_MINUTE_NUM {
		achain.NextBlock.AddEndOfMinuteMarker(pli.Ack.Type)
	}
}

// build Genesis blocks
func buildGenesisBlocks() error {

	/*
		// Send an End of Minute message to the Factoid component to create a genesis block
		eomMsg := &wire.MsgInt_EOM{
			EOM_Type:         wire.FORCE_FACTOID_GENESIS_REBUILD,
			NextDBlockHeight: 0,
		}
		outCtlMsgQueue <- eomMsg
	*/

	// Allocate the first two dbentries for ECBlock and Factoid block
	dchain.AddDBEntry(&common.DBEntry{}) // AdminBlock
	dchain.AddDBEntry(&common.DBEntry{}) // ECBlock
	dchain.AddDBEntry(&common.DBEntry{}) // Factoid block

	// Entry Credit Chain
	cBlock := newEntryCreditBlock(ecchain)
	fmt.Printf("buildGenesisBlocks: cBlock=%s\n", spew.Sdump(cBlock))
	dchain.AddECBlockToDBEntry(cBlock)
	exportECChain(ecchain)

	// Admin chain
	aBlock := newAdminBlock(achain)
	fmt.Printf("buildGenesisBlocks: aBlock=%s\n", spew.Sdump(aBlock))
	dchain.AddABlockToDBEntry(aBlock)
	exportAChain(achain)

    // Simplecoin Genesis Address
    scBlock := newFactoidBlock(scchain)
    data, _ := scBlock.MarshalBinary()
    fmt.Println("\n\n ",common.Sha(data).String(),"\n\n")
    dchain.AddSCBlockToDBEntry(scBlock)
	exportSCChain(scchain)
    
    
	// Directory Block chain
	util.Trace("in buildGenesisBlocks")
	dbBlock := newDirectoryBlock(dchain)

	// Check block hash if genesis block
	if dbBlock.DBHash.String() != common.GENESIS_DIR_BLOCK_HASH {
        
		panic("\nGenesis block hash expected: "+common.GENESIS_DIR_BLOCK_HASH+
		    "\nGenesis block hash found:    " + dbBlock.DBHash.String()+"\n")
	}

	exportDChain(dchain)

	// place an anchor into btc
	placeAnchor(dbBlock)

	return nil
}

// build blocks from all process lists
func buildBlocks() error {
	util.Trace()

	// Allocate the first three dbentries for Admin block, ECBlock and Factoid block
	dchain.AddDBEntry(&common.DBEntry{}) // AdminBlock
	dchain.AddDBEntry(&common.DBEntry{}) // ECBlock
    dchain.AddDBEntry(&common.DBEntry{}) // Simplecoin	

	if plMgr != nil && plMgr.MyProcessList.IsValid() {
		buildFromProcessList(plMgr.MyProcessList)
	}

	// Entry Credit Chain
	ecBlock := newEntryCreditBlock(ecchain)
	dchain.AddECBlockToDBEntry(ecBlock)
	exportECChain(ecchain)

	// Admin chain
	aBlock := newAdminBlock(achain)
	//fmt.Printf("buildGenesisBlocks: aBlock=%s\n", spew.Sdump(aBlock))
	dchain.AddABlockToDBEntry(aBlock)
	exportAChain(achain)
	
	// Factoid chain
	fBlock := newFactoidBlock(scchain)
	//fmt.Printf("buildGenesisBlocks: aBlock=%s\n", spew.Sdump(aBlock))
	dchain.AddSCBlockToDBEntry(fBlock)
	exportSCChain(scchain)	

	// sort the echains by chain id
	var keys []string
	for k := range chainIDMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Entry Chains
	for _, k := range keys {
		chain := chainIDMap[k]
		eblock := newEntryBlock(chain)
		if eblock != nil {
			dchain.AddEBlockToDBEntry(eblock)
		}
		exportEChain(chain)
	}

	// Directory Block chain
	util.Trace("in buildBlocks")
	dbBlock := newDirectoryBlock(dchain)
	// Check block hash if genesis block here??

	// Generate the inventory vector and relay it.
	binary, _ := dbBlock.MarshalBinary()
	commonHash := common.Sha(binary)
	hash, _ := wire.NewShaHash(commonHash.Bytes())
	outMsgQueue <- (&wire.MsgInt_DirBlock{hash})

	exportDChain(dchain)

	// re-initialize the process lit manager
	initProcessListMgr()

	// Initialize timer for the new dblock
	if nodeMode == common.SERVER_NODE {
		timer := &BlockTimer{
			nextDBlockHeight: dchain.NextBlockHeight,
			inCtlMsgQueue:    inCtlMsgQueue,
		}
		go timer.StartBlockTimer()
	}

	// place an anchor into btc
	placeAnchor(dbBlock)

	return nil
}

// Sign the directory block
func SignDirectoryBlock() error {
	// Only Servers can write the anchor to Bitcoin network
	if nodeMode == common.SERVER_NODE && dchain.NextBlockHeight > 0 {
		// get the previous directory block from db
		dbBlock, _ := db.FetchDBlockByHeight(dchain.NextBlockHeight - 1)
		dbHeaderBytes, _ := dbBlock.Header.MarshalBinary()
		identityChainID := common.NewHash() // 0 ID for milestone 1
		sig := serverPrivKey.Sign(dbHeaderBytes)
		achain.NextBlock.AddABEntry(common.NewDBSignatureEntry(identityChainID, sig))
	}

	return nil
}

// Place an anchor into btc
func placeAnchor(dbBlock *common.DirectoryBlock) error {
	util.Trace()
	// Only Servers can write the anchor to Bitcoin network
	if nodeMode == common.SERVER_NODE && dbBlock != nil {
		// todo: need to make anchor as a go routine, independent of factomd
		// same as blockmanager to btcd
		go anchor.SendRawTransactionToBTC(dbBlock.KeyMR, uint64(dbBlock.Header.BlockHeight))
	}
	return nil
}

// build blocks from a process lists
func buildFromProcessList(pl *consensus.ProcessList) error {
	for _, pli := range pl.GetPLItems() {
		if pli.Ack.Type == wire.ACK_COMMIT_CHAIN {
			buildCommitChain(pli.Msg.(*wire.MsgCommitChain))
		} else if pli.Ack.Type == wire.ACK_COMMIT_ENTRY {
			buildCommitEntry(pli.Msg.(*wire.MsgCommitEntry))
		} else if pli.Ack.Type == wire.ACK_REVEAL_CHAIN {
			buildRevealChain(pli.Msg.(*wire.MsgRevealChain))
		} else if pli.Ack.Type == wire.ACK_REVEAL_ENTRY {
			buildRevealEntry(pli.Msg.(*wire.MsgRevealEntry))
		} else if wire.END_MINUTE_1 <= pli.Ack.Type && pli.Ack.Type <= wire.END_MINUTE_10 {
			buildEndOfMinute(pl, pli)
		}
	}

	return nil
}

func newEntryBlock(chain *common.EChain) *common.EBlock {

	// acquire the last block
	block := chain.NextBlock

	if len(block.EBEntries) < 1 {
		//log.Println("No new entry found. No block created for chain: "  + common.EncodeChainID(chain.ChainID))
		return nil
	}

	// Create the block and add a new block for new coming entries

	block.Header.DBHeight = dchain.NextBlockHeight
	block.Header.EntryCount = uint32(len(block.EBEntries))
	block.Header.StartTime = dchain.NextBlock.Header.StartTime

	if devNet {
		block.Header.NetworkID = common.NETWORK_ID_TEST
	} else {
		block.Header.NetworkID = common.NETWORK_ID_EB
	}

	// Create the Entry Block Boday Merkle Root from EB Entries
	hashes := make([]*common.Hash, 0, len(block.EBEntries))
	for _, entry := range block.EBEntries {
		hashes = append(hashes, entry.EntryHash)
	}
	merkle := common.BuildMerkleTreeStore(hashes)
	block.Header.BodyMR = merkle[len(merkle)-1]

	// Create the Entry Block Key Merkle Root from the hash of Header and the Body Merkle Root
	hashes = make([]*common.Hash, 0, 2)
	binaryEBHeader, _ := block.Header.MarshalBinary()
	hashes = append(hashes, common.Sha(binaryEBHeader))
	hashes = append(hashes, block.Header.BodyMR)
	merkle = common.BuildMerkleTreeStore(hashes)
	block.MerkleRoot = merkle[len(merkle)-1] // MerkleRoot is not marshalized in Entry Block
	fmt.Println("block.MerkleRoot:%v", block.MerkleRoot.String())
	blkhash, _ := common.CreateHash(block)
	block.EBHash = blkhash
	log.Println("blkhash:%v", blkhash.Bytes())

	block.IsSealed = true
	chain.NextBlockHeight++
	chain.NextBlock, _ = common.CreateBlock(chain, block, 10)

	//Store the block in db
	db.ProcessEBlockBatch(block)
	log.Println("EntryBlock: block" + strconv.FormatUint(uint64(block.Header.EBHeight), 10) + " created for chain: " + chain.ChainID.String())
	return block
}

func newEntryCreditBlock(chain *common.ECChain) *common.ECBlock {

	// acquire the last block
	block := chain.NextBlock

	if chain.NextBlockHeight != dchain.NextBlockHeight {
		panic("Entry Credit Block height does not match Directory Block height:" + string(dchain.NextBlockHeight))
	}

	block.BuildHeader()

	// Create the block and add a new block for new coming entries
	chain.BlockMutex.Lock()
	chain.NextBlockHeight++
	chain.NextBlock = common.NextECBlock(block)
	chain.BlockMutex.Unlock()

	//Store the block in db
	db.ProcessECBlockBatch(block)
	log.Println("EntryCreditBlock: block" + strconv.FormatUint(uint64(block.Header.DBHeight), 10) + " created for chain: " + chain.ChainID.String())

	return block
}

func newAdminBlock(chain *common.AdminChain) *common.AdminBlock {

	// acquire the last block
	block := chain.NextBlock

	if chain.NextBlockHeight != dchain.NextBlockHeight {
		panic("Admin Block height does not match Directory Block height:" + string(dchain.NextBlockHeight))
	}

	block.Header.EntryCount = uint32(len(block.ABEntries))
	block.Header.BodySize = uint32(block.MarshalledSize() - block.Header.MarshalledSize())
	block.BuildABHash()

	// Create the block and add a new block for new coming entries
	chain.BlockMutex.Lock()
	chain.NextBlockHeight++
	chain.NextBlock, _ = common.CreateAdminBlock(chain, block, 10)
	chain.BlockMutex.Unlock()

	//Store the block in db
	db.ProcessABlockBatch(block)
	log.Println("Admin Block: block" + strconv.FormatUint(uint64(block.Header.DBHeight), 10) + " created for chain: " + chain.ChainID.String())

	return block
}

func newFactoidBlock(chain *common.SCChain) block.ISCBlock {

	// acquire the last block
	currentBlock := chain.NextBlock

	if chain.NextBlockHeight != dchain.NextBlockHeight {
		panic("Factoid Block height does not match Directory Block height:" + strconv.Itoa(int(dchain.NextBlockHeight)))
	}

	//block.BuildHeader()

	// Create the block and add a new block for new coming entries
	chain.BlockMutex.Lock()
	chain.NextBlockHeight++
	chain.NextBlock = block.NewSCBlock(FactoshisPerCredit, chain.NextBlockHeight)
	chain.BlockMutex.Unlock()

	//Store the block in db
	db.ProcessSCBlockBatch(currentBlock)
	log.Println("Factoid chain: block" + " created for chain: " + chain.ChainID.String())

	return currentBlock
}

func newDirectoryBlock(chain *common.DChain) *common.DirectoryBlock {
	util.Trace("**** new Dir Block")
	// acquire the last block
	block := chain.NextBlock

	if devNet {
		block.Header.NetworkID = common.NETWORK_ID_TEST
	} else {
		block.Header.NetworkID = common.NETWORK_ID_EB
	}

	// Create the block add a new block for new coming entries
	chain.BlockMutex.Lock()
	block.Header.EntryCount = uint32(len(block.DBEntries))
	// Calculate Merkle Root for FBlock and store it in header
	if block.Header.BodyMR == nil {
		block.Header.BodyMR, _ = block.BuildBodyMR()
    //  Factoid1 block not in the right place...    
	}
	block.IsSealed = true
	chain.AddDBlockToDChain(block)
	chain.NextBlockHeight++
	chain.NextBlock, _ = common.CreateDBlock(chain, block, 10)
	chain.BlockMutex.Unlock()

	block.DBHash, _ = common.CreateHash(block)
	block.BuildKeyMerkleRoot()

	//Store the block in db
	db.ProcessDBlockBatch(block)

	// Initialize the dirBlockInfo obj in db
	db.InsertDirBlockInfo(common.NewDirBlockInfoFromDBlock(block))
	anchor.UpdateDirBlockInfoMap(common.NewDirBlockInfoFromDBlock(block))

	log.Println("DirectoryBlock: block" + strconv.FormatUint(uint64(block.Header.BlockHeight), 10) + " created for directory block chain: " + chain.ChainID.String())

	// To be improved in milestone 2
	SignDirectoryBlock()

	return block
}

func GetEntryCreditBalance(pubKey *[32]byte) (int32, error) {

	return eCreditMap[pubKey], nil
}

// Validate dir chain from genesis block
func validateDChain(c *common.DChain) error {

	if uint32(len(c.Blocks)) != c.NextBlockHeight {
		return errors.New("Dir chain doesn't have an expected Next Block ID: " + strconv.Itoa(int(c.NextBlockHeight)))
	}

	//prevBlk := c.Blocks[0]
	prevMR, prevBlkHash, err := validateDBlock(c, c.Blocks[0])
	if err != nil {
		return err
	}

	//validate the genesis block
	if prevBlkHash == nil || prevBlkHash.String() != common.GENESIS_DIR_BLOCK_HASH {
		panic("Genesis dir block is not as expected: " + prevBlkHash.String())
	}

	for i := 1; i < len(c.Blocks); i++ {
		if !prevBlkHash.IsSameAs(c.Blocks[i].Header.PrevBlockHash) {
			return errors.New("Previous block hash not matching for Dir block: " + strconv.Itoa(i))
		}
		if !prevMR.IsSameAs(c.Blocks[i].Header.PrevKeyMR) {
			return errors.New("Previous merkle root not matching for Dir block: " + strconv.Itoa(i))
		}
		mr, dblkHash, err := validateDBlock(c, c.Blocks[i])
		if err != nil {
			c.Blocks[i].IsValidated = false
			return err
		}

		prevMR = mr
		prevBlkHash = dblkHash
		c.Blocks[i].IsValidated = true
	}

	return nil
}

// Validate a dir block
func validateDBlock(c *common.DChain, b *common.DirectoryBlock) (merkleRoot *common.Hash, dbHash *common.Hash, err error) {

	bodyMR, err := b.BuildBodyMR()
	if err != nil {
		return nil, nil, err
	}

	if !b.Header.BodyMR.IsSameAs(bodyMR) {
        fmt.Printf("\n\nERROR!!!!!! !b.Header.BodyMR.IsSameAs(bodyMR) fails.\n\n") 
//		return nil, nil, errors.New("Invalid body MR for dir block: " + string(b.Header.BlockHeight))
	}

	for _, dbEntry := range b.DBEntries {
		switch dbEntry.ChainID.String() {
		case ecchain.ChainID.String():
			err := validateCBlockByMR(dbEntry.MerkleRoot)
			if err != nil {
				return nil, nil, err
			}
        case achain.ChainID.String():
            err := validateABlockByMR(dbEntry.MerkleRoot)
            if err != nil {
                return nil, nil, err
            }
        case scchain.ChainID.String():
            err := validateSCBlockByMR(dbEntry.MerkleRoot)
            if err != nil {
                return nil, nil, err
            }
		case wire.FChainID.String():
			err := validateFBlockByMR(dbEntry.MerkleRoot)
			if err != nil {
				return nil, nil, err
			}
		default:
			err := validateEBlockByMR(dbEntry.ChainID, dbEntry.MerkleRoot)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	b.DBHash, _ = common.CreateHash(b)
	b.BuildKeyMerkleRoot()

	return b.KeyMR, b.DBHash, nil
}

func validateFBlockByMR(mr *common.Hash) error {
	// Call BTCD side for factoid block validation??

	return nil
}

func validateCBlockByMR(mr *common.Hash) error {
	cb, _ := db.FetchECBlockByHash(mr)

	if cb == nil {
		return errors.New("Entry Credit block not found in db for merkle root: " + mr.String())
	}

	return nil
}

// Validate Admin Block by merkle root
func validateABlockByMR(mr *common.Hash) error {
	b, _ := db.FetchABlockByHash(mr)

	if b == nil {
		return errors.New("Admin block not found in db for merkle root: " + mr.String())
	}

	return nil
}

// Validate SCBlock by merkle root
func validateSCBlockByMR(mr *common.Hash) error {
    b, _ := db.FetchSCBlockByHash(mr)
    
    if b == nil {
        return errors.New("Simple Coin block not found in db for merkle root: " + mr.String())
    }
    
    return nil
}

func validateEBlockByMR(cid *common.Hash, mr *common.Hash) error {

	eb, _ := db.FetchEBlockByMR(mr)

	if eb == nil {
		return errors.New("Entry block not found in db for merkle root: " + mr.String())
	}

	eb.BuildMerkleRoot()

	if !mr.IsSameAs(eb.MerkleRoot) {
		return errors.New("Entry block's merkle root does not match with: " + mr.String())
	}

	for _, ebEntry := range eb.EBEntries {
		entry, _ := db.FetchEntryByHash(ebEntry.EntryHash)
		if entry == nil {
			return errors.New("Entry not found in db for entry hash: " + ebEntry.EntryHash.String())
		}
	}

	return nil
}

func exportDChain(chain *common.DChain) {
	if len(chain.Blocks) == 0 || procLog.Level() < factomlog.Info {
		//log.Println("no blocks to save for chain: " + string (*chain.ChainID))
		return
	}

	for _, block := range chain.Blocks {
		//the open block is not saved
		if block == nil || block.IsSealed == false {
			continue
		}

		data, err := block.MarshalBinary()
		if err != nil {
			panic(err)
		}

		strChainID := chain.ChainID.String()
		if fileNotExists(dataStorePath + strChainID) {
			err := os.MkdirAll(dataStorePath+strChainID, 0777)
			if err == nil {
				log.Println("Created directory " + dataStorePath + strChainID)
			} else {
				log.Println(err)
			}
		}
		err = ioutil.WriteFile(fmt.Sprintf(dataStorePath+strChainID+"/store.%09d.block", block.Header.BlockHeight), data, 0777)
		if err != nil {
			panic(err)
		}
	}
}

func exportEChain(chain *common.EChain) {
	if procLog.Level() < factomlog.Info {
		return
	}

	eBlocks, _ := db.FetchAllEBlocksByChain(chain.ChainID)
	sort.Sort(util.ByEBlockIDAccending(*eBlocks))

	for _, block := range *eBlocks {

		data, err := block.MarshalBinary()
		if err != nil {
			panic(err)
		}

		strChainID := chain.ChainID.String()
		if fileNotExists(dataStorePath + strChainID) {
			err := os.MkdirAll(dataStorePath+strChainID, 0777)
			if err == nil {
				log.Println("Created directory " + dataStorePath + strChainID)
			} else {
				log.Println(err)
			}
		}

		err = ioutil.WriteFile(fmt.Sprintf(dataStorePath+strChainID+"/store.%09d.block", block.Header.DBHeight), data, 0777)
		if err != nil {
			panic(err)
		}
	}
}

func exportECChain(chain *common.ECChain) {
	if procLog.Level() < factomlog.Info {
		return
	}
	// get all ecBlocks from db
	ecBlocks, _ := db.FetchAllECBlocks()
	sort.Sort(util.ByECBlockIDAccending(ecBlocks))

	for _, block := range ecBlocks {
		data, err := block.MarshalBinary()
		if err != nil {
			panic(err)
		}

		strChainID := chain.ChainID.String()
		if fileNotExists(dataStorePath + strChainID) {
			err := os.MkdirAll(dataStorePath+strChainID, 0777)
			if err == nil {
				log.Println("Created directory " + dataStorePath + strChainID)
			} else {
				log.Println(err)
			}
		}
		err = ioutil.WriteFile(fmt.Sprintf(dataStorePath+strChainID+"/store.%09d.block", block.Header.DBHeight), data, 0777)
		if err != nil {
			panic(err)
		}
	}
}

func exportAChain(chain *common.AdminChain) {
    if procLog.Level() < factomlog.Info {
        return
    }
    // get all aBlocks from db
    aBlocks, _ := db.FetchAllABlocks()
    sort.Sort(util.ByABlockIDAccending(aBlocks))
    
    for _, block := range aBlocks {
        
        data, err := block.MarshalBinary()
        if err != nil {
            panic(err)
        }
        
        strChainID := chain.ChainID.String()
        if fileNotExists(dataStorePath + strChainID) {
            err := os.MkdirAll(dataStorePath+strChainID, 0777)
            if err == nil {
                log.Println("Created directory " + dataStorePath + strChainID)
            } else {
                log.Println(err)
            }
        }
        err = ioutil.WriteFile(fmt.Sprintf(dataStorePath+strChainID+"/store.%09d.block", block.Header.DBHeight), data, 0777)
        if err != nil {
            panic(err)
        }
    }
}


func exportSCChain(chain *common.SCChain) {
	if procLog.Level() < factomlog.Info {
		return
	}
	// get all aBlocks from db
	scBlocks, _ := db.FetchAllSCBlocks()
	sort.Sort(util.BySCBlockIDAccending(scBlocks))

	for _, block := range scBlocks {

		data, err := block.MarshalBinary()
		if err != nil {
			panic(err)
		}

		strChainID := chain.ChainID.String()
		if fileNotExists(dataStorePath + strChainID) {
			err := os.MkdirAll(dataStorePath+strChainID, 0777)
			if err == nil {
				log.Println("Created directory " + dataStorePath + strChainID)
			} else {
				log.Println(err)
			}
		}
		err = ioutil.WriteFile(fmt.Sprintf(dataStorePath+strChainID+"/store.%09d.block", block.GetDBHeight()), data, 0777)
		if err != nil {
			panic(err)
		}
	}
}

func initDChain() {
	dchain = new(common.DChain)

	//Initialize the Directory Block Chain ID
	dchain.ChainID = new(common.Hash)
	barray := common.D_CHAINID
	dchain.ChainID.SetBytes(barray)

	// get all dBlocks from db
	dBlocks, _ := db.FetchAllDBlocks()
	sort.Sort(util.ByDBlockIDAccending(dBlocks))

	//fmt.Printf("initDChain: dBlocks=%s\n", spew.Sdump(dBlocks))

	dchain.Blocks = make([]*common.DirectoryBlock, len(dBlocks), len(dBlocks)+1)

	for i := 0; i < len(dBlocks); i = i + 1 {
		if dBlocks[i].Header.BlockHeight != uint32(i) {
			panic("Error in initializing dChain:" + dchain.ChainID.String())
		}
		dBlocks[i].Chain = dchain
		dBlocks[i].IsSealed = true
		dBlocks[i].IsSavedInDB = true
		dchain.Blocks[i] = &dBlocks[i]
	}

	// double check the block ids
	for i := 0; i < len(dchain.Blocks); i = i + 1 {
		if uint32(i) != dchain.Blocks[i].Header.BlockHeight {
			panic(errors.New("BlockID does not equal index for chain:" + dchain.ChainID.String() + " block:" + fmt.Sprintf("%v", dchain.Blocks[i].Header.BlockHeight)))
		}
	}

	//Create an empty block and append to the chain
	if len(dchain.Blocks) == 0 {
		dchain.NextBlockHeight = 0
		dchain.NextBlock, _ = common.CreateDBlock(dchain, nil, 10)
		// Update dir block height cache in db
		h, _ := common.HexToHash(common.GENESIS_DIR_BLOCK_HASH)
		db.UpdateBlockHeightCache( 0, h)		
	} else {
		dchain.NextBlockHeight = uint32(len(dchain.Blocks))
		dchain.NextBlock, _ = common.CreateDBlock(dchain, dchain.Blocks[len(dchain.Blocks)-1], 10)
		// Update dir block height cache in db
		db.UpdateBlockHeightCache(dchain.NextBlockHeight -1, dchain.NextBlock.Header.PrevBlockHash)		
	}

	exportDChain(dchain)

	//Double check the sealed flag
	if dchain.NextBlock.IsSealed == true {
		panic("dchain.Blocks[dchain.NextBlockID].IsSealed for chain:" + dchain.ChainID.String())
	}


}

func initECChain() {

	eCreditMap = make(map[*[32]byte]int32)

	//Initialize the Entry Credit Chain ID
	ecchain = common.NewECChain()

	// get all ecBlocks from db
	ecBlocks, _ := db.FetchAllECBlocks()
	sort.Sort(util.ByECBlockIDAccending(ecBlocks))

	for i, v := range ecBlocks {
		if v.Header.DBHeight != uint32(i) {
			panic("Error in initializing dChain:" + ecchain.ChainID.String() + " DBHeight:" + strconv.Itoa(int(v.Header.DBHeight)) + " i:" + strconv.Itoa(i))
		}

		// Calculate the EC balance for each account
		initializeECreditMap(&v)
	}

	//Create an empty block and append to the chain
	if len(ecBlocks) == 0 || dchain.NextBlockHeight == 0 {
		ecchain.NextBlockHeight = 0
		ecchain.NextBlock = common.NewECBlock()
	} else {
		// Entry Credit Chain should have the same height as the dir chain
		ecchain.NextBlockHeight = dchain.NextBlockHeight
		ecchain.NextBlock = common.NextECBlock(&ecBlocks[ecchain.NextBlockHeight-1])
	}

	// create a backup copy before processing entries
	copyCreditMap(eCreditMap, eCreditMapBackup)
	exportECChain(ecchain)

	// ONly for debugging
	//printCChain()
	//printCreditMap()
	//printPaidEntryMap()

}

func initAChain() {

	//Initialize the Admin Chain ID
	achain = new(common.AdminChain)
	achain.ChainID = new(common.Hash)
	achain.ChainID.SetBytes(common.ADMIN_CHAINID)

	// get all aBlocks from db
	aBlocks, _ := db.FetchAllABlocks()
	sort.Sort(util.ByABlockIDAccending(aBlocks))

	fmt.Printf("initAChain: aBlocks=%s\n", spew.Sdump(aBlocks))

	// double check the block ids
	for i := 0; i < len(aBlocks); i = i + 1 {
		//if uint32(i) != aBlocks[i].Header.DBHeight {
		//	panic(errors.New("BlockID does not equal index for chain:" + achain.ChainID.String() + " block:" + fmt.Sprintf("%v", aBlocks[i].Header.DBHeight)))
		//}
	}

	//Create an empty block and append to the chain
	if len(aBlocks) == 0 || dchain.NextBlockHeight == 0 {
		achain.NextBlockHeight = 0
		achain.NextBlock, _ = common.CreateAdminBlock(achain, nil, 10)

	} else {
		// Entry Credit Chain should have the same height as the dir chain
		achain.NextBlockHeight = dchain.NextBlockHeight
		achain.NextBlock, _ = common.CreateAdminBlock(achain, &aBlocks[achain.NextBlockHeight-1], 10)
	}

	exportAChain(achain)

}

func initSCChain() {
    
    //Initialize the Admin Chain ID
    scchain = new(common.SCChain)
    scchain.ChainID = new(common.Hash)
    scchain.ChainID.SetBytes(sc.FACTOID_CHAINID)
    
    // get all aBlocks from db
    scBlocks, _ := db.FetchAllSCBlocks()
    sort.Sort(util.BySCBlockIDAccending(scBlocks))
        
    // double check the block ids
    for i := 0; i < len(scBlocks); i = i + 1 {
        //if uint32(i) != aBlocks[i].Header.DBHeight {
        //	panic(errors.New("BlockID does not equal index for chain:" + achain.ChainID.String() + " block:" + fmt.Sprintf("%v", aBlocks[i].Header.DBHeight)))
        //}
    }
    
    //Create an empty block and append to the chain
    if len(scBlocks) == 0 || dchain.NextBlockHeight == 0 {
        scchain.NextBlockHeight = 0
        
        // THIS IS IN TWO PLACES HERE! THEY NEED TO MATCH!
        scchain.NextBlock = block.GetGenesisBlock(1000000,10,200000000000)  
        data, _ := scchain.NextBlock.MarshalBinary()
        fmt.Println("\n\n ",common.Sha(data).String(),"\n\n")
    } else {
        // Entry Credit Chain should have the same height as the dir chain
        scchain.NextBlockHeight = dchain.NextBlockHeight
        scchain.NextBlock = block.NewSCBlock(FactoshisPerCredit, dchain.NextBlockHeight)
    }
    
    exportSCChain(scchain)
    
}


func initEChains() {

	chainIDMap = make(map[string]*common.EChain)

	chains, err := db.FetchAllChains()

	if err != nil {
		panic(err)
	}

	for _, chain := range chains {
		var newChain = chain
		chainIDMap[newChain.ChainID.String()] = &newChain
		exportEChain(&chain)
	}

}

func initializeECreditMap(block *common.ECBlock) {
	for _, entry := range block.Body.Entries {
		// Only process: ECIDChainCommit, ECIDEntryCommit, ECIDBalanceIncrease
		switch entry.ECID() {
		case common.ECIDChainCommit:
			e := entry.(*common.CommitChain)
			eCreditMap[e.ECPubKey] += int32(e.Credits)
		case common.ECIDEntryCommit:
			e := entry.(*common.CommitEntry)
			eCreditMap[e.ECPubKey] += int32(e.Credits)
		case common.ECIDBalanceIncrease:
			e := entry.(*common.IncreaseBalance)
			eCreditMap[e.ECPubKey] += int32(e.Credits)
		}
	}
}

// Initialize server private key and server public key for milestone 1
func initServerKeys() {
	if nodeMode == common.SERVER_NODE {
		var err error
		serverPrivKey, err = common.NewPrivateKeyFromHex(serverPrivKeyHex)
		if err != nil {
			panic("Cannot parse Server Private Key from configuration file: " + err.Error())
		}
	} else {
		serverPubKey = common.PubKeyFromString(common.SERVER_PUB_KEY)
	}
}

func initProcessListMgr() {
	plMgr = consensus.NewProcessListMgr(dchain.NextBlockHeight, 1, 10)

}

func getPrePaidChainKey(entryHash *common.Hash, chainIDHash *common.Hash) string {
	return chainIDHash.String() + entryHash.String()
}

func copyCreditMap(
	originalMap map[*[32]byte]int32,
	newMap map[*[32]byte]int32) {
	newMap = make(map[*[32]byte]int32)

	// copy every element from the original map
	for k, v := range originalMap {
		newMap[k] = v
	}

}

func printCreditMap() {
	fmt.Println("eCreditMap:")
	for key := range eCreditMap {
		fmt.Println("Key: %x Value %d\n", key, eCreditMap[key])
	}
}
