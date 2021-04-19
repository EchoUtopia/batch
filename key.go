package batch

import (
	"github.com/EchoUtopia/zerror/v2"
	"strconv"
)

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

func (keys *Keys) FromStringList(from []string) {
	list := make(Keys, 0, len(from))
	for _, v := range from {
		list = append(list, StringKey(v))
	}
	*keys = list
}
func (keys *Keys) FromInt64List(from []int64) {
	list := make(Keys, 0, len(from))
	for _, v := range from {
		list = append(list, Int64Key(v))
	}
	*keys = list
}

func (keys Keys) ToStringList() ([]string, error) {
	list := make([]string, 0, len(keys))
	for _, v := range keys {
		stringKey, ok := v.(StringKey)
		if !ok {
			return nil, zerror.BadRequest.New()
		}
		list = append(list, string(stringKey))
	}
	return list, nil
}

func (keys Keys) ToInt64List() ([]int64, error) {
	list := make([]int64, 0, len(keys))
	for _, v := range keys {
		int64Key, ok := v.(Int64Key)
		if !ok {
			return nil, zerror.BadRequest.New()
		}
		list = append(list, int64(int64Key))
	}
	return list, nil
}
