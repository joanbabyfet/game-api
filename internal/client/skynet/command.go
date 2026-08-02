package skynet

import "errors"

// Message ID
//
// Client  -> Skynet
// Skynet  -> Client
//
// 建议预留区间，方便以后扩充。

const (
	// 保持该顺序 1000 User 2000 Wallet 3000 Slot 4000 Jackpot

	CmdLogin 		uint16 = 1001
	CmdKick 		uint16 = 1002

    CmdSpin         uint16 = 2001
	CmdFreeSpin     uint16 = 2002

    CmdPing         uint16 = 9999
	
	// Response（可选）
	MsgOK    uint16 = 9000
	MsgError uint16 = 9001
)

// Default Config
const (

	// Header
	PacketHeaderSize = 10

	// 4MB
	DefaultMaxPacketSize = 4 * 1024 * 1024

	DefaultReadBuffer  = 64 * 1024
	DefaultWriteBuffer = 64 * 1024
)

// Timeout
const (

	DefaultDialTimeout = 5

	DefaultReadTimeout = 10

	DefaultWriteTimeout = 10
)

// Error

var (

	// Connection
	ErrClosed = errors.New("skynet: connection closed")

	ErrTimeout = errors.New("skynet: timeout")

	ErrReconnect = errors.New("skynet: reconnect failed")

	// Packet
	ErrPacketTooLarge = errors.New("skynet: packet too large")

	ErrInvalidPacket = errors.New("skynet: invalid packet")

	ErrInvalidLength = errors.New("skynet: invalid packet length")

	// Protocol
	ErrUnknownMessage = errors.New("skynet: unknown message")

	ErrSequenceNotFound = errors.New("skynet: sequence not found")
)