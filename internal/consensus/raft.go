package consensus

import (
	"context"
	"fmt"
	"math/rand"
	"net/rpc"
	"strings"
	"sync"
	"time"
)

type RaftNode struct {
	mu              sync.Mutex
	id              string
	currentTerm     int
	votedFor        string
	log             []LogEntry
	commitIndex     int
	lastApplied     int
	nextIndex       map[string]int
	matchIndex      map[string]int
	state           string            // follower, leader, candidate
	peers           map[string]string // node_ids <> addresses
	electionTimeout time.Duration
	heartbeatTicker *time.Ticker
	//voteCount       int
}

type LogEntry struct {
	Term    int
	Command string
}

type AppendEntriesArgs struct {
	Term         int
	LeaderId     string
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}

func NewRaftNode(id string, address string, peerAddrs []string) *RaftNode {

	peers := make(map[string]string)

	// Convert peer list to map
	for i, addr := range peerAddrs {
		if addr != "" && addr != address { // Skip empty addresses and self
			peerId := fmt.Sprintf("node%d", i+1)
			if strings.HasPrefix(addr, ":") {
				// If only port is provided, add localhost
				addr = "localhost" + addr
			}
			peers[peerId] = addr
		}
	}
	rn := &RaftNode{
		id:              id,
		currentTerm:     0,
		votedFor:        "",
		log:             make([]LogEntry, 0),
		commitIndex:     -1,
		lastApplied:     -1,
		nextIndex:       make(map[string]int),
		matchIndex:      make(map[string]int),
		state:           "follower",
		peers:           peers,
		electionTimeout: time.Duration(150+rand.Intn(150)) * time.Millisecond,
	}

	rn.log = append(rn.log, LogEntry{Term: 0, Command: "nil"})
	return rn
}

func (rn *RaftNode) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	reply.Success = false
	reply.Term = rn.currentTerm

	if args.Term < rn.currentTerm {
		return nil
	}

	rn.resetElectionTimeout()

	if args.Term > rn.currentTerm {
		rn.currentTerm = args.Term
		rn.state = "follower"
		rn.votedFor = ""
	}

	//check if log contains an entry at prevLogIndex with prevLogTerm
	if args.PrevLogIndex > -1 {
		if len(rn.log) <= args.PrevLogIndex || rn.log[args.PrevLogIndex].Term != args.PrevLogTerm {
			return nil
		}
	}

	// Append new entries
	for i, entry := range args.Entries {
		index := args.PrevLogIndex + 1 + i
		if index < len(rn.log) {
			if rn.log[index].Term != entry.Term {
				rn.log = rn.log[:index]
				rn.log = append(rn.log, entry)
			}
		} else {
			rn.log = append(rn.log, entry)
		}
	}

	//update commit index
	if args.LeaderCommit > rn.commitIndex {
		rn.commitIndex = min(args.LeaderCommit, len(rn.log)-1)
	}

	reply.Success = true
	return nil
}

// send beat to all peers to maintain Leader state
func (rn *RaftNode) sendHeartbeat() {
	rn.mu.Lock()
	if rn.state != "leader" {
		rn.mu.Unlock()
		return
	}
	term := rn.currentTerm
	rn.mu.Unlock()

	for peerId, peerAddr := range rn.peers {
		go func(id, addr string) {
			rn.mu.Lock()
			prevLogIndex := rn.nextIndex[id] - 1
			prevLogTerm := -1
			if prevLogIndex >= 0 && prevLogIndex < len(rn.log) {
				prevLogTerm = rn.log[prevLogIndex].Term
			}
			entries := rn.log[rn.nextIndex[id]:]
			args := &AppendEntriesArgs{
				Term:         term,
				LeaderId:     rn.id,
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  prevLogTerm,
				Entries:      entries,
				LeaderCommit: rn.commitIndex,
			}
			rn.mu.Unlock()

			var reply AppendEntriesReply
			client, err := rpc.DialHTTP("tcp", addr)
			if err != nil {
				fmt.Printf("Error connecting to peer^^ %s: %v\n", id, err)
				return
			}
			defer client.Close()

			err = client.Call("RaftNode.AppendEntries", args, &reply)
			if err != nil {
				fmt.Printf("Error calling AppendEntries on peer %s: %v\n", id, err)
				return
			}

			rn.mu.Lock()
			defer rn.mu.Unlock()

			if reply.Term > rn.currentTerm {
				rn.currentTerm = reply.Term
				rn.state = "follower"
				rn.votedFor = ""
				return
			}

			if reply.Success {
				rn.nextIndex[id] = prevLogIndex + len(entries) + 1
				rn.matchIndex[id] = rn.nextIndex[id] - 1
			} else {
				rn.nextIndex[id]--
			}
		}(peerId, peerAddr)
	}
}

