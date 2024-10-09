package partitioning

import "hash/fnv"

type Partitioner struct {
	nodeCount int
}

func NewPartioner(nodeCount int) *Partitioner {
	return &Partitioner{nodeCount: nodeCount}
}

func (p *Partitioner) GetNodeForKey(key string) int {
		/* (Fowler-Noll-Vo) is a fast, non-cryptographic hash function */
		hash := fnv.New32a() // 32 byte hash in big-endian byte order
		hash.Write([]byte(key))
		return int(hash.Sum32()) % p.nodeCount
}


