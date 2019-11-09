package bytepal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func testBitAccess(t *testing.T, out Writer, expected ...int) {
	fun := out.BitAccess()
	fun(1, 1)
	assert.Equal(t, expected[0], out.Size())
	fun(2, 1)
	fun(5, 1)
	assert.Equal(t, byte(161), out.Payload()[0])
	assert.Equal(t, expected[1], out.Size())
	fun(1, 1)
	assert.Equal(t, byte(128), out.Payload()[1])
	assert.Equal(t, expected[2], out.Size())
	fun(15, 3)
	assert.Equal(t, expected[3], out.Size())
}

func TestBitWriter_BitAccess2(t *testing.T) {
	out := NewExpandableWriterWithCap(0)
	testBitAccess(t, out, 1, 1, 2, 3)
	out = NewExpandableWriter()
	testBitAccess(t, out, 1, 1, 2, 3)
}
func TestBitWriter_BitAccess3(t *testing.T) {
	out := NewFixedWriter(0)
	testBitAccess(t, out, 1, 1, 2, 3)
	out = NewFixedWriter(3)
	testBitAccess(t, out, 3, 3, 3, 3)
}
func TestBitWriter_BitAccessAndWrite(t *testing.T) {
	out := NewExpandableWriterWithCap(0)
	bits := out.BitAccess()
	bits(1, 0)
	bits(8, 0)
	bits(11, 2047)
	assert.Len(t, out.Payload(), 3)
	out.Write([]byte{0})
	assert.Len(t, out.Payload(), 4)
}

func TestBitWriter_BitAccessAndWrite2(t *testing.T) {
	out := NewExpandableWriterWithCap(0)
	bits := out.BitAccess()
	bits(11, 1)
	bits(1, 1)
	bits(1, 1)
	v := -2
	bits(5, uint(v))
	bits(5, 2)
	assert.Equal(t, 3, out.Size())

	in := NewReader(out.Payload())
	b := in.ReadBits()
	assert.Equal(t, uint(1), b(11))
	assert.Equal(t, uint(1), b(1))
	assert.Equal(t, uint(1), b(1))
	assert.Equal(t, uint(30), b(5))
	assert.Equal(t, uint(2), b(5))
}

func TestFixedWriter_Write(t *testing.T) {
	out := NewFixedWriter(4)
	test := []byte{0, 1, 0, 1}
	out.Write(test)
	assert.Len(t, out.Payload(), len(test))

	for i := range test {
		assert.Equal(t, test[i], out.Payload()[i])
	}
}
func TestExpandableWriter_Write(t *testing.T) {
	out := NewExpandableWriterWithCap(4)
	test := []byte{0, 1, 0, 1}
	out.Write(test)

	assert.Len(t, out.Payload(), len(test))

	for i := range test {
		assert.Equal(t, test[i], out.Payload()[i])
	}
}
func TestExpandableWriter_Write2(t *testing.T) {
	out := NewExpandableWriter()
	test := []byte{0, 1, 0, 1}
	out.Write(test)

	assert.Len(t, out.Payload(), len(test))

	for i := range test {
		assert.Equal(t, test[i], out.Payload()[i])
	}
}
func BenchmarkFixedWriter_Write(b *testing.B) {
	test := []byte(TestString)
	out := NewFixedWriter(len(test) * b.N)
	for i := 0; i < b.N; i++ {
		out.Write(test)
	}
}
func BenchmarkExpandableWriter_Write(b *testing.B) {
	test := []byte(TestString)
	out := NewExpandableWriterWithCap(len(test) * b.N)
	for i := 0; i < b.N; i++ {
		out.Write(test)
	}
}
func BenchmarkExpandableWriter_Write2(b *testing.B) {
	test := []byte(TestString)
	out := NewExpandableWriter()
	for i := 0; i < b.N; i++ {
		out.Write(test)
	}
}

