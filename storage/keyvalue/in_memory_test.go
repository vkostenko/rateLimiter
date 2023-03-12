package keyvalue

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInMemory(t *testing.T) {
	keyA := "a"
	keyB := "b"

	m := NewInMemory()

	require.False(t, m.Exists(keyA))
	require.Nil(t, m.Get(keyA)) // default value when key doesn't exist

	m.Delete(keyA) // delete non-existing key, nothing happens

	m.Set(keyA, 5)
	m.Set(keyB, 3)
	require.True(t, m.Exists(keyA))
	require.Equal(t, 5, m.Get(keyA))

	m.Delete(keyA)
	require.False(t, m.Exists(keyA))
	require.Nil(t, m.Get(keyA)) // default value when key doesn't exist

	require.True(t, m.Exists(keyB))
	require.Equal(t, 3, m.Get(keyB))

	m.Set(keyB, 15) // rewrite value
	require.True(t, m.Exists(keyB))
	require.Equal(t, 15, m.Get(keyB))
}
