package helper

import "sync"

type SyncMap[K comparable, V any] struct {
	m  map[K]V
	rw *sync.RWMutex
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return NewSyncMapSized[K, V](0)
}

func NewSyncMapSized[K comparable, V any](size int) *SyncMap[K, V] {
	return &SyncMap[K, V]{
		m:  make(map[K]V, size),
		rw: &sync.RWMutex{},
	}
}

func (sm *SyncMap[K, V]) Get(key K) (V, bool) {
	sm.rw.RLock()
	defer sm.rw.RUnlock()
	value, ok := sm.m[key]

	return value, ok
}

func (sm *SyncMap[K, V]) Set(key K, value V) {
	sm.rw.Lock()
	defer sm.rw.Unlock()

	sm.m[key] = value
}

func (sm *SyncMap[K, V]) Delete(key K) {
	sm.rw.Lock()
	defer sm.rw.Unlock()

	delete(sm.m, key)
}

func (sm *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	sm.rw.RLock()
	defer sm.rw.RUnlock()

	for k, v := range sm.m {
		if !f(k, v) {
			break
		}
	}
}

func (sm *SyncMap[K, V]) Len() int {
	sm.rw.RLock()
	defer sm.rw.RUnlock()

	return len(sm.m)
}
