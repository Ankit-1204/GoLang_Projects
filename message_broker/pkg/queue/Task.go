package queue

import (
	"net"
	"sync"
	"time"
)

type Task struct {
	id          int
	topic       string
	status      string
	data        []byte
	createdTime time.Time
}

type Consumer struct {
	con   net.Conn
	topic string
}
type Queue struct {
	id      int
	mu      sync.Mutex
	Consu   []*Consumer
	Message []Task
	topic   string
}
