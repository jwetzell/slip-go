package slip

import (
	"errors"
	"io"
)

type Reader struct {
	pendingBytes  []byte
	r             io.Reader
	messages      [][]byte
	underlyingEOF bool
}

func NewReader(r io.Reader) *Reader {
	return &Reader{pendingBytes: []byte{}, r: r}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	readBuffer := make([]byte, len(p))
	n, err = r.r.Read(readBuffer)
	if err != nil {
		if errors.Is(err, io.EOF) {
			r.underlyingEOF = true
		} else {
			return 0, err
		}
	}
	bytes := readBuffer[:n]
	escapeNext := false
	for _, packetByte := range bytes {

		if packetByte == esc {
			escapeNext = true
			continue
		}

		if escapeNext {
			switch packetByte {
			case esc_end:
				r.pendingBytes = append(r.pendingBytes, end)
			case esc_esc:
				r.pendingBytes = append(r.pendingBytes, esc)
			}
			escapeNext = false
		} else if packetByte == end {
			if len(r.pendingBytes) == 0 {
				// probably opening END byte, can discard
				continue
			} else {
				r.messages = append(r.messages, r.pendingBytes)
			}
			r.pendingBytes = []byte{}
		} else {
			r.pendingBytes = append(r.pendingBytes, packetByte)
		}
	}

	if len(r.messages) == 0 && r.underlyingEOF {
		return 0, io.EOF
	}

	if len(r.messages) > 0 {
		copy(p, r.messages[0])
		n = len(r.messages[0])
		r.messages = r.messages[1:]
		return n, nil
	}
	return 0, nil
}
