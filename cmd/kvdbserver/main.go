package main

import (
	"flag"
	"kvdb/internal/node"
	"log"
	"strings"
)

func main() {
	var (
		id			= flag.String("id", "", "Node ID")
		address		= flag.String("address", ":8000", "Node address")
		peers		= flag.String("peers","", "list of peer addresses")
	)
	flag.Parse()

	if *id == "" {
		log.Fatal("Node ID is required")
	}


	peerList := strings.Split(*peers, ",")

	if *peers == "" {
		peerList = []string{}
	}

	//Convert String to map

	node := node.NewNode(*id, *address, peerList)
	if err := node.Start(*address); err != nil {
		log.Fatal(err)
	}


	log.Printf("Node %s started on %s", *id, *address)
	select {} 
}
