package tracker

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/chrispritchard/gotorrent/internal/bencode"
	"github.com/chrispritchard/gotorrent/internal/torrent"
)

func escape(data []byte) string {
	return url.QueryEscape(string(data))
}

func CallTracker(metadata torrent.TorrentMetadata) (TrackerResponse, error) {
	id := make([]byte, 20)
	_, err := rand.Read(id)
	if err != nil {
		return nil_info, err
	}

	port := 6881

	keys := fmt.Sprintf("info_hash=%s&peer_id=%s&port=%d&uploaded=%d&downloaded=%d&left=%d&event=%s&compact=1",
		escape(metadata.InfoHash[:]), escape(id), port, 0, 0, metadata.Length, "started")

	url := metadata.Announcers[0] + "?" + keys
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil_info, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil_info, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil_info, err
	}

	return parse_tracker_response(body)
}

func parse_tracker_response(data []byte) (TrackerResponse, error) {
	decoded, _, err := bencode.Decode(data)
	if err != nil {
		return nil_info, err
	}
	root, ok := decoded.(map[string]any)
	if !ok {
		return nil_info, fmt.Errorf("invalid tracker response - not a bencode dict")
	}

	failure, err := bencode.Get[string](root, "failure reason")
	if err == nil {
		return nil_info, fmt.Errorf("tracker returned failure: %s", failure)
	}

	interval, err := bencode.Get[int](root, "interval")
	if err != nil {
		return nil_info, fmt.Errorf("invalid tracker response - missing interval")
	}

	peers_compact, err1 := bencode.Get[string](root, "peers")
	peers_full, err2 := bencode.Get[[]any](root, "peers")
	if err1 != nil && err2 != nil {
		return nil_info, fmt.Errorf("invalid tracker response - missing peers")
	}

	var peers []PeerInfo
	if err1 == nil {
		peers, err = parse_compact_peers(peers_compact)
		if err != nil {
			return nil_info, fmt.Errorf("invalid tracker response - invalid compact peer response: %v", err)
		}
	} else {
		peers, err = parse_full_peers(peers_full)
		if err != nil {
			return nil_info, fmt.Errorf("invalid tracker response - invalid peer response: %v", err)
		}
	}

	return TrackerResponse{
		Peers:    peers,
		Interval: interval,
	}, nil
}

func parse_full_peers(peers_full []any) ([]PeerInfo, error) {
	panic("unimplemented")
}

func parse_compact_peers(peers_compact string) ([]PeerInfo, error) {
	panic("unimplemented")
}
