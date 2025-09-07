package main

import (
	store "kvStore/internal"
	wal "kvStore/internal"
	"log"
)

func main() {
	kv := store.CreateStore()
	w, err := wal.CreateLogFile("data")
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	err = w.Recover(func(rec wal.LogFile) {
		switch rec.Operation {
		case wal.OpSet:
			kv.Set(rec.Key, rec.Data)
		case wal.OpDelete:
			kv.Delete(rec.Key)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}
