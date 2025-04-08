package udp

func (r *ForwarderRequest) ReleaseReqPacket() {
	r.pkt.DecRef()
	r.pkt = nil
}
