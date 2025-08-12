package cache

import (
	"reflect"
	"sync"

	"github.com/renxzen/golidator/internal/fieldinfo"
)

type TypeInfo struct {
	Type   reflect.Type
	Fields []fieldinfo.Info
}

type TypeCache struct {
	mu    sync.RWMutex
	cache map[reflect.Type]*TypeInfo
}

func NewTypeCache() *TypeCache {
	return &TypeCache{
		cache: make(map[reflect.Type]*TypeInfo),
	}
}

func (tc *TypeCache) Get(t reflect.Type) *TypeInfo {
	tc.mu.RLock()
	info, exists := tc.cache[t]
	tc.mu.RUnlock()
	if exists {
		return info
	}

	tc.mu.Lock()
	defer tc.mu.Unlock()

	if info, exists := tc.cache[t]; exists {
		return info
	}

	computedInfo := tc.computeTypeInfo(t)
	tc.cache[t] = computedInfo
	return computedInfo
}

func (tc *TypeCache) computeTypeInfo(t reflect.Type) *TypeInfo {
	numField := t.NumField()
	info := &TypeInfo{
		Type:   t,
		Fields: make([]fieldinfo.Info, numField),
	}

	dummyValue := reflect.New(t).Elem()

	for i := 0; i < numField; i++ {
		fieldInfo := fieldinfo.ExtractInfo(dummyValue, i)
		fieldInfo.Value = reflect.Value{}
		info.Fields[i] = fieldInfo
	}

	return info
}

func (tc *TypeCache) GetWithValues(t reflect.Type, value reflect.Value) *TypeInfo {
	typeInfo := tc.Get(t)

	result := &TypeInfo{
		Type:   typeInfo.Type,
		Fields: make([]fieldinfo.Info, len(typeInfo.Fields)),
	}

	for i, fieldInfo := range typeInfo.Fields {
		fieldInfo.Value = value.Field(fieldInfo.Index)
		result.Fields[i] = fieldInfo
	}

	return result
}

func (tc *TypeCache) Clear() {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.cache = make(map[reflect.Type]*TypeInfo)
}

func (tc *TypeCache) Size() int {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return len(tc.cache)
}
