package btc

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"time"
)

// Define constants
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
	svcBitmask := vm.NodeNetworkBitmask(0x1)
	addressBytes := vm.NetAddrAsBytes(&svcBitmask, &vm.AddrRecv)

	buffer := make([]byte, 0)
	buffer = append(buffer, Int32ToBytes(vm.ProtocolVersion)...)
	buffer = append(buffer, Uint64ToBytesLE(svcBitmask)...)
	buffer = append(buffer, Int64ToBytes(vm.Timestamp)...)
	buffer = append(buffer, addressBytes...)
	buffer = append(buffer, make([]byte, 26)...) // addrFrom
	buffer = append(buffer, Uint64ToBytesLE(vm.Nonce)...)
	buffer = append(buffer, byte(0)) // user agent
	buffer = append(buffer, Int32ToBytes(vm.StartHeight)...)
	buffer = append(buffer, byte(0)) // relay

	return buffer, nil
}

// NodeNetworkBitmask generates the bitmask.
func (vm *VersionMessage) NodeNetworkBitmask(nodeNet uint64) uint64 {
	return nodeNet & 0xffffffffffffffff
}

// NetAddrAsBytes returns the network address as bytes.
func (vm *VersionMessage) NetAddrAsBytes(nodeBitmask *uint64, address *net.TCPAddr) []byte {
	buffer := make([]byte, 0)
	buffer = append(buffer, Uint64ToBytesLE(*nodeBitmask)...)
	buffer = append(buffer, address.IP.To4()...)
	buffer = append(buffer, Uint16ToBytesBE(uint16(address.Port))...)
	return buffer
}

// FromBytes converts bytes to a VersionMessage.
func VersionMessageFromBytes(data []byte) (*VersionMessage, error) {
	vm := &VersionMessage{}

	vm.ProtocolVersion = BytesToInt32(data[:4])
	vm.Service = BytesToUint64(data[4:12])
	vm.Timestamp = BytesToInt64(data[12:20])

	recvAddr, err := NetAddrFromBytes(data[20:46])
	if err != nil {
		return nil, err
	}
	vm.AddrRecv = *recvAddr

	fromAddr, err := NetAddrFromBytes(data[46:72])
	if err != nil {
		return nil, err
	}
	vm.AddrFrom = *fromAddr

	vm.Nonce = BytesToUint64(data[72:80])

	userAgent, _, err := ReadVarStr(data[80:])
	if err != nil {
		return nil, err
	}
	vm.UserAgent = userAgent

	vm.StartHeight = BytesToInt32(data[len(data)-5 : len(data)-1])

	return vm, nil
}

// NetAddrFromBytes converts bytes to a network address.
func NetAddrFromBytes(data []byte) (*net.TCPAddr, error) {
	if len(data) < 18 {
		return nil, errors.New("invalid data length for network address")
	}

	var ip net.IP
	if len(data) == 18 {
		ip = net.IP(data[0:16])
	} else {
		ip = net.IP(data[0:4])
	}
	port := binary.BigEndian.Uint16(data[len(data)-2:])

	return &net.TCPAddr{IP: ip, Port: int(port)}, nil
}

// Uint64ToBytesLE converts uint64 to bytes in little endian.
func Uint64ToBytesLE(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)
	return b
}

// Uint64ToBytesBE converts uint64 to bytes in big-endian.
func Uint64ToBytesBE(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

func Uint32ToBytesBE(n uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return b
}

// Uint16ToBytesBE converts uint16 to bytes in big endian.
func Uint16ToBytesBE(n uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return b
}

// Int32ToBytes converts int32 to bytes.
func Int32ToBytes(value int32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(value))
	return buf
}

// BytesToInt32 converts bytes to int32.
func BytesToInt32(data []byte) int32 {
	return int32(binary.LittleEndian.Uint32(data))
}

// Int64ToBytes converts int64 to bytes.
func Int64ToBytes(value int64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(value))
	return buf
}

// BytesToInt64 converts bytes to int64.
func BytesToInt64(data []byte) int64 {
	return int64(binary.LittleEndian.Uint64(data))
}

// Uint64ToBytes converts uint64 to bytes.
func Uint64ToBytes(value uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, value)
	return buf
}

// BytesToUint64 converts bytes to uint64.
func BytesToUint64(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}

// StringToVarStr converts a string to a variable length string.
func StringToVarStr(value string) []byte {
	length := len(value)
	buf := make([]byte, 0, length+1)
	buf = append(buf, byte(length))
	buf = append(buf, []byte(value)...)
	return buf
}

// ReadVarStr reads a variable length string from bytes.
func ReadVarStr(data []byte) (string, []byte, error) {
	if len(data) == 0 {
		return "", nil, errors.New("unexpected end of data")
	}
	length := int(data[0])
	if len(data) < length+1 {
		return "", nil, errors.New("unexpected end of data")
	}
	str := string(data[1 : length+1])
	remaining := data[length+1:]
	return str, remaining, nil
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
