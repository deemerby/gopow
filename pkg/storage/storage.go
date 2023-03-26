package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type MemoryStore struct {
	dataMap map[int]Data
	mx      sync.RWMutex
}

// Data
type Data struct {
	CreatedTime time.Time
	Expiration  int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		dataMap: make(map[int]Data, 0),
	}
}

// Add - add rand client data to mememory
func (c *MemoryStore) Add(key int) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.dataMap[key] = Data{
		CreatedTime: time.Now().UTC(),
		Expiration:  viper.GetInt("hashcash.duration"),
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
	createdTime := value.CreatedTime

	if curTime.After(createdTime.Add(time.Second * time.Duration(value.Expiration))) {
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
