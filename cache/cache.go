package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

// Cache Object
type Cache struct {
	client *memcache.Client
}

// Item cached
type Item struct {
	Key        string
	Value      []byte
	Expiration int32
}

// NewCache new cache Object
func NewCache(host ...string) *Cache {
	client := memcache.New(host...)
	cache := &Cache{client}
	return cache
}

// Set value to cache
func (j *Cache) Set(key string, value []byte) error {
	return j.client.Set(&memcache.Item{Key: key, Value: value})
}

// Add value to cache if not in cache
func (j *Cache) Add(key string, value []byte) error {
	return j.client.Add(&memcache.Item{Key: key, Value: value})
}

// Delete value from cache
func (j *Cache) Delete(key string) error {
	return j.client.Delete(key)
}

// DeleteAll value from cache
func (j *Cache) DeleteAll() error {
	return j.client.DeleteAll()
}

// FlushAll value from cache
func (j *Cache) FlushAll() error {
	return j.client.FlushAll()
}

// Get value from cache
func (j *Cache) Get(key string) (*Item, error) {
	it, err := j.client.Get(key)
	if err != nil {
		return nil, err
	}
	return &Item{it.Key, it.Value, it.Expiration}, nil
}

// Decrement value from cache
func (j *Cache) Decrement(key string, delta uint64) (newValue uint64, err error) {
	return j.client.Decrement(key, delta)
}

// Increment value from cache
func (j *Cache) Increment(key string, delta uint64) (newValue uint64, err error) {
	return j.client.Increment(key, delta)
}

// Touch value from cache
func (j *Cache) Touch(key string, seconds int32) (err error) {
	return j.client.Touch(key, seconds)
}

// CompareAndSwap value from cache
func (j *Cache) CompareAndSwap(item *Item) error {
	it, err := j.client.Get(item.Key)
	if err != nil {
		return err
	}
	it.Value = item.Value
	it.Expiration = item.Expiration
	return j.client.CompareAndSwap(it)
}
