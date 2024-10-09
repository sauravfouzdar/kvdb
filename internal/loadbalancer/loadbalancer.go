package loadbalancer

import (
		"sync"
		"sync/atomic"
)

type LoadBalancer struct {
	nodes []string
	index uint64
	mu    sync.RWMutex
}

func NewLoadBalancer(nodes []string) *LoadBalancer {
		return &LoadBalancer{
			nodes: nodes,
		}
}

func (lb *LoadBalancer) NextNode() string {
		lb.mu.RLock()
		defer lb.mu.RUnlock()

		index := atomic.AddUint64(&lb.index, 1)
		return lb.nodes[index%uint64(len(lb.nodes))]

}

func (lb *LoadBalancer) AddNode(node string) {
		lb.mu.Lock()
		defer lb.mu.Unlock()
		lb.nodes = append(lb.nodes, node)
}

func (lb *LoadBalancer) RemoveNode(node string) {
		lb.mu.Lock()
		defer lb.mu.Unlock()
		for i, n := range lb.nodes {
			if n == node {
				lb.nodes = append(lb.nodes[:i], lb.nodes[i+1:]...)
				break
			}
		}
}

