package torrent

import (
    "net/url"
    "strconv"
)

func (m *Metainfo) BuildTrackerURL(peerId, port string) (string, error) {
    base, err := url.Parse(m.Announce)
    if err != nil {
        return "", err
    }
    params := url.Values{
        "info_hash": []string{string(m.InfoHash)},
        "peer_id": []string{peerId},
        "port": []string{port},
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
