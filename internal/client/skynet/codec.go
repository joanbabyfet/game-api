package skynet

import (
	"encoding/binary"
	"fmt"
	"io"
)

//
// Packet Type
//

const (
	PacketRequest uint8 = 0
	PacketResponse uint8 = 1
	PacketError uint8 = 2
)

//
// Packet
//
// +----------+--------+------------+------------+--------------+
// | Length   | Type   | CmdID      | SeqID      | Payload      |
// | uint32BE | uint8  | uint16BE   | uint32BE   | protobuf     |
// +----------+--------+------------+------------+--------------+
//
// Length = Type + CmdID + SeqID + Payload
//

type Packet struct {
	Type    uint8
	CmdID   uint16
	SeqID   uint32
	Payload []byte
}

//
// Header
//
// Type(1) + CmdID(2) + SeqID(4)
//

const (
	LengthSize = 4
	TypeSize   = 1
	CmdSize    = 2
	SeqSize    = 4

	HeaderSize = TypeSize + CmdSize + SeqSize
)

//
// Encode
//
// Length(4) + Type(1) + CmdID(2) + SeqID(4) + Payload
//

func Encode(p *Packet) ([]byte, error) {

	if p == nil {
		return nil, ErrInvalidPacket
	}

	bodyLen := HeaderSize + len(p.Payload)

	if bodyLen > DefaultMaxPacketSize {
		return nil, ErrPacketTooLarge
	}

	buf := make([]byte, LengthSize+bodyLen)

	// Length
	binary.BigEndian.PutUint32(buf[0:4], uint32(bodyLen))

	// Type
	buf[4] = p.Type

	// CmdID
	binary.BigEndian.PutUint16(buf[5:7], p.CmdID)

	// SeqID
	binary.BigEndian.PutUint32(buf[7:11], p.SeqID)

	// Payload
	copy(buf[LengthSize+HeaderSize:], p.Payload)

	return buf, nil
}

//
// Decode
//

func Decode(data []byte) (*Packet, error) {

	if len(data) < LengthSize+HeaderSize {
		return nil, ErrInvalidPacket
	}

	bodyLen := binary.BigEndian.Uint32(data[:4])

	if bodyLen < HeaderSize {
		return nil, ErrInvalidLength
	}

	if int(bodyLen)+LengthSize != len(data) {
		return nil, fmt.Errorf(
			"%w: expect=%d actual=%d",
			ErrInvalidLength,
			bodyLen+LengthSize,
			len(data),
		)
	}

	packetType := data[4]

	cmdID := binary.BigEndian.Uint16(data[5:7])

	seqID := binary.BigEndian.Uint32(data[7:11])

	payload := make([]byte, int(bodyLen)-HeaderSize)

	copy(payload, data[11:])

	return &Packet{
		Type:    packetType,
		CmdID:   cmdID,
		SeqID:   seqID,
		Payload: payload,
	}, nil
}

//
// ReadPacket
//

func ReadPacket(r io.Reader) (*Packet, error) {

	var lengthBuf [LengthSize]byte

	if _, err := io.ReadFull(r, lengthBuf[:]); err != nil {
		return nil, err
	}

	bodyLen := binary.BigEndian.Uint32(lengthBuf[:])

	if bodyLen < HeaderSize {
		return nil, ErrInvalidLength
	}

	if bodyLen > DefaultMaxPacketSize {
		return nil, ErrPacketTooLarge
	}

	body := make([]byte, bodyLen)

	if _, err := io.ReadFull(r, body); err != nil {
		return nil, err
	}

	packet := make([]byte, LengthSize+len(body))

	copy(packet[:LengthSize], lengthBuf[:])
	copy(packet[LengthSize:], body)

	return Decode(packet)
}

//
// WritePacket
//

func WritePacket(w io.Writer, p *Packet) error {

	data, err := Encode(p)
	if err != nil {
		return err
	}

	total := 0

	for total < len(data) {

		n, err := w.Write(data[total:])
		if err != nil {
			return err
		}

		total += n
	}

	return nil
}