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
	subs map[string]*queue.Queue
	mu   sync.Mutex
}

func Makecoordinator() *Coordinator {
	c := Coordinator{}
	c.subs = make(map[string]*queue.Queue)
	return &c
}

func (c *Coordinator) Deliver(topic string) {
	for msg := range c.subs[topic].Message {
		// not locked since .Message is a channel
		c.subs[topic].Mu.Lock()
		for _, consumer := range c.subs[topic].Consu {
			select {
			case consumer <- msg:
			default:
			}
		}
		c.subs[topic].Mu.Unlock()
	}
}

func (c *Coordinator) Publish(income queue.Incoming) {
	topic := income.Topic
	c.mu.Lock()
	_, ok := c.subs[topic]
	c.mu.Unlock()
	if !ok {
		uid := uuid.New().String()
		msgQueue := queue.Queue{Id: uid, Consu: make([]chan *queue.Incoming, 0), Message: make(chan *queue.Incoming, 50), Topic: topic}
		c.mu.Lock()
		c.subs[topic] = &msgQueue
		c.mu.Unlock()
		go c.Deliver(income.Topic)
	}
	c.subs[topic].Message <- &income

}
func (c *Coordinator) Subscribe(topic string) chan *queue.Incoming {
	conChan := make(chan *queue.Incoming, 10)
	c.mu.Lock()
	_, ok := c.subs[topic]
	if !ok {
		uid := uuid.New().String()
		msgQueue := queue.Queue{Id: uid, Consu: make([]chan *queue.Incoming, 0), Message: make(chan *queue.Incoming, 50), Topic: topic}

		c.subs[topic] = &msgQueue

		go c.Deliver(topic)
	}
	c.subs[topic].Consu = append(c.subs[topic].Consu, conChan)
	c.mu.Unlock()
	return conChan

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
				lchannel := c.Subscribe(topic)

				for msg := range lchannel {
					sublistener.Write(msg.Data)
				}
			}(subListener)
		}
	}()
	go func() {
		for {
			pubListener, _ := pub.Accept()
			go func(pubListener net.Conn) {
				defer pubListener.Close()
				pubReader := bufio.NewReader(pubListener)
				for {
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
