package keyvalue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTypeDecorator(t *testing.T) {
	type storageStruct struct {
		val  int
		time *time.Time
	}

	m := NewInMemory()
	mm := NewStorageValueTypeDecorator[storageStruct](m)

	key := "a"
	tt := time.Now()
	data := storageStruct{
		val:  5,
		time: &tt,
	}

	require.False(t, mm.Exists(key))
	require.NotNil(t, mm.Get(key))
	require.Empty(t, mm.Get(key)) // default value when key doesn't exist

	mm.Delete(key) // delete non-existing key, nothing happens

	mm.Set(key, data)
	require.True(t, mm.Exists(key))
	val := mm.Get(key)
	require.Equal(t, data, val)

	val.val = 6
	mm.Set(key, val)

	val = mm.Get(key)
	require.Equal(t, 6, val.val)

	mm.Delete(key)
	require.False(t, mm.Exists(key))
	require.NotNil(t, mm.Get(key))
	require.Empty(t, mm.Get(key)) // default value when key doesn't exist
}
