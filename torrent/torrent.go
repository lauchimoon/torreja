package torrent

import (
    "errors"
//    "fmt"
    "os"
    "github.com/lauchimoon/torreja/bencode"
)

const (
    modeSingleFile = iota
    modeMultiFile
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

    data, err := getInfo(decoded)
    if err != nil {
        return nil, err
    }
    metainfo.Info = data

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

func getInfo(decoded map[string]any) (info, error) {
    i := info{}
    dataRaw, ok := decoded["info"]
    if !ok {
        return info{}, errors.New("no 'info' dictionary found")
    }
    data, ok := dataRaw.(map[string]any)
    if !ok {
        return info{}, errors.New("failed to parse 'info' dictionary")
    }

    name, ok := data["name"]
    if !ok {
        return info{}, errors.New("failed to get name of file/output directory")
    }
    i.Name, ok = name.(string)
    if !ok {
        return info{}, errors.New("failed to parse name as string")
    }

    pieceLength, ok := data["piece length"]
    if !ok {
        return info{}, errors.New("failed to get piece length")
    }
    i.PieceLength, ok = pieceLength.(int64)
    if !ok {
        return info{}, errors.New("failed to parse piece length as int64")
    }

    getField(data, "private", &i.Private)
    pieces, ok := data["pieces"]
    if !ok {
        return info{}, errors.New("failed to get pieces")
    }
    hashes, err := parsePieces(pieces)
    if err != nil {
        return info{}, err
    }
    i.Pieces = hashes

    mode := modeMultiFile
    _, ok = data["files"]
    if !ok {
        mode = modeSingleFile
    }

    i.Files, err = getFiles(data, mode)
    return i, nil
}

func parsePieces(pieces any) ([]hash, error) {
    return nil, nil
}

func getFiles(data map[string]any, mode int) ([]map[string]file, error) {
    return nil, nil
}
