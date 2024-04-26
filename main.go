package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"net"
)

const (
	ServerPort      = 18445
	ProtocolVersion = 60002
)

// Handshake performs a network handshake with a Bitcoin node
func Handshake() error {
	// Connect to the Bitcoin node
	address := &net.TCPAddr{
		IP:   TCPAddress,
		Port: ServerPort,
	}

	// Connect to the address
	conn, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		log.Println("Error connecting to address:", err)
		return err
	}
	defer conn.Close()

	// Create and send the version message
	versionMsg := NewVersionMessage(ProtocolVersion, net.TCPAddr{IP: TCPAddress, Port: ServerPort})
	payload, err := versionMsg.ToRawMessage()
	if err != nil {
		return err
	}
	versCheck := versionMsg.CalculateSHA256()
	checksum := binary.BigEndian.Uint32(versCheck[:4])
	btcMessage := NewBtcMessage(Regtest.MagicValue(), versionMsg.Command(), checksum, payload)
	btcBytes, err := json.Marshal(btcMessage)
	if err != nil {
		return err
	}

	if _, err := conn.Write(btcBytes); err != nil {
		return err
	}

	log.Println("Waiting for server response...")
	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		log.Println("Error reading from connection:", err)
		return err
	}
	// Check if the received command matches the expected version command
	if btcMessage.Command != versionMsg.Command() {
		log.Println("Received unexpected command:", btcMessage.Command)
		return errors.New("wrong command")
	}

	log.Println("Connection established. Received message:", btcMessage.Command)
	return nil
}

func main() {
	if err := Handshake(); err != nil {
		log.Fatal("Error during handshake:", err)
	}
}
