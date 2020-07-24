package bytepal

import (
	"bytes"
	"encoding/binary"
	"io"
)

var bitMask [32]uint

// Wrapper that will read incremental bytes of an array into variables
type Reader struct {
	bytes        []byte
	currentIndex int
}

// Create a Reader from a existing byte array with endianess set to BigEndian.
func NewReader(bytes []byte) *Reader {
	return &Reader{
		bytes: bytes,
	}
}

// Reads bytes and creates a new reader.
func ReadBytes(reader io.Reader, size int) (*Reader, error) {
	bytes := make([]byte, size)
	_, err := reader.Read(bytes)
	if err != nil {
		return nil, err
	}
	return &Reader{
		bytes: bytes,
	}, nil
}

// Payload returns the underlying byte array
func (b *Reader) Payload() []byte {
	return b.bytes
}

// Seek sets the reading index to the desired position. NOTE: No out of bounds checks are performed.
func (b *Reader) Seek(position int) {
	b.currentIndex = position
}

// Reads a single byte off the array and increments the index pointer
func (b *Reader) ReadUInt8() uint8 {
	//defer lastResult.Inc(1)
	b.currentIndex++
	return b.bytes[b.currentIndex-1]
}

// ReadSlice reads the given number of bytes and returns a slice of the payload containing those
func (b *Reader) ReadSlice(size int) []byte {
	data := b.bytes[b.currentIndex:b.currentIndex+size]
	b.currentIndex += size
	return data
}

func (b *Reader) ReadBytes(payload []byte) {
	for i := range payload {
		payload[i] = b.bytes[b.currentIndex+i]
	}
	b.currentIndex += len(payload)
}

// Reads a twos byte off the array and increments the index pointer
func (b *Reader) ReadLEUInt16() uint16 {
	b.currentIndex += 2
	return binary.LittleEndian.Uint16(b.bytes[b.currentIndex-2 : b.currentIndex])
}

// Reads a twos byte off the array and increments the index pointer
func (b *Reader) ReadUInt16() uint16 {
	b.currentIndex += 2
	return binary.BigEndian.Uint16(b.bytes[b.currentIndex-2 : b.currentIndex])
}

// ReadUMedium reads a 24bit unsigned value
func (b *Reader) ReadUMedium() uint32 {
	b.currentIndex += 3
	return uint32(b.bytes[b.currentIndex-3]) << 16 | uint32(b.bytes[b.currentIndex-2]) << 8 | uint32(b.bytes[b.currentIndex-1])
}

// ReadBigSmart attempts to read either a short or int based on the next value.
func (b *Reader) ReadBigSmart() uint32 {
	if int8(b.bytes[b.currentIndex]) >= 0 {
		return uint32(b.ReadUInt16())
	}
	return b.ReadUInt32()
}

// Reads a twos byte off the array and increments the index pointer
func (b *Reader) ReadLEUInt32() uint32 {
	b.currentIndex += 4
	return binary.LittleEndian.Uint32(b.bytes[b.currentIndex-4 : b.currentIndex])
}

// Reads a twos byte off the array and increments the index pointer
func (b *Reader) ReadUInt32() uint32 {
	b.currentIndex += 4
	return binary.BigEndian.Uint32(b.bytes[b.currentIndex-4 : b.currentIndex])
}

// Reads a twos byte off the array and increments the index pointer
func (b *Reader) ReadLEUInt64() uint64 {
	b.currentIndex += 8
	return binary.LittleEndian.Uint64(b.bytes[b.currentIndex-8 : b.currentIndex])
}

// Reads a twos byte off the array and increments the index pointer
func (b *Reader) ReadUInt64() uint64 {
	b.currentIndex += 8
	return binary.BigEndian.Uint64(b.bytes[b.currentIndex-8 : b.currentIndex])
}

// Remaining bytes available to be read.
func (b *Reader) Remaining() int {
	return len(b.bytes) - b.currentIndex
}

// Continuously reads bytes until the deliminiter character is read.
func (b *Reader) ReadString(delim byte) string {
	if index := bytes.IndexByte(b.bytes[b.currentIndex:], delim); index > 0 {
		end := index + b.currentIndex
		data := b.bytes[b.currentIndex:end]
		b.currentIndex = end + 1
		return string(data)
	}
	return ""
}

// Increments the index pointer by the amount
func (b *Reader) Inc(amt int) {
	b.currentIndex += amt
}

// ReadBits returns a function that will continuously read bits off of the array.
func (b *Reader) ReadBits() func(uint) uint {
	bitPosition := uint(b.currentIndex * 8)
	return func(numBits uint) uint {
		bytePos := bitPosition >> 3
		bitOffset := 8 - (bitPosition & 7)
		bitPosition += numBits

		value := uint(0)
		for ; numBits > bitOffset; bitOffset = 8 {
			value += uint(b.bytes[bytePos]) & bitMask[bitOffset] << (numBits - bitOffset)
			bytePos++
			numBits -= bitOffset
		}
		if numBits == bitOffset {
			value += uint(b.bytes[bytePos]) & bitMask[bitOffset]
		} else {
			value += uint(b.bytes[bytePos]) >> (bitOffset - numBits) & bitMask[numBits]
		}
		b.currentIndex = int(bitPosition+7) / 8
		return value
	}
}
