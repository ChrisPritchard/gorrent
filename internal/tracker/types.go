package tracker

type PeerInfo struct {
	Id   string
	IP   int
	Port int
}

type TrackerResponse struct {
	Peers    []PeerInfo
	Interval int
}

var nil_info TrackerResponse
