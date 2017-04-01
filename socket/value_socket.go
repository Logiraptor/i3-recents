package socket

import (
	"encoding/json"
	"io"
)

type ValueSocket struct {
	sock *Socket
}

func NewValueSocket(conn io.ReadWriter) *ValueSocket {
	return &ValueSocket{sock: NewSocket(conn)}
}

func (v *ValueSocket) Write(messageType uint32, val interface{}) error {
	buf, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return v.sock.Write(Message{
		Type:    messageType,
		Payload: string(buf),
	})
}

func (v *ValueSocket) Read(val interface{}) (uint32, error) {
	message, err := v.sock.Read()
	if err != nil {
		return 0, err
	}
	return message.Type, json.Unmarshal([]byte(message.Payload), val)
}
