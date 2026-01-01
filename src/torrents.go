package main

import "fmt"

// Decodes a torrent file into the relevant properties for further downloading

type TorrentMetadata struct {
	Announcers  []string
	InfoHash    []byte
	Name        string
	PieceLength int
	Pieces      []string
	Length      int
	Files       []TorrentFile
}

type TorrentFile struct {
	Path   string
	Length int
}

var nil_file TorrentMetadata

func parse_torrent_file(file_data []byte) (TorrentMetadata, error) {
	decoded, _, err := decode_from_bencoded(file_data)
	if err != nil {
		return nil_file, err
	}

	root, ok := decoded.(map[string]any)
	if !ok {
		return nil_file, fmt.Errorf("invalid torrent: root is not a dict")
	}

	announce, err := get_val[string](root, "announce")
	if err != nil {
		return nil_file, fmt.Errorf("invalid torrent: %A", err)
	}
	announcers := []string{announce}

	announce_list, err := get_string_list(root, "announce-list")
	if err == nil {
		announcers = append(announcers, announce_list...)
	}

	info, err := get_val[map[string]any](root, "info")
	if err != nil {
		return nil_file, fmt.Errorf("invalid torrent: %A", err)
	}

	name, err := get_val[string](info, "name")
	if err != nil {
		return nil_file, fmt.Errorf("invalid torrent: %A", err)
	}

	piece_length, err := get_val[int](info, "piece length")
	if err != nil {
		return nil_file, fmt.Errorf("invalid torrent: %A", err)
	}

	pieces, err := get_val[string](info, "pieces")
	if err != nil {
		return nil_file, fmt.Errorf("invalid torrent: %A", err)
	}
	pieces_parsed := []string{}
	for i := 0; i < len(pieces)/piece_length; i += piece_length {
		pieces_parsed = append(pieces_parsed, pieces[i*piece_length:i+1*piece_length])
	}

	return TorrentMetadata{
		Announcers:  announcers,
		InfoHash:    []byte{},
		Name:        name,
		PieceLength: piece_length,
		Pieces:      pieces_parsed,
		Length:      0,
		Files:       []TorrentFile{},
	}, nil
}

func get_val[T any](m map[string]any, key string) (T, error) {
	var nilT T
	val, exists := m[key]
	if !exists {
		return nilT, fmt.Errorf("key %s was not in map", key)
	}
	res, ok := val.(T)
	if !ok {
		return nilT, fmt.Errorf("key %s's value was an invalid type: %A", key, val)
	}
	return res, nil
}

func get_string_list(m map[string]any, key string) ([]string, error) {
	list, err := get_val[[]any](m, key)
	if err != nil {
		return nil, err
	}
	results := []string{}
	for _, v := range list {
		s, ok := v.([]byte)
		if !ok {
			return nil, fmt.Errorf("a non-string value was in the list: %A", v)
		}
		results = append(results, string(s))
	}
	return results, nil
}
