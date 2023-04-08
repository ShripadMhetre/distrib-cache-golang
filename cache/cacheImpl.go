package cache

import (
	"fmt"
	"time"
)

type CacheImpl struct {
	dataStore map[string][]byte
}

func New() *CacheImpl {
	return &CacheImpl{
		dataStore: make(map[string][]byte),
	}
}

func (c *CacheImpl) Delete(key []byte) error {
	delete(c.dataStore, string(key))

	return nil
}

func (c *CacheImpl) Has(key []byte) bool {
	_, ok := c.dataStore[string(key)]

	return ok
}

func (c *CacheImpl) Get(key []byte) ([]byte, error) {
	keyStr := string(key)
	val, ok := c.dataStore[keyStr]
	if !ok {
		return nil, fmt.Errorf("key: %s not found", keyStr)
	}

	return val, nil
}

func (c *CacheImpl) Set(key, value []byte, ttl time.Duration) error {
	c.dataStore[string(key)] = value

	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			delete(c.dataStore, string(key))
		}()
	}

	return nil
}
