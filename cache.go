package batch

import (
	"context"
	"sync"
)

type Cache interface {
	Get(ctx context.Context, k Key) (res interface{}, ok bool)
	Set(ctx context.Context, k Key, res interface{})
}

type defaultCache struct {
	sync.RWMutex
	data map[Key]interface{}
}

func DefaultCache() Cache {
	dc := &defaultCache{
		data: make(map[Key]interface{}),
	}
	return dc
}

func (c *defaultCache) Get(_ context.Context, k Key) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	res, ok := c.data[k]
	return res, ok
}

func (c *defaultCache) Set(_ context.Context, k Key, res interface{}) {
	c.Lock()
	defer c.Unlock()
	c.data[k] = res
}
