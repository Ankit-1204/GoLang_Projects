package queue

import (
	"net"
	"sync"
	"time"
)

type Task struct {
	Topic       string
	Status      string
	Data        []byte
	CreatedTime time.Time
}

type Consumer struct {
	Con   net.Conn
	Topic string
}
type Queue struct {
	Id      string
	Mu      sync.Mutex
	Consu   []*Consumer
	Message []Task
	Topic   string
}

type Incoming struct {
	Topic string
	Data  []byte
}
