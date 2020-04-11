package chache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func getter(key string) ([]byte, error) {
	return []byte(key), nil
}

func TestGetter(t *testing.T) {
	f := GetterFunc(getter)
	expect := []byte("key")

	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatal("callback failed")
	}
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))

	testCache := NewGroup("scores", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Printf("Get key %s from source.", key)
		if v, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key]++
			return []byte(v), nil
		}

		return nil, fmt.Errorf("Can't load key %s from source", key)

	}))

	for k, v := range db {
		if view, err := testCache.Get(k); err != nil || view.String() != v {
			t.Fatal("Failed to get value from cache")
		}

		if _, err := testCache.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s missed", k)
		}
	}

	if view, err := testCache.Get("unknown"); err == nil {
		t.Fatalf("The value of unknown should be empty, but got %s", view)
	}
}
