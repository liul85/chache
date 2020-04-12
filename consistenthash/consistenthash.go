package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//Hash function
type Hash func(data []byte) uint32

//Map maintain all the keys
type Map struct {
	hashFunc    Hash
	vNodeCount  int
	nodeHashes  []int
	nodeHashMap map[int]string
}

//New create a new Map
func New(vNodeCount int, fn Hash) *Map {
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}

	return &Map{
		vNodeCount:  vNodeCount,
		hashFunc:    fn,
		nodeHashMap: make(map[int]string),
	}
}

//Add allow map to add new nodes
func (m *Map) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < m.vNodeCount; i++ {
			hash := int(m.hashFunc([]byte(strconv.Itoa(i) + node)))
			m.nodeHashes = append(m.nodeHashes, hash)
			m.nodeHashMap[hash] = node
		}
	}
	sort.Ints(m.nodeHashes)
}

//Get returns close node with a given key
func (m *Map) Get(key string) string {
	if len(m.nodeHashes) == 0 {
		return ""
	}

	hash := int(m.hashFunc([]byte(key)))

	idx := sort.Search(len(m.nodeHashes), func(i int) bool {
		return m.nodeHashes[i] >= hash
	})

	return m.nodeHashMap[m.nodeHashes[idx%len(m.nodeHashes)]]
}
