package message

import (
    "encoding/binary"
    "errors"
    "fmt"
    "io"
)

const (
    IdChoke = iota
    IdUnchoke
    IdInterested
    IdNotInterested
    IdHave
    IdBitfield
    IdRequest
    IdPiece
    IdCancel
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
    messageLen := binary.BigEndian.Uint32(lengthBuf)
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

func ParseHave(m *Message) (int, error) {
    if m.Id != IdHave {
        return 0, fmt.Errorf("expected have (id %d), got %d", IdHave, m.Id)
    }
    if len(m.Payload) != 4 {
        return 0, fmt.Errorf("expected payload of length 4, got length %d", len(m.Payload))
    }
    idx := int(binary.BigEndian.Uint32(m.Payload))
    return idx, nil
}

func ParsePiece(idx int, buf []byte, m *Message) (int64, error) {
    if m.Id != IdPiece {
        return 0, fmt.Errorf("expected piece (id %d), got %d", IdPiece, m.Id)
    }
    if len(m.Payload) < 8 {
        return 0, fmt.Errorf("payload is too short")
    }
    parsedIdx := int(binary.BigEndian.Uint32(m.Payload[0:4]))
    if parsedIdx != idx {
        return 0, fmt.Errorf("expected index %v, got %v", idx, parsedIdx)
    }
    begin := int(binary.BigEndian.Uint32(m.Payload[4:8]))
    if begin >= len(buf) {
        return 0, fmt.Errorf("begin offset is too high")
    }
    data := m.Payload[8:]
    if begin+len(data) > len(buf) {
        return 0, fmt.Errorf("data is too long")
    }
    copy(buf[begin:], data)
    return int64(len(data)), nil
}

func FormatHave(idx int) *Message {
    buf := make([]byte, 4)
    binary.BigEndian.PutUint32(buf, uint32(idx))
    return &Message{
        Id: IdHave,
        Payload: buf,
    }
}

func FormatRequest(idx int, requestedBytes, blockSize int64) *Message {
    buf := make([]byte, 12)
    binary.BigEndian.PutUint32(buf[0:4], uint32(idx))
    binary.BigEndian.PutUint32(buf[4:8], uint32(requestedBytes))
    binary.BigEndian.PutUint32(buf[8:12], uint32(blockSize))
    return &Message{
        Id: IdRequest,
        Payload: buf,
    }
}
