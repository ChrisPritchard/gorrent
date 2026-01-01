package main

import (
	"fmt"
	"os"
)

func main() {
	file := "c:\\users\\chris\\onedrive\\desktop\\test.torrent"
	d, _ := os.ReadFile(file)
	res, rem, err := decode_from_bencoded(d)
	fmt.Println(res)
	fmt.Println(rem)
	fmt.Println(err)

	torrent, err := parse_torrent_file(d)
	fmt.Println(torrent)
	fmt.Println(err)
}
