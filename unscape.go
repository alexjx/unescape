package unescape

import (
	"bytes"
	"io"
	"unicode/utf8"
)

const Marker = '\x1B'

func isTerminator(c rune) bool {
	return (c >= 0x40 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a)
}

// Unescaper is used to remove ANSI escape sequences from a stream of bytes.
type Unescaper struct {
	Forward io.Writer

	ansi       bool
	ansiseq    bytes.Buffer
	lastseq    bytes.Buffer
	seqchanged bool
	runeBuf    []byte
}

// Write implements io.Writer, it consumes the bytes and removes ANSI escape
func (w *Unescaper) Write(b []byte) (int, error) {
	for _, c := range string(b) {
		if c == Marker {
			// ANSI escape sequence
			w.ansi = true
			w.seqchanged = true
			_, _ = w.ansiseq.WriteRune(c)
		} else if w.ansi {
			_, _ = w.ansiseq.WriteRune(c)
			if isTerminator(c) {
				// ANSI sequence terminated
				w.ansi = false

				if bytes.HasSuffix(w.ansiseq.Bytes(), []byte("[0m")) {
					// reset sequence
					w.lastseq.Reset()
					w.seqchanged = false
				} else if c == 'm' {
					// color code
					_, _ = w.lastseq.Write(w.ansiseq.Bytes())
				}
			}
		} else {
			_, err := w.writeRune(c)
			if err != nil {
				return 0, err
			}
		}
	}

	return len(b), nil
}

func (w *Unescaper) writeRune(r rune) (int, error) {
	if w.runeBuf == nil {
		w.runeBuf = make([]byte, utf8.UTFMax)
	}
	n := utf8.EncodeRune(w.runeBuf, r)
	return w.Forward.Write(w.runeBuf[:n])
}
