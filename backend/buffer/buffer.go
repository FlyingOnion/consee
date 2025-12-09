// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package buffer

import (
	"bytes"
	"io"
	"strconv"
	"time"
	"unicode/utf8"
)

type Buffer struct {
	data []byte
}

func NewBuffer(b []byte) *Buffer {
	return &Buffer{data: b}
}

// Reader creates an io.Reader for Read.
func (b Buffer) Reader() io.Reader { return bytes.NewReader(b.BytesCopy()) }

// ReadCloser creates an io.ReadCloser, but Close does nothing.
func (b Buffer) ReadCloser() io.ReadCloser { return io.NopCloser(bytes.NewReader(b.BytesCopy())) }

// Write implements io.Writer, and it will write data to buffer directly, without considering any conditions.
// If you need chaining call, use WriteBytes instead.
func (b *Buffer) Write(data []byte) (int, error) {
	b.data = append(b.data, data...)
	return len(data), nil
}

// String implements Stringer.
func (b Buffer) String() string { return string(b.data) }

func (b Buffer) Len() int { return len(b.data) }

func (b Buffer) Cap() int { return cap(b.data) }

func (b *Buffer) Reset() {
	b.data = b.data[:0]
}

// Bytes returns b.data. It is NOT a copy, so the result will be changed after b modified.
func (b Buffer) Bytes() []byte { return b.data }

func (b Buffer) BytesCopy() []byte {
	bs := make([]byte, b.Len())
	copy(bs, b.data)
	return bs
}

func (b *Buffer) WriteInt(i int) *Buffer {
	return b.WriteInt64(int64(i))
}

func (b *Buffer) WriteInt8(i int8) *Buffer {
	return b.WriteInt64(int64(i))
}

func (b *Buffer) WriteInt16(i int16) *Buffer {
	return b.WriteInt64(int64(i))
}

func (b *Buffer) WriteInt32(i int32) *Buffer {
	return b.WriteInt64(int64(i))
}

func (b *Buffer) WriteInt64(i int64) *Buffer {
	b.data = strconv.AppendInt(b.data, i, 10)
	return b
}

func (b *Buffer) WriteUint(i uint) *Buffer {
	return b.WriteUint64(uint64(i))
}

func (b *Buffer) WriteUint8(i uint8) *Buffer {
	return b.WriteUint64(uint64(i))
}

func (b *Buffer) WriteUint16(i uint16) *Buffer {
	return b.WriteUint64(uint64(i))
}

func (b *Buffer) WriteUint32(i uint32) *Buffer {
	return b.WriteUint64(uint64(i))
}

func (b *Buffer) WriteUint64(i uint64) *Buffer {
	b.data = strconv.AppendUint(b.data, i, 10)
	return b
}

func (b *Buffer) WriteBytes(s []byte) *Buffer {
	b.data = append(b.data, s...)
	return b
}

func (b *Buffer) WriteByte(c byte) *Buffer {
	b.data = append(b.data, c)
	return b
}

func (b *Buffer) WriteRune(r rune) *Buffer {
	var br [4]byte
	n := utf8.EncodeRune(br[:], r)
	b.data = append(b.data, br[:n]...)
	return b
}

func (b *Buffer) WriteString(s string) *Buffer {
	b.data = append(b.data, s...)
	return b
}

func isInsecureCharacter(c byte) bool {
	return c < 32 || c == '\\' || c == '"'
}

var hex = "0123456789abcdef"

func (b *Buffer) WriteJsonSafeBytes(s []byte) *Buffer {
	return b.WriteJsonSafeString(string(s))
}

func (b *Buffer) WriteJsonSafeString(s string) *Buffer {
	return b.writeJsonSafeString(s)
}

// WriteJsonSafeString converts strings with json syntax to safe strings (by adding a '\\' prefix)
func (b *Buffer) writeJsonSafeString(s string) *Buffer {
	i := 0
	for j := 0; j < len(s); {
		if c := s[j]; c < utf8.RuneSelf {
			if !isInsecureCharacter(c) {
				j++
				continue
			}
			if i < j {
				b.WriteString(s[i:j])
			}
			b.WriteByte('\\')
			switch c {
			case '\\', '"':
				b.WriteByte(c)
			case '\n':
				b.WriteByte('n')
			case '\r':
				b.WriteByte('r')
			case '\t':
				b.WriteByte('t')
			default:
				b.WriteString("u00").WriteByte(hex[c>>4]).WriteByte(hex[c&0xF])
			}
			j++
			i = j
			continue
		}
		c, size := utf8.DecodeRuneInString(s[j:])
		if c == utf8.RuneError && size == 1 {
			if i < j {
				b.WriteString(s[i:j])
			}
			b.WriteString(`\ufffd`)
			j += size
			i = j
			continue
		}
		if c == '\u2028' || c == '\u2029' {
			if i < j {
				b.WriteString(s[i:j])
			}
			b.WriteString(`\u202`).WriteByte(hex[c&0xF])
			j += size
			i = j
			continue
		}
		j += size
	}
	if i < len(s) {
		b.WriteString(s[i:])
	}
	return b
}

func (b *Buffer) WriteFloat32(f float32, format byte, precision int) *Buffer {
	b.data = strconv.AppendFloat(b.data, float64(f), format, precision, 32)
	return b
}

func (b *Buffer) WriteFloat64(f float64, format byte, precision int) *Buffer {
	b.data = strconv.AppendFloat(b.data, f, format, precision, 64)
	return b

}

// WritePointer appends the hexadecimal address of p with a "0x" prefix.
// uintptr is like uint, its bitsize depends on the operating system.
func (b *Buffer) WritePointer(p uintptr) *Buffer {
	return b.WriteString("0x").WriteUint(uint(p))
}

func (b *Buffer) WriteBool(v bool) *Buffer {
	b.data = strconv.AppendBool(b.data, v)
	return b
}

func (b *Buffer) WriteTime(t time.Time, format string) *Buffer {
	b.data = t.AppendFormat(b.data, format)
	return b
}

func (b *Buffer) WriteDuration(d time.Duration) *Buffer {
	return b.WriteString(d.String())
}
