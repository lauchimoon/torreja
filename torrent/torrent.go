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
    announce, ok := decoded["announce"]
    if !ok {
        return nil, errors.New("no 'announce' found.")
    }
    metainfo.Announce, ok = announce.(string)
    if !ok {
        return nil, errors.New("failed to parse 'announce' as string.")
    }

    // optional fields:
    // announce-list
    // creation date
    // comment
    // created by
    // encoding
    // If they're not found, there's no problem.
    metainfo.AnnounceList = getAnnounceList(decoded)
    getField(decoded, "creation date", &metainfo.CreationDate)
    getField(decoded, "comment", &metainfo.Comment)
    getField(decoded, "created by", &metainfo.CreatedBy)
    getField(decoded, "encoding", &metainfo.Encoding)

    return &metainfo, nil
}

func getField[T any](decoded map[string]any, field string, target *T) {
    if v, ok := decoded[field]; ok {
        if typedVal, ok := v.(T); ok {
            *target = typedVal
        }
    }
}

func getAnnounceList(decoded map[string]any) []string {
    announceList := []string{}
    list, ok := decoded["announce-list"].([]any)
    if !ok {
        return nil
    }
    for _, elem := range list {
        tracker := elem.([]any)[0].(string)
        announceList = append(announceList, tracker)
    }
    return announceList
}
