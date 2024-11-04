package consensus

import (
	"testing"
	"time"
)

func TestNewRaftNode(t *testing.T) {
	peerAddrs := []string{":8002", ":8003"}
	rn := NewRaftNode("node0", ":8000", peerAddrs)

	if rn.id != "node0" {
		t.Errorf("Expected node id to be 'node0', got %s", rn.id)
	}
	if rn.state != "follower" {
		t.Errorf("Expected initial state to be 'follower', got %s", rn.state)
	}
	if len(rn.peers) != 3 {
		t.Errorf("Expected 3 peers, got %d", len(rn.peers))
	}
}

func TestAppendEntries(t *testing.T) {
	rn := NewRaftNode("node0", ":8000", []string{})
	rn.currentTerm = 1
	rn.log = append(rn.log, LogEntry{Term: 1, Command: "set x 1"})

	args := &AppendEntriesArgs{
		Term:         1,
		LeaderId:     "node1",
		PrevLogIndex: 0,
		PrevLogTerm:  1,
		Entries:      []LogEntry{{Term: 1, Command: "set y 2"}},
		LeaderCommit: 1,
	}
	reply := &AppendEntriesReply{}

	err := rn.AppendEntries(args, reply)
	if err != nil {
		t.Fatalf("AppendEntries returned error: %v", err)
	}
	if !reply.Success {
		t.Errorf("Expected success, got failure")
	}
	if len(rn.log) != 2 {
		t.Errorf("Expected log length to be 2, got %d", len(rn.log))
	}
	if rn.commitIndex != 1 {
		t.Errorf("Expected commitIndex to be 1, got %d", rn.commitIndex)
	}
}

func TestElectionTimeoutReset(t *testing.T) {
	rn := NewRaftNode("node0", ":8000", []string{})
	initialTimeout := rn.electionTimeout
	rn.resetElectionTimeout()
	if rn.electionTimeout == initialTimeout {
		t.Errorf("Expected election timeout to be reset, but it was not")
	}
}

func TestBecomeFollower(t *testing.T) {
	rn := NewRaftNode("node0", ":8000", []string{})
	rn.state = "leader"
	rn.becomeFollower(2)
	if rn.state != "follower" {
		t.Errorf("Expected state to be 'follower', got %s", rn.state)
	}
	if rn.currentTerm != 2 {
		t.Errorf("Expected currentTerm to be 2, got %d", rn.currentTerm)
	}
	if rn.votedFor != "" {
		t.Errorf("Expected votedFor to be empty, got %s", rn.votedFor)
	}
}

func TestBecomeLeader(t *testing.T) {
	rn := NewRaftNode("node0", ":8000", []string{":8001", ":8002"})
	rn.state = "candidate"
	rn.becomeLeader()
	if rn.state != "leader" {
		t.Errorf("Expected state to be 'leader', got %s", rn.state)
	}
	for peerId := range rn.peers {
		if rn.nextIndex[peerId] != 0 {
			t.Errorf("Expected nextIndex for peer %s to be 0, got %d", peerId, rn.nextIndex[peerId])
		}
		if rn.matchIndex[peerId] != 0 {
			t.Errorf("Expected matchIndex for peer %s to be 0, got %d", peerId, rn.matchIndex[peerId])
		}
	}
}

func TestIsLogUpToDate(t *testing.T) {
	rn := NewRaftNode("node0", ":8000", []string{})
	rn.log = append(rn.log, LogEntry{Term: 1, Command: "set x 1"})

	if !rn.isLogUpToDate(0, 1) {
		t.Errorf("Expected log to be up to date")
	}
	if rn.isLogUpToDate(1, 1) {
		t.Errorf("Expected log to not be up to date")
	}
}

func TestRunConsensus(t *testing.T) {
	rn := NewRaftNode("node0", ":8000", []string{})
	go rn.RunConsensus()
	time.Sleep(1 * time.Second)
	if rn.state != "candidate" && rn.state != "leader" {
		t.Errorf("Expected state to be 'candidate' or 'leader', got %s", rn.state)
	}
}