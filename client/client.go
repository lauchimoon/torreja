package client

import (
    "bytes"
    "fmt"
    "net"
    "time"

    "github.com/lauchimoon/torreja/handshake"
    "github.com/lauchimoon/torreja/message"
    "github.com/lauchimoon/torreja/peers"
    bf "github.com/lauchimoon/torreja/bitfield"
)

type Client struct {
    Conn     net.Conn
    Choked   bool
    Bitfield bf.Bitfield
    peer     peers.Peer
    infoHash [20]byte
    peerId   string
}

func New(peer peers.Peer, peerId string, infoHash [20]byte) (*Client, error) {
    conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
    if err != nil {
        return nil, err
    }

    _, err = completeHandshake(conn, infoHash, peerId)
    if err != nil {
        return nil, err
    }

    bitField, err := receiveBitfield(conn)
    if err != nil {
        return nil, err
    }

    return &Client{
        Conn: conn,
        Choked: true,
        Bitfield: bitField,
        peer: peer,
        infoHash: infoHash,
        peerId: peerId,
    }, nil
}

func completeHandshake(conn net.Conn, infoHash [20]byte, peerId string) (*handshake.Handshake, error) {
    conn.SetDeadline(time.Now().Add(5*time.Second))
    defer conn.SetDeadline(time.Time{})

    hs := handshake.New(infoHash, peerId)
    _, err := conn.Write(hs.Serialize())
    if err != nil {
        return nil, err
    }

    res, err := handshake.Read(conn)
    if err != nil {
        return nil, err
    }
    if !bytes.Equal(res.InfoHash[:], infoHash[:]) {
        return nil, fmt.Errorf("expected infohash %x but got %x", res.InfoHash, infoHash)
    }
    return res, nil
}

func receiveBitfield(conn net.Conn) (bf.Bitfield, error) {
    conn.SetDeadline(time.Now().Add(5*time.Second))
    defer conn.SetDeadline(time.Time{})

    msg, err := message.Read(conn)
    if err != nil {
        return nil, err
    }
    if msg.Id != message.IdBitfield {
        return nil, fmt.Errorf("expected bitfield message, got id %d", msg.Id)
    }
    return msg.Payload, nil
}

func (c *Client) Read() (*message.Message, error) {
    return nil, nil
}

func (c *Client) SendUnchoked() {
}

func (c *Client) SendInterested() {
}

func (c *Client) SendHave(idx int) {
}

func (c *Client) SendRequest(idx, requestedBytes, blockSize int) error {
    return nil
}
