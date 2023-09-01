package serialization

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestNumberSerializationAndDeserialization(t *testing.T) {
	testNumber := Number(42)

	var buf bytes.Buffer
	err := testNumber.SerializeToBinaryStream(&buf)
	if err != nil {
		t.Fatalf("Error serializing: %v", err)
	}

	var deserializedNumber Num
	err = deserializedNumber.DeserializeFromBinaryStream(&buf)
	if err != nil {
		t.Fatalf("Error deserializing: %v", err)
	}

	if deserializedNumber != testNumber {
		t.Fatalf("Deserialized number %d doesn't match original number %d", deserializedNumber, testNumber)
	}
}

func TestNumberSerializedSize(t *testing.T) {
	testNumber := Number(42)
	expectedSize := uint64(binary.Size(testNumber))

	serializedSize := testNumber.SerializedSize()

	if serializedSize != expectedSize {
		t.Fatalf("Serialized size mismatch: got %d, expected %d", serializedSize, expectedSize)
	}
}
