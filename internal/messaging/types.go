package messaging

import "encoding/binary"

type PeerMessageType int

const (
	MSG_CHOKE PeerMessageType = iota
	MSG_UNCHOKE
	MSG_INTERESTED
	MSG_NOTINTERESTED
	MSG_HAVE
	MSG_BITFIELD
	MSG_REQUEST
	MSG_PIECE
	MSG_CANCEL
)

type Received struct {
	Kind PeerMessageType
	Data []byte
}

func (r *Received) AsPiece() (int, int, []byte) {
	index := binary.BigEndian.Uint32(r.Data[0:4])
	begin := binary.BigEndian.Uint32(r.Data[4:8])
	piece := r.Data[8:]
	return int(index), int(begin), piece
}

var nil_received Received
