package bytepal

import (
	"math"
)

// bitWriter is the underlying base structure of this writer, shares all common characteristics as a bit is the lowest.
type bitWriter struct {
	bytes        []byte
	currentIndex int
}

// SetCurrentWrite moves the current index that will be written.
func (w *bitWriter) SetCurrentWrite(index int) {
	w.currentIndex = index
}

// Size returns the size of the payload
func (w *bitWriter) Size() int {
	return len(w.bytes)
}

// Payload returns the byte buffer inside the FixedWriter
func (a *bitWriter) Payload() []byte {
	return a.bytes
}

// BitAccess allows writing bits to the array, expands the size if capacity is exceeded.
//	If a normal write is executed, a new BitAccess must be called to get the appropriate bitPosition.
func (a *bitWriter) BitAccess() func(uint, uint) {
	bitPosition := uint(a.currentIndex * 8)
	return func(numBits, value uint) {
		bytePos := bitPosition >> 3
		bitOffset := 8 - (bitPosition & 7)

		if availBits := uint(len(a.bytes)*8) - bitPosition; numBits > availBits {
			addSize := make([]byte, int(math.Ceil(float64(numBits-availBits)/8)))
			a.bytes = append(a.bytes, addSize...)
		}
		bitPosition += numBits

		// Check if numBits is greater than int8, if so write value and continue.
		for ; numBits > bitOffset; bitOffset = 8 {
			a.bytes[bytePos] &= ^byte(bitMask[bitOffset])
			a.bytes[bytePos] |= byte((value >> (numBits - bitOffset)) & bitMask[bitOffset])
			bytePos++
			numBits -= bitOffset
		}
		if numBits == bitOffset {
			a.bytes[bytePos] &= ^byte(bitMask[bitOffset])
			a.bytes[bytePos] |= byte(value & bitMask[bitOffset])
		} else {
			a.bytes[bytePos] &= ^(byte(bitMask[numBits] << (bitOffset - numBits)))
			a.bytes[bytePos] |= byte(value & bitMask[numBits] << (bitOffset - numBits))
		}
		a.currentIndex = int(bitPosition+7) / 8
	}
}

type Writer interface {
	Size() int
	Payload() []byte
	BitAccess() func(uint, uint)

	Write([]uint8)
	WriteUInt8(uint8)
	WriteInt16(int16)
	WriteLEInt16(int16)
	WriteInt32(int32)
	WriteLEInt32(int32)
	WriteInt64(int64)
	WriteLEInt64(int64)
	WriteString(string, byte)
}

// FixedWriter represents a fixed buffer size with functions to write data to its buffer
var _ Writer = &FixedWriter{}

type FixedWriter struct {
	*bitWriter
}

// NewFixedWriter allocates an Array with a fixed size and values already initialized.
//	NOTE: See benchmarks in terms of comparison between FixedWriter & ExpandableWriter to determine which you want.
func NewFixedWriter(size int) Writer {
	return &FixedWriter{
		&bitWriter{
			bytes: make([]byte, size),
		},
	}
}

// Writes a byte onto the buffer
func (a *FixedWriter) WriteUInt8(v uint8) {
	a.bytes[a.currentIndex] = v
	a.currentIndex++
}

// WriteUInt16 writes two bytes to the buffer
func (a *FixedWriter) WriteInt16(v int16) {
	a.bytes[a.currentIndex] = byte(v >> 8)
	a.bytes[a.currentIndex+1] = byte(v)
	a.currentIndex += 2
}

// WriteInt16 writes two bytes in little endian to the buffer
func (a *FixedWriter) WriteLEInt16(v int16) {
	a.bytes[a.currentIndex+1] = byte(v >> 8)
	a.bytes[a.currentIndex] = byte(v)
	a.currentIndex += 2
}

// WriteLEInt32 writes a int32 to the buffer in little Endian order
func (a *FixedWriter) WriteLEInt32(v int32) {
	a.bytes[a.currentIndex+3] = byte(v >> 24)
	a.bytes[a.currentIndex+2] = byte(v >> 16)
	a.bytes[a.currentIndex+1] = byte(v >> 8)
	a.bytes[a.currentIndex] = byte(v)
	a.currentIndex += 4
}

// WriteInt32 writes a int32 to the buffer in Big Endian order
func (a *FixedWriter) WriteInt32(v int32) {
	a.bytes[a.currentIndex+3] = byte(v)
	a.bytes[a.currentIndex+2] = byte(v >> 8)
	a.bytes[a.currentIndex+1] = byte(v >> 16)
	a.bytes[a.currentIndex] = byte(v >> 24)
	a.currentIndex += 4
}

// WriteLEInt64 writes a int64 to the buffer in little Endian order
func (a *FixedWriter) WriteLEInt64(v int64) {
	a.bytes[a.currentIndex+7] = byte(v >> 56)
	a.bytes[a.currentIndex+6] = byte(v >> 48)
	a.bytes[a.currentIndex+5] = byte(v >> 40)
	a.bytes[a.currentIndex+4] = byte(v >> 32)
	a.bytes[a.currentIndex+3] = byte(v >> 24)
	a.bytes[a.currentIndex+2] = byte(v >> 16)
	a.bytes[a.currentIndex+1] = byte(v >> 8)
	a.bytes[a.currentIndex] = byte(v)
	a.currentIndex += 8
}

