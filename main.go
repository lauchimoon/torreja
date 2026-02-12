package main

import (
    "fmt"
    "os"
    "github.com/lauchimoon/torreja/torrent"
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
}