func TestFixedWriter_WriteUInt8(t *testing.T) {
	out := NewFixedWriter(1)
	out.WriteUInt8(1)
	assert.Equal(t, byte(1), out.Payload()[0])
}
func TestExpandableWriter_WriteUInt8(t *testing.T) {
	out := NewExpandableWriterWithCap(1)
	out.WriteUInt8(1)
	assert.Equal(t, byte(1), out.Payload()[0])
}
func TestExpandableWriter_WriteUInt82(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteUInt8(1)
	assert.Equal(t, byte(1), out.Payload()[0])
}
func BenchmarkFixedWriter_WriteUInt8(b *testing.B) {
	out := NewFixedWriter(b.N)
	for i := 0; i < b.N; i++ {
		out.WriteUInt8(1)
	}
}
func BenchmarkExpandableWriter_WriteUInt8(b *testing.B) {
	out := NewExpandableWriterWithCap(b.N)
	for i := 0; i < b.N; i++ {
		out.WriteUInt8(1)
	}
}
func BenchmarkExpandableWriter_WriteUInt82(b *testing.B) {
	out := NewExpandableWriter()
	for i := 0; i < b.N; i++ {
		out.WriteUInt8(1)
	}
}

func TestFixedWriter_WriteInt16(t *testing.T) {
	out := NewFixedWriter(2)
	out.WriteInt16(512)

	assert.Equal(t, byte(2), out.Payload()[0])
	assert.Equal(t, byte(0), out.Payload()[1])
}
func TestExpandableWriter_WriteInt16(t *testing.T) {
	out := NewExpandableWriterWithCap(2)
	out.WriteInt16(512)

	require.Len(t, out.Payload(), 2)
	assert.Equal(t, byte(2), out.Payload()[0])
	assert.Equal(t, byte(0), out.Payload()[1])
}
func TestExpandableWriter_WriteInt162(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteInt16(512)

	require.Len(t, out.Payload(), 2)
	assert.Equal(t, byte(2), out.Payload()[0])
	assert.Equal(t, byte(0), out.Payload()[1])
}
func BenchmarkFixedWriter_WriteInt16(b *testing.B) {
	out := NewFixedWriter(b.N * 2)
	for i := 0; i < b.N; i++ {
		out.WriteInt16(2)
	}
}
func BenchmarkExpandableWriter_WriteInt16(b *testing.B) {
	out := NewExpandableWriterWithCap(b.N * 2)
	for i := 0; i < b.N; i++ {
		out.WriteInt16(2)
	}
}
func BenchmarkExpandableWriter_WriteInt162(b *testing.B) {
	out := NewExpandableWriter()
	for i := 0; i < b.N; i++ {
		out.WriteInt16(2)
	}
}

func TestFixedWriter_WriteLEInt16(t *testing.T) {
	out := NewFixedWriter(2)
	out.WriteLEInt16(512)


	if out.Payload()[1] != 2 && out.Payload()[0] != 0 {
		fmt.Println(out.Payload())
		t.FailNow()
	}
}
func TestExpandableWriter_WriteLEInt16(t *testing.T) {
	out := NewExpandableWriterWithCap(2)
	out.WriteLEInt16(512)

	require.Len(t, out.Payload(), 2)
	assert.Equal(t, byte(0), out.Payload()[0])
	assert.Equal(t, byte(2), out.Payload()[1])
}
func TestExpandableWriter_WriteLEInt162(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteLEInt16(512)

	require.Len(t, out.Payload(), 2)
	assert.Equal(t, byte(0), out.Payload()[0])
	assert.Equal(t, byte(2), out.Payload()[1])
}
func BenchmarkFixedWriter_WriteLEInt16(b *testing.B) {
	out := NewFixedWriter(b.N * 2)
	for i := 0; i < b.N; i++ {
		out.WriteLEInt16(2)
	}
}
func BenchmarkExpandableWriter_WriteLeInt16(b *testing.B) {
	out := NewExpandableWriterWithCap(b.N * 2)
	for i := 0; i < b.N; i++ {
		out.WriteLEInt16(2)
	}
}
func BenchmarkExpandableWriter_WriteLEInt162(b *testing.B) {
	out := NewExpandableWriter()
	for i := 0; i < b.N; i++ {
		out.WriteLEInt16(2)
	}
}

