package replication

import "sync"

type Leader struct {
		followers map[string]Follower
		mu        sync.RWMutex
}

type Follower interface {
		Replicate(key, value string) error 
}

func NewLeader() *Leader {
		return &Leader{
			followers: make(map[string]Follower),
		}
}

func (ldr *Leader) AddFollower(id string, follower Follower) {
		ldr.mu.Lock()
		defer ldr.mu.Unlock()
		ldr.followers[id] = follower
}

func (ldr *Leader) RemoveFollower(id string) {
		ldr.mu.Lock()
		defer ldr.mu.Unlock()
		delete(ldr.followers, id)
}

func (ldr *Leader) ReplicateToFollowers(key, value string) []error {
		ldr.mu.RLock()
		defer ldr.mu.RUnlock()

		var errors []error 
		for _, follower := range ldr.followers {
			if err := follower.Replicate(key, value); err != nil {
				errors = append(errors, err)
			}
		}
		return errors
}






