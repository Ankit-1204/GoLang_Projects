package main

import (
	"fmt"
	store "kvStore/internal"
)

func main() {
	kv := store.CreateStore()

	kv.Set("lang", "Go")
	if v, ok := kv.Get("lang"); ok {
		fmt.Println("Lang:", v)
	}
}
