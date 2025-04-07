package gonet

import (
	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/waiter"
)

type UDPReadChan struct {
	Ep        tcpip.Endpoint
	WaitEntry waiter.Entry
	NotifyCh  chan struct{}
	DealineCh <-chan struct{}
}

func (c *UDPConn) BeginReadChan() *UDPReadChan {
	waitEntry, notifyCh := waiter.NewChannelEntry(waiter.ReadableEvents)
	c.wq.EventRegister(&waitEntry)
	return &UDPReadChan{
		WaitEntry: waitEntry,
		NotifyCh:  notifyCh,
		DealineCh: c.readCancel(),
	}
}

func (c *UDPConn) EndReadChan(rc *UDPReadChan) {
	c.wq.EventUnregister(&rc.WaitEntry)
}

func (c *UDPConn) EP() tcpip.Endpoint {
	return c.ep
}
