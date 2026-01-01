package main

import (
	"reflect"
	"testing"
)

func TestParseStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		want     any
		want_rem []byte
		want_err bool
	}{
		{
			name:     "basic parse",
			input:    []byte("4:spam"),
			want:     "spam",
			want_rem: []byte{},
			want_err: false,
		},

		{
			name:     "remainder returned",
			input:    []byte("4:spamtest"),
			want:     "spam",
			want_rem: []byte("test"),
			want_err: false,
		},

		{
			name:     "longer parse",
			input:    []byte("10:abcdefghij"),
			want:     "abcdefghij",
			want_rem: []byte{},
			want_err: false,
		},

		{
			name:     "bad length",
			input:    []byte("02:aa"),
			want:     nil,
			want_rem: nil,
			want_err: true,
		},

		{
			name:     "wrong length",
			input:    []byte("2:a"),
			want:     nil,
			want_rem: nil,
			want_err: true,
		},

		{
			name:     "invalid header",
			input:    []byte("4aspam"),
			want:     nil,
			want_rem: nil,
			want_err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, rem, err := parse(tt.input)
			if (err != nil) != tt.want_err {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.want_err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(rem, tt.want_rem) {
				t.Errorf("parse() = %v, want remainder %v", rem, tt.want_rem)
			}
		})
	}
}
