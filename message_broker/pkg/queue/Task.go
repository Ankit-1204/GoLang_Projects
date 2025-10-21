package queue

import (
	"net"
	"sync"
)

type Task struct {
	Topic  string
	Status string
	Data   []byte
}

type Consumer struct {
	Con   net.Conn
	Topic string
}
type Queue struct {
	Id      string
	Mu      sync.RWMutex
	Consu   []chan *Incoming
	Message chan *Incoming
	Topic   string
}

type Incoming struct {
	Topic string
	Data  []byte
}