func (rn *RaftNode) startHeartBeat() {
	rn.heartbeatTicker = time.NewTicker(50 * time.Millisecond)
	go func() {
		for range rn.heartbeatTicker.C {
			rn.sendHeartbeat()
		}
	}()
}

func (rn *RaftNode) stopHeartbeat() {
	if rn.heartbeatTicker != nil {
		rn.heartbeatTicker.Stop()
	}
}

func (rn *RaftNode) resetElectionTimeout() {
	rn.electionTimeout = time.Duration(150+rand.Intn(150)) * time.Millisecond
}

func (rn *RaftNode) startElection() {
	rn.mu.Lock()
	rn.currentTerm++
	rn.state = "candidate"
	rn.votedFor = rn.id
	currentTerm := rn.currentTerm
	PrevLogIndex := len(rn.log) - 1
	PrevLogTerm := rn.log[PrevLogIndex].Term
	rn.resetElectionTimeout()
	rn.mu.Unlock()

	votes := 1
	voteChan := make(chan bool, len(rn.peers))

	for peerId, peerAddr := range rn.peers {
		go func(id, addr string) {
			args := &AppendEntriesArgs{
				Term:         currentTerm,
				LeaderId:     rn.id,
				PrevLogIndex: PrevLogIndex,
				PrevLogTerm:  PrevLogTerm,
			}
			var reply AppendEntriesReply

			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				client, err := rpc.DialHTTP("tcp", addr)
				if err != nil {
					done <- err
					return
				}
				defer client.Close()
				done <- client.Call("RaftNode.requestVote", args, &reply)
			}()

			select {
			case <-ctx.Done():
				voteChan <- false
			case err := <-done:
				if err != nil {
					voteChan <- false
				} else {
					rn.mu.Lock()
					if reply.Term > rn.currentTerm {
						rn.becomeFollower(reply.Term)
						rn.mu.Unlock()
						voteChan <- false
					} else {
						rn.mu.Unlock()
						voteChan <- reply.Success
					}
				}
			}
		}(peerId, peerAddr)
	}

	go func() {
		for i := 0; i < len(rn.peers); i++ {
			if <-voteChan {
				votes++
				if votes > len(rn.peers)/2 {
					rn.mu.Lock()
					if rn.state == "candidate" {
						rn.becomeLeader()
					}
					rn.mu.Unlock()
				}
			}
		}
		// If we reach here, we didn't win election
		rn.mu.Lock()
		if rn.state == "candidate" {
			rn.becomeFollower(rn.currentTerm)
		}
		rn.mu.Unlock()
	}()
}

func (rn *RaftNode) requestVote(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	// Implement RPC to get vote from peer
	// if vote granted, increment voteCount
	// if voteCount > len(peers)/2, become Modi
	rn.mu.Lock()
	defer rn.mu.Unlock()

	reply.Term = rn.currentTerm
	reply.Success = false

	if args.Term < rn.currentTerm {
		return nil
	}

	if args.Term > rn.currentTerm {
		rn.becomeFollower(args.Term)
	}

	if (rn.votedFor == "" || rn.votedFor == args.LeaderId) && rn.isLogUpToDate(args.PrevLogIndex, args.PrevLogTerm) {
		reply.Success = true
		rn.votedFor = args.LeaderId
		rn.resetElectionTimeout()
	}

	return nil
}

func (rn *RaftNode) becomeFollower(term int) {
	rn.state = "follower"
	rn.currentTerm = term
	rn.votedFor = ""
	rn.resetElectionTimeout()
	if rn.heartbeatTicker != nil {
		rn.heartbeatTicker.Stop()
		rn.heartbeatTicker = nil
	}
}

func (rn *RaftNode) becomeLeader() {
	rn.state = "leader"
	for peerId := range rn.peers {
		rn.nextIndex[peerId] = len(rn.log)
		rn.matchIndex[peerId] = 0
	}
	rn.startHeartBeat()
}

func (rn *RaftNode) isLogUpToDate(candidatePrevLogIndex, candidatePrevLogTerm int) bool {
	PrevLogIndex := len(rn.log) - 1
	PrevLogTerm := rn.log[PrevLogIndex].Term

	if candidatePrevLogIndex > PrevLogTerm {
		return true
	}
	if candidatePrevLogTerm == PrevLogTerm && candidatePrevLogIndex >= PrevLogIndex {
		return true
	}
	return false
}

func (rn *RaftNode) RunConsensus() {
	for {
		if rn.state != "Leader" {
			time.Sleep(time.Second * time.Duration(3+rand.Intn(3)))
			rn.startElection()
		} else {
			time.Sleep(time.Millisecond * 100)
			rn.sendHeartbeat()
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
