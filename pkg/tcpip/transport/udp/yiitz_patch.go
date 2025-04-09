package udp

func (r *ForwarderRequest) ReleaseReqPacket() {
	r.pkt.DecRef()
	r.pkt = nil
}

func (e *endpoint) SetNotifyReadFunc(notifyReadFunc func()) {
	e.notifyReadFunc = notifyReadFunc
}

type ISetNotifyReadFunc interface {
	SetNotifyReadFunc(notifyReadFunc func())
}
