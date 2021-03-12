package batch

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
	"time"
)

var bFn = func(ctx context.Context, keys Keys) (map[Key]interface{}, error) {
	out := map[Key]interface{}{}
	for _, v := range keys {
		out[v] = v
	}
	return out, nil
}

func testBatch(t *testing.T, b *Batch) {
	cnt := 0
	for i := 0; i < b.batchSpan; i++ {
		ki := Int64Key(i)
		b.Do(ki, func(r interface{}) {
			require.Equal(t, true, r.(Int64Key) == ki)
			cnt++
		})
	}

	b.Do(Int64Key(53), func(r interface{}) {
		require.Equal(t, true, r.(Int64Key) == Int64Key(53))
		cnt++
	})
	b.Do(Int64Key(40), func(r interface{}) {
		require.Equal(t, true, r.(Int64Key) == Int64Key(40))
		cnt++
	})
	err := b.Flush()
	require.Equal(t, err, nil)
	require.Equal(t, cnt, b.batchSpan+2)
}

func TestBatch(t *testing.T) {
	b := NewBatch(context.TODO(), bFn, WithBatchSpan(50), WithCache(DefaultCache()))
	testBatch(t, b)
	b = NewBatch(context.TODO(), bFn, WithBatchSpan(50))
	testBatch(t, b)
}

func TestError(t *testing.T) {
	bFn := func(ctx context.Context, keys Keys) (map[Key]interface{}, error) {
		return nil, errors.New(`err`)
	}
	b := NewBatch(context.TODO(), bFn, WithBatchSpan(3))
	for i := 0; i < 5; i++ {
		ki := Int64Key(i)
		b.Do(ki)
	}
	err := b.Flush()
	require.NotNil(t, err)
}

func TestCtxDone(t *testing.T) {
	bFn := func(ctx context.Context, keys Keys) (map[Key]interface{}, error) {
		results := map[Key]interface{}{}
		for _, v := range keys {
			results[v] = v
		}
		return results, nil
	}
	cnt := 0
	hook := func(_ interface{}) {
		cnt++
	}
	ctx, cancel := context.WithCancel(context.Background())
	b := NewBatch(ctx, bFn, WithBatchSpan(3))
	for i := 0; i < 5; i++ {
		ki := Int64Key(i)
		b.Do(ki, hook)
	}
	time.Sleep(time.Millisecond)
	cancel()
	time.Sleep(time.Millisecond)
	err := b.Flush()
	require.NotNil(t, err, err)
	require.Equal(t, 3, cnt)
}

func ExampleBatch() {
	cnt := 0
	bFn := func(ctx context.Context, keys Keys) (map[Key]interface{}, error) {
		cnt += len(keys)
		return nil, nil
	}
	b := NewBatch(context.TODO(), bFn, WithBatchSpan(10))
	for i := 0; i < 23; i++ {
		b.Do(Int64Key(i))
	}
	if err := b.Flush(); err != nil {
		panic(err)
	}
	fmt.Println(cnt)
	// Output:
	// 23
}

func ExampleDataLoader() {
	bFn := func(ctx context.Context, keys Keys) (map[Key]interface{}, error) {
		results := map[Key]interface{}{}
		for _, v := range keys {
			results[v] = v
		}
		return results, nil
	}
	b := NewBatch(context.TODO(), bFn, WithBatchSpan(10), WithCache(DefaultCache()))
	list := []int{}
	for i := 0; i < 23; i++ {
		b.Do(Int64Key(i), func(res interface{}) {
			list = append(list, int(res.(Int64Key)))
		})
	}
	if err := b.Flush(); err != nil {
		panic(err)
	}
	sort.Ints(list)
	fmt.Println(len(list))
	fmt.Println(list)
	// Output:
	// 23
	// [0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22]
}
