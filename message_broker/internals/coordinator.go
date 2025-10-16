package internals

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"queue/pkg/queue"
	"sync"

	"github.com/google/uuid"
)

type Coordinator struct {
	subs map[string][]*queue.Queue
	mu   sync.Mutex
}

func Makecoordinator() *Coordinator {
	c := Coordinator{}
	c.subs = make(map[string][]*queue.Queue)
	return &c
}

func Deliver(task *queue.Task, q *queue.Queue) {
	q.Mu.Lock()
	q.Message = append(q.Message, task)
	q.Mu.Unlock()

	for {
		q.Mu.Lock()
		if len(q.Message) == 0 {
			return
		}
		t := q.Message[0]
		_, err := q.Consu.Con.Write(t.Data)
		if err == nil {
			q.Message = q.Message[1:]
		}
		q.Mu.Unlock()
	}

}

func (c *Coordinator) Publish(income queue.Incoming) {
	topic := income.Topic
	data := income.Data
	task := queue.Task{Topic: topic, Status: "pending", Data: data}
	for _, q := range c.subs[topic] {
		go Deliver(&task, q)
	}

}
func (c *Coordinator) Subscribe(topic string, subListener net.Conn) {
	consumer := queue.Consumer{Con: subListener, Topic: topic}
	uid := uuid.New().String()
	msgList := make([]*queue.Task, 0)
	msgQueue := queue.Queue{Id: uid, Consu: &consumer, Message: msgList, Topic: topic}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.subs[topic] = append(c.subs[topic], &msgQueue)

}
func main() {
	c := Makecoordinator()
	sub, _ := net.Listen("tcp", "127.0.0.1:8000/sub")
	pub, _ := net.Listen("tcp", "127.0.0.1:3333/pub")
	go func() {
		for {
			subListener, _ := sub.Accept()
			// sub will hold till it recieves request
			go func(sublistener net.Conn) {
				subReader := bufio.NewReader(subListener)
				// basically \n so that it listens until first delim (line here)
				topic, _ := subReader.ReadString('\n')
				c.Subscribe(topic, sublistener)
				scanner := bufio.NewScanner(sublistener)
				for scanner.Scan() {
				}
			}(subListener)
		}
	}()
	go func() {
		for {
			pubListener, _ := pub.Accept()
			go func(pubListener net.Conn) {
				pubReader := bufio.NewReader(pubListener)
				data, _ := pubReader.ReadBytes('\n')
				var income queue.Incoming
				err := json.Unmarshal(data, &income)
				if err != nil {
					log.Panic(err)
				}
				c.Publish(income)
			}(pubListener)
		}
	}()
}
