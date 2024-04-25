package btc

import (
	"encoding/binary"
	"errors"
)

// BtcMessage represents a Bitcoin message.
type BtcMessage struct {
	Magic    uint32
	Command  [12]byte
	Length   uint32
	Checksum uint32
	Payload  []byte
}

// NewBtcMessage creates a new BtcMessage instance.
func NewBtcMessage(magic uint32, command [12]byte, checksum uint32, payload []byte) *BtcMessage {
	return &BtcMessage{
		Magic:    magic,
		Command:  command,
		Length:   uint32(len(payload)),
		Checksum: checksum,
		Payload:  payload,
	}
}

// ToBytes converts the message to bytes.
func (msg *BtcMessage) ToBytes() []byte {
	buffer := make([]byte, 0, 4+12+4+4+len(msg.Payload))
	buffer = append(buffer, make([]byte, 4)...) // Placeholder for magic
	buffer = append(buffer, msg.Command[:]...)
	buffer = append(buffer, make([]byte, 4)...) // Placeholder for length
	buffer = append(buffer, make([]byte, 4)...) // Placeholder for checksum
	buffer = append(buffer, msg.Payload...)
	binary.LittleEndian.PutUint32(buffer[0:4], msg.Magic)
	binary.LittleEndian.PutUint32(buffer[16:20], msg.Length)
	binary.LittleEndian.PutUint32(buffer[20:24], msg.Checksum)
	return buffer
}

// FromBytes creates a BtcMessage from bytes.
func (msg *BtcMessage) FromBytes(data []byte) error {
	if len(data) < 24 {
		return errors.New("insufficient data for message")
	}

	msg.Magic = binary.LittleEndian.Uint32(data[0:4])
	copy(msg.Command[:], data[4:16])
	msg.Length = binary.LittleEndian.Uint32(data[16:20])
	msg.Checksum = binary.LittleEndian.Uint32(data[20:24])

	if len(data) < int(24+msg.Length) {
		return errors.New("insufficient data for payload")
	}

	msg.Payload = make([]byte, msg.Length)
	copy(msg.Payload, data[24:24+msg.Length])

	return nil
}
