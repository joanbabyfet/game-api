package skynet

import (
	"context"
	"fmt"
	"game-api/proto/rpcpb"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"
)

// Client 是 Skynet TCP Client。
// transport.go 会基于此结构实现连接、收发、重连等逻辑。
type Client struct {
	addr string

	// TCP
	conn net.Conn

	// 状态
	closed atomic.Bool

	// Sequence
	sequence atomic.Uint32

	// 并发控制
	mu sync.RWMutex

	// 等待响应
	pending sync.Map // map[uint32]chan *Packet

	// 写队列（transport.go 使用）
	writeCh chan *Packet

	// 关闭信号
	closeCh chan struct{}

	// 配置
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration

	maxPacketSize uint32
}

// New 创建 Client，但不会主动连接。
// 调用 Connect()（在 transport.go 中实现）后建立连接。
func New(addr string) *Client {
	c := &Client{
		addr:          addr,
		writeCh:       make(chan *Packet, 1024),
		closeCh:       make(chan struct{}),
		dialTimeout:   DefaultDialTimeout * time.Second,
		readTimeout:   DefaultReadTimeout * time.Second,
		writeTimeout:  DefaultWriteTimeout * time.Second,
		maxPacketSize: DefaultMaxPacketSize,
	}

	c.sequence.Store(0)

	return c
}

// Addr 返回服务器地址。
func (c *Client) Addr() string {
	return c.addr
}

// IsClosed 返回连接是否关闭。
func (c *Client) IsClosed() bool {
	return c.closed.Load()
}

// NextSequence 返回下一个请求序号。
func (c *Client) NextSequence() uint32 {
	return c.sequence.Add(1)
}

// RegisterPending 注册等待响应。
func (c *Client) RegisterPending(seq uint32) chan *Packet {
	ch := make(chan *Packet, 1)
	c.pending.Store(seq, ch)
	return ch
}

// RemovePending 删除等待响应。
func (c *Client) RemovePending(seq uint32) {
	c.pending.Delete(seq)
}

// FindPending 查找等待响应。
func (c *Client) FindPending(seq uint32) (chan *Packet, bool) {
	v, ok := c.pending.Load(seq)
	if !ok {
		return nil, false
	}
	ch, ok := v.(chan *Packet)
	return ch, ok
}

// Close 标记客户端关闭。
// transport.go 会扩展真正关闭 conn 等逻辑。
func (c *Client) Close() error {
	if c.closed.CompareAndSwap(false, true) {
		close(c.closeCh)

		c.mu.Lock()
		if c.conn != nil {
			_ = c.conn.Close()
			c.conn = nil
		}
		c.mu.Unlock()
	}
	return nil
}

// 像调用本地函数一样，去调用另一台机器（或另一个进程）的函数
// Call 发起一次同步 RPC 调用。
// 请求与响应均使用 Protobuf 编解码。
func (c *Client) Call(ctx context.Context, cmdID uint16, req proto.Message, resp proto.Message) error {

	payload, err := proto.Marshal(req)
    if err != nil {
        return err
    }

	packet := &Packet{
		Type:    PacketRequest,
		CmdID:   cmdID,
		Payload: payload,
	}

	reply, err := c.Request(ctx, packet)
	if err != nil {
		return err
	}

	switch reply.Type {

	case PacketResponse:

		return proto.Unmarshal(reply.Payload, resp)

	case PacketError:

		var rpcErr rpcpb.Error

		if err := proto.Unmarshal(reply.Payload, &rpcErr); err != nil {
			return err
		}

		fmt.Printf("rpcErr.Code = %d\n", rpcErr.Code)
		fmt.Printf("rpcErr.Msg  = %s\n", rpcErr.Msg)

		return &Error{
			Code: int(rpcErr.Code),
			Msg:  rpcErr.Msg,
		}

	default:

		return ErrInvalidPacket
	}
}