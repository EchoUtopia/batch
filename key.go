package batch

import "strconv"

type Key interface {
	String() string
}

type StringKey string

func (k StringKey) String() string {
	return string(k)
}

type Int64Key int64

func (k Int64Key) String() string {
	return strconv.FormatInt(int64(k), 10)
}

type Keys []Key

func KeysFromStrings(keys []string) Keys {
	list := make(Keys, 0, len(keys))
	for _, v := range keys {
		list = append(list, StringKey(v))
	}
	return list
}
func KeysFromInt64s(keys []int64) Keys {
	list := make(Keys, 0, len(keys))
	for _, v := range keys {
		list = append(list, Int64Key(v))
	}
	return list
}
