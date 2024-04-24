package main

import (
	"fmt"
	"log"
	"net"
)

const (
	PORT       = "8333"
	IP         = "127.0.0.1"
	ClientPort = ":8000"
)

// BitCoinNode contains IP and port of the
type BitCoinNode struct {
	IP   string
	Port string
}

// Handshake performs a network handshake with a Bitcoin node
func Handshake(node BitCoinNode) error {
	strPing := "PING"
	connString := fmt.Sprintf("%s:%s", node.IP, node.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", connString)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(strPing))
	if err != nil {
		return err
	}

	log.Println("Write to Server", strPing)
	reply := make([]byte, 1024)

	_, err = conn.Read(reply)
	if err != nil {
		return err
	}

	log.Println("Reply from Node Server\n", string(reply))
	conn.Close()
	return err
}

func main() {
	clientListener, err := net.Listen("tcp", ClientPort)
	if err != nil {
		panic(err)
	}
	defer clientListener.Close()

	log.Println("Client Listening...")
	// Start the Bitcoin node handshake
	node := BitCoinNode{
		IP:   IP,
		Port: PORT,
	}
	if err := Handshake(node); err != nil {
		panic(err)
	}
}
