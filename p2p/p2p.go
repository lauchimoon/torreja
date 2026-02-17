package p2p

import (
    "bytes"
    "crypto/sha1"
    "fmt"
    "log"
    "time"

    "github.com/lauchimoon/torreja/client"
    "github.com/lauchimoon/torreja/peers"
    "github.com/lauchimoon/torreja/message"
)

const MaxBlockSize = 16384
const MaxPipelined = 5

type Torrent struct {
    Peers       []peers.Peer
    PeerId      string
    InfoHash    [20]byte
    PieceHashes [][20]byte
    PieceLength int
    Length      int
    Name        string
}

type pieceWork struct {
    idx    int
    hash   [20]byte
    length int
}

type pieceResult struct {
    idx int
    buf []byte
}

type pieceProgress struct {
    idx        int
    client     *client.Client
    buf        []byte
    downloaded int
    requested  int
    pipelined  int
}

func (t *Torrent) Download() ([]byte, error) {
    workQueue := make(chan *pieceWork, len(t.PieceHashes))
    result := make(chan *pieceResult)
    for index, hash := range t.PieceHashes {
        length := t.calculatePieceSize(index)
        workQueue <- &pieceWork{index, hash, length}
    }

    for _, peer := range t.Peers {
        go t.startDownload(peer, workQueue, result)
    }

    buf := make([]byte, t.Length)
    donePieces := 0
    for donePieces < len(t.PieceHashes) {
        res := <- result
        begin, end := t.calculateBoundsForPiece(res.idx)
        copy(buf[begin:end], res.buf)
        donePieces++
    }
    close(workQueue)

    return buf, nil
}

func (t *Torrent) calculatePieceSize(idx int) int {
    begin, end := t.calculateBoundsForPiece(idx)
    return end - begin
}

func (t *Torrent) calculateBoundsForPiece(idx int) (int, int) {
    begin := idx*t.PieceLength
    end := begin + t.PieceLength
    if end > t.PieceLength {
        end = t.PieceLength
    }
    return begin, end
}

func (t *Torrent) startDownload(peer peers.Peer, workQueue chan *pieceWork, result chan *pieceResult) {
    c, err := client.New(peer, t.PeerId, t.InfoHash)
    if err != nil {
        log.Printf("could not perform handshake with IP %s.\n", peer.Ip)
        return
    }
    defer c.Conn.Close()
    log.Printf("connection with %s successful.\n", peer.Ip)

    c.SendUnchoked()
    c.SendInterested()

    for worker := range workQueue {
        if !c.Bitfield.HasPiece(worker.idx) {
            workQueue <- worker
            continue
        }

        buf, err := attemptDownload(c, worker)
        if err != nil {
            log.Println("exiting:", err)
            workQueue <- worker
            return
        }
        err = checkIntegrity(worker, buf)
        if err != nil {
            log.Printf("piece %d failed integrity check\n", worker.idx)
            workQueue <- worker
            continue
        }

        c.SendHave(worker.idx)
        result <- &pieceResult{worker.idx, buf}
    }
}

func attemptDownload(c *client.Client, worker *pieceWork) ([]byte, error) {
    state := pieceProgress{
        idx: worker.idx,
        client: c,
        buf: make([]byte, worker.length),
    }

    c.Conn.SetDeadline(time.Now().Add(30*time.Second))
    defer c.Conn.SetDeadline(time.Time{})

    for state.downloaded < worker.length {
        if !state.client.Choked {
            for state.pipelined < MaxPipelined && state.requested < worker.length {
                blockSize := MaxBlockSize
                if worker.length - state.requested < blockSize {
                    blockSize = worker.length - state.requested
                }

                err := c.SendRequest(worker.idx, state.requested, blockSize)
                if err != nil {
                    return nil, err
                }
                state.pipelined++
                state.requested += blockSize
            }
        }

        err := state.readMessage()
        if err != nil {
            return nil, err
        }
    }
    return state.buf, nil
}

func checkIntegrity(worker *pieceWork, buf []byte) error {
    hash := sha1.Sum(buf)
    if !bytes.Equal(hash[:], worker.hash[:]) {
        return fmt.Errorf("index %d failed integrity check", worker.idx)
    }
    return nil
}

func (p *pieceProgress) readMessage() error {
    msg, err := p.client.Read()
    if err != nil {
        return err
    }
    switch msg.Id {
    case message.IdChoke:
        p.client.Choked = true
    case message.IdUnchoke:
        p.client.Choked = false
    case message.IdHave:
        idx, err := message.ParseHave(msg)
        if err != nil {
            return err
        }
        p.client.Bitfield.SetPiece(idx)
    case message.IdPiece:
        n, err := message.ParsePiece(p.idx, p.buf, msg)
        if err != nil {
            return err
        }
        p.downloaded += n
        p.pipelined--
    }
    return nil
}
