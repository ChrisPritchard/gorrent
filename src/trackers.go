package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type PeerInfo struct {
	Id   string
	IP   int
	Port int
}

type TrackerResponse struct {
	Peers    []PeerInfo
	Interval int
}

func escape(data []byte) string {
	return url.QueryEscape(string(data))
}

func call_tracker(metadata TorrentMetadata) (TrackerResponse, error) {
	var nil_info TrackerResponse

	id := make([]byte, 20)
	_, err := rand.Read(id)
	if err != nil {
		return nil_info, err
	}

	port := 6881

	keys := fmt.Sprintf("info_hash=%s&peer_id=%s&port=%d&uploaded=%d&downloaded=%d&left=%d&event=%s",
		escape(metadata.InfoHash[:]), escape(id), port, 0, 0, metadata.Length, "started")

	url := metadata.Announcers[0] + "?" + keys
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, _ := client.Do(req)

	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	fmt.Println(s)

	return nil_info, nil
}
