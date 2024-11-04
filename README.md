# kvdb
A Naive attempt to implement Consensus using Raft on a in-memory key-value pair distributed database.

## Start the project
```
go run cmd/kvdbserver/main.go -id node1 -address :8000 -peers :8002,:8003


```

## Run tests
```
go test ./...

```

## To do list
- [x] Write unit test for all internal modules
- [x] Write data to disk(lowpriority)
- [x] Concurrent read/write, test with multiple clients etc..
- [x] Update readme
- [ ] API client
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
- [Raft Consensus Algorithm](https://raft.github.io/raft.pdf) for quick understanding go to `Page no. 4`
- [Raft Visualization](http://thesecretlivesofdata.com/raft/) 
