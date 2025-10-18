package internals

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Consumer struct {
	id    int
	topic string
}

func MakeConsumer(topic string) *Consumer {
	id := os.Getpid()
	fmt.Println(topic)
	c := Consumer{id: id, topic: topic}
	return &c
}

func (c *Consumer) Subscribe() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}
	fmt.Println("here")
	defer conn.Close()
	_, err = conn.Write([]byte(c.topic + "\n"))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	reader := bufio.NewReader(conn)
	for {

		data, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("error retrieving msg")
			break
		}
		fmt.Println("here")
		fmt.Println(string(data))
	}
}
