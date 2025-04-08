package gonet

import (
	"errors"
	"io"
	"net"

	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/waiter"
)

func (c *UDPConn) ReadTo(w io.Writer) (int, net.Addr, error) {
	deadline := c.readCancel()

	var addr tcpip.FullAddress
	n, err := commonReadTo(w, c.ep, c.wq, deadline, &addr, c)
	if err != nil {
		return 0, nil, err
	}
	return n, fullToUDPAddr(addr), nil
}

func commonReadTo(w io.Writer, ep tcpip.Endpoint, wq *waiter.Queue, deadline <-chan struct{}, addr *tcpip.FullAddress, errorer opErrorer) (int, error) {
	select {
	case <-deadline:
		return 0, errorer.newOpError("read", &timeoutError{})
	default:
	}

	opts := tcpip.ReadOptions{NeedRemoteAddr: addr != nil}
	res, err := ep.Read(w, opts)

	if _, ok := err.(*tcpip.ErrWouldBlock); ok {
		// Create wait queue entry that notifies a channel.
		waitEntry, notifyCh := waiter.NewChannelEntry(waiter.ReadableEvents)
		wq.EventRegister(&waitEntry)
		defer wq.EventUnregister(&waitEntry)
		for {
			res, err = ep.Read(w, opts)
			if _, ok := err.(*tcpip.ErrWouldBlock); !ok {
				break
			}
			select {
			case <-deadline:
				return 0, errorer.newOpError("read", &timeoutError{})
			case <-notifyCh:
			}
		}
	}

	if _, ok := err.(*tcpip.ErrClosedForReceive); ok {
		return 0, io.EOF
	}

	if err != nil {
		return 0, errorer.newOpError("read", errors.New(err.String()))
	}

	if addr != nil {
		*addr = res.RemoteAddr
	}
	return res.Count, nil
}
