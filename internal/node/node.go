/* handles network communication. */
package node

import (
	"encoding/gob"
	"kvdb/internal/consensus"
	"kvdb/internal/partitioning"
	"kvdb/internal/replication"
	"kvdb/internal/storage"
	"log"
	"net"
)


type Node struct {
	id				string
	storage			*storage.Storage
	replication		*replication.Leader
	partitioner		*partitioning.Partitioner
	raft			*consensus.RaftNode
	listener		net.Listener
}

func NewNode(id string, address string, peers []string) *Node {
		return &Node{
				id:				id,
				storage: 		storage.NewStorage(),
				replication:	replication.NewLeader(),
				partitioner:	partitioning.NewPartioner(len(peers) + 1),
				raft:			consensus.NewRaftNode(id, address, peers),
		}
}

func (n *Node) Start(address string) error {
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return err 
		}
		n.listener = listener


		go n.raft.RunConsensus()
		go n.acceptConnections()

		
		return nil
}

func (n *Node) acceptConnections() {
		for {
			conn, err := n.listener.Accept()
			if err != nil {
				log.Println("Following error occured while accepting conn: ",err)
				continue
			}
			go n.handleConnection(conn)
		}
}

func (n *Node) handleConnection(conn net.Conn) {
		defer conn.Close()
		decoder := gob.NewDecoder(conn)
		encoder := gob.NewEncoder(conn)


		var request struct {
			Operation	string
			Key			string 
			Value		string
		}

		if err := decoder.Decode(&request); err != nil {
			log.Println("An error occured: ", err)
			return
		}

		switch request.Operation {
		case "GET":
			value, ok := n.storage.Get(request.Key)
			encoder.Encode(map[string]interface{}{"ok":ok, "value":value})
		case "SET":
			log.Println("Setting key--> ", request.Key, " value: ", request.Value)
			n.storage.Set(request.Key, request.Value)
			n.replication.ReplicateToFollowers(request.Key, request.Value)
			encoder.Encode(map[string]interface{}{"ok":true})
		case "DELETE":
			n.storage.Delete(request.Key)
			n.replication.ReplicateToFollowers(request.Key,"")
			encoder.Encode(map[string]interface{}{"ok":true})
		}
}

