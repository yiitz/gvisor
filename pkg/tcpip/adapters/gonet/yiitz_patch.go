package gonet

import (
	"gvisor.dev/gvisor/pkg/tcpip"
)

func (c *UDPConn) GetEp() tcpip.Endpoint {
	return c.ep
}
