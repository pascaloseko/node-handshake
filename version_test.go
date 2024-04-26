package main

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVersionMessage(t *testing.T) {
	addrRecv := net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8000}
	vm := NewVersionMessage(70015, addrRecv)

	assert.Equal(t, int32(70015), vm.ProtocolVersion)
	assert.Equal(t, uint64(0x1), vm.Service)
	assert.Equal(t, calculateTimestamp(), vm.Timestamp)
	assert.Equal(t, addrRecv, vm.AddrRecv)
	assert.Equal(t, net.TCPAddr{IP: TCPAddress, Port: 8000}, vm.AddrFrom)
	assert.NotEqual(t, uint64(0), vm.Nonce)
	assert.Equal(t, "", vm.UserAgent)
	assert.Equal(t, int32(1), vm.StartHeight)
}

func TestCalculateSHA256(t *testing.T) {
	addrRecv := net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8000}
	vm := NewVersionMessage(70015, addrRecv)
	sha256 := vm.CalculateSHA256()

	assert.Len(t, sha256, 32)
}

func TestToRawMessage(t *testing.T) {
	addrRecv := net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8000}
	vm := NewVersionMessage(70015, addrRecv)
	rawMessage, err := vm.ToRawMessage()

	assert.NoError(t, err)
	assert.Len(t, rawMessage, 72)
}

func TestNodeNetworkBitmask(t *testing.T) {
	vm := NewVersionMessage(70015, net.TCPAddr{})
	bitmask := vm.nodeNetworkBitmask(0x1)

	assert.Equal(t, uint64(0x1), bitmask)
}

func TestNetAddrAsBytes(t *testing.T) {
	addrRecv := net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8000}
	vm := NewVersionMessage(70015, addrRecv)
	var nodeBitmask uint64 = 0x1
	addressBytes := vm.netAddrAsBytes(&nodeBitmask, &addrRecv)

	assert.Len(t, addressBytes, 14)
}

func TestUint64ToBytesLE(t *testing.T) {
	n := uint64(1234567890)
	bytes := uint64ToBytesLE(n)

	assert.Len(t, bytes, 8)
	assert.Equal(t, []byte{0xd2, 0x2, 0x96, 0x49, 0x0, 0x0, 0x0, 0x0}, bytes)
}

func TestUint16ToBytesBE(t *testing.T) {
	n := uint16(12345)
	bytes := uint16ToBytesBE(n)

	assert.Len(t, bytes, 2)
	assert.Equal(t, []byte{0x30, 0x39}, bytes)
}

func TestInt32ToBytes(t *testing.T) {
	value := int32(12345)
	bytes := int32ToBytesLE(value)

	assert.Len(t, bytes, 4)
	assert.Equal(t, []byte{0x39, 0x30, 0x0, 0x0}, bytes)
}

func TestInt64ToBytes(t *testing.T) {
	value := int64(1234567890)
	bytes := int64ToBytesLE(value)

	assert.Len(t, bytes, 8)
	assert.Equal(t, []byte{0xd2, 0x2, 0x96, 0x49, 0x0, 0x0, 0x0, 0x0}, bytes)
}

func TestCalculateTimestamp(t *testing.T) {
	timestamp := calculateTimestamp()

	assert.True(t, timestamp > 0)
}

func TestGenerateRandomNonce(t *testing.T) {
	nonce := generateRandomNonce()

	assert.NotEqual(t, uint64(0), nonce)
}

func TestCommand(t *testing.T) {
	vm := NewVersionMessage(70015, net.TCPAddr{})
	command := vm.Command()

	assert.Equal(t, [12]byte{'v', 'e', 'r', 's', 'i', 'o', 'n', 0x00, 0x00, 0x00, 0x00, 0x00}, command)
}
