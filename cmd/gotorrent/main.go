package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chrispritchard/gotorrent/internal/bitfields"
	"github.com/chrispritchard/gotorrent/internal/messaging"
	"github.com/chrispritchard/gotorrent/internal/peer"
	"github.com/chrispritchard/gotorrent/internal/torrent"
	"github.com/chrispritchard/gotorrent/internal/tracker"
	"github.com/chrispritchard/gotorrent/internal/util"
)

func main() {
	file := "c:\\users\\chris\\onedrive\\desktop\\test.torrent"
	if _, err := os.Stat("ScreenToGif.exe"); err == nil {
		os.Remove("ScreenToGif.exe")
	}

	err := try_download(file)
	if err != nil {
		fmt.Printf("unable to download via torrent file: %v", err)
		os.Exit(1)
	}
}

func try_download(torrent_file_path string) error {

	metadata, err := parse_torrent(torrent_file_path)
	tracker_info, err := tracker.CallTracker(metadata)
	if err != nil {
		return fmt.Errorf("failed to register with tracker: %v", err)
	}

	local_field, err := get_local_bit_field(metadata)
	if err != nil {
		return fmt.Errorf("failed to register with tracker: %v", err)
	}

	out_file, err := establish_outfile(metadata)
	if err != nil {
		return fmt.Errorf("failed to create/read out_file: %v", err)
	}
	defer out_file.Close()

	return start_state_machine(metadata, tracker_info, local_field, out_file)
}

func parse_torrent(torrent_file_path string) (torrent.TorrentMetadata, error) {
	var nil_result torrent.TorrentMetadata
	d, err := os.ReadFile(torrent_file_path)
	if err != nil {
		return nil_result, fmt.Errorf("unable to read file at path %s: %v", torrent_file_path, err)
	}

	torrent, err := torrent.ParseTorrentFile(d)
	if err != nil {
		return nil_result, fmt.Errorf("unable to parse torrent file: %v", err)
	}

	return torrent, nil
}

func get_local_bit_field(metadata torrent.TorrentMetadata) (bitfields.BitField, error) {
	return bitfields.CreateBlankBitfield(len(metadata.Pieces)), nil // TODO: evaluate existing file
}

func establish_outfile(metadata torrent.TorrentMetadata) (*os.File, error) {
	out_file, err := os.Create(metadata.Name) // assuming a single file with no directory info
	if err != nil {
		return nil, err
	}

	err = out_file.Truncate(int64(metadata.Length)) // create full size file
	if err != nil {
		out_file.Close()
		return nil, err
	}

	return out_file, nil
}

func start_state_machine(metadata torrent.TorrentMetadata, tracker_info tracker.TrackerResponse, local_field bitfields.BitField, out_file *os.File) error {
	ctx := context.Background()
	defer ctx.Done()

	received_channel := make(chan messaging.Received)
	error_channel := make(chan error)

	peers := connect_to_peers(metadata, tracker_info, local_field)
	if len(peers) == 0 {
		return fmt.Errorf("failed to connect to a peer")
	}

	for _, p := range peers {
		p.StartReceiving(ctx, received_channel, error_channel)
		defer p.Close()
	}

	// start requesting pieces

	mu := sync.Mutex{}
	requested := map[int]map[int]struct{}{}
	partials := peer.CreatePartialPieces(metadata.Pieces, metadata.PieceLength, metadata.Length)
	pipeline := make(chan int, 5) // concurrent requests
	for range 5 {
		pipeline <- 1
	}
	go start_requesting_pieces(ctx, peers, partials, pipeline, &mu, requested)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	finished_pieces := 0
	for {
		select {
		case <-ticker.C:
			print_status(partials, requested)
		case received := <-received_channel:
			piece_finished := handle_received(received, &mu, requested, partials, out_file)
			if piece_finished {
				finished_pieces++
				if finished_pieces == len(metadata.Pieces) {
					fmt.Println("done")
					return nil
				}
			}
			pipeline <- 1
		case err := <-error_channel:
			return err
		}
	}
}

func handle_received(received messaging.Received, mu *sync.Mutex, requested map[int]map[int]struct{}, partials []peer.PartialPiece, out_file *os.File) (piece_finished bool) {
	piece_finished = false
	if received.Kind == messaging.MSG_PIECE {
		index := binary.BigEndian.Uint32(received.Data[0:4])
		begin := binary.BigEndian.Uint32(received.Data[4:8])
		piece := received.Data[8:]

		mu.Lock()
		p := requested[int(index)]
		delete(p, int(begin))
		if len(p) == 0 {
			delete(requested, int(index))
		}
		mu.Unlock()

		partials[index].Set(int(begin), piece)
		fmt.Printf("piece %d block received\n", index)
		if partials[index].Valid() {
			partials[index].WritePiece(out_file)
			piece_finished = true
			fmt.Printf("piece %d finished\n", index)
		}
	} else {
		fmt.Printf("received an unhandled kind: %d\n", received.Kind)
	}
	return
}

func print_status(partials []peer.PartialPiece, requested map[int]map[int]struct{}) {
	for i, p := range partials {
		if !p.Valid() {
			fmt.Printf("partial %d is invalid\n", i)
			fmt.Printf("\tmissing: %v\n", p.Missing())
		}
	}
	fmt.Printf("requested:\n")
	for k, v := range requested {
		var indices strings.Builder
		for k2 := range v {
			indices.WriteString(strconv.Itoa(k2) + " ")
		}
		fmt.Printf("\t%d: %s\n", k, indices.String())
	}
	fmt.Println()
}

func start_requesting_pieces(ctx context.Context, peers []peer.PeerHandler, partials []peer.PartialPiece, pipeline <-chan int, mu *sync.Mutex, requested map[int]map[int]struct{}) error {
	for i, p := range partials {
		for j := range p.Length() {
			select {
			case <-ctx.Done():
				return nil
			case <-pipeline:
				err := peers[0].RequestPieceBlock(i, j*peer.BLOCK_SIZE, p.BlockSize(j))
				mu.Lock()
				if e, ok := requested[i]; ok {
					e[j*peer.BLOCK_SIZE] = struct{}{}
				} else {
					requested[i] = map[int]struct{}{}
					requested[i][j*peer.BLOCK_SIZE] = struct{}{}
				}
				mu.Unlock()
				fmt.Printf("requested block %d of piece %d\n", j, i)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func connect_to_peers(metadata torrent.TorrentMetadata, tracker_response tracker.TrackerResponse, local_bitfield bitfields.BitField) []peer.PeerHandler {
	ops := make([]util.Op[peer.PeerHandler], len(tracker_response.Peers))
	for i, p := range tracker_response.Peers {
		local_p := p
		ops[i] = func() (peer.PeerHandler, error) {
			return peer.ConnectToPeer(local_p, metadata.InfoHash[:], tracker_response.LocalID, local_bitfield)
		}
	}

	conns, _ := util.Concurrent(ops, 20)
	return conns
}
