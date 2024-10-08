// reference: https://github.com/uber-go/zap/blob/master/buffer/buffer.go
package logx

import (
	"slices"
	"strconv"
	"time"
	"unsafe"
)

// Buffer is a thin wrapper around a byte slice. It's intended to be pooled, so
// the only way to construct one is via a Pool.
type Buffer struct {
	bs []byte
}

func NewBuffer(buf []byte) *Buffer { return &Buffer{bs: buf} }

// AppendByte writes a single byte to the Buffer.
func (b *Buffer) AppendByte(v byte) {
	b.bs = append(b.bs, v)
}

// AppendBytes writes the given slice of bytes to the Buffer.
func (b *Buffer) AppendBytes(v []byte) {
	b.bs = append(b.bs, v...)
}

// AppendString writes a string to the Buffer.
func (b *Buffer) AppendString(s string) {
	b.bs = append(b.bs, s...)
}

// AppendInt appends an integer to the underlying buffer (assuming base 10).
func (b *Buffer) AppendInt(i int64) {
	b.bs = strconv.AppendInt(b.bs, i, 10)
}

// AppendTime appends the time formatted using the specified layout.
func (b *Buffer) AppendTime(t time.Time, layout string) {
	b.bs = t.AppendFormat(b.bs, layout)
}

// AppendUint appends an unsigned integer to the underlying buffer (assuming
// base 10).
func (b *Buffer) AppendUint(i uint64) {
	b.bs = strconv.AppendUint(b.bs, i, 10)
}

// AppendBool appends a bool to the underlying buffer.
func (b *Buffer) AppendBool(v bool) {
	b.bs = strconv.AppendBool(b.bs, v)
}

// AppendFloat appends a float to the underlying buffer. It doesn't quote NaN
// or +/- Inf.
func (b *Buffer) AppendFloat(f float64, bitSize int) {
	b.bs = strconv.AppendFloat(b.bs, f, 'f', -1, bitSize)
}

func (b *Buffer) AppendQuote(s string) {
	b.bs = strconv.AppendQuote(b.bs, s)
}

func (b *Buffer) TryGrow(size int) { b.bs = slices.Grow(b.bs, size) }

// Reset resets the underlying byte slice. Subsequent writes re-use the slice's
// backing array.
func (b *Buffer) Reset() { b.bs = b.bs[:0] }

// Len returns the length of the underlying byte slice.
func (b *Buffer) Len() int { return len(b.bs) }

// Cap returns the capacity of the underlying byte slice.
func (b *Buffer) Cap() int { return cap(b.bs) }

// Bytes returns a mutable reference to the underlying byte slice.
func (b *Buffer) Bytes() []byte { return b.bs }

func (b *Buffer) String() string { return unsafe.String(&b.bs[0], len(b.bs)) }

func (b *Buffer) Clone() *Buffer { return &Buffer{bs: slices.Clone(b.bs)} }

// Write implements io.Writer.
func (b *Buffer) Write(bs []byte) (int, error) {
	b.bs = append(b.bs, bs...)
	return len(bs), nil
}

// WriteByte writes a single byte to the Buffer.
//
// Error returned is always nil, function signature is compatible
// with bytes.Buffer and bufio.Writer
func (b *Buffer) WriteByte(v byte) error {
	b.AppendByte(v)
	return nil
}

// WriteString writes a string to the Buffer.
//
// Error returned is always nil, function signature is compatible
// with bytes.Buffer and bufio.Writer
func (b *Buffer) WriteString(s string) (int, error) {
	b.AppendString(s)
	return len(s), nil
}

// TrimNewline trims any final "\n" byte from the end of the buffer.
func (b *Buffer) TrimNewline() {
	if i := len(b.bs) - 1; i >= 0 {
		if b.bs[i] == '\n' {
			b.bs = b.bs[:i]
		}
	}
}
