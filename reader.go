package slip

import (
	"bytes"
	"io"
	"sync"
)

type Reader struct {
	mu sync.Mutex
	r  io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	packet := bytes.Buffer{}

	readBuf := make([]byte, 1)

	// NOTE(jwetzell): one byte reads seems odd but not sure how else to read as little as possible
	for {
		n, err := r.r.Read(readBuf)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}

		switch readBuf[0] {
		case end:
			if packet.Len() > 0 {
				if packet.Len() > len(p) {
					return 0, io.ErrShortBuffer
				}
				copy(p, packet.Bytes())
				return packet.Len(), nil
			}
		case esc:
			n, err := r.r.Read(readBuf)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				return 0, nil
			}
			switch readBuf[0] {
			case esc_end:
				packet.WriteByte(end)
			case esc_esc:
				packet.WriteByte(esc)
			default:
				return 0, io.ErrUnexpectedEOF
			}
		default:
			packet.WriteByte(readBuf[0])
		}
	}
}
