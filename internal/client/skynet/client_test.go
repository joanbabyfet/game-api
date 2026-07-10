package skynet

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func TestEncodeDecode(t *testing.T) {
	src := &Packet{
		Type:    PacketResponse,
		CmdID:   CmdBalance,
		SeqID:   100,
		Payload: []byte("hello"),
	}

	data, err := Encode(src)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	dst, err := Decode(data)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if dst.Type != src.Type {
		t.Fatal("packet type mismatch")
	}

	if dst.CmdID != src.CmdID {
		t.Fatalf("message id mismatch")
	}

	if dst.SeqID != src.SeqID {
		t.Fatalf("sequence mismatch")
	}

	if !bytes.Equal(dst.Payload, src.Payload) {
		t.Fatalf("payload mismatch")
	}
}

func TestReadWritePacket(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	want := &Packet{
		Type:    PacketResponse,
		CmdID:   CmdBet,
		SeqID:   88,
		Payload: []byte("protobuf-bytes"),
	}

	go func() {
		if err := WritePacket(server, want); err != nil {
			t.Errorf("write failed: %v", err)
		}
	}()

	got, err := ReadPacket(client)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if got.Type != want.Type {
		t.Fatal("packet type mismatch")
	}

	if got.CmdID != want.CmdID {
		t.Fatalf("message id mismatch")
	}

	if got.SeqID != want.SeqID {
		t.Fatalf("sequence mismatch")
	}

	if !bytes.Equal(got.Payload, want.Payload) {
		t.Fatalf("payload mismatch")
	}
}

func TestPending(t *testing.T) {
	c := New("127.0.0.1:9000")

	seq := c.NextSequence()

	ch := c.RegisterPending(seq)

	go func() {
		ch <- &Packet{
			Type:  PacketResponse,
			CmdID: MsgOK,
			SeqID: seq,
		}
	}()

	select {
	case p := <-ch:
		if p.Type != PacketResponse {
			t.Fatal("packet type mismatch")
		}

		if p.SeqID != seq {
			t.Fatalf("unexpected sequence")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}

	c.RemovePending(seq)
}

func TestClose(t *testing.T) {
	c := New("127.0.0.1:9000")

	if c.IsClosed() {
		t.Fatal("client should be open")
	}

	if err := c.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	if !c.IsClosed() {
		t.Fatal("client should be closed")
	}
}

func TestSequence(t *testing.T) {
	c := New("127.0.0.1:9000")

	a := c.NextSequence()
	b := c.NextSequence()

	if b != a+1 {
		t.Fatalf("sequence should increase")
	}
}
