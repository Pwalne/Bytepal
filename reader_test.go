package bytebuf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TestString = "hello2umydarling"
	Delim      = 0
)

func TestReadString(t *testing.T) {
	data := append([]byte(TestString), Delim)
	reader := NewReader(data)
	assert.Equal(t, TestString, reader.ReadString(Delim))
}
func TestReadString2(t *testing.T) {
	data := append([]byte(TestString), Delim)
	data = append(data, data...)
	reader := NewReader(data)
	assert.Equal(t, TestString, reader.ReadString(Delim))
	assert.Equal(t, TestString, reader.ReadString(Delim))
}
func BenchmarkReadBuffer_ReadString(b *testing.B) {
	data := make([]byte, 0)
	for n := 0; n < b.N; n++ {
		data = append(data, []byte(TestString)...)
		data = append(data, Delim)
	}
	reader := NewReader(data)
	for i := 0; i < b.N; i++ {
		_ = reader.ReadString(Delim)
	}
}

func TestReadBuffer_ReadByte(t *testing.T) {
	reader := NewReader([]byte{0x5, 0x8})
	assert.Equal(t, uint8(0x5), reader.ReadUInt8())
	assert.Equal(t, uint8(0x8), reader.ReadUInt8())
}
func BenchmarkReadBuffer_ReadByte(b *testing.B) {
	data := make([]byte,b.N)
	for n := 0; n < b.N; n++ {
		data = append(data, 5)
	}
	reader := NewReader(data)
	for i := 0; i < b.N; i++ {
		_ = reader.ReadUInt8()
	}
}

func TestReadBuffer_ReadInt16LE(t *testing.T) {
	reader := NewReader([]byte{0x12, 0x2E})
	assert.Equal(t, uint16(11794), reader.ReadLEUInt16())
}
func BenchmarkReadBuffer_ReadInt16LE(b *testing.B) {
	data := make([]byte, 4 * b.N)
	for n := 0; n < b.N; n++ {
		data = append(data, 5, 6, 10, 5)
	}
	reader := NewReader(data)
	for i := 0; i < b.N; i++ {
		_ = reader.ReadLEUInt16()
	}
}

func TestReadBuffer_ReadIntLE(t *testing.T) {
	reader := NewReader([]byte{0xDB, 0xAC, 0xCD, 0xB})
	assert.Equal(t, uint32(198028507), reader.ReadLEUInt32())
}
func BenchmarkReadBuffer_ReadIntLE(b *testing.B) {
	data := make([]byte, 4 * b.N)
	for n := 0; n < b.N; n++ {
		data = append(data, 5, 6, 10, 5)
	}
	reader := NewReader(data)
	for i := 0; i < b.N; i++ {
		_ = reader.ReadLEUInt32()
	}
}

func TestReadBuffer_ReadInt64LE(t *testing.T) {
	reader := NewReader([]byte{0xCA, 0xBD, 0xFD, 0xCF, 0xCA, 0xBD, 0xCA, 0x3D})
	assert.Equal(t, uint64(4452579860379712970), reader.ReadLEUInt64())
}
func BenchmarkReadBuffer_ReadInt64LE(b *testing.B) {
	reader := NewReader([]byte(TestString))
	for i := 0; i < b.N; i++ {
		_ = reader.ReadLEUInt64()
		reader.currentIndex = 0
	}
}

func TestArrayReader_ReadBits(t *testing.T) {
	data := []byte{
		161, 128, 3, 0, 0,
	}
	in := NewReader(data)
	r := in.ReadBits()
	assert.Equal(t, uint(1), r(1))
	assert.Equal(t, uint(1), r(2))
	assert.Equal(t, uint(1), r(5))
	assert.Equal(t, uint(1), r(1))
	assert.Equal(t, uint(3), r(15))
}
