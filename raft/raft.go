package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"log"
	"math/rand"
	"sync"
	"time"
)
import "github.com/roessland/distributed-systems/labrpc"

// A peer's electionTimeout is set randomly
// to a duration in the interval [T, 2T] and never changes.
const baseElectionTimeout time.Duration = time.Duration(300) * time.Millisecond

// import "bytes"
// import "github.com/roessland/distributed/systems/labgob"

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in Lab 3 you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh; at that point you can add fields to
// ApplyMsg, but set CommandValid to false for these other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

// Server states. Initial state is Follower.
// See Figure 4 under §5.1 in extended Raft paper.

// Follower:
//	- Times out, starts election and becomes Candidate
// Candidate:
//  - Discovers current leader or new term and becomes Follower
//  - Times out, becomes Candidate again in a new election
//	- Receives majority vote and becomes leader

type ServerState int

const (
	Follower = iota + 1
	Candidate
	Leader
)

type Command string

type LogEntry struct {
	Term    int
	Command Command
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	// Initially Follower
	state ServerState

	// Update when receiving AppendEntries from current leader and when
	// granting vote to a candidate
	lastHeartbeat time.Time

	// Persistent state -- updated on stable storage before responding to RPCs
	currentTerm int        // latest term peer has seen
	votedFor    int        // candidateId that received vote in current term (or 0 if none)
	log         []LogEntry // commands for state machine

	// Volatile state
	commitIndex int // index of highest log entry known to be committed
	lastApplied int // index of highest executed log entry

	// Volatile state on leaders
	// (Reinitialized after election)
	nextIndex  []int // Index of last log entry to send to peers
	matchIndex []int // index of highest log entry known to be replicated on peers
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (term int, isLeader bool) {
	log.Print(rf.me, "GetState requesting lock")
	rf.mu.Lock()
	log.Print(rf.me, "GetState got lock")
	term = rf.currentTerm
	isLeader = rf.state == Leader
	log.Print(rf.me, "GetState returning lock")
	rf.mu.Unlock()
	return
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term         int // Candidate's term
	CandidateId  int // Candidate requesting vote
	LastLogIndex int // Index of candidate's last log entry
	LastLogTerm  int // Term of candidate's last log entry
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Term        int  // currentTerm, for candidate to update itself
	VoteGranted bool // true means candidate received vote
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	log.Print(rf.me, "RequestVote: requesting lock")
	rf.mu.Lock()
	log.Print(rf.me, "RequestVote: got lock")
	defer func() {
		log.Print(rf.me, "RequestVote: returning lock")
		rf.mu.Unlock()
	}()

	// §5.1 in extended Raft paper
	if args.Term < rf.currentTerm {
		log.Print(rf.me, " does NOT vote for ", args.CandidateId,
			" since his term is ", args.Term, " and my term is ", rf.currentTerm)
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return
	}
	// §5.2, §5.4 in extended Raft paper
	canSendVote := rf.votedFor == -1 || rf.votedFor == args.CandidateId
	myLastLogTerm := rf.log[len(rf.log)-1].Term
	candidateLogIsUpToDate := args.LastLogTerm >= myLastLogTerm && args.LastLogIndex >= len(rf.log)-1
	if canSendVote && candidateLogIsUpToDate {
		log.Println(rf.me, " votes for for ", args.CandidateId)
		reply.Term = rf.currentTerm
		reply.VoteGranted = true
		rf.votedFor = args.CandidateId
		return
	}

	// Otherwise don't grant vote
	log.Print(rf.me, " does NOTT vote for ", args.CandidateId,
		", canSendVote is ", canSendVote,
		", and candidateLogIsUpToDate is ", candidateLogIsUpToDate,
		", and myLastLogTerm is ", myLastLogTerm,
		", already voted for ", rf.votedFor,
		", and my log is ", rf.log)
	reply.VoteGranted = false
	reply.Term = rf.currentTerm
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

type AppendEntriesArgs struct {
	Term         int        // leader's term
	LeaderId     int        // so follower can redirect clients
	PrevLogIndex int        // index of log entry immediately preceding new ones
	PrevLogTerm  int        // term of PrevLogIndex entry
	Entries      []LogEntry // log entries to store (empty for heartbeat)
	LeaderCommit int        // leader's commitIndex
}

type AppendEntriesReply struct {
	Term    int  // currentTerm, for leader to update itself
	Success bool // true if follower contained entry matching PrevLogIndex and PrevLogTerm
}

// Invoked by leader to replicate log entries. (§5.3)
// Also used as heartbeat. (§5.2)
func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	log.Print(rf.me, "AppendEntries: requesting lock")
	rf.mu.Lock()
	log.Print(rf.me, "AppendEntries: got lock")
	defer func() {
		log.Println(rf.me, "AppendEntries: returning lock")
		rf.mu.Unlock()
	}()
	rf.lastHeartbeat = time.Now()

	// (7). If RPC request or response contains term T > currentTerm: set
	// currentTerm = T, convert to follower (§5.1)
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.votedFor = -1
		rf.state = Follower
	}

	// 1. Reply false if term < currentTerm (§5.1)
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.Success = false
		return
	}
	// 2. Reply false if log doesn't contain an entry at PrevLogIndex whose
	// term matches PrevLogTerm (§5.3)
	if rf.log[args.PrevLogIndex].Term != args.PrevLogTerm {
		reply.Term = rf.currentTerm
		reply.Success = false
		return
	}
	// 3. If an existing entry conflicts with a new one (same index but
	// different terms), delete the existing entry and all that follow it
	// (§5.3)
	// TODO
	// 4. Append any new entries not already in the log
	// TODO
	// 5. If LeaderCommit > commitIndex,
	// set commitIndex = min(LeaderCommit, index of last new entry)
	if args.LeaderCommit > rf.commitIndex {
		// TODO
	}
	// (6). If commitIndex > lastApplied: increment lastApplied, apply
	// log[lastApplied] to state machine (§5.3)
	if rf.commitIndex > rf.lastApplied {
		// TODO apply to state machine
	}
}