func TestFixedWriter_WriteInt32(t *testing.T) {
	out := NewFixedWriter(4)
	out.WriteInt32(512)

	require.Len(t, out.Payload(), 4)
	assert.Equal(t, byte(0), out.Payload()[0])
	assert.Equal(t, byte(2), out.Payload()[2])
}
func TestExpandableWriter_WriteInt32(t *testing.T) {
	out := NewExpandableWriterWithCap(4)
	out.WriteInt32(512)

	require.Len(t, out.Payload(), 4)
	assert.Equal(t, byte(0), out.Payload()[0])
	assert.Equal(t, byte(2), out.Payload()[2])
}
func TestExpandableWriter_WriteInt322(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteInt32(512)

	require.Len(t, out.Payload(), 4)
	assert.Equal(t, byte(0), out.Payload()[0])
	assert.Equal(t, byte(2), out.Payload()[2])
}
func BenchmarkFixedWriter_WriteInt32(b *testing.B) {
	out := NewFixedWriter(b.N * 4)
	for i := 0; i < b.N; i++ {
		out.WriteInt32(2)
	}
}
func BenchmarkExpandableWriter_WriteInt32(b *testing.B) {
	out := NewExpandableWriterWithCap(b.N * 4)
	for i := 0; i < b.N; i++ {
		out.WriteInt32(2)
	}
}
func BenchmarkExpandableWriter_WriteInt322(b *testing.B) {
	out := NewExpandableWriter()
	for i := 0; i < b.N; i++ {
		out.WriteInt32(2)
	}
}

func TestFixedWriter_WriteLEInt32(t *testing.T) {
	out := NewFixedWriter(4)
	out.WriteLEInt32(512)

	require.Len(t, out.Payload(), 4)
	assert.Equal(t, byte(0), out.Payload()[0])
	assert.Equal(t, byte(2), out.Payload()[1])
}
func TestExpandableWriter_WriteLEInt32(t *testing.T) {
	out := NewExpandableWriterWithCap(4)
	out.WriteLEInt32(512)

	require.Len(t, out.Payload(), 4)
	assert.Equal(t, byte(0), out.Payload()[0])
	assert.Equal(t, byte(2), out.Payload()[1])
}
func TestExpandableWriter_WriteLEInt322(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteLEInt32(512)

	require.Len(t, out.Payload(), 4)
	assert.Equal(t, byte(0), out.Payload()[0])
	assert.Equal(t, byte(2), out.Payload()[1])
}
func BenchmarkFixedWriter_WriteLEInt32(b *testing.B) {
	out := NewFixedWriter(b.N * 4)
	for i := 0; i < b.N; i++ {
		out.WriteLEInt32(2)
	}
}
func BenchmarkExpandableWriter_WriteLEInt32(b *testing.B) {
	out := NewExpandableWriterWithCap(b.N * 4)
	for i := 0; i < b.N; i++ {
		out.WriteLEInt32(2)
	}
}
func BenchmarkExpandableWriter_WriteLEInt322(b *testing.B) {
	out := NewExpandableWriter()
	for i := 0; i < b.N; i++ {
		out.WriteLEInt32(2)
	}
}

func TestFixedWriter_WriteInt64(t *testing.T) {
	out := NewFixedWriter(8)
	out.WriteInt64(144115188075856000)

	assert.Equal(t, byte(2), out.Payload()[0])
	assert.Equal(t, byte(128), out.Payload()[7])
}
func TestExpandableWriter_WriteInt64(t *testing.T) {
	out := NewExpandableWriterWithCap(8)
	out.WriteInt64(144115188075856000)

	assert.Equal(t, byte(2), out.Payload()[0])
	assert.Equal(t, byte(128), out.Payload()[7])
}
func TestExpandableWriter_WriteInt642(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteInt64(144115188075856000)

	assert.Equal(t, byte(2), out.Payload()[0])
	assert.Equal(t, byte(128), out.Payload()[7])
}
func BenchmarkFixedWriter_WriteInt64(b *testing.B) {
	out := NewFixedWriter(b.N * 8)
	for i := 0; i < b.N; i++ {
		out.WriteInt64(144115188075856000)
	}
}
func BenchmarkExpandableWriter_WriteInt64(b *testing.B) {
	out := NewExpandableWriterWithCap(b.N * 8)
	for i := 0; i < b.N; i++ {
		out.WriteInt64(144115188075856000)
	}
}
func BenchmarkExpandableWriter_WriteInt642(b *testing.B) {
	out := NewExpandableWriter()
	for i := 0; i < b.N; i++ {
		out.WriteInt64(144115188075856000)
	}
}

