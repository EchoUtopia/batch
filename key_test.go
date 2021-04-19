package batch

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeys_FromInt64List(t *testing.T) {
	var keys Keys
	keys.FromInt64List([]int64{1, 2, 3})
	require.ElementsMatch(t, Keys{Int64Key(1), Int64Key(2), Int64Key(3)}, keys)
}

func TestKeys_FromStringList(t *testing.T){
	var keys Keys
	keys.FromStringList([]string{`1`, `2`})
	require.ElementsMatch(t, Keys{StringKey(`1`), StringKey(`2`)}, keys)
}

func TestStringKeysToStringList(t *testing.T){
	keys := Keys{
		StringKey(`1`),
		StringKey(`2`),
	}
	list, err := keys.ToStringList()
	require.Nil(t, err)
	require.ElementsMatch(t, []string{`1`, `2`}, list)
	_, err = keys.ToInt64List()
	require.NotNil(t, err)
}

func TestStringKeysToInt64List(t *testing.T){
	keys := Keys{
		Int64Key(1),
		Int64Key(2),
	}
	list, err := keys.ToInt64List()
	require.Nil(t, err)
	require.ElementsMatch(t, []int64{1,2}, list)
	_, err = keys.ToStringList()
	require.NotNil(t, err)
	keys = append(keys, StringKey(`1`))

	list, err = keys.ToInt64List()
	require.NotNil(t, err)
}