func (rf *Raft) sendAppendEntries(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).

	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}

func (rf *Raft) startElection() {
	log.Println("Peer ", rf.me, " didn't get heartbeats and is starting an election")
}

func (rf *Raft) holdElections() {
	electionTimeout := time.Duration((1.0 + rand.Float64()) * float64(baseElectionTimeout))
	for {
		time.Sleep(time.Duration(10) * time.Millisecond)
		//log.Print(rf.me, "holdElections requesting lock")
		rf.mu.Lock()
		//log.Print(rf.me, "holdElections got lock")
		state := rf.state
		timeSinceLastHeartbeat := time.Since(rf.lastHeartbeat)
		//log.Print(rf.me, "holdElections returning lock")
		rf.mu.Unlock()
		if state == Leader {
			continue // Only followers and candidates can start a new election
		}
		if timeSinceLastHeartbeat > electionTimeout {
			log.Print(rf.me, "holdElections#2 requestion lock")
			rf.mu.Lock()
			log.Println(rf.me, "holdElections#2 got lock")
			rf.state = Candidate
			rf.currentTerm++
			rf.votedFor = rf.me
			rf.lastHeartbeat = time.Now()
			requestVoteArgs := RequestVoteArgs{
				Term:         rf.currentTerm,
				CandidateId:  rf.me,
				LastLogIndex: len(rf.log) - 1,
				LastLogTerm:  rf.log[len(rf.log)-1].Term,
			}
			gotVotesChan := make(chan bool)
			for i := 0; i < len(rf.peers); i++ {
				if i == rf.me {
					continue
				}
				go func(i int) {
					reply := RequestVoteReply{}
					gotReply := rf.peers[i].Call("Raft.RequestVote", &requestVoteArgs, &reply)
					if gotReply {
						// TODO: if newer term in reply, convert to follower
						rf.mu.Lock()
						if reply.Term > rf.currentTerm {
							rf.state = Follower
							rf.votedFor = -1
							rf.currentTerm = reply.Term
						}
						rf.mu.Unlock()
						if reply.VoteGranted {
							gotVotesChan <- true
						}
					}
				}(i)
			}
			rf.mu.Unlock()

			// Votes requested. Now we wait.
			electionWon := false
			timeout := time.After(electionTimeout)
			totalVotes := 1 // Votes for itself
		loop:
			for {
				select {
				case <-gotVotesChan:
					totalVotes++
					log.Println(rf.me, " we got a vote hurray! we now have ", totalVotes)
					if totalVotes == len(rf.peers)/2+1 {
						electionWon = true
						break loop
					}

				case <-timeout:
					log.Println(rf.me, "Election timed out")
					electionWon = false
					break loop
				}
			}

			rf.mu.Lock()
			if rf.state == Candidate && electionWon {
				log.Print(rf.me, ": I won the election, hurray!")
				rf.state = Leader
				// If won election, initialize leader state
				for i, _ := range rf.matchIndex {
					rf.nextIndex[i] = len(rf.log)
					rf.matchIndex[i] = 0
				}
			}
			log.Println(rf.me, "holdElections#2 returning lock")
			rf.mu.Unlock()
		}
	}
}

