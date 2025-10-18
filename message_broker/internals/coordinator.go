package internals

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"queue/pkg/queue"
	"strings"
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
			q.Mu.Unlock()
			return
		}
		t := q.Message[0]
		q.Mu.Unlock()
		fmt.Println(t.Data)
		_, err := q.Consu.Con.Write(t.Data)
		q.Mu.Lock()
		if err == nil {
			q.Message = q.Message[1:]
		}
		q.Mu.Unlock()
	}

}

func (c *Coordinator) Publish(income queue.Incoming) {
	topic := income.Topic
	data := string(income.Data) + "\n"
	c.mu.Lock()
	fmt.Println(topic)
	qlist, ok := c.subs[topic]
	fmt.Println(c.subs)
	if !ok {
		sub := make([]*queue.Queue, len(qlist))
		c.subs[topic] = sub
	}
	sublist := make([]*queue.Queue, len(qlist))
	copy(sublist, qlist)
	c.mu.Unlock()
	fmt.Println(qlist)
	for _, q := range sublist {

		task := &queue.Task{Topic: topic, Status: "pending", Data: []byte(data)}
		fmt.Println(*task)
		go Deliver(task, q)
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
func Run() {
	c := Makecoordinator()
	sub, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	pub, err := net.Listen("tcp", "127.0.0.1:3333")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("here")
	go func() {
		for {
			subListener, _ := sub.Accept()
			// sub will hold till it recieves request
			go func(sublistener net.Conn) {
				subReader := bufio.NewReader(sublistener)
				// basically \n so that it listens until first delim (line here)

				topic, err := subReader.ReadString('\n')
				if err == io.EOF {
					fmt.Println("conn lost sub")
					return
				}
				topic = strings.TrimSuffix(topic, "\n")
				fmt.Println(topic)
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
				for {
					pubReader := bufio.NewReader(pubListener)
					data, err := pubReader.ReadBytes('\n')
					if err == io.EOF {
						fmt.Println("conn lost pub")
						return
					}
					var income queue.Incoming
					err = json.Unmarshal(data, &income)
					if err != nil {
						fmt.Println(err)
						break
					}
					c.Publish(income)
				}
			}(pubListener)
		}
	}()
}
