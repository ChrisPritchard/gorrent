package peer

import (
	"net"

	"github.com/chrispritchard/gotorrent/internal/tracker"
)

type PeerDetails struct {
	tracker.PeerInfo

	Conn net.Conn

	Chocked    bool
	Interested bool

	Has      BitField
	Requests map[int]struct{}
}

func NewPeerCommunication(p tracker.PeerInfo, conn net.Conn) PeerDetails {
	return PeerDetails{
		PeerInfo:   p,
		Conn:       conn,
		Chocked:    true,
		Interested: false,
		Has:        BitField{},
		Requests:   make(map[int]struct{}),
	}
}

type PeerManager struct {
	Have BitField

	Peers      []PeerDetails
	Requesting map[int]struct{}
}
