package message

import (
    "encoding/binary"
    "errors"
    "io"
)

const (
    MessageIdChoke = iota
    MessageIdUnchoke
    MessageIdInterested
    MessageIdNotInterested
    MessageIdHave
    MessageIdBitfield
    MessageIdRequest
    MessageIdPiece
    MessageIdCancel
)

type Message struct {
    Id      int
    Payload []byte
}

func (m *Message) Serialize() []byte {
    if m == nil {
        return make([]byte, 4)
    }
    length := uint32(len(m.Payload) + 1)
    buf := make([]byte, 4+length)
    binary.BigEndian.PutUint32(buf[0:4], length)
    buf[4] = byte(m.Id)
    copy(buf[5:], m.Payload)
    return buf
}

func Read(r io.Reader) (*Message, error) {
    lengthBuf := make([]byte, 4)
    _, err := io.ReadFull(r, lengthBuf)
    if err != nil {
        return nil, err
    }

    messageLen := binary.BigEndian.Uint32(lengthBuf[0:4])
    if messageLen == 0 {
        return nil, errors.New("length of message cannot be 0")
    }

    messageBuf := make([]byte, messageLen)
    _, err = io.ReadFull(r, messageBuf)
    if err != nil {
        return nil, err
    }

    return &Message{
        Id: int(messageBuf[0]),
        Payload: messageBuf[1:],
    }, nil
}