// WriteInt64 writes a int64 to the buffer in Big Endian order
func (a *FixedWriter) WriteInt64(v int64) {
	a.bytes[a.currentIndex] = byte(v >> 56)
	a.bytes[a.currentIndex+1] = byte(v >> 48)
	a.bytes[a.currentIndex+2] = byte(v >> 40)
	a.bytes[a.currentIndex+3] = byte(v >> 32)
	a.bytes[a.currentIndex+4] = byte(v >> 24)
	a.bytes[a.currentIndex+5] = byte(v >> 16)
	a.bytes[a.currentIndex+6] = byte(v >> 8)
	a.bytes[a.currentIndex+7] = byte(v)
	a.currentIndex += 8
}

// Write adds all the bytes to the payload
func (a *FixedWriter) Write(v []byte) {
	copy(a.bytes[a.currentIndex:], v)
	a.currentIndex += len(v)
}

// WriteString writes a sequence of characters (string) followed by a delimiter byte
func (a *FixedWriter) WriteString(value string, delim byte) {
	copy(a.bytes[a.currentIndex:], value)
	i := len(value)
	a.bytes[a.currentIndex+i] = delim
	a.currentIndex += 1 + i
}

// ExpandableWriter allows the array to grow past its capacity
var _ Writer = &ExpandableWriter{}

type ExpandableWriter struct {
	*bitWriter
}

// NewExpandableWriterWithCap makes a Writer interface with a initial soft capacity limit.
//	NOTE: See benchmarks in terms of comparison between FixedWriter & ExpandableWriter to determine which you want.
func NewExpandableWriterWithCap(size int) Writer {
	return &ExpandableWriter{
		&bitWriter{
			bytes: make([]byte, 0, size),
		},
	}
}

// NewExpandableWriter makes a Writer interface that will expand as more is written to it.
//	NOTE: See benchmarks in terms of comparison between FixedWriter & ExpandableWriter to determine which you want.
func NewExpandableWriter() Writer {
	return &ExpandableWriter{
		&bitWriter{
			bytes: make([]byte, 0),
		},
	}
}

// Writes a byte onto the buffer
func (a *ExpandableWriter) WriteUInt8(v uint8) {
	a.bytes = append(a.bytes, v)
	a.currentIndex++
}

// WriteUInt16 writes two bytes to the buffer
func (a *ExpandableWriter) WriteInt16(v int16) {
	a.bytes = append(a.bytes, byte(v>>8), byte(v))
	a.currentIndex += 2
}

// WriteUInt16 writes two bytes in little endian to the buffer
func (a *ExpandableWriter) WriteLEInt16(v int16) {
	a.bytes = append(a.bytes, byte(v), byte(v>>8))
	a.currentIndex += 2
}

// WriteInt32 writes a integer to the byte buffer
func (a *ExpandableWriter) WriteInt32(v int32) {
	a.bytes = append(a.bytes,
		byte(v>>24), byte(v>>16), byte(v>>8), byte(v),
	)
	a.currentIndex += 4
}

// WriteLEInt32 writes a integer to the byte buffer in little Endian
func (a *ExpandableWriter) WriteLEInt32(v int32) {
	a.bytes = append(a.bytes,
		byte(v), byte(v>>8), byte(v>>16), byte(v>>24),
	)
	a.currentIndex += 4
}

// WriteInt64 writes a int64 to the buffer in Big Endian order
func (a *ExpandableWriter) WriteInt64(v int64) {
	a.bytes = append(a.bytes,
		byte(v>>56), byte(v>>48), byte(v>>40),
		byte(v>>32), byte(v>>24), byte(v>>16),
		byte(v>>8), byte(v),
	)
	a.currentIndex += 8
}

// WriteLEInt64 writes a int64 to the buffer in Little Endian order
func (a *ExpandableWriter) WriteLEInt64(v int64) {
	a.bytes = append(a.bytes,
		byte(v), byte(v>>8), byte(v>>16), byte(v>>24),
		byte(v>>32), byte(v>>40), byte(v>>48),
		byte(v>>56),
	)
	a.currentIndex += 8
}

// Write adds all the bytes to the payload
func (a *ExpandableWriter) Write(v []byte) {
	a.bytes = append(a.bytes, v...)
	a.currentIndex += len(v)
}

// AppendString writes a sequence of characters (string) followed by a delimiter byte
//	NOTE: This will go beyond the size of the buffered Array. If not desired, use WriteString instead.
func (a *ExpandableWriter) WriteString(value string, delim byte) {
	a.bytes = append(a.bytes, []byte(value)...)
	a.bytes = append(a.bytes, delim)
}

func init() {
	for i := range bitMask {
		bitMask[i] = (1 << uint(i)) - 1
	}
}
