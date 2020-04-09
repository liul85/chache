package chache

import (
	"fmt"
	"reflect"
	"testing"
)

type Str string

func (s Str) Len() int64 {
	return int64(len(s))
}

func TestNewAdd(t *testing.T) {
	lru := New(int64(10), nil)
	lru.Add("key1", Str("value1"))

	if v, ok := lru.Get("key1"); !ok || string(v.(Str)) != "value1" {
		t.Fatal("Can't get item from cache with key 'key1'")
	}

	fmt.Println(lru.Len())
}

func TestShouldUpdateExistingItemWhenAdd(t *testing.T) {
	lru := New(int64(10), nil)
	lru.Add("key1", Str("value1"))
	lru.Add("key1", Str("value2"))

	if _, ok := lru.Get("key1"); !ok {
		t.Fatal("Can't get cache from cache with key 'key1'")
	}

	v, _ := lru.Get("key1")

	if string(v.(Str)) != "value2" {
		t.Fatal("value in cache was not updated")
	}
}

func TestShouldDeleteOldElementWhenExceedCapacity(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, Str(v1))
	lru.Add(k2, Str(v2))
	lru.Add(k3, Str(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatal("cache did not remove old element")
	}
}

func TestShouldInvokeCallbackFuncWhenRemoveOldElement(t *testing.T) {
	evictedKeys := make([]string, 0)
	callback := func(key string, value Value) {
		evictedKeys = append(evictedKeys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", Str("123456"))
	lru.Add("k2", Str("k2"))
	lru.Add("k3", Str("k3"))
	lru.Add("k4", Str("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, evictedKeys) {
		t.Fatal("Callback fun was not invoked")
	}
}
