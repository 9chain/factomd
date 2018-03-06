package elections

import (
	"fmt"
	"time"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/state"
)

var _ = fmt.Print
var _ = time.Tick

type FaultId struct {
	Dbheight int
	Minute   int
	Round    int
}

type Elections struct {
	FedID     interfaces.IHash
	Name      string
	Sync      []bool // List of servers that have Synced
	Federated []interfaces.IServer
	Audit     []interfaces.IServer
	FPriority []interfaces.IHash
	APriority []interfaces.IHash
	DBHeight  int               // Height of this election
	Minute    int               // Minute of this election (-1 for a DBSig)
	VMIndex   int               // VMIndex of this election
	Msgs      []interfaces.IMsg // Messages we are collecting in this election.  Look here for what's missing.
	Input     interfaces.IQueue
	Output    interfaces.IQueue
	Round     []int
	Electing  int // This is the federated Server index that we are looking to replace
	State     interfaces.IState
	feedback  []string
	VName     string
	Msg       interfaces.IMsg // The missing message as supplied by the volunteer
	Ack       interfaces.IMsg // The missing ack for the message as supplied by the volunteer

	Sigs [][]interfaces.IHash // Signatures from the Federated Servers for a given round.

	Adapter interfaces.IElectionAdapter

	Timeout time.Duration
	EndFault chan FaultId
	PrevElection FaultId
}

func (e *Elections) AdapterStatus() string {
	if e.Adapter != nil {
		return e.Adapter.Status()
	}
	return ""
}

// Add the given sig list to the list of signatures for the given round.
func (e *Elections) AddSigs(round int, sigs []interfaces.IHash) {
	for len(e.Sigs) <= round {
		e.Sigs = append(e.Sigs)
	}
}

func (e *Elections) NewFeedback() {
	e.feedback = make([]string, len(e.Federated)+len(e.Audit))
}

func (e *Elections) FeedBackStr(v string, fed bool, index int) string {

	if !fed {
		index = index + len(e.Federated)
	}

	// If I have no feedback, then get some.
	if e.feedback == nil || len(e.feedback) == 0 {
		e.NewFeedback()
	}

	// Add the status if it is in my known range.
	if index >= 0 && index < len(e.feedback) {
		e.feedback[index] = v
	}

	// Make a string of the status.
	r := ""
	for _, v := range e.feedback {
		r = r + fmt.Sprintf("%4s ", v)
	}
	if e.Msg != nil {
		r = r + " " + e.VName
		r = r + " " + e.Msg.String()
	}
	return r
}

func (e *Elections) String() string {
	str := fmt.Sprintf("eee %10s %s  dbht %d\n", e.State.GetFactomNodeName(), e.Name, e.DBHeight)
	str += fmt.Sprintf("eee %10s  %s\n", e.State.GetFactomNodeName(), "Federated Servers")
	for _, s := range e.Federated {
		str += fmt.Sprintf("eee %10s     %x\n", e.State.GetFactomNodeName(), s.GetChainID().Bytes())
	}
	str += fmt.Sprintf("eee %10s  %s\n", e.State.GetFactomNodeName(), "Audit Servers")
	for _, s := range e.Audit {
		str += fmt.Sprintf("eee %10s     %x\n", e.State.GetFactomNodeName(), s.GetChainID().Bytes())
	}
	return str
}

func (e *Elections) Print() {
	fmt.Println(e.String())
}

// Returns the index of the given server. -1 if it isn't a Federated Server
func (e *Elections) LeaderIndex(server interfaces.IHash) int {
	for i, b := range e.Federated {
		if server.IsSameAs(b.GetChainID()) {
			return i
		}
	}
	return -1
}

// Returns the index of the given server. -1 if it isn't a Audit Server
func (e *Elections) AuditIndex(server interfaces.IHash) int {
	for i, b := range e.Audit {
		if server.IsSameAs(b.GetChainID()) {
			return i
		}
	}
	return -1
}

func (e *Elections) AuditPriority() int {
	// Get the priority order list of audit servers in the priority order
	for len(e.Round) <= e.Electing {
		e.Round = append(e.Round, 0)
	}
	e.APriority = Order(e.Audit, e.DBHeight, e.Minute, e.Electing, e.Round[e.Electing])
	auditIdx := MaxIdx(e.APriority)
	return auditIdx
}

// Runs the main loop for elections for this instance of factomd
func Run(s *state.State) {
	e := new(Elections)
	s.Elections = e
	e.State = s
	e.Name = s.FactomNodeName
	e.Input = s.ElectionsQueue()
	e.Output = s.InMsgQueue()
	e.Electing = -1

	e.Timeout = 10 * time.Second

	// Actually run the elections
	for {
		msg := e.Input.BlockingDequeue().(interfaces.IElectionMsg)
		s.LogMessage("Election", "exec", msg.(interfaces.IMsg))
		msg.ElectionProcess(s, e)
	}
}