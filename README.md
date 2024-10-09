# kvdb
A distributed key-value store built on top of Raft consensus algorithm. 


## To do list
- [x] Write unit test for all intenal modules
- [x] Implement Failover/Recovery, Consensus algo
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
