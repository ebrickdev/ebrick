package cache

import (
	"context"
	"crypto"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/ebrickdev/ebrick/cache/store"
)

var DefaultCache Cache

// Cache is an interface for a cache with methods to get, set, delete, invalidate, and clear cache entries.
type Cache interface {
	Get(ctx context.Context, key any) (any, error)
	GetWithTTL(ctx context.Context, key any) (any, time.Duration, error)
	Set(ctx context.Context, key any, object any, options ...store.Option) error
	Delete(ctx context.Context, key any) error
	Invalidate(ctx context.Context, options ...store.InvalidateOption) error
	Clear(ctx context.Context) error
	GetType() string
	GetStore() store.Store
	GetStats() *Stats
	getCacheKey(key any) string
}

// CacheKeyGenerator is an interface for generating cache keys.
type CacheKeyGenerator interface {
	GetCacheKey() string
}

// Stats holds statistics about cache operations.
type Stats struct {
	Hits              int
	Miss              int
	SetSuccess        int
	SetError          int
	DeleteSuccess     int
	DeleteError       int
	InvalidateSuccess int
	InvalidateError   int
	ClearSuccess      int
	ClearError        int
}

// cache is an implementation of the Cache interface.
type cache struct {
	store    store.Store
	stats    *Stats
	statsMtx sync.Mutex
}

// New creates a new cache instance with the given store.
func New(store store.Store) Cache {
	return &cache{
		store: store,
		stats: &Stats{},
	}
}

// Get retrieves an item from the cache by key.
func (c *cache) Get(ctx context.Context, key any) (any, error) {
	cacheKey := c.getCacheKey(key)
	value, err := c.store.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// GetStats returns the current cache statistics.
func (c *cache) GetStats() *Stats {
	c.statsMtx.Lock()
	defer c.statsMtx.Unlock()
	stats := *c.stats
	return &stats
}

// GetStore returns the underlying store of the cache.
func (c *cache) GetStore() store.Store {
	return c.store
}

// GetType returns the type of the underlying store.
func (c *cache) GetType() string {
	return c.store.GetType()
}

// GetWithTTL retrieves an item from the cache by key along with its TTL.
func (c *cache) GetWithTTL(ctx context.Context, key any) (any, time.Duration, error) {
	val, ttl, err := c.store.GetWithTTL(ctx, key)

	c.statsMtx.Lock()
	defer c.statsMtx.Unlock()
	if err == nil {
		c.stats.Hits++
	} else {
		c.stats.Miss++
	}

	return val, ttl, err
}

// Invalidate invalidates cache entries based on the given options.
func (c *cache) Invalidate(ctx context.Context, options ...store.InvalidateOption) error {
	err := c.store.Invalidate(ctx, options...)
	c.statsMtx.Lock()
	defer c.statsMtx.Unlock()
	if err == nil {
		c.stats.InvalidateSuccess++
	} else {
		c.stats.InvalidateError++
	}
	return err
}

// Set adds an item to the cache with the given key and options.
func (c *cache) Set(ctx context.Context, key any, object any, options ...store.Option) error {
	err := c.store.Set(ctx, key, object, options...)

	c.statsMtx.Lock()
	defer c.statsMtx.Unlock()
	if err == nil {
		c.stats.SetSuccess++
	} else {
		c.stats.SetError++
	}
	return err
}

// Clear clears all items from the cache.
func (c *cache) Clear(ctx context.Context) error {
	err := c.store.Clear(ctx)

	c.statsMtx.Lock()
	defer c.statsMtx.Unlock()
	if err == nil {
		c.stats.ClearSuccess++
	} else {
		c.stats.ClearError++
	}

	return err
}

// Delete removes an item from the cache by key.
func (c *cache) Delete(ctx context.Context, key any) error {
	err := c.store.Delete(ctx, key)

	c.statsMtx.Lock()
	defer c.statsMtx.Unlock()
	if err == nil {
		c.stats.DeleteSuccess++
	} else {
		c.stats.DeleteError++
	}

	return err
}

// getCacheKey returns the cache key for the given key object by returning
// the key if type is string or by computing a checksum of key structure
// if its type is other than string.
func (c *cache) getCacheKey(key any) string {
	switch v := key.(type) {
	case string:
		return v
	case CacheKeyGenerator:
		return v.GetCacheKey()
	default:
		return checksum(key)
	}
}

// checksum hashes a given object into a string.
func checksum(object any) string {
	digester := crypto.MD5.New()
	fmt.Fprint(digester, reflect.TypeOf(object))
	fmt.Fprint(digester, object)
	hash := digester.Sum(nil)

	return fmt.Sprintf("%x", hash)
}
