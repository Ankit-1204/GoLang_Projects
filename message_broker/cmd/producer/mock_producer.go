package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"queue/pkg/queue"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3333")
	if err != nil {
		fmt.Println("Error connecting to coordinator:", err)
		os.Exit(1)
	}
	defer conn.Close()
	i := 0
	for {
		i++
		msg := queue.Incoming{
			Topic: "prod1",
			Data:  []byte(fmt.Sprintf("This msg %d \n", i)),
		}
		data, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			break
		}
		_, err = conn.Write(append(data, '\n'))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}

		fmt.Printf("Sent message %d\n", i)
		time.Sleep(3 * time.Second)
	}
}
