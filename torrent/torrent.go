package torrent

import (
    "errors"
    "os"
    "github.com/lauchimoon/torreja/bencode"
)

type file struct {
    Length int64
    MD5Sum string
    Path []string
}

type hash [20]byte

type info struct {
    PieceLength int64
    Pieces []hash
    Private int64

    Name string
    Files []map[string]file
}

type Metainfo struct {
    Info info
    Announce string
    AnnounceList []string
    CreationDate int64
    Comment string
    CreatedBy string
    Encoding string
}

func New(torrentFilePath string) (*Metainfo, error) {
    fileContents, err := os.ReadFile(torrentFilePath)
    if err != nil {
        return nil, err
    }
    decoded, err := bencode.Decode(string(fileContents))
    if err != nil {
        return nil, err
    }

    metainfo := Metainfo{}
    announce, ok := decoded["announce"].(string)
    if !ok {
        return nil, errors.New("no 'announce' found")
    }
    metainfo.Announce = announce

    // optional fields:
    // announce-list
    // creation date
    // comment
    // created by
    // encoding
    // If they're not found, there's no problem.
    announceList := getAnnounceList(decoded)
    metainfo.AnnounceList = announceList

    creationDate, ok := decoded["creation date"]
    if ok {
        metainfo.CreationDate = creationDate.(int64)
    }
    comment, ok := decoded["comment"]
    if ok {
        metainfo.Comment = comment.(string)
    }
    createdBy, ok := decoded["created by"]
    if ok {
        metainfo.CreatedBy = createdBy.(string)
    }
    encoding, ok := decoded["encoding"]
    if ok {
        metainfo.Encoding = encoding.(string)
    }

    return &metainfo, nil
}

func getAnnounceList(decoded map[string]interface{}) []string {
    announceList := []string{}
    list, ok := decoded["announce-list"].([]interface{})
    if !ok {
        return nil
    }
    for _, elem := range list {
        tracker := elem.([]interface{})[0].(string)
        announceList = append(announceList, tracker)
    }
    return announceList
}
