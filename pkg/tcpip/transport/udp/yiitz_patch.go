package udp

func (r *ForwarderRequest) ReleaseReqPacket() {
	r.pkt.DecRef()
	r.pkt = nil
}

func (e *endpoint) SetReadChan(readChan chan<- any, readChanMsg any) {
	e.readChan = readChan
	e.readChanMsg = readChanMsg
}

type SetReadChaner interface {
	SetReadChan(readChan chan<- any, readChanMsg any)
}
