// MIT License
//
// Copyright (c) 2021 Xiantu Li
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package errors

import (
	"unicode/utf8"
	"unsafe"
)

const (
	hex = "0123456789abcdef"
)

var (
	// 标识 JSON string字符是否需要需要的转义
	safeSet = func() (safeSet [utf8.RuneSelf]bool) {
		for i := '\u0000'; i <= '\u007f'; i++ {
			if i < ' ' || i == '"' || i == '\\' {
				safeSet[i] = false
				continue
			}
			safeSet[i] = true
		}
		return
	}()
)

type writeBuffer struct {
	buf []byte
}

func NewWriteBuffer(n int) (buf *writeBuffer) { //nolint:bgolint
	return &writeBuffer{
		buf: make([]byte, 0, n),
	}
}

func (buf *writeBuffer) Bytes() []byte { return buf.buf }

func (buf *writeBuffer) String() string {
	return *(*string)(unsafe.Pointer(&buf.buf))
}

func (buf *writeBuffer) Grow(n int) {
	if cap(buf.buf) == 0 {
		buf.buf = make([]byte, 0, n)
		return
	}
	if cap(buf.buf)-len(buf.buf) >= n {
		return
	}
	bs := buf.buf
	buf.buf = make([]byte, len(bs), n+len(bs))
	copy(buf.buf, bs)
}
func (buf *writeBuffer) Write(p []byte) {
	buf.buf = append(buf.buf, p...)
}
func (buf *writeBuffer) WriteString(s string) {
	buf.buf = append(buf.buf, s...)
}
func (buf *writeBuffer) WriteByte(c byte) {
	buf.buf = append(buf.buf, c)
}

// WriteEscape 抄 std json 库
func (buf *writeBuffer) WriteEscape(src string) {
	start := 0
	for i := 0; i < len(src); {
		if c := src[i]; c < utf8.RuneSelf {
			if safeSet[c] {
				i++
				continue
			}
			if start < i {
				buf.WriteString(src[start:i])
			}
			buf.WriteByte('\\')
			switch c {
			case '\\', '"':
				buf.WriteByte(c)
			case '\n':
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('t')
			default:
				buf.WriteString(`u00`)
				buf.WriteByte(hex[c>>4])
				buf.WriteByte(hex[c&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(src[i:])
		if c == utf8.RuneError && size == 1 {
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				buf.WriteString(src[start:i])
			}
			// to: \u202 hex[c&0xF]
			buf.WriteString(`\u202`)
			buf.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		// 跳过无效字符
		i += size
	}
	if start < len(src) {
		buf.WriteString(src[start:])
	}
}

func countEscape(str string) (l int, escape bool) {
	start := 0
	for i := 0; i < len(str); {
		if buf := str[i]; buf < utf8.RuneSelf {
			if safeSet[buf] {
				i++
				continue
			}
			escape = true
			l += i - start
			switch buf {
			case '\\', '"', '\n', '\r', '\t':
				l += 2
			default:
				l += 6
			}
			i++
			start = i
			continue
		}
		escape = true
		c, size := utf8.DecodeRuneInString(str[i:])
		if c == utf8.RuneError && size == 1 {
			// to: \ufffd
			i += size
			start = i
			l += 6
			continue
		}
		if c == '\u2028' || c == '\u2029' {
			// to: \u202 hex[c&0xF]
			i += size
			start = i
			l += 6
			continue
		}
		i += size
	}
	l += len(str) - start
	return
}
