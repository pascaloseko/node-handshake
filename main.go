package main

import (
	"encoding/binary"
	"errors"
	"log"
	"net"

	"github.com/pascaloseko/node-handshake/btc"
)

const (
	SERVER_PORT      = 18445
	CLIENT_PORT      = ":8000"
	PROTOCOL_VERSION = 60002
)

// Handshake performs a network handshake with a Bitcoin node
func Handshake() error {
	// Connect to the Bitcoin node
	address := &net.TCPAddr{
		IP:   btc.TCPAddress,
		Port: SERVER_PORT,
	}

	// Connect to the address
	conn, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		log.Println("Error connecting to address:", err)
		return err
	}
	defer conn.Close()

	// Create and send the version message
	versionMsg := btc.NewVersionMessage(PROTOCOL_VERSION, net.TCPAddr{IP: btc.TCPAddress, Port: SERVER_PORT})
	payload, err := versionMsg.ToRawMessage()
	if err != nil {
		return err
	}
	versCheck := versionMsg.CalculateSHA256()
	checksum := binary.BigEndian.Uint32(versCheck[:4])
	btcMessage := btc.NewBtcMessage(btc.Regtest.MagicValue(), versionMsg.Command(), checksum, payload)

	if _, err := conn.Write(btcMessage.ToBytes()); err != nil {
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
