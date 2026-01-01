package main

import (
	"fmt"
	"strconv"
)

func parse(data []byte) (any, []byte, error) {
	data_len := len(data)

	if data_len == 0 {
		return nil, nil, nil
	}

	// if data[0] == 'd' {
	// 	return parse_dict(data)
	// } else if data[0] == 'l' {
	// 	return parse_list(data)
	// }

	i := 0
	for data[i] >= '0' && data[i] <= '9' {
		if data_len <= i {
			return nil, nil, fmt.Errorf("unrecognised start token")
		}
		i++
	}

	if i == 0 {
		return nil, nil, fmt.Errorf("unrecognised start token")
	} else if data[0] == '0' {
		return nil, nil, fmt.Errorf("invalid string length - starts with 0")
	}

	length, _ := strconv.Atoi(string(data[0:i]))
	if data_len < i+1+length {
		return nil, nil, fmt.Errorf("invalid string length - string len does not match length header")
	}
	if data[i] != ':' {
		return nil, nil, fmt.Errorf("invalid header, missing separator colon")
	}

	s := data[i+1 : i+1+length]
	return string(s), data[i+1+length:], nil
}
