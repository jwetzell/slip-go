package slip_test

import (
	"bytes"
	"testing"

	"github.com/jwetzell/slip-go"
	"github.com/stretchr/testify/require"
)

func TestGoodWriter(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "one message",
			input:    []byte{'h', 'i'},
			expected: []byte{0xc0, 'h', 'i', 0xc0},
		},
		{
			name:     "message w/escaping",
			input:    []byte{'b', 0xc0, 'e'},
			expected: []byte{0xc0, 'b', 0xdb, 0xdc, 'e', 0xc0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(nil)
			w := slip.NewWriter(buffer)
			_, err := w.Write(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, buffer.Bytes())
		})
	}
}
