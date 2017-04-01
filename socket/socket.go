package socket

import (
	"encoding/binary"
	"io"
)

const MagicString = "i3-ipc"

var magicStringBytes [6]byte

func init() {
	copy(magicStringBytes[:], MagicString)
}

type Header struct {
	MagicString   [6]byte
	PayloadLength uint32
	MessageType   uint32
}

type Message struct {
	Type    uint32
	Payload string
}

type Socket struct {
	conn io.ReadWriter
}

func NewSocket(conn io.ReadWriter) *Socket {
	return &Socket{
		conn: conn,
	}
}

type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) Write(buf []byte) (int, error) {
	if ew.err != nil {
		return 0, ew.err
	}
	var n int
	n, ew.err = ew.w.Write(buf)
	return n, ew.err
}

func (s *Socket) Write(message Message) error {
	wr := &errWriter{w: s.conn}
	binary.Write(wr, nativeOrder, Header{
		MagicString:   magicStringBytes,
		MessageType:   message.Type,
		PayloadLength: uint32(len(message.Payload)),
	})
	io.WriteString(wr, message.Payload)
	return wr.err
}

func (s *Socket) Read() (Message, error) {
	var header Header
	err := binary.Read(s.conn, nativeOrder, &header)
	if err != nil {
		return Message{}, err
	}
	buf := make([]byte, header.PayloadLength)
	_, err = io.ReadFull(s.conn, buf)
	if err != nil {
		return Message{}, err
	}
	return Message{
		Type:    header.MessageType,
		Payload: string(buf),
	}, nil
}
