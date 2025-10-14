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
	subs map[string]*queue.Queue
	mu   sync.Mutex
}

func Makecoordinator() *Coordinator {
	c := Coordinator{}
	c.subs = make(map[string]*queue.Queue)
	return &c
}

func Publish(income queue.Incoming) {
	topic := income.Topic
	data := income.Data

}
func (c *Coordinator) Subscribe(topic string, subListener net.Conn) {
	consumer := queue.Consumer{Con: subListener, Topic: topic}
	_, ok := c.subs[topic]
	if !ok {
		uid := uuid.New().String()
		msgList := make([]queue.Task, 0)
		consumerList := make([]*queue.Consumer, 0)
		c.subs[topic] = &queue.Queue{Id: uid, Message: msgList, Topic: topic, Consu: consumerList}
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.subs[topic].Consu = append(c.subs[topic].Consu, &consumer)

}
func main() {
	c := Makecoordinator()
	sub, _ := net.Listen("tcp", "127.0.0.1:3333/sub")
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
				Publish(income)
			}(pubListener)
		}
	}()
}
