package main

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
