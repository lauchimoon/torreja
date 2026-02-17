package handshake

import (
    "errors"
    "io"
)

type Handshake struct {
    Pstr     string
    InfoHash [20]byte
    PeerId   string
}

func New(infoHash [20]byte, peerId string) *Handshake {
    return &Handshake{
        Pstr: "BitTorrent protocol",
        InfoHash: infoHash,
        PeerId: peerId,
    }
}

func (hs *Handshake) Serialize() []byte {
    const pstrByte = 1
    const numExtensions = 8
    const infoHashLen = 20
    peerIdLen := len(hs.PeerId)
    pstrLen := len(hs.Pstr)
    total := pstrByte + pstrLen + numExtensions + infoHashLen + peerIdLen
    buf := make([]byte, total)

    buf[0] = byte(pstrLen)
    curr := 1
    curr += copy(buf[curr:], hs.Pstr)
    curr += copy(buf[curr:], make([]byte, 8))
    curr += copy(buf[curr:], hs.InfoHash[:])
    curr += copy(buf[curr:], hs.PeerId)

    return buf
}

func Read(r io.Reader) (*Handshake, error) {
    lengthBuf := make([]byte, 1)
    _, err := io.ReadFull(r, lengthBuf)
    if err != nil {
        return nil, err
    }

    pStrLen := int(lengthBuf[0])
    if pStrLen == 0 {
        return nil, errors.New("length of protocol string cannot be 0")
    }

    // 48: numExtensions + infoHashLen + peerIdLen,
    // which is 8 + 20 + 20
    handshakeBuf := make([]byte, pStrLen+48)
    _, err = io.ReadFull(r, handshakeBuf)
    if err != nil {
        return nil, err
    }

    var infoHash, peerId [20]byte

    copy(infoHash[:], handshakeBuf[pStrLen+8:pStrLen+8+20])
    copy(peerId[:], handshakeBuf[pStrLen+8+20:])
    return &Handshake{
        Pstr: string(handshakeBuf[0:pStrLen]),
        InfoHash: infoHash,
        PeerId: string(peerId[:]),
    }, nil
}
