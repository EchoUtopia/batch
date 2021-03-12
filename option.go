package batch

type Option func(*Batch)

// WithBatchSpan sets the batch capacity. Default is 0 (unbounded).
func WithBatchSpan(c int) Option {
	return func(b *Batch) {
		b.batchSpan = c
	}
}

func WithCache(c Cache) Option {
	return func(b *Batch) {
		b.cache = c
	}
}
