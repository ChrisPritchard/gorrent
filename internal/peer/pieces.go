package peer

import (
	"encoding/binary"
	"net"
)

func RequestPiecePart(conn net.Conn, piece_index, start_byte, length uint) error {
	to_send := make([]byte, 12)
	binary.BigEndian.PutUint32(to_send[:4], uint32(piece_index))
	binary.BigEndian.PutUint32(to_send[4:8], uint32(start_byte))
	binary.BigEndian.PutUint32(to_send[8:], uint32(length))
	return send_message(conn, request, to_send)
}
