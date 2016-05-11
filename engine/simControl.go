// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package engine

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/wsapi"
)

var _ = fmt.Print

func SimControl(listenTo int) {

	var _ = time.Sleep

	var summary bool
	var watchPL bool
	var watchMessages bool

	for {

		l := make([]byte, 100)
		var err error
		if _, err = os.Stdin.Read(l); err != nil {
			l = []byte("no command") // This is a hack to handle running in the background. (Eg: as a detatched process)
			// Being unable to read from StdIn gives error, this pretends like "no command" was typed, which causes nothing (unlike simply hitting return)
		}

		// This splits up the command at anycodepoint that is not a letter, number of punctuation, so usually by spaces.
		parseFunc := func(c rune) bool {
			return !unicode.IsLetter(c) && !unicode.IsNumber(c) && !unicode.IsPunct(c)
		}
		// cmd is not a list of the parameters, much like command line args show up in args[]
		cmd := strings.FieldsFunc(string(l), parseFunc)
		if 0 == len(cmd) {
			cmd = []string{"h"}
		}
		b := string(cmd[0])
		v, err := strconv.Atoi(string(b))
		if err == nil && v >= 0 && v < len(fnodes) && fnodes[listenTo].State != nil {
			listenTo = v
			os.Stderr.WriteString(fmt.Sprintf("Switching to Node %d\n", listenTo))
			wsapi.SetState(fnodes[listenTo].State)
		} else {
			// fmt.Printf("Parsing command, found %d elements.  The first element is: %+v / %s \n Full command: %+v\n", len(cmd), b[0], string(b), cmd)
			switch {
			case 's' == b[0]:
				summary = !summary
				if summary {
					os.Stderr.WriteString("--Print Summary On--\n")
					go printSummary(&summary, &listenTo)
				} else {
					os.Stderr.WriteString("--Print Summary Off--\n")
				}
			case 'p' == b[0]:
				watchPL = !watchPL
				if watchPL {
					os.Stderr.WriteString("--Print Process Lists On--\n")
					go printProcessList(&watchPL, &listenTo)
				} else {
					os.Stderr.WriteString("--Print Process Lists Off--\n")
				}
			case 'a' == b[0]:
				mLog.all = false
				for _, fnode := range fnodes {
					fnode.State.SetOut(false)
				}
				if listenTo < 0 || listenTo > len(fnodes) {
					fmt.Println("Select a node first")
					break
				}
				f := fnodes[listenTo]
				fmt.Println("-----------------------------", f.State.FactomNodeName, "--------------------------------------", string(b[:len(b)]))
				if len(b) < 2 {
					break
				}
				ht, err := strconv.Atoi(string(b[1:]))
				if err != nil {
					fmt.Println(err, "Dump Adminblock block with an  where n = blockheight, i.e. 'a10'")
				} else {
					msg, err := f.State.LoadDBState(uint32(ht))
					if err == nil && msg != nil {
						dsmsg := msg.(*messages.DBStateMsg)
						ABlock := dsmsg.AdminBlock
						fmt.Println(ABlock.String())
					} else {
						fmt.Println("Error: ", err, msg)
					}
				}
			case 'f' == b[0]:
				mLog.all = false
				for _, fnode := range fnodes {
					fnode.State.SetOut(false)
				}
				if listenTo < 0 || listenTo > len(fnodes) {
					fmt.Println("Select a node first")
					break
				}
				f := fnodes[listenTo]
				fmt.Println("-----------------------------", f.State.FactomNodeName, "--------------------------------------", string(b[:len(b)]))
				if len(b) < 2 {
					break
				}
				ht, err := strconv.Atoi(string(b[1:]))
				if err != nil {
					fmt.Println(err, "Dump Factoid block with fn  where n = blockheight, i.e. 'f10'")
				} else {
					msg, err := f.State.LoadDBState(uint32(ht))
					if err == nil && msg != nil {
						dsmsg := msg.(*messages.DBStateMsg)
						FBlock := dsmsg.FactoidBlock
						fmt.Printf(FBlock.String())
					} else {
						fmt.Println("Error: ", err, msg)
					}
				}
			case 'd' == b[0]:
				mLog.all = false
				for _, fnode := range fnodes {
					fnode.State.SetOut(false)
				}
				if listenTo < 0 || listenTo > len(fnodes) {
					fmt.Println("Select a node first")
					break
				}
				f := fnodes[listenTo]
				fmt.Println("-----------------------------", f.State.FactomNodeName, "--------------------------------------", string(b[:len(b)]))
				if len(b) < 2 {
					break
				}
				ht, err := strconv.Atoi(string(b[1:]))
				if err != nil {
					fmt.Println(err, "Dump Directory block with dn  where n = blockheight, i.e. 'd10'")
				} else {
					msg, err := f.State.LoadDBState(uint32(ht))
					if err == nil && msg != nil {
						dsmsg := msg.(*messages.DBStateMsg)
						DBlock := dsmsg.DirectoryBlock
						fmt.Printf(DBlock.String())
					} else {
						fmt.Println("Error: ", err, msg)
					}
				}
			case 'x' == b[0]:

				if listenTo >= 0 && listenTo < len(fnodes) {
					f := fnodes[listenTo]
					v := f.State.GetNetStateOff()
					if v {
						os.Stderr.WriteString("Bring "+f.State.FactomNodeName+" Back onto the network\n")
					}else{
						os.Stderr.WriteString("Take  "+f.State.FactomNodeName+" off the network\n")
					}
					f.State.SetNetStateOff(!v)
				}

			case 'm' == b[0]:
				watchMessages = !watchMessages
				if watchMessages {
					os.Stderr.WriteString("--Print Messages On--\n")
					go printMessages(&watchMessages, &listenTo)
				} else {
					os.Stderr.WriteString("--Print Messages Off--\n")
				}
			case 'l' == b[0]:
				msg := messages.NewAddServerMsg(fnodes[listenTo].State, 0)
				fnodes[listenTo].State.InMsgQueue() <- msg
				os.Stderr.WriteString(fmt.Sprintln("Attempting to make", fnodes[listenTo].State.GetFactomNodeName(), "a Leader"))
				fallthrough
			case 'n' == b[0]:
				fnodes[listenTo].State.SetOut(false)
				listenTo++
				if listenTo >= len(fnodes) {
					listenTo = 0
				}
				fnodes[listenTo].State.SetOut(true)
				os.Stderr.WriteString(fmt.Sprint("\r\nSwitching to Node ", listenTo, "\r\n"))
				wsapi.SetState(fnodes[listenTo].State)

			case 'h' == b[0]:
				os.Stderr.WriteString("-------------------------------------------------------------------------------\n")
				os.Stderr.WriteString("h or ENTER    Shows this help\n")
				os.Stderr.WriteString("aN            Show Admin block     N. Indicate node eg:\"a5\" to shows blocks for that node.\n")
				os.Stderr.WriteString("fN            Show Factoid block   N. Indicate node eg:\"f5\" to shows blocks for that node.\n")
				os.Stderr.WriteString("dN            Show Directory block N. Indicate node eg:\"d5\" to shows blocks for that node.\n")
				os.Stderr.WriteString("m             Show Messages as they are passed through the simulator.\n")
				os.Stderr.WriteString("s             Show the state of all nodes as their state changes in the simulator.\n")
				os.Stderr.WriteString("p             Show the process lists and directory block states as they change.\n")
				os.Stderr.WriteString("n             Change the focus to the next node.\n")
				os.Stderr.WriteString("l             Make focused node the Leader.\n")
				os.Stderr.WriteString("x             Take the given node out of the netork or bring an offline node back in.\n")
				os.Stderr.WriteString("h or <enter>  Show help\n")
				os.Stderr.WriteString("\n")
				os.Stderr.WriteString("Most commands are case insensitive.\n")
				os.Stderr.WriteString("-------------------------------------------------------------------------------\n\n")
			// -- add node (and give its connections or topology)
			// TODO JAYJAY Need to make an option that causes the p2p network to print out all messsages it gets and sends, for easier debugging.

			default:
			}
		}
	}

}

