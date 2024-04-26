package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"net"
	"time"
)

// TCPAddress Define constants
var TCPAddress = net.IPv4(127, 0, 0, 1)

// VersionMessage struct definition
type VersionMessage struct {
	ProtocolVersion int32
	Service         uint64
	Timestamp       int64
	AddrRecv        net.TCPAddr
	AddrFrom        net.TCPAddr
	Nonce           uint64
	UserAgent       string
	StartHeight     int32
}

// NewVersionMessage creates a new VersionMessage instance.
func NewVersionMessage(protocolVersion int32, addrRecv net.TCPAddr) *VersionMessage {
	timestamp := calculateTimestamp()
	return &VersionMessage{
		ProtocolVersion: protocolVersion,
		Service:         0x1,
		Timestamp:       timestamp,
		AddrRecv:        addrRecv,
		AddrFrom:        net.TCPAddr{IP: TCPAddress, Port: 8000},
		Nonce:           generateRandomNonce(),
		UserAgent:       "",
		StartHeight:     1,
	}
}

// CalculateSHA256 calculates the SHA256 checksum of the message.
func (vm *VersionMessage) CalculateSHA256() [32]byte {
	payload, _ := vm.ToRawMessage()

	// Calculate the first hash
	h1 := sha256.New()
	h1.Write(payload)
	hash1 := h1.Sum(nil)

	// Calculate the second hash
	h2 := sha256.New()
	h2.Write(hash1)
	hash2 := h2.Sum(nil)

	var result [32]byte
	copy(result[:], hash2)

	return result
}

// ToRawMessage converts the VersionMessage to a raw message.
func (vm *VersionMessage) ToRawMessage() ([]byte, error) {
	log.Println(vm)
	svcBitmask := vm.nodeNetworkBitmask(0x1)
	addressBytes := vm.netAddrAsBytes(&svcBitmask, &vm.AddrRecv)
	zeroByte := make([]byte, 0)

	buffer := make([]byte, 0)
	buffer = append(buffer, int32ToBytesLE(vm.ProtocolVersion)...)
	buffer = append(buffer, uint64ToBytesLE(svcBitmask)...)
	buffer = append(buffer, int64ToBytesLE(vm.Timestamp)...)
	buffer = append(buffer, addressBytes...)
	buffer = append(buffer, make([]byte, 26)...) // addrFrom
	buffer = append(buffer, uint64ToBytesLE(vm.Nonce)...)
	buffer = append(buffer, zeroByte...) // user agent
	buffer = append(buffer, int32ToBytesLE(vm.StartHeight)...)
	buffer = append(buffer, zeroByte...) // relay

	return buffer, nil
}

// nodeNetworkBitmask generates the bitmask.
func (vm *VersionMessage) nodeNetworkBitmask(nodeNet uint64) uint64 {
	return nodeNet & 0xffffffffffffffff
}

// netAddrAsBytes returns the network address as bytes.
func (vm *VersionMessage) netAddrAsBytes(nodeBitmask *uint64, address *net.TCPAddr) []byte {
	buffer := make([]byte, 0)
	buffer = append(buffer, uint64ToBytesLE(*nodeBitmask)...)
	buffer = append(buffer, address.IP.To4()...)
	buffer = append(buffer, uint16ToBytesBE(uint16(address.Port))...)
	return buffer
}

// uint64ToBytesLE converts uint64 to bytes in little endian.
func uint64ToBytesLE(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)
	return b
}

// uint16ToBytesBE converts uint16 to bytes in big endian.
func uint16ToBytesBE(n uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return b
}

// int32ToBytes converts int32 to bytes.
func int32ToBytesLE(value int32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(value))
	return buf
}

// int64ToBytes converts int64 to bytes.
func int64ToBytesLE(value int64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(value))
	return buf
}

// calculateTimestamp calculates the current UNIX timestamp.
func calculateTimestamp() int64 {
	return time.Now().Unix()
}

// generateRandomNonce generates a random nonce value.
func generateRandomNonce() uint64 {
	var nonce uint64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &nonce); err != nil {
		log.Fatal("Error generating random nonce:", err)
	}
	return nonce
}

// Command returns the command bytes for the version message.
func (vm *VersionMessage) Command() [12]byte {
	return [12]byte{'v', 'e', 'r', 's', 'i', 'o', 'n', 0x00, 0x00, 0x00, 0x00, 0x00}
}
