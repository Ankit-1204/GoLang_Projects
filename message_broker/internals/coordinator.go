package internals

import (
	"io/ioutil"
	"net"
	"net/http"
	"queue/pkg/queue"
)

type Coordinator struct {
	subs map[string][]*queue.Queue
}

func Makecoordinator() *Coordinator {
	c := Coordinator{}
	c.subs = make(map[string][]*queue.Queue)
	return &c
}

func Publish(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
	}
}
func main() {
	c := Makecoordinator()
	ln, _ := net.Listen("tcp", "127.0.0.1:3333")
	for {

	}
}
