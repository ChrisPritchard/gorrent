# Gorrent

Small Golang BitTorrent client.

> Gorrent: **Go** To**rrent**: too many other 'gotorrent' projects out there :D

Idea #1 from <https://codecrafters.io/blog/programming-project-ideas>, and built against the BitTorrent specification: <https://www.bittorrent.org/beps/bep_0003.html>

> Presently only downloading is supported - if a torrent is loaded for a file that is fully completed, the application will panic as 'seeding' isn't supported (yet?)

> Only torrent files are supported, only over tcp and unencrypted (e.g. no magnet links, no utorrent protocol, no TLS)

```
Usage: gorrent [options] <torrent-file>
  -v    enable verbose output
exit status 1
```

## Components

- gorrent/main.go: gets a torrent file from the arguments, parses it, creates or reads local files, then initiates a parallel process of requesting pieces and receiving them from peersfrom the tracker
- bencode: contains methods to parse the bencoded torrent file, and bencoded responses
- bitfields: contains a type used to represent available pieces of a torrent - this is a long bit array where each positive bit represents a held piece. these fields are exchanged with peers
- downloading: a manager of local files, local bit fields and remote peers that makes requests for pieces, cancels requests, and receives requests for writing to the local files
- messaging: helper methods for the inter-peer communication structure, including message types and tcp conn management
- out_files: a manager for local files: abstracts single vs multi-file torrent structures away from the communication primitives (which are just pieces and offsets). writes received data to the correct files at the correct locations, and also maintains the local bitfield
- peer: types for talking to peers, including a handler manages the connection
- terminal: some utility methods for presenting status and progress bars in the terminal, mostly using escape codes
- torrent_files: contains types and methods for parsing torrent files into useful structs
- tracker: communication with trackers, registering as a peer and finding other peers
- util: at present, just some useful concurrency functions

## LLM Use disclaimer

LLMs only generated unit tests, all other code is hand written.
