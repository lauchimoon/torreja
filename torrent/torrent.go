package torrent

import (
    "crypto/sha1"
    "errors"
    "os"
    "strings"

    "github.com/lauchimoon/torreja/bencode"
)

const (
    modeSingleFile = iota
    modeMultiFile
)

type file struct {
    Length int64
    MD5Sum string
    Path string
}

type hash [20]byte

type info struct {
    PieceLength int64
    Pieces []hash
    Private int64

    Name string
    Files []file
}

type Metainfo struct {
    Info info
    InfoHash hash
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

    iHash, err := getInfoHash(decoded)
    if err != nil {
        return nil, err
    }
    metainfo.InfoHash = iHash

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

    name, ok := data["name"]
    if !ok {
        return info{}, errors.New("failed to get name of file")
    }
    i.Name, ok = name.(string)
    if !ok {
        return info{}, errors.New("failed to parse name as string")
    }

    i.Files, err = getFiles(data, mode, i.Name)
    if err != nil {
        return info{}, err
    }

    return i, nil
}

func parsePieces(pieces any) ([]hash, error) {
    str, ok := pieces.(string)
    if !ok {
        return nil, errors.New("could not parse pieces as a string")
    }
    piecesSlice := []byte(str)
    piecesSize := len(piecesSlice)
    if piecesSize % 20 != 0 {
        return nil, errors.New("size of pieces string must be a multiple of 20")
    }

    hashes := []hash{}
    for i := 0; i < piecesSize; i += 20 {
        var chunk [20]byte
        for j := 0; j < 20; j++ {
            chunk[j] = piecesSlice[i+j]
        }
        hashes = append(hashes, chunk)
    }
    return hashes, nil
}

func getFiles(data map[string]any, mode int, name string) ([]file, error) {
    if mode == modeSingleFile {
        f, err := getSingleFile(data, name)
        return f, err
    }

    files, err := getMultiFile(data)
    return files, err
}

func getSingleFile(data map[string]any, name string) ([]file, error) {
    f := file{}
    lengthRaw, ok := data["length"]
    if !ok {
        return nil, errors.New("failed to find file length")
    }
    length, ok := lengthRaw.(int64)
    if !ok {
        return nil, errors.New("failed to parse file length as int64")
    }
    getField(data, "md5sum", &f.MD5Sum)
    f.Length = length
    f.Path = name
    return []file{f}, nil
}

func getMultiFile(data map[string]any) ([]file, error) {
    files := []file{}

    // We already made sure it exists in getInfo
    filesRaw, ok := data["files"].([]any)
    if !ok {
        return nil, errors.New("failed to parse as list of files")
    }

    for _, elem := range filesRaw {
        fRaw, ok := elem.(map[string]any)
        if !ok {
            return nil, errors.New("failed to parse file")
        }
        lengthRaw, ok := fRaw["length"]
        if !ok {
            return nil, errors.New("failed to find file length")
        }
        length, ok := lengthRaw.(int64)
        if !ok {
            return nil, errors.New("failed to parse file length as int64")
        }

        pathListRaw, ok := fRaw["path"]
        if !ok {
            return nil, errors.New("failed to find file path")
        }
        parts, ok := pathListRaw.([]any)
        if !ok {
            return nil, errors.New("failed to parse file path as []string")
        }

        pathList := []string{}
        for _, part := range parts {
            p, ok := part.(string)
            if !ok {
                return nil, errors.New("failed to parse part of file path as string")
            }
            pathList = append(pathList, p)
        }

        path := strings.Join(pathList, "/")
        f := file{}
        f.Length = length
        f.Path = path
        getField(fRaw, "md5sum", &f.MD5Sum)
        files = append(files, f)
    }

    return files, nil
}

func getInfoHash(data map[string]any) (hash, error) {
    buf := bencode.Encode(data["info"])
    h := sha1.Sum([]byte(buf))
    return h, nil
}
