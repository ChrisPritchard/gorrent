package main

import (
	"fmt"
	"os"
)

func main() {
	file := "c:\\users\\chris\\onedrive\\desktop\\test.torrent"

	err := try_download(file)
	if err != nil {
		fmt.Printf("unable to download from torrent file: %v", err)
		os.Exit(1)
	}
}

func try_download(torrent_file_path string) error {
	d, err := os.ReadFile(torrent_file_path)
	if err != nil {
		return fmt.Errorf("unable to read file at path %s: %v", torrent_file_path, err)
	}

	torrent, err := parse_torrent_file(d)
	if err != nil {
		return fmt.Errorf("unable to parse torrent file: %v", err)
	}

	peer_info, err := call_tracker(torrent)
	if err != nil {
		return fmt.Errorf("failed to register with tracker: %v", err)
	}

	fmt.Println(peer_info)
	// handshake with peers
	// download and write file

	return nil
}
