package main

import (
	"fmt"
	"os"
)

func main() {
	file := "c:\\users\\chris\\onedrive\\desktop\\test.torrent"
	d, _ := os.ReadFile(file)
	res, rem, err := parse(d)
	fmt.Println(res)
	fmt.Println(rem)
	fmt.Println(err)
}
