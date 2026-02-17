package main

import (
    "os"
    "github.com/lauchimoon/torreja/torrent"
)

func main() {
    torr, err := torrent.New(os.Args[1])
    if err != nil {
        panic(err)
    }

    err = torr.Download(os.Args[2])
    if err != nil {
        panic(err)
    }
}
