package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	testMap := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	testMap.Add("6", "4", "2")

	testKeys := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testKeys {
		if testMap.Get(k) != v {
			t.Errorf("did not get value %s for key %s", v, k)
		}
	}

	testMap.Add("8")

	testKeys["27"] = "8"

	for k, v := range testKeys {
		if testMap.Get(k) != v {
			t.Errorf("did not get value %s for key %s", v, k)
		}
	}
}
