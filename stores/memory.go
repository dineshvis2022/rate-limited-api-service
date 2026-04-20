package stores

import (
	"maps"
	"sync"
)

type MemoryStore struct {
	mu    sync.RWMutex // protects concurrent access
	stats map[string]int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{stats: make(map[string]int)}
}

// Increment increases request count for a user (thread-safe)
func (m *MemoryStore) Increment(userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stats[userID]++
	return nil
}

// GetStats returns a copy to avoid race conditions
func (m *MemoryStore) GetStats() (map[string]int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make(map[string]int)
	maps.Copy(res, m.stats)

	// for k, v := range m.stats {
	// 	// res[k] = v
	// }
	return res, nil
}
