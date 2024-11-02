package encoding

import (
	"strconv"
	"testing"
)

func TestEncodeString(t *testing.T) {
	var tests = []struct {
		str, want string
	}{
		{"https://yahoo.com", "c88f320dec138ba5ab0a5f990ff082ba"},
		{"https://google.com", "99999ebcfdb78df077ad2727fd00969f"},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			ans := EncodeMd5(tt.str)
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}

func TestEncodeInt64(t *testing.T) {
	var tests = []struct {
		id, want string
	}{
		{"1", "1"},
		{"61", "z"},
		{"62", "10"},
		{"1000", "G8"},
		{"1000000000", "15ftgG"},
		{"4000000000000", "18QB6MKG"},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			ans := Encode(stoi(tt.id))
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}

func stoi(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1
	}
	return i
}
