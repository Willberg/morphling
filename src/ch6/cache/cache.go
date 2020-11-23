// 缓存
package main

import (
	"fmt"
	"sync"
)

var cache = struct {
	sync.Mutex
	mapping map[string]string
}{
	mapping: make(map[string]string),
}

func lookup(key string) string {
	return cache.mapping[key]
}

func add(key, val string) {
	cache.Lock()
	cache.mapping[key] = val
	cache.Unlock()
}

func main() {
	add("1", "1+")
	fmt.Println(lookup("1"))
}
