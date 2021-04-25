package batch

import (
	"context"
)

type HookFunc func(res interface{})

// if you dont want return result, return nil, err
type BatchFunc func(context.Context, Keys) (map[Key]interface{}, error)

type Batch struct {
	batchFn   BatchFunc
	batchSpan int
	reqs      chan *batchRequest
	cache     Cache
	errChan   chan error
	ctx       context.Context
}

type batchRequest struct {
	key   Key
	hooks []HookFunc
}

func NewBatch(ctx context.Context, batchFn BatchFunc, opts ...Option) *Batch {
	batch := &Batch{
		batchFn: batchFn,
		ctx:     ctx,
		errChan: make(chan error, 1),
	}
	// Apply options
	for _, apply := range opts {
		apply(batch)
	}
	batch.reqs = make(chan *batchRequest, batch.batchSpan*2)
	batch.start()
	return batch
}

func (b *Batch) Flush() error {
	//println(`flush`)
	b.reqs <- nil
	return <-b.errChan
}

func (b *Batch) done(err error) {
	b.errChan <- err
}

func (b *Batch) start() {
	go func() {
		keys := make(Keys, 0, b.batchSpan)
		hooksMap := make(map[string][]HookFunc, b.batchSpan)
		for {
			select {
			case <-b.ctx.Done():
				b.done(b.ctx.Err())
				return
			case req := <-b.reqs:
				if req != nil {
					if b.cache != nil {
						res, ok := b.cache.Get(b.ctx, req.key)
						if ok {
							for _, hook := range req.hooks {
								hook(res)
							}
							continue
						}
					}
					keys = append(keys, req.key)
					if req.hooks != nil {
						hooksMap[req.key.String()] = append(hooksMap[req.key.String()], req.hooks...)
					}
				}
				if (len(keys) == b.batchSpan || req == nil) && len(keys) > 0 {
					results, err := b.batchFn(b.ctx, keys)
					if err != nil {
						b.done(err)
						return
					}
					for k, v := range results {
						if b.cache != nil {
							b.cache.Set(b.ctx, k, v)
						}
						hooks := hooksMap[k.String()]
						if len(hooks) > 0 {
							for _, hook := range hooks {
								hook(v)
							}
						}
					}
					for _, k := range keys {
						delete(hooksMap, k.String())
					}
					keys = keys[:0]
				}
				if req == nil {
					b.done(nil)
					return
				}
			}
		}
	}()
}

// pass key, nil if you dont want hooks
// hook funcs will be executed in another goroutine
func (b *Batch) Do(key Key, hooks ...HookFunc) {
	b.reqs <- &batchRequest{
		key:   key,
		hooks: hooks,
	}
}
