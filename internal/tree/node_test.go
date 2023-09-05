package tree

import (
	"bytes"
	"testing"
)

func TestNodeSerializeDeserialize(t *testing.T) {
	node := Node{
		ID:   42,
		Word: [NodeWordSize]byte{1, 2, 3, 4},
		Seed: 12,
	}

	var buffer bytes.Buffer
	err := node.SerializeToBinaryStream(&buffer)
	if err != nil {
		t.Fatalf("Failed to serialize node: %v", err)
	}

	var node2 Node
	node2, err = node2.DeserializeFromBinaryStream(&buffer)

	if err != nil {
		t.Fatalf("Failed to deserialize node: %v", err)
	}

	if node2.ID != node.ID {
		t.Fatalf("ID did not match. Got %d, expected %d", node2.ID, node.ID)
	}

	if node2.Word != node.Word {
		t.Fatalf("Word did not match. Got %d, expected %d", node2.Word, node.Word)
	}

	if node2.Seed != node.Seed {
		t.Fatalf("Seed did not match. Got %d, expected %d", node2.Seed, node.Seed)
	}
}

func TestNodeSerializedSize(t *testing.T) {
	node := Node{
		ID:   42,
		Word: [NodeWordSize]byte{1, 2, 3, 4},
		Seed: 12,
	}

	size := node.StrideLength()

	var buffer bytes.Buffer
	err := node.SerializeToBinaryStream(&buffer)
	if err != nil {
		t.Fatalf("Failed to serialize node: %v", err)
	}

	if size != uint64(len(buffer.Bytes())) {
		t.Fatalf("StrideLength() did not return the correct size. Got %d, expected %d", size, len(buffer.Bytes()))
	}

	println("Node serialized size:", size)
}

func TestNewNode(t *testing.T) {
	node := NewNode(42, [NodeWordSize]byte{1, 2, 3, 4}, 12)

	if node.ID != 42 {
		t.Fatalf("ID did not match. Got %d, expected %d", node.ID, 42)
	}

	if node.Word != [NodeWordSize]byte{1, 2, 3, 4} {
		t.Fatalf("Word did not match. Got %d, expected %d", node.Word, [NodeWordSize]byte{1, 2, 3, 4})
	}

	if node.Seed != 12 {
		t.Fatalf("Seed did not match. Got %d, expected %d", node.Seed, 12)
	}
}
