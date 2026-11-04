package slip_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/jwetzell/slip-go"
	"github.com/stretchr/testify/require"
)

func TestGoodReader(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected [][]byte
	}{
		{
			name:     "one message",
			input:    []byte{0xc0, 'h', 'i', 0xc0},
			expected: [][]byte{{'h', 'i'}},
		},
		{
			name:     "multiple messages",
			input:    []byte{0xc0, 'h', 'i', 0xc0, 0xc0, 'b', 'y', 'e', 0xc0},
			expected: [][]byte{{'h', 'i'}, {'b', 'y', 'e'}},
		},
		{
			name:     "message w/escaping",
			input:    []byte{0xc0, 'b', 0xdb, 0xdc, 'e', 0xc0},
			expected: [][]byte{{'b', 0xc0, 'e'}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := slip.NewReader(bytes.NewReader(tc.input))
			messages := [][]byte{}
			for {
				buf := make([]byte, 1024)
				n, err := r.Read(buf)
				if err != nil {
					if !errors.Is(err, io.EOF) {
						t.Fatalf("unexpected error: %v", err)
					}
					break
				}
				messages = append(messages, buf[:n])
			}
			require.Equal(t, tc.expected, messages)
		})
	}
}

func TestBadReader(t *testing.T) {
	testCases := []struct {
		name      string
		input     []byte
		buffer    []byte
		errString string
	}{
		{
			name:      "not enough bytes",
			input:     []byte{0xc0, 'h', 'i', 0xc0},
			buffer:    make([]byte, 1),
			errString: "short buffer",
		},
		{
			name:      "EOF",
			input:     []byte{0xc0, 'h', 'i'},
			buffer:    make([]byte, 2),
			errString: "EOF",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := slip.NewReader(bytes.NewReader(tc.input))
			messages := [][]byte{}
			didError := false
			for {
				n, err := r.Read(tc.buffer)
				require.Equal(t, 0, n)
				if err != nil {
					require.Contains(t, err.Error(), tc.errString)
					didError = true
					break
				}
				messages = append(messages, tc.buffer[:n])
			}
			require.Empty(t, messages)
			require.True(t, didError)
		})
	}
}