func TestFixedWriter_WriteLELEInt64(t *testing.T) {
	out := NewFixedWriter(8)
	out.WriteLEInt64(144115188075856000)

	assert.Equal(t, byte(2), out.Payload()[7])
	assert.Equal(t, byte(128), out.Payload()[0])
}
func TestExpandableWriter_WriteLEInt64(t *testing.T) {
	out := NewExpandableWriterWithCap(8)
	out.WriteLEInt64(144115188075856000)

	assert.Equal(t, byte(2), out.Payload()[7])
	assert.Equal(t, byte(128), out.Payload()[0])
}
func TestExpandableWriter_WriteLEInt642(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteLEInt64(144115188075856000)

	assert.Equal(t, byte(2), out.Payload()[7])
	assert.Equal(t, byte(128), out.Payload()[0])
}
func BenchmarkFixedWriter_WriteLEInt64(b *testing.B) {
	out := NewFixedWriter(b.N * 8)
	for i := 0; i < b.N; i++ {
		out.WriteLEInt64(144115188075856000)
	}
}
func BenchmarkExpandableWriter_WriteLEInt64(b *testing.B) {
	out := NewExpandableWriterWithCap(b.N * 8)
	for i := 0; i < b.N; i++ {
		out.WriteLEInt64(144115188075856000)
	}
}
func BenchmarkExpandableWriter_WriteLEInt642(b *testing.B) {
	out := NewExpandableWriter()
	for i := 0; i < b.N; i++ {
		out.WriteLEInt64(144115188075856000)
	}
}

func TestFixedWriter_WriteString(t *testing.T) {
	out := NewFixedWriter(len(TestString) + 1)
	out.WriteString(TestString, Delim)
	require.Len(t, out.Payload(), len(TestString) + 1)
}
func TestExpandableWriter_WriteString(t *testing.T) {
	out := NewExpandableWriterWithCap(len(TestString) + 1)
	out.WriteString(TestString, Delim)
	require.Len(t, out.Payload(), len(TestString) + 1)

	out2 := NewExpandableWriterWithCap(len(TestString) + 1)
	out2.WriteString(TestString, Delim)
	require.Len(t, out2.Payload(), len(TestString) + 1)

	for i := range out.Payload() {
		assert.Equal(t, out.Payload()[i], out2.Payload()[i])
	}
}
func TestExpandableWriter_WriteString2(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteString(TestString, Delim)
	require.Len(t, out.Payload(), len(TestString) + 1)

	out2 := NewExpandableWriterWithCap(len(TestString) + 1)
	out2.WriteString(TestString, Delim)
	require.Len(t, out2.Payload(), len(TestString) + 1)

	for i := range out.Payload() {
		assert.Equal(t, out.Payload()[i], out2.Payload()[i])
	}
}
func BenchmarkFixedWriter_WriteString(b *testing.B) {
	size := len(TestString) + 1
	out := NewFixedWriter(size * b.N)
	for i := 0; i < b.N; i++ {
		out.WriteString(TestString, Delim)
	}
}
func BenchmarkExpandableWriter_WriteString(b *testing.B) {
	size := len(TestString) + 1
	out := NewExpandableWriterWithCap(size * b.N)
	for i := 0; i < b.N; i++ {
		out.WriteString(TestString, Delim)
	}
}
func BenchmarkExpandableWriter_WriteString2(b *testing.B) {
	out := NewExpandableWriter()
	for i := 0; i < b.N; i++ {
		out.WriteString(TestString, Delim)
	}
}