func (rf *Raft) sendHeartbeats() {
	for {
		// Minimum time between heartbeats
		time.Sleep(50 * time.Millisecond)

		// Remember not to do anything blocking within the critical section
		//log.Print(rf.me, "sendHeartbeats requesting lock")
		rf.mu.Lock()
		//log.Print(rf.me, "sendHeartbeats got lock")

		// Only leaders send heartbeats
		if rf.state != Leader {
			//log.Print(rf.me, " sendHeartbeats: was not leader, so not sending anything")
			rf.mu.Unlock()
			continue
		}

		log.Println(rf.me, " thinks it should send heartbeats")

		// Attempt sending a heartbeat to each peer.
		for i, _ := range rf.peers {
			// Don't send heartbeat to self
			if i == rf.me {
				continue
			}
			// Non-blocking RPC call by using a goroutine
			go func(i int) {
				//log.Print(rf.me, "heartbeat goroutine requesting lock")
				rf.mu.Lock()
				//log.Print(rf.me, "heartbeat goroutine got lock")
				heartbeat := AppendEntriesArgs{
					Term:         rf.currentTerm,
					LeaderId:     rf.me,
					PrevLogIndex: rf.nextIndex[i] - 1,
					PrevLogTerm:  rf.log[rf.nextIndex[i]-1].Term,
					Entries:      []LogEntry{},
					LeaderCommit: rf.commitIndex,
				}
				//log.Print(rf.me, "heartbeat goroutine returning lock")
				rf.mu.Unlock()
				reply := AppendEntriesReply{Term: 0, Success: false}
				gotReply := rf.peers[i].Call("Raft.AppendEntries", &heartbeat, &reply)
				if gotReply {
					//log.Print(rf.me, "heartbeat response requesting lock")
					rf.mu.Lock()
					//log.Print(rf.me, "heartbeat response got lock")
					if reply.Term > rf.currentTerm {
						rf.currentTerm = reply.Term
						rf.state = Follower
						rf.votedFor = -1
					}
					if reply.Success {
						// TODO update nextIndex and matchIndex for follower (§5.3)
					} else {
						rf.nextIndex[i]--
						if rf.nextIndex[i] == 0 {
							rf.nextIndex[i] = 1
						}
					}
					log.Print(rf.me, "heartbeat response returning lock")
					rf.mu.Unlock()
				}
			}(i)
		}
		rf.mu.Unlock()
	}
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.state = Follower
	rf.votedFor = -1
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.log = []LogEntry{{0, ""}} // no-op first log entry
	rf.lastHeartbeat = time.Now()
	rf.matchIndex = make([]int, len(peers))
	rf.nextIndex = make([]int, len(peers))

	// Your initialization code here (2A, 2B, 2C).

	//Modify Make() to create a
	// background goroutine that will kick off leader election periodically by
	// sending out RequestVote RPCs when it hasn't heard from another peer for
	// a while. This way a peer will learn who is the leader, if there is
	// already a leader, or become the leader itself.
	go rf.holdElections()
	go rf.sendHeartbeats()

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	return rf
}
