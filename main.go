package main

import (
    "fmt"
    "github.com/lauchimoon/torreja/bencode"
)

func main() {
    var l []interface{}
    l = append(l, "a")
    l = append(l, "b")
    d := make(map[string]interface{})
    d["cow"] = "moo"
    d["spam"] = "eggs"
    d["spammy"] = l

    bc := bencode.Encode(d)
    fmt.Println(bc)
    fmt.Println(bencode.Decode(bc))
}
