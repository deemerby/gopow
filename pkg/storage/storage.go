package storage

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type MemoryStore struct {
	dataMap  map[int]Data
	mx       sync.RWMutex
	duration time.Duration
}

// Data
type Data struct {
	ExpiratioTime time.Time
}

func NewMemoryStore(ctx context.Context, duration time.Duration) *MemoryStore {
	ms := &MemoryStore{
		dataMap:  make(map[int]Data, 0),
		duration: duration,
	}
	go ms.Watcher(ctx)

	return ms
}

// Add - add rand client data to mememory
func (c *MemoryStore) Add(key int) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.dataMap[key] = Data{
		ExpiratioTime: time.Now().UTC().Add(c.duration),
	}
	return nil
}

// Get - check data of client in memory
func (c *MemoryStore) Get(key int) error {
	c.mx.RLock()
	defer c.mx.RUnlock()
	value, ok := c.dataMap[key]

	if !ok {
		return fmt.Errorf("required value was not found or was automatically deleted after expiration")
	}

	curTime := time.Now().UTC()

	if curTime.After(value.ExpiratioTime) {
		return fmt.Errorf("value was expired")
	}

	return nil
}

// Delete - delete key from memory
func (c *MemoryStore) Delete(key int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.dataMap, key)
}

// Add - add rand client data to mememory
func (c *MemoryStore) Watcher(ctx context.Context) {
	for {
		c.mx.Lock()
		curTime := time.Now().UTC()
		for k, v := range c.dataMap {
			if curTime.After(v.ExpiratioTime) {
				delete(c.dataMap, k)
			}
		}
		c.mx.Unlock()

		select {
		case <-ctx.Done():
			return
		case <-time.After(c.duration * 2):
		}
	}
}
