package internals

import (
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Coordinator struct {
	client *redis.Client
}

func Makecoordinator() {
	c := Coordinator{}
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})
	c.client = client
	http.HandleFunc("/enqueue", Enqueue)
	err := http.ListenAndServe("127.0.0.1:3333", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Enqueue(w http.ResponseWriter, r *http.Request) {

}
