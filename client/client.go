package client

import (
    "bytes"
    "fmt"
    "net"
    "time"

    "github.com/lauchimoon/torreja/handshake"
    "github.com/lauchimoon/torreja/peers"
)

type Client struct {
    Conn     net.Conn
    Choked   bool
    peer     peers.Peer
    infoHash [20]byte
    peerId   string
}

func New(peer peers.Peer, peerId string, infoHash [20]byte) (*Client, error) {
    conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
    if err != nil {
        return nil, err
    }

    hs, err := completeHandshake(conn, infoHash, peerId)
    fmt.Println(hs)
    return &Client{
        Conn: conn,
        Choked: true,
    }, nil
}

func completeHandshake(conn net.Conn, infoHash [20]byte, peerId string) (*handshake.Handshake, error) {
    conn.SetDeadline(time.Now().Add(7*time.Second))
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
