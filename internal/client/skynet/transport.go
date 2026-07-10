package skynet

import (
	"context"
	"errors"
	"net"
	"time"
)

// Connect 建立 TCP 连接并启动读写协程。
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	conn, err := net.DialTimeout("tcp", c.addr, c.dialTimeout)
	if err != nil {
		return err
	}

	c.conn = conn
	c.closed.Store(false)

	go c.readLoop()
	go c.writeLoop()

	return nil
}

// Reconnect 重连。
func (c *Client) Reconnect() error {
	_ = c.Close()

	// 重新初始化关闭通道
	c.closeCh = make(chan struct{})

	return c.Connect()
}

// Send 异步发送。
func (c *Client) Send(p *Packet) error {
	if c.IsClosed() {
		return ErrClosed
	}

	select {
	case c.writeCh <- p:
		return nil
	case <-c.closeCh:
		return ErrClosed
	}
}

// Request 发送请求并等待响应。
func (c *Client) Request(ctx context.Context, p *Packet) (*Packet, error) {
	if p == nil {
		return nil, ErrInvalidPacket
	}

	if ctx == nil {
		ctx = context.Background()
	}

	// 自动分配 SeqID
	p.SeqID = c.NextSequence()

	// 注册等待响应
	ch := c.RegisterPending(p.SeqID)
	defer c.RemovePending(p.SeqID)

	// 发送数据
	if err := c.Send(p); err != nil {
		return nil, err
	}

	// 等待响应
	select {

	case resp := <-ch:
		return resp, nil

	case <-ctx.Done():
		return nil, ctx.Err()

	case <-c.closeCh:
		return nil, ErrClosed
	}
}

func (c *Client) writeLoop() {
	for {
		select {
		case <-c.closeCh:
			return

		case p := <-c.writeCh:
			// fmt.Printf("SEND type=%d cmd=%d seq=%d payload=%d\n",
			// 	p.Type,
			// 	p.CmdID,
			// 	p.SeqID,
			// 	len(p.Payload),
			// )

			c.mu.RLock()
			conn := c.conn
			c.mu.RUnlock()

			if conn == nil {
				continue
			}

			_ = conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))

			if err := WritePacket(conn, p); err != nil {
				go c.Reconnect()
				return
			}
		}
	}
}

func (c *Client) readLoop() {
	for {

		select {
		case <-c.closeCh:
			return
		default:
		}

		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()

		if conn == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		_ = conn.SetReadDeadline(time.Now().Add(c.readTimeout))

		packet, err := ReadPacket(conn)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			go c.Reconnect()
			return
		}

		// 当前仅处理 Response Packet。
		// Error Packet 将在 RPC V2 下一阶段实现。
		switch packet.Type {

			case PacketResponse:
				// 正常响应

			case PacketError:
				// RPC Error Packet

			default:
				continue
		}

		ch, ok := c.FindPending(packet.SeqID)
		if !ok {
			continue
		}

		select {
			case ch <- packet:
			default:
		}
	}
}