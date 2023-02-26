package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func New() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Has(key []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.data[string(key)]
	return ok
}

func (c *Cache) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, string(key))
	return nil
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()

	defer c.lock.RUnlock()

	log.Printf("GET %s\n", string(key))

	keyStr := string(key)

	val, ok := c.data[keyStr]
	if !ok {
		return nil, fmt.Errorf("Key (%s) not found", keyStr)
	}
	log.Printf("GET %s = %s\n", string(key), string(val))
	return val, nil
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Printf("SET %s to %s\n", string(key), string(value))

	// goroutine para verificar se o tempo do dado salvo no cache foi expirado
	go func() {
		<-time.After(ttl) // Boolean que retorna se o tempo do dado ja expirou ou não
		delete(c.data, string(key))
	}()
	c.data[string(key)] = value
	return nil
}
