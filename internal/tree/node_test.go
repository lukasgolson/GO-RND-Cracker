package tree

import (
	"awesomeProject/internal/serialization"
	"testing"
)

func TestNodeSerializeDeserialize(t *testing.T) {
	node1 := node{
		ID:   42,
		Word: [NodeWordSize]byte{1, 2, 3, 4},
		Seed: 12,
	}

	buffer := make([]byte, node1.StrideLength())
	err := node1.SerializeToBinaryStream(buffer)
	if err != nil {
		t.Fatalf("Failed to serialize node: %v", err)
	}

	var node2 node
	node2, err = node2.DeserializeFromBinaryStream(buffer)

	if err != nil {
		t.Fatalf("Failed to deserialize node: %v", err)
	}

	if node2.ID != node1.ID {
		t.Fatalf("ID did not match. Got %d, expected %d", node2.ID, node1.ID)
	}

	if node2.Word != node1.Word {
		t.Fatalf("Word did not match. Got %d, expected %d", node2.Word, node1.Word)
	}

	if node2.Seed != node1.Seed {
		t.Fatalf("Seed did not match. Got %d, expected %d", node2.Seed, node1.Seed)
	}
}

func TestNodeSerializedSize(t *testing.T) {
	node := node{
		ID:   42,
		Word: [NodeWordSize]byte{1, 2, 3, 4},
		Seed: 12,
	}

	size := node.StrideLength()

	buffer := make([]byte, node.StrideLength())
	err := node.SerializeToBinaryStream(buffer)
	if err != nil {
		t.Fatalf("Failed to serialize node: %v", err)
	}

	if size != serialization.Length(len(buffer)) {
		t.Fatalf("StrideLength() did not return the correct size. Got %d, expected %d", size, len(buffer))
	}

	println("node serialized size:", size)
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
