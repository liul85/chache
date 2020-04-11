package chache

import (
	"fmt"
	"log"
	"sync"
)

//Getter defines interface of how get data from source
type Getter interface {
	Get(key string) ([]byte, error)
}

//GetterFunc implements Getter
type GetterFunc func(key string) ([]byte, error)

//Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.Mutex
	groups = make(map[string]*Group)
)

//NewGroup create new group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}

	mu.Lock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}

	groups[name] = g
	return g
}

//GetGroup returns group with specified namespace
func GetGroup(name string) *Group {
	return groups[name]
}

//Get returns value of a key
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("cache hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)

	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}

	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