func printSummary(summary *bool, listenTo *int) {
	out := ""
	for {
		if *summary {
			prt := ""
			for _, f := range fnodes {
				f.State.SetOut(false)
				prt = prt + fmt.Sprintf("%8s %s\n", f.State.FactomNodeName, f.State.ShortString())
			}
			if *listenTo >= 0 && *listenTo < len(fnodes) {
				state := fnodes[*listenTo].State
				prt = prt + fmt.Sprintf("   %s\n", fnodes[*listenTo].State.GetFactomNodeName())
				prt = prt + fmt.Sprintf("      FollowerMsgQueue       %d\n", len(state.FollowerMsgQueue()))
				prt = prt + fmt.Sprintf("      InMsgQueue             %d\n", len(state.InMsgQueue()))
				prt = prt + fmt.Sprintf("      LeaderMsgQueue         %d\n", len(state.LeaderMsgQueue()))
				prt = prt + fmt.Sprintf("      stall Queue            %d\n", len(state.Stall()))
				prt = prt + fmt.Sprintf("      TimerMsgQueue          %d\n", len(state.TimerMsgQueue()))
				prt = prt + fmt.Sprintf("      NetworkOutMsgQueue     %d\n", len(state.NetworkOutMsgQueue()))
				prt = prt + fmt.Sprintf("      NetworkInvalidMsgQueue %d\n", len(state.NetworkInvalidMsgQueue()))
			}
			if prt != out {
				fmt.Println(prt)
				out = prt
			}
		} else {
			return
		}
		time.Sleep(time.Second * 2)
	}
}

func printProcessList(watchPL *bool, listenTo *int) {
	out := ""
	for {
		if *watchPL {
			fnode := fnodes[*listenTo]
			nprt := fnode.State.DBStates.String()
			b := fnode.State.GetHighestRecordedBlock()
			nprt = nprt + fnode.State.ProcessLists.String()
			pl := fnode.State.ProcessLists.Get(b)
			nprt = nprt + pl.PrintMap()

			if out != nprt {
				fmt.Println(nprt)
				out = nprt
			}

		} else {
			return
		}
		time.Sleep(time.Second)
	}
}

func printMessages(Messages *bool, listenTo *int) {
	fmt.Println("Printing Messages")
	for {
		if *Messages {
			fnode := fnodes[*listenTo]
			fnode.MLog.PrtMsgs(fnode.State)
		} else {
			fmt.Println("Done Printing Messages!")
			return
		}
		time.Sleep(2 * time.Second)
	}
}
