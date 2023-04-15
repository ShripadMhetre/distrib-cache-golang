package cache

import (
	"fmt"
	"sync"
	"time"
)

type CacheImpl struct {
	lock      sync.RWMutex
	dataStore map[string][]byte
}

func New() *CacheImpl {
	return &CacheImpl{
		dataStore: make(map[string][]byte),
	}
}

func (c *CacheImpl) Get(key []byte) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	keyStr := string(key)
	val, ok := c.dataStore[keyStr]
	if !ok {
		return nil, fmt.Errorf("key: %s not found", keyStr)
	}

	return val, nil
}

func (c *CacheImpl) Set(key, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.dataStore[string(key)] = value

	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			delete(c.dataStore, string(key))
		}()
	}

	return nil
}

func (c *CacheImpl) Exists(key []byte) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.dataStore[string(key)]

	return ok, nil
}

func (c *CacheImpl) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.dataStore, string(key))

	return nil
}
