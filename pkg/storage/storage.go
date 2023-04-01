package storage

import (
	"fmt"
	"sync"
	"time"
)

type MemoryStore struct {
	dataMap  map[int]Data
	mx       sync.RWMutex
	duration int
}

// Data
type Data struct {
	ExpiratioTime time.Time
}

func NewMemoryStore(duration int) *MemoryStore {
	return &MemoryStore{
		dataMap:  make(map[int]Data, 0),
		duration: duration,
	}
}

// Add - add rand client data to mememory
func (c *MemoryStore) Add(key int) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.dataMap[key] = Data{
		ExpiratioTime: time.Now().UTC().Add(time.Second * time.Duration(c.duration)),
	}
	return nil
}

// Get - check data of client in memory
func (c *MemoryStore) Get(key int) error {
	c.mx.RLock()
	defer c.mx.RUnlock()
	value, ok := c.dataMap[key]

	if !ok {
		return fmt.Errorf("required values not found")
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
