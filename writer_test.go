package bytebuf

import (
	"fmt"
	"testing"
)

func testBitAccess(t *testing.T, out Writer, expected ...int) {
	fun := out.BitAccess()
	fun(1, 1)
	if out.Size() != expected[0] {
		fmt.Printf("Incorrect payload size. Should be %d, got %d\n", expected[0], out.Size())
		t.Fail()
	}
	fun(2, 1)
	fun(5, 1)
	if (out.Payload())[0] != 161 {
		t.Fail()
	}
	if out.Size() != expected[1] {
		fmt.Printf("Incorrect payload size. Should be %d, got %d\n", expected[1], out.Size())
		t.Fail()
	}
	fun(1, 1)
	if (out.Payload())[1] != 128 {
		t.Fail()
	}
	if out.Size() != expected[2] {
		fmt.Printf("Incorrect payload size. Should be %d, got %d\n", expected[2], out.Size())
		t.Fail()
	}
	fun(15, 3)
	if out.Size() != expected[3] {
		fmt.Printf("Incorrect payload size. Should be %d, got %d\n", expected[3], out.Size())
		t.Fail()
	}
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
	if out.Size() != 3 {
		fmt.Println("Incorrect buffer size of", out.Size())
		t.FailNow()
	}
	out.Write([]byte{0})
	if out.Size() != 4 {
		fmt.Println("Incorrect buffer size of", out.Size())
		t.FailNow()
	}
}

