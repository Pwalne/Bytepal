package bytebuf

import (
	"fmt"
	"testing"
)

const (
	TestString = "hello2umydarling"
	Delim      = 0
)

func TestReadString(t *testing.T) {
	data := append([]byte(TestString), Delim)
	reader := NewReader(data)
	if value := reader.ReadString(Delim); value != TestString {
		fmt.Printf("Did not decode string correctly, got %s\n", value)
		t.FailNow()
	}
}
func TestReadString2(t *testing.T) {
	data := append([]byte(TestString), Delim)
	data = append(data, data...)
	reader := NewReader(data)
	if value := reader.ReadString(Delim); value != TestString {
		fmt.Printf("Did not decode string correctly, got %s\n", value)
		t.FailNow()
	}
	if value := reader.ReadString(Delim); value != TestString {
		fmt.Printf("Did not decode second string correctly, got %s\n", value)
		t.FailNow()
	}
}
func BenchmarkReadBuffer_ReadString(b *testing.B) {
	data := append([]byte(TestString), Delim)
	reader := NewReader(data)
	for i := 0; i < b.N; i++ {
		_ = reader.ReadString(Delim)
		reader.currentIndex = 0
	}
}

func TestReadBuffer_ReadByteLE(t *testing.T) {
	reader := NewReader([]byte{0x5, 0x8})
	if reader.ReadUInt8() != 0x5 {
		t.FailNow()
	}
	if reader.ReadUInt8() != 0x8 {
		t.FailNow()
	}
}
func BenchmarkReadBuffer_ReadByteLE(b *testing.B) {
	reader := NewReader([]byte(TestString))
	for i := 0; i < b.N; i++ {
		_ = reader.ReadUInt8()
		reader.currentIndex = 0
	}
}

func TestReadBuffer_ReadInt16LE(t *testing.T) {
	reader := NewReader([]byte{0x12, 0x2E})
	if v := reader.ReadLEUInt16(); v != 11794 {
		fmt.Printf("Value: %d\n", v)
		t.FailNow()
	}
}
func BenchmarkReadBuffer_ReadInt16LE(b *testing.B) {
	reader := NewReader([]byte(TestString))
	for i := 0; i < b.N; i++ {
		_ = reader.ReadLEUInt16()
		reader.currentIndex = 0
	}
}

func TestReadBuffer_ReadIntLE(t *testing.T) {
	reader := NewReader([]byte{0xDB, 0xAC, 0xCD, 0xB})
	if v := reader.ReadLEUInt32(); v != 198028507 {
		fmt.Printf("Value: %d\n", v)
		t.FailNow()
	}
}
func BenchmarkReadBuffer_ReadIntLE(b *testing.B) {
	reader := NewReader([]byte(TestString))
	for i := 0; i < b.N; i++ {
		_ = reader.ReadLEUInt32()
		reader.currentIndex = 0
	}
}

func TestReadBuffer_ReadInt64LE(t *testing.T) {
	reader := NewReader([]byte{0xCA, 0xBD, 0xFD, 0xCF, 0xCA, 0xBD, 0xCA, 0x3D})
	if v := reader.ReadLEUInt64(); v != 4452579860379712970 {
		fmt.Printf("Value: %d\n", v)
		t.FailNow()
	}
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
	if a := r(1); a != 1 {
		fmt.Println("bit 1: ", a)
		t.Fail()
	}
	if a := r(2); a != 1 {
		fmt.Println("Bit 3", a)
		t.Fail()
	}
	if a := r(5); a != 1 {
		fmt.Println("Bit 8:", a)
		t.Fail()
	}
	if a := r(1); a != 1 {
		fmt.Println("bit 9: ", a)
		t.Fail()
	}
	if a := r(15); a != 3 {
		fmt.Println("bit 24: ", a)
		t.Fail()
	}
}
