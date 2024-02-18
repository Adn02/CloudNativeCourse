package cache

import "errors"

type Cacher[K comparable, V any] interface {
	Get(key K) (value V, err error)
	Put(key K, value V) (err error)
}

// Concrete LRU cache
type lruCache[K comparable, V any] struct {
	size      int
	remaining int
	cache     map[K]V
	queue     []K
}

// Constructor
func NewCacher[K comparable, V any](size int) Cacher[K, V] {
	return &lruCache[K, V]{size: size, remaining: size, cache: make(map[K]V), queue: make([]K, 0)}
}

func (c *lruCache[K, V]) Get(key K) (value V, err error) {
	value, exist := c.cache[key] // exist: true or false based on if key exists in map

	if exist {
		c.deleteFromQueue(key)
		c.queue = append(c.queue, key) // Adds key to the tail of queue
	} else {
		err = errors.New("key does not exist in cache")
	}

	return value, err

}

func (c *lruCache[K, V]) Put(key K, value V) (err error) {
	_, exists := c.cache[key]

	// Case 1: Key already exists
	if exists {
		c.deleteFromQueue(key)
	}

	// Case 2: Cache is full capacity
	if !exists && c.remaining == 0 {
		victim := c.queue[0]      // Get least recently used key from queue (head)
		delete(c.cache, victim)   // Remove key from cache
		c.deleteFromQueue(victim) // Remove key from queue
		c.remaining++
	}
	// Default: Add key-value pair to map and end of queue
	c.cache[key] = value
	c.queue = append(c.queue, key)
	c.remaining--

	return nil
}

// Helper method to delete all occurrences of a key from the queue
func (c *lruCache[K, V]) deleteFromQueue(key K) {
	newQueue := make([]K, 0, c.size)
	for _, k := range c.queue {
		if k != key {
			newQueue = append(newQueue, k)
		}
	}
	c.queue = newQueue
}