func TestFixedWriter_Write(t *testing.T) {
	out := NewFixedWriter(4)
	test := []byte{0, 1, 0, 1}
	out.Write(test)

	if out.Size() != len(test) {
		fmt.Printf("Incorrect payload size: Got %d, expected %d\n", out.Size(), len(test))
	}

	for i := range test {
		if test[i] != out.Payload()[i] {
			fmt.Printf("Got %d Expected %d.\n", out.Payload()[i], test[i])
			t.FailNow()
		}
	}
}
func TestExpandableWriter_Write(t *testing.T) {
	out := NewExpandableWriterWithCap(4)
	test := []byte{0, 1, 0, 1}
	out.Write(test)

	if out.Size() != len(test) {
		fmt.Printf("Incorrect payload size: Got %d, expected %d\n", out.Size(), len(test))
	}

	for i := range test {
		if test[i] != out.Payload()[i] {
			fmt.Printf("Got %d Expected %d.\n", out.Payload()[i], test[i])
			t.FailNow()
		}
	}
}
func TestExpandableWriter_Write2(t *testing.T) {
	out := NewExpandableWriter()
	test := []byte{0, 1, 0, 1}
	out.Write(test)

	if out.Size() != len(test) {
		fmt.Printf("Incorrect payload size: Got %d, expected %d\n", out.Size(), len(test))
	}

	for i := range test {
		if test[i] != out.Payload()[i] {
			fmt.Printf("Got %d Expected %d.\n", out.Payload()[i], test[i])
			t.FailNow()
		}
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
	if out.Payload()[0] != 1 {
		t.FailNow()
	}
}
func TestExpandableWriter_WriteUInt8(t *testing.T) {
	out := NewExpandableWriterWithCap(1)
	out.WriteUInt8(1)
	if out.Payload()[0] != 1 {
		t.FailNow()
	}
}
func TestExpandableWriter_WriteUInt82(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteUInt8(1)
	if out.Payload()[0] != 1 {
		t.FailNow()
	}
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

	if out.Payload()[0] != 2 && out.Payload()[1] != 0{
		fmt.Println(out.Payload())
		t.FailNow()
	}
}
func TestExpandableWriter_WriteInt16(t *testing.T) {
	out := NewExpandableWriterWithCap(2)
	out.WriteInt16(512)

	if out.Size() != 2 {
		t.FailNow()
	}
	if out.Payload()[0] != 2 && out.Payload()[1] != 0{
		fmt.Println(out.Payload())
		t.FailNow()
	}
}
func TestExpandableWriter_WriteInt162(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteInt16(512)

	if out.Size() != 2 {
		t.FailNow()
	}
	if out.Payload()[0] != 2 && out.Payload()[1] != 0{
		fmt.Println(out.Payload())
		t.FailNow()
	}
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

	if out.Payload()[1] != 2 && out.Payload()[0] != 0{
		fmt.Println(out.Payload())
		t.FailNow()
	}
}
func TestExpandableWriter_WriteLEInt16(t *testing.T) {
	out := NewExpandableWriterWithCap(2)
	out.WriteLEInt16(512)

	if out.Size() != 2 {
		t.FailNow()
	}
	if out.Payload()[1] != 2 && out.Payload()[0] != 0{
		fmt.Println(out.Payload())
		t.FailNow()
	}
}
func TestExpandableWriter_WriteLEInt162(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteLEInt16(512)

	if out.Size() != 2 {
		t.FailNow()
	}
	if out.Payload()[1] != 2 && out.Payload()[0] != 0{
		fmt.Println(out.Payload())
		t.FailNow()
	}
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

	if out.Payload()[1] != 2 && out.Payload()[0] != 0{
		fmt.Println(out.Payload())
		t.FailNow()
	}
}
func TestExpandableWriter_WriteInt32(t *testing.T) {
	out := NewExpandableWriterWithCap(4)
	out.WriteInt32(512)

	if out.Size() != 4 {
		t.FailNow()
	}
	if out.Payload()[1] != 2 && out.Payload()[0] != 0{
		fmt.Println(out.Payload())
		t.FailNow()
	}
}
func TestExpandableWriter_WriteInt322(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteInt32(512)

	if out.Size() != 4 {
		t.FailNow()
	}
	if out.Payload()[1] != 2 && out.Payload()[0] != 0{
		fmt.Println(out.Payload())
		t.FailNow()
	}
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

	if out.Payload()[1] != 2 && out.Payload()[0] != 0{
		fmt.Println(out.Payload())
		t.FailNow()
	}
}
func TestExpandableWriter_WriteLEInt32(t *testing.T) {
	out := NewExpandableWriterWithCap(4)
	out.WriteLEInt32(512)

	if out.Size() != 4 {
		t.FailNow()
	}
	if out.Payload()[1] != 2 && out.Payload()[0] != 0{
		fmt.Println(out.Payload())
		t.FailNow()
	}
}
func TestExpandableWriter_WriteLEInt322(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteLEInt32(512)

	if out.Size() != 4 {
		t.FailNow()
	}
	if out.Payload()[1] != 2 && out.Payload()[0] != 0{
		fmt.Println(out.Payload())
		t.FailNow()
	}
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

	if out.Payload()[0] != 2 && out.Payload()[7] != 128{
		fmt.Println(out.Payload())
		t.FailNow()
	}
}
func TestExpandableWriter_WriteInt64(t *testing.T) {
	out := NewExpandableWriterWithCap(8)
	out.WriteInt64(144115188075856000)

	if out.Payload()[0] != 2 && out.Payload()[7] != 128{
		fmt.Println(out.Payload())
		t.FailNow()
	}
}
func TestExpandableWriter_WriteInt642(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteInt64(144115188075856000)

	if out.Payload()[0] != 2 && out.Payload()[7] != 128{
		fmt.Println(out.Payload())
		t.FailNow()
	}
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

	if out.Payload()[7] != 2 && out.Payload()[0] != 128{
		fmt.Println(out.Payload())
		t.FailNow()
	}
}
func TestExpandableWriter_WriteLEInt64(t *testing.T) {
	out := NewExpandableWriterWithCap(8)
	out.WriteLEInt64(144115188075856000)

	if out.Payload()[7] != 2 && out.Payload()[0] != 128{
		fmt.Println(out.Payload())
		t.FailNow()
	}
}
func TestExpandableWriter_WriteLEInt642(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteLEInt64(144115188075856000)

	if out.Payload()[7] != 2 && out.Payload()[0] != 128{
		fmt.Println(out.Payload())
		t.FailNow()
	}
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
	if out.Size() != len(TestString)+1 {
		fmt.Printf("Size: %d Expected :%d\n", out.Size(), len(TestString)+1)
		t.Fail()
	}
}
func TestExpandableWriter_WriteString(t *testing.T) {
	out := NewExpandableWriterWithCap(len(TestString) + 1)
	out.WriteString(TestString, Delim)
	if out.Size() != len(TestString)+1 {
		fmt.Printf("Size: %d Expected :%d\n", out.Size(), len(TestString)+1)
		t.Fail()
	}

	out2 := NewExpandableWriterWithCap(len(TestString) + 1)
	out2.WriteString(TestString, Delim)
	if out2.Size() != len(TestString)+1 {
		fmt.Printf("Size: %d Expected :%d\n", out2.Size(), len(TestString)+1)
		t.Fail()
	}

	for i := range out.Payload() {
		if out.Payload()[i] != out2.Payload()[i] {
			fmt.Printf("Array [%s] != [%s]", string(out.Payload()), string(out2.Payload()))
			t.FailNow()
		}
	}
}
func TestExpandableWriter_WriteString2(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteString(TestString, Delim)
	if out.Size() != len(TestString)+1 {
		fmt.Printf("Size: %d Expected :%d\n", out.Size(), len(TestString)+1)
		t.Fail()
	}

	out2 := NewExpandableWriterWithCap(len(TestString) + 1)
	out2.WriteString(TestString, Delim)
	if out2.Size() != len(TestString)+1 {
		fmt.Printf("Size: %d Expected :%d\n", out2.Size(), len(TestString)+1)
		t.Fail()
	}

	for i := range out.Payload() {
		if out.Payload()[i] != out2.Payload()[i] {
			fmt.Printf("Array [%s] != [%s]", string(out.Payload()), string(out2.Payload()))
			t.FailNow()
		}
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
