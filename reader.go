package bytebuf

import (
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

// Reads a single byte off the array and increments the index pointer
func (b *Reader) ReadUInt8() uint8 {
	//defer lastResult.Inc(1)
	b.currentIndex++
	return b.bytes[b.currentIndex-1]
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
	cur := b.currentIndex
	for ; b.currentIndex < len(b.bytes); b.currentIndex++ {
		if b.bytes[b.currentIndex] == delim {
			b.currentIndex++
			return string(b.bytes[cur : b.currentIndex-1])
		}
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