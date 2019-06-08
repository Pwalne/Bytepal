package bytebuf

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
			a.bytes[bytePos] |= byte((value >> (numBits - bitOffset)) & bitMask[numBits])
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
	WriteUInt8A(uint8)
	WriteUInt16(uint16)
	WriteUInt16A(uint16)
	WriteLEUInt16(uint16)
	WriteLEUInt16A(uint16)
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
	//defer lastResult.Inc(1)
	a.bytes[a.currentIndex] = v
	a.currentIndex++
}

// Writes a byte onto the buffer but appends 128 to the value
func (a *FixedWriter) WriteUInt8A(v uint8) {
	//defer lastResult.Inc(1)
	a.bytes[a.currentIndex] = v + 128
	a.currentIndex++
}

// WriteUInt16 writes two bytes to the buffer
func (a *FixedWriter) WriteUInt16(v uint16) {
	//defer lastResult.Inc(1)
	a.bytes[a.currentIndex] = byte(v >> 8)
	a.bytes[a.currentIndex+1] = byte(v)
	a.currentIndex += 2
}

// WriteUInt16 writes two bytes to the buffer, but appends the second with 128
func (a *FixedWriter) WriteUInt16A(v uint16) {
	//defer lastResult.Inc(1)
	a.bytes[a.currentIndex] = byte(v >> 8)
	a.bytes[a.currentIndex+1] = byte(v + 128)
	a.currentIndex += 2
}

// WriteUInt16 writes two bytes in little endian to the buffer
func (a *FixedWriter) WriteLEUInt16(v uint16) {
	//defer lastResult.Inc(1)
	a.bytes[a.currentIndex+1] = byte(v >> 8)
	a.bytes[a.currentIndex] = byte(v)
	a.currentIndex += 2
}

// WriteUInt16 writes two bytes in little endian to the buffer, but appends the second with 128
func (a *FixedWriter) WriteLEUInt16A(v uint16) {
	//defer lastResult.Inc(1)
	a.bytes[a.currentIndex+1] = byte(v >> 8)
	a.bytes[a.currentIndex] = byte(v + 128)
	a.currentIndex += 2
}

// Write adds all the bytes to the payload
func (a *FixedWriter) Write(v []uint8) {
	copy(a.bytes[a.currentIndex:], v)
	a.currentIndex += len(v)
}

// WriteString writes a sequence of characters (string) followed by a delimiter byte
func (a *FixedWriter) WriteString(value string, delim byte) {
	copy(a.bytes[a.currentIndex:], value)
	i := len(value)
	a.bytes[a.currentIndex + i] = delim
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

// Writes a byte onto the buffer but appends 128 to the value
func (a *ExpandableWriter) WriteUInt8A(v uint8) {
	a.bytes = append(a.bytes, v+128)
	a.currentIndex++
}

// WriteUInt16 writes two bytes to the buffer
func (a *ExpandableWriter) WriteUInt16(v uint16) {
	a.bytes = append(a.bytes, byte(v>>8), byte(v))
	a.currentIndex += 2
}

// WriteUInt16 writes two bytes to the buffer, but appends the second with 128
func (a *ExpandableWriter) WriteUInt16A(v uint16) {
	a.bytes = append(a.bytes, byte(v>>8), byte(v+128))
	a.currentIndex += 2
}

// WriteUInt16 writes two bytes in little endian to the buffer
func (a *ExpandableWriter) WriteLEUInt16(v uint16) {
	a.bytes = append(a.bytes, byte(v), byte(v>>8))
	a.currentIndex += 2
}

// WriteUInt16 writes two bytes in little endian to the buffer, but appends the second with 128
func (a *ExpandableWriter) WriteLEUInt16A(v uint16) {
	a.bytes = append(a.bytes, byte(v+128), byte(v>>8))
	a.currentIndex += 2
}

// Write adds all the bytes to the payload
func (a *ExpandableWriter) Write(v []uint8) {
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
