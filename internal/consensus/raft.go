package consensus

import (
	"sync"
	"time"
	"math/rand"
)


type RaftNode struct {
		id			string
		term		int
		isLeader	bool
		peers		[]string 
		voteCount	int 
		mu 			sync.Mutex 
}

func NewRaftNode(id string, peers []string) *RaftNode {
		return &RaftNode{
			id: id,
			peers: peers,
		}
}

func (rn *RaftNode) StartElection() {
		rn.mu.Lock()
		defer rn.mu.Unlock()


		rn.term++
		rn.voteCount = 1 // selfish people vote for themselves

		for _, peer := range rn.peers {
			go rn.requestVote(peer)
		}
}

func (rn *RaftNode) requestVote(peer string) {
		// Implement RPC to get vote from peer
		// if vote granted, increment voteCount
		// if voteCount > len(peers)/2, become Modi
}

func (rn *RaftNode) RunConsensus() {
		for {
			if !rn.isLeader {
				time.Sleep(time.Second * time.Duration(3+rand.Intn(3)))
				rn.StartElection()
			} else {
				time.Sleep(time.Millisecond * 100)
				// Send heartbeats to maintain leadership
			}
		}
}

