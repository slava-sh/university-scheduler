package main

import (
	"bufio"
	"bytes"
	"io"
)

type FastReader struct {
	r   *bufio.Reader
	buf string
}

func NewFastReader(r io.Reader) *FastReader {
	return &FastReader{
		r: bufio.NewReader(r),
	}
}

func (r *FastReader) advance() {
	for len(r.buf) != 0 && r.buf[0] == ' ' {
		r.buf = r.buf[1:]
	}
	var err error
	for len(r.buf) == 0 && err != io.EOF {
		buf := new(bytes.Buffer)
		for {
			var chunk []byte
			var more bool
			chunk, more, err = r.r.ReadLine()
			buf.Write(chunk)
			if !more {
				break
			}
		}
		r.buf = buf.String()
	}
}

func (r *FastReader) NextLine() string {
	r.advance()
	line := r.buf
	r.buf = ""
	return line
}

func (r *FastReader) NextWord() string {
	r.advance()
	var word string
	wordStart := 0
	for i := 0; i < len(r.buf); i++ {
		b := r.buf[i]
		if b == ' ' {
			if i != wordStart {
				word = r.buf[wordStart:i]
				r.buf = r.buf[i:]
				break
			}
			wordStart = i + 1
		}
	}
	if len(word) == 0 {
		word = r.buf
		r.buf = ""
	}
	return word
}

func (r *FastReader) NextInt() int {
	return parseInt(r.NextWord())
}

func parseInt(word string) int {
	sign := 1
	if word[0] == '-' {
		sign = -1
		word = word[1:]
	}
	result := 0
	for i := 0; i < len(word); i++ {
		result = result*10 + int(word[i]) - '0'
	}
	result *= sign
	return result
}
