package bitfields

import (
	"reflect"
	"testing"
)

func TestCreateBlankBitfield(t *testing.T) {
	tests := []struct {
		name   string
		length int
		want   BitField
	}{
		{
			name:   "zero length",
			length: 0,
			want:   BitField{Data: []byte{}, Length: 0},
		},
		{
			name:   "single bit",
			length: 1,
			want:   BitField{Data: []byte{0}, Length: 1},
		},
		{
			name:   "eight bits",
			length: 8,
			want:   BitField{Data: []byte{0}, Length: 8},
		},
		{
			name:   "nine bits",
			length: 9,
			want:   BitField{Data: []byte{0, 0}, Length: 9},
		},
		{
			name:   "sixteen bits",
			length: 16,
			want:   BitField{Data: []byte{0, 0}, Length: 16},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateBlankBitfield(tt.length)
			if !reflect.DeepEqual(got.Data, tt.want.Data) {
				t.Errorf("CreateBlankBitfield().Data = %v, want %v", got.Data, tt.want.Data)
			}
			if got.Length != tt.want.Length {
				t.Errorf("CreateBlankBitfield().Length = %v, want %v", got.Length, tt.want.Length)
			}
		})
	}
}

func TestNewBitfield(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		length int
		want   BitField
	}{
		{
			name:   "empty bitfield",
			data:   []byte{},
			length: 0,
			want:   BitField{Data: []byte{}, Length: 0},
		},
		{
			name:   "single byte bitfield",
			data:   []byte{0xFF},
			length: 8,
			want:   BitField{Data: []byte{0xFF}, Length: 8},
		},
		{
			name:   "multi byte bitfield",
			data:   []byte{0xAA, 0x55},
			length: 16,
			want:   BitField{Data: []byte{0xAA, 0x55}, Length: 16},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBitfield(tt.data, tt.length)
			if !reflect.DeepEqual(got.Data, tt.want.Data) {
				t.Errorf("NewBitfield().Data = %v, want %v", got.Data, tt.want.Data)
			}
			if got.Length != tt.want.Length {
				t.Errorf("NewBitfield().Length = %v, want %v", got.Length, tt.want.Length)
			}
		})
	}
}

func TestBitfieldSet(t *testing.T) {
	tests := []struct {
		name     string
		bf       BitField
		index    uint
		want     []byte
		want_err bool
	}{
		{
			name:     "single byte set first bit",
			bf:       NewBitfield([]byte{0}, 8),
			index:    0,
			want:     []byte{0b10000000},
			want_err: false,
		},
		{
			name:     "single byte set middle bit",
			bf:       NewBitfield([]byte{0}, 8),
			index:    3,
			want:     []byte{0b00010000},
			want_err: false,
		},
		{
			name:     "single byte set last bit",
			bf:       NewBitfield([]byte{0}, 8),
			index:    7,
			want:     []byte{0b00000001},
			want_err: false,
		},
		{
			name:     "single byte already set",
			bf:       NewBitfield([]byte{0b00010000}, 8),
			index:    3,
			want:     []byte{0b00010000},
			want_err: false,
		},
		{
			name:     "single byte out of range",
			bf:       NewBitfield([]byte{0}, 8),
			index:    8,
			want:     []byte{0},
			want_err: true,
		},
		{
			name:     "multi byte set",
			bf:       NewBitfield([]byte{0, 0, 0}, 24),
			index:    12,
			want:     []byte{0, 0b00001000, 0},
			want_err: false,
		},
		{
			name:     "multi byte out of range",
			bf:       NewBitfield([]byte{0, 0, 0}, 24),
			index:    24,
			want:     []byte{0, 0, 0},
			want_err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := tt.bf
			err := bf.Set(tt.index)
			if (err != nil) != tt.want_err {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.want_err)
				return
			}
			if err == nil && !reflect.DeepEqual(bf.Data, tt.want) {
				t.Errorf("Set() = %v, want %v", bf.Data, tt.want)
			}
		})
	}
}

func TestBitfieldGet(t *testing.T) {
	tests := []struct {
		name  string
		bf    BitField
		index int
		want  bool
	}{
		{
			name:  "single byte get true first bit",
			bf:    NewBitfield([]byte{0b10000000}, 8),
			index: 0,
			want:  true,
		},
		{
			name:  "single byte get true middle bit",
			bf:    NewBitfield([]byte{0b00010000}, 8),
			index: 3,
			want:  true,
		},
		{
			name:  "single byte get true last bit",
			bf:    NewBitfield([]byte{0b00000001}, 8),
			index: 7,
			want:  true,
		},
		{
			name:  "single byte get false",
			bf:    NewBitfield([]byte{0}, 8),
			index: 3,
			want:  false,
		},
		{
			name:  "single byte out of range",
			bf:    NewBitfield([]byte{0}, 8),
			index: 8,
			want:  false,
		},
		{
			name:  "single byte negative index",
			bf:    NewBitfield([]byte{0}, 8),
			index: -1,
			want:  false,
		},
		{
			name:  "multi byte get true",
			bf:    NewBitfield([]byte{0, 0b00010000, 0}, 24),
			index: 11,
			want:  true,
		},
		{
			name:  "multi byte get false",
			bf:    NewBitfield([]byte{0, 0, 0}, 24),
			index: 12,
			want:  false,
		},
		{
			name:  "multi byte out of range",
			bf:    NewBitfield([]byte{0, 0, 0}, 24),
			index: 24,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bf.Get(tt.index)
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitfieldBitString(t *testing.T) {
	tests := []struct {
		name string
		bf   BitField
		want string
	}{
		{
			name: "empty bitfield",
			bf:   NewBitfield([]byte{}, 0),
			want: "",
		},
		{
			name: "single byte all zeros",
			bf:   NewBitfield([]byte{0}, 8),
			want: "00000000",
		},
		{
			name: "single byte all ones",
			bf:   NewBitfield([]byte{0xFF}, 8),
			want: "11111111",
		},
		{
			name: "single byte pattern",
			bf:   NewBitfield([]byte{0b10101010}, 8),
			want: "10101010",
		},
		{
			name: "partial byte",
			bf:   NewBitfield([]byte{0b11110000}, 4),
			want: "1111",
		},
		{
			name: "multi byte pattern",
			bf:   NewBitfield([]byte{0xAA, 0x55}, 16),
			want: "1010101001010101",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bf.BitString()
			if got != tt.want {
				t.Errorf("BitString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitfieldIncomplete(t *testing.T) {
	tests := []struct {
		name string
		bf   BitField
		want bool
	}{
		{
			name: "empty bitfield",
			bf:   NewBitfield([]byte{}, 0),
			want: false,
		},
		{
			name: "single byte all zeros",
			bf:   NewBitfield([]byte{0}, 8),
			want: true,
		},
		{
			name: "single byte all ones",
			bf:   NewBitfield([]byte{0xFF}, 8),
			want: false,
		},
		{
			name: "single byte partial",
			bf:   NewBitfield([]byte{0b00010000}, 8),
			want: true,
		},
		{
			name: "multi byte all ones",
			bf:   NewBitfield([]byte{0xFF, 0xFF}, 16),
			want: false,
		},
		{
			name: "multi byte partial",
			bf:   NewBitfield([]byte{0xFF, 0x00}, 16),
			want: true,
		},
		{
			name: "partial length",
			bf:   NewBitfield([]byte{0xFF}, 4),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bf.Incomplete()
			if got != tt.want {
				t.Errorf("Incomplete() = %v, want %v", got, tt.want)
			}
		})
	}
}
