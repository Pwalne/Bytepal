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

func TestFixedWriter_WriteUInt8A(t *testing.T) {
	out := NewFixedWriter(1)
	out.WriteUInt8A(2)

	if out.Payload()[0] != 130 {
		t.FailNow()
	}
}
func TestExpandableWriter_WriteUInt8A(t *testing.T) {
	out := NewExpandableWriterWithCap(1)
	out.WriteUInt8A(2)

	if out.Size() != 1 {
		t.FailNow()
	}
	if out.Payload()[0] != 130 {
		t.FailNow()
	}
}
func TestExpandableWriter_WriteUInt8A2(t *testing.T) {
	out := NewExpandableWriter()
	out.WriteUInt8A(2)

	if out.Size() != 1 {
		t.FailNow()
	}
	if out.Payload()[0] != 130 {
		t.FailNow()
	}
}
func BenchmarkFixedWriter_WriteUInt8A(b *testing.B) {
	out := NewFixedWriter(b.N)
	for i := 0; i < b.N; i++ {
		out.WriteUInt8A(2)
	}
}
func BenchmarkExpandableWriter_WriteUInt8A(b *testing.B) {
	out := NewExpandableWriterWithCap(b.N)
	for i := 0; i < b.N; i++ {
		out.WriteUInt8A(2)
	}
}
func BenchmarkExpandableWriter_WriteUInt8A2(b *testing.B) {
	out := NewExpandableWriter()
	for i := 0; i < b.N; i++ {
		out.WriteUInt8A(2)
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
