package keyvalue

import (
	"sync"
)

type InMemory struct {
	mux  *sync.RWMutex
	data map[string]interface{}
}

func NewInMemory() *InMemory {
	return &InMemory{
		mux:  &sync.RWMutex{},
		data: map[string]interface{}{},
	}
}

func (m *InMemory) Exists(key string) bool {
	m.mux.RLock()
	defer m.mux.RUnlock()

	_, ok := m.data[key]

	return ok
}

func (m *InMemory) Get(key string) interface{} {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.data[key]
}

func (m *InMemory) Set(key string, value interface{}) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.data[key] = value
}

func (m *InMemory) Delete(key string) {
	delete(m.data, key)
}
