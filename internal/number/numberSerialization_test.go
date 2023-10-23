package number

import (
	"bytes"
	"encoding/binary"
	"github.com/lukasgolson/FileArray/serialization"
	"testing"
)

func TestNumber_SerializeToBinaryStream(t *testing.T) {

	testNumber := NewNumber(42)

	buf := make([]byte, testNumber.StrideLength())

	err := testNumber.SerializeToBinaryStream(buf)
	if err != nil {
		t.Fatalf("Error serializing: %v", err)
	}

	expectedBytes := []byte{42, 0, 0, 0, 0, 0, 0, 0}

	if !bytes.Equal(buf, expectedBytes) {
		t.Fatalf("Serialized bytes don't match expected bytes. Got %v, expected %v", buf, expectedBytes)
	}
}

func TestNumber_DeserializeFromBinaryStream(t *testing.T) {
	testNumber := NewNumber(42)

	buf := make([]byte, testNumber.StrideLength())
	err := testNumber.SerializeToBinaryStream(buf)
	if err != nil {
		t.Fatalf("Error serializing: %v", err)
	}

	deserializedNumber := NewNumber(28)
	deserializedNumber, err = deserializedNumber.DeserializeFromBinaryStream(buf)
	if err != nil {
		t.Fatalf("Error deserializing: %v", err)
	}

	if deserializedNumber != testNumber {
		t.Fatalf("Deserialized number %d doesn't match original number %d", deserializedNumber, testNumber)
	}
}

func TestNumberSerializationAndDeserialization(t *testing.T) {
	testNumber := NewNumber(42)

	buf := make([]byte, testNumber.StrideLength())
	err := testNumber.SerializeToBinaryStream(buf)
	if err != nil {
		t.Fatalf("Error serializing: %v", err)
	}

	var deserializedNumber Number
	deserializedNumber, err = deserializedNumber.DeserializeFromBinaryStream(buf)
	if err != nil {
		t.Fatalf("Error deserializing: %v", err)
	}

	if deserializedNumber != testNumber {
		t.Fatalf("Deserialized number %d doesn't match original number %d", deserializedNumber, testNumber)
	}
}

func TestNumberSerializedSize(t *testing.T) {
	testNumber := NewNumber(42)
	expectedSize := serialization.Length(binary.Size(testNumber))

	serializedSize := testNumber.StrideLength()

	if serializedSize != expectedSize {
		t.Fatalf("Serialized size mismatch: got %d, expected %d", serializedSize, expectedSize)
	}
}
