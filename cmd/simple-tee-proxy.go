package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

type multipleHosts []string

func (i *multipleHosts) String() string {
	return fmt.Sprintf("%s", *i)
}

func (i *multipleHosts) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var hostingSide string
	var forwardSide string
	var shadowSide multipleHosts

	flag.StringVar(&hostingSide, "hosting", "0.0.0.0:7700", "What the external address is of the proxy")
	flag.StringVar(&forwardSide, "forward", "localhost:8800", "Where to forward all tcp connections to")
	flag.Var(&shadowSide, "shadow", "Where to send a copy of all connections to, the replies will be ignored, and the proxy will continue working regardless of the shadow side of things")

	flag.Parse()

	server, err := net.Listen("tcp", hostingSide)
	if err != nil {
		log.Fatalf("Failed to start listening on: %v, error: %v", hostingSide, err)
	}
	log.Printf("Listening on %v\n", hostingSide)

	for {
		newConnection, err := server.Accept()
		if err != nil {
			log.Fatalf("Failed to accept new connections: %v", err)
		}
		go forward(newConnection, forwardSide, shadowSide)
	}
}

func forward(source net.Conn, forwardSide string, shadowSide []string) {
	defer source.Close()

	forwardQueue := make(chan []byte, 512)
	shadowQueues := make([]chan []byte, len(shadowSide))
	for s := range shadowQueues {
		shadowQueues[s] = make(chan []byte, 512)
	}

	go func() {
		for {
			// a fresh buffer for every read, since we pass the buffer to the proxied connection
			buffer := make([]byte, 16*1024)
			read, err := source.Read(buffer)
			if err != nil {
				close(forwardQueue)
				for _, q := range shadowQueues {
					close(q)
				}
				return
			}
			if read > 0 {
				forwardQueue <- buffer[:read]
				for _, q := range shadowQueues {
					q <- buffer[:read]
				}
			}
		}

	}()

	replyQueue := make(chan []byte, 512)

	go connectBackend(forwardSide, forwardQueue, replyQueue)

	for i := range shadowQueues {
		go connectBackend(shadowSide[i], shadowQueues[i], nil)
	}

	for r := range replyQueue {
		source.Write(r)
	}
}

func connectBackend(target string, incoming <-chan []byte, outgoing chan<- []byte) {
	con, err := net.Dial("tcp", target)
	if err != nil {
		if outgoing != nil {
			log.Printf("Backend %v is not available: %v\n", target, err)
			close(outgoing)
		}
		// otherwise we are the shadow side, so not important if we can open or not
		// either way, just clear the incoming queue until it's closed
		for range incoming {
			// just drop it all,
		}
		return
	}
	if outgoing != nil {
		go func() {
			defer close(outgoing)
			for {
				buffer := make([]byte, 16*1024)
				read, err := con.Read(buffer)
				if err != nil {
					return
				}
				outgoing <- buffer[:read]
			}
		}()
	}
	connectionClosed := false
	for b := range incoming {
		written := 0
		for !connectionClosed && written < len(b) {
			w, err := con.Write(b[written:])
			if err != nil {
				connectionClosed = true
			}
			written += w
		}
	}
}
