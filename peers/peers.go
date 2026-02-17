package peers

import (
    "net"
    "fmt"
)

type Peer struct {
    Ip   net.IP
    Port int64
}

func (p Peer) String() string {
    return fmt.Sprintf("%s:%v", p.Ip.String(), p.Port)
}
