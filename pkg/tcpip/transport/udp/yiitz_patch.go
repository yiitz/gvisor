package udp

func (r *ForwarderRequest) Complete() {
	if r.pkt == nil {
		return
	}
	r.pkt.DecRef()
	r.pkt = nil
}

func (e *endpoint) SetNotifyReadFunc(notifyReadFunc func()) {
	e.notifyReadFunc = notifyReadFunc
}

type ISetNotifyReadFunc interface {
	SetNotifyReadFunc(notifyReadFunc func())
}
