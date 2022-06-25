package main

import "sync"

func main() {
	mu := new(sync.Mutex) // want "unwrap_mutex"
	mu.Lock()
	defer mu.Unlock()
}
