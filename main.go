package main

import (
    "fmt"
    "os"
    "github.com/lauchimoon/torreja/torrent"
    "github.com/lauchimoon/torreja/client"
)

func main() {
    torr, err := torrent.New(os.Args[1])
    if err != nil {
        panic(err)
    }

    fmt.Println("Announce:", torr.Announce)
    fmt.Println("Announce list:", torr.AnnounceList)
    fmt.Println("Creation date (UNIX timestamp):", torr.CreationDate)
    fmt.Printf("Comment: '%s'\n", torr.Comment)
    fmt.Printf("Created by: '%s'\n", torr.CreatedBy)
    fmt.Printf("Encoding: '%s'\n", torr.Encoding)

    fmt.Printf("\nName: '%s'\n", torr.Info.Name)
    fmt.Println("Piece length:", torr.Info.PieceLength)
    fmt.Println("Private:", torr.Info.Private)
    fmt.Println("Files:", torr.Info.Files)

    peerList, err := torr.RequestPeers("something1something1", 6881)
    if err != nil {
        panic(err)
    }
    fmt.Println("Peers found:", peerList)

    i := 0
    fmt.Printf("Inspect peer %s\n", peerList[i])
    c, err := client.New(peerList[i], "something1something1", torr.InfoHash)
    for err != nil {
        i++
        fmt.Printf("Inspect peer %s\n", peerList[i])
        c, err = client.New(peerList[i], "something1something1", torr.InfoHash)
    }
    fmt.Println(peerList[i], c)
}
