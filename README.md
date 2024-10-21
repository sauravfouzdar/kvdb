# kvdb
A naive attempt to implement a distributed key-value database using Raft. 

## Start the project
```
go run cmd/kvdbserver/main.go -id node1 -address :8000 -peers :8002,:8003

```

## To do list
- [x] Write unit test for all internal modules
- [x] Write data to disk(lowpriority)
- [x] Concurrent read/write, test with multiple clients etc..
- [x] Update readme
- [X] API client
- [X] code cleanup
- [ ] check for compatibility with 32bit systems(data types etc.)


## Project Structure
```
   kvdb/
   ├── cmd/
   │   └── kvdbserver/
   │       └── main.go
   ├── internal/
   │   ├── storage/
   │   │   └── storage.go
   │   ├── replication/
   │   │   └── replication.go
   │   ├── partitioning/
   │   │   └── partitioning.go
   │   ├── consensus/
   │   │   └── raft.go
   │   ├── loadbalancer/
   │   │   └── loadbalancer.go
   │   └── node/
   │       └── node.go
   ├── pkg/
   │   └── client/
   │       └── client.go
   └── go.mod
```


## References
- [Raft Consensus Algorithm](https://raft.github.io/raft.pdf)
- [Raft Visualization](http://thesecretlivesofdata.com/raft/)