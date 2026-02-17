package torrent

import (
    "errors"
    "io"
    "net"
    "net/http"
    "net/url"
    "strconv"
    "time"

    "github.com/lauchimoon/torreja/bencode"
    "github.com/lauchimoon/torreja/peers"
)

func (m *Metainfo) buildTrackerURL(peerId string, port int64) (string, error) {
    base, err := url.Parse(m.Announce)
    if err != nil {
        return "", err
    }
    params := url.Values{
        "info_hash": []string{string(m.InfoHash[:])},
        "peer_id": []string{peerId},
        "port": []string{strconv.FormatInt(port, 10)},
        "uploaded": []string{"0"},
        "downloaded": []string{"0"},
        "left": []string{strconv.FormatInt(m.getTotalLength(), 10)},
    }
    base.RawQuery = params.Encode()
    return base.String(), nil
}

func (m *Metainfo) getTotalLength() int64 {
    var length int64
    for _, f := range m.Info.Files {
        length += f.Length
    }
    return length
}

func (m *Metainfo) RequestPeers(peerId string, port int64) ([]peers.Peer, error) {
    url, err := m.buildTrackerURL(peerId, port)
    if err != nil {
        return nil, err
    }

    client := &http.Client{Timeout: 15*time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    decoded, err := bencode.Decode(string(bodyBytes))
    if err != nil {
        return nil, err
    }
    return parsePeers(decoded)
}

func parsePeers(decoded map[string]any) ([]peers.Peer, error) {
    peersRaw, ok := decoded["peers"]
    if !ok {
        return nil, errors.New("failed to find peers to connect to")
    }
    peersList, ok := peersRaw.([]any)
    if !ok {
        return nil, errors.New("failed to parse peers as list")
    }

    list := []peers.Peer{}
    for _, peerRaw := range peersList {
        peer, ok := peerRaw.(map[string]any)
        if !ok {
            return nil, errors.New("failed to parse peer as dictionary")
        }
        list = append(list, peers.Peer{
            Ip: net.ParseIP(peer["ip"].(string)),
            Port: peer["port"].(int64),
        })
    }

    return list, nil
}
