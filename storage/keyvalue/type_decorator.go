package keyvalue

type Storage interface {
	Exists(key string) bool
	Get(key string) interface{}
	Set(key string, value interface{})
	Delete(key string)
}

type StorageValueTypeDecorator[V any] struct {
	storage Storage
}

func NewStorageValueTypeDecorator[V any](storage Storage) *StorageValueTypeDecorator[V] {
	return &StorageValueTypeDecorator[V]{
		storage: storage,
	}
}

func (v *StorageValueTypeDecorator[V]) Exists(key string) bool {
	return v.storage.Exists(key)
}

func (v *StorageValueTypeDecorator[V]) Get(key string) V {
	val := v.storage.Get(key)
	if val == nil {
		var emptyVar V

		return emptyVar
	}

	return val.(V)
}

func (v *StorageValueTypeDecorator[V]) Set(key string, value V) {
	v.storage.Set(key, value)
}

func (v *StorageValueTypeDecorator[V]) Delete(key string) {
	v.storage.Delete(key)
}
