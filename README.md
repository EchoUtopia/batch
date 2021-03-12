## Batch

this simple lib helps you with batch process your tasks, and you can also use it as dataloader


## Example:

#### dataloader
```go

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
		if err := b.Do(Int64Key(i), func(res interface{}) {
			list = append(list, int(res.(Int64Key)))
		}); err != nil {
			panic(err)
		}
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
```
#### batch

```go
func ExampleBatch() {
	cnt := 0
	bFn := func(ctx context.Context, keys Keys) (map[Key]interface{}, error) {
		cnt += len(keys)
		return nil, nil
	}
	b := NewBatch(context.TODO(), bFn, WithBatchSpan(10))
	for i := 0; i < 23; i++ {
		if err := b.Do(Int64Key(i)); err != nil {
			panic(err)
		}
	}
	if err := b.Flush(); err != nil {
		panic(err)
	}
	fmt.Println(cnt)
	// Output:
	// 23
}
```

