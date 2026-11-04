package slip

import "io"

type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	var encodedBytes = []byte{end}

	for _, byteToEncode := range p {
		switch byteToEncode {
		case end:
			encodedBytes = append(encodedBytes, esc, esc_end)
		case esc:
			encodedBytes = append(encodedBytes, esc, esc_esc)
		default:
			encodedBytes = append(encodedBytes, byteToEncode)
		}
	}

	encodedBytes = append(encodedBytes, end)
	return w.w.Write(encodedBytes)
}
