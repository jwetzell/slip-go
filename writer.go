package slip

import (
	"bytes"
	"io"
	"sync"
)

type Writer struct {
	mu sync.Mutex
	w  io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	var encodedBytes = bytes.Buffer{}

	encodedBytes.WriteByte(end)

	for _, byteToEncode := range p {
		switch byteToEncode {
		case end:
			encodedBytes.Write([]byte{esc, esc_end})
		case esc:
			encodedBytes.Write([]byte{esc, esc_esc})
		default:
			encodedBytes.WriteByte(byteToEncode)
		}
	}

	encodedBytes.WriteByte(end)
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.w.Write(encodedBytes.Bytes())
}
