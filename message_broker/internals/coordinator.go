package internals

import (
	"bufio"
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

func Publish() {
	
}
func Subscribe(){

}
func main() {
	c := Makecoordinator()
	sub, _ := net.Listen("tcp", "127.0.0.1:3333/sub")
	pub, _ := net.Listen("tcp", "127.0.0.1:3333/pub")
	go func(){
		for {
			subListener,err:=sub.Accept()
			// sub will hold till it recieves request
			go func(sublistener net.Conn){
				subReader:=bufio.NewReader(subListener)
				// basically \n so that it listens until first delim (line here)
				topic,_:=reader.ReadString("\n") 
				
			}
		}
	}
	go func(){
		for{
			pubListener,err:=pub.Accept()
			go func(pubListener net.Conn){
				pubReader:=bufio.NewReader(pubListener)
				

			}
		}
	}
}
