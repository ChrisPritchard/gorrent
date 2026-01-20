package messaging

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
