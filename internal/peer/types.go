package peer

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math"
	"net"
	"os"

	. "github.com/chrispritchard/gotorrent/internal/bitfields"
	"github.com/chrispritchard/gotorrent/internal/tracker"
)

var BLOCK_SIZE = 1 << 14

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

type PartialPiece struct {
	hash   string
	offset int
	blocks []bool
	data   []byte
}

func CreatePartialPiece(hash string, offset, full_length int) PartialPiece {
	return PartialPiece{
		hash:   hash,
		offset: offset,
		blocks: make([]bool, int(math.Ceil(float64(full_length)/float64(BLOCK_SIZE)))),
		data:   make([]byte, full_length),
	}
}

func (pp *PartialPiece) Length() int {
	return len(pp.blocks)
}

func (pp *PartialPiece) Set(offset int, data []byte) error {
	block_index := offset / BLOCK_SIZE
	if block_index < 0 || block_index >= len(pp.blocks) {
		return fmt.Errorf("invalid block index, out of range")
	}
	if len(data) > BLOCK_SIZE {
		return fmt.Errorf("data is too large for a single block")
	}
	// if pp.blocks[block_index] {
	// 	panic("existing block redownloaded")
	// }
	pp.blocks[block_index] = true
	target := pp.data[block_index*BLOCK_SIZE:]
	if len(target) < len(data) {
		return fmt.Errorf("data is too large for the target location") // should only be possible for the last block if truncated
	}
	copy(pp.data[block_index*BLOCK_SIZE:], data)
	return nil
}

func (pp *PartialPiece) Valid() bool {
	for _, b := range pp.blocks {
		if !b {
			return false
		}
	}
	hash := sha1.Sum(pp.data)
	return string(hash[:]) == pp.hash
}

func (pp *PartialPiece) WritePiece(file *os.File) error {
	if !pp.Valid() {
		return fmt.Errorf("piece is not valid")
	}
	_, err := file.Seek(int64(pp.offset), io.SeekStart)
	if err != nil {
		return err
	}
	_, err = file.Write(pp.data)
	return err
}
