package common

import "net"

type CustomIPAddr struct {
	IPAddr net.IPAddr
	Loc    string
	Note   string
}
