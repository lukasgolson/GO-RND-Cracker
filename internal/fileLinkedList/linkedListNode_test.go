package fileLinkedList

import (
	"awesomeProject/internal/number"
	"awesomeProject/internal/serialization"
	"testing"
)

func TestNodeSerializeDeserialize(t *testing.T) {
	node1 := linkedListNode[number.Number]{NextOffset: 43, Item: *new(number.Number)}

	buffer := make([]byte, node1.StrideLength())
	err := node1.SerializeToBinaryStream(buffer)
	if err != nil {
		t.Fatalf("Failed to serialize edge: %v", err)
	}

	var node2 linkedListNode[number.Number]
	node2, err = node2.DeserializeFromBinaryStream(buffer)

	if err != nil {
		t.Fatalf("Failed to deserialize edge: %v", err)
	}

	if node2.NextOffset != node1.NextOffset {
		t.Fatalf("NextOffset did not match. Got %d, expected %d", node2.NextOffset, node1.NextOffset)
	}

	if node2.Item != node1.Item {
		t.Fatalf("Item did not match. Got %d, expected %d", node2.Item, node1.Item)
	}

}

func TestEdgeSerializedSize(t *testing.T) {
	node := linkedListNode[number.Number]{NextOffset: 43, Item: *new(number.Number)}

	size := node.StrideLength()

	buffer := make([]byte, node.StrideLength())
	err := node.SerializeToBinaryStream(buffer)
	if err != nil {
		t.Fatalf("Failed to serialize node: %v", err)
	}

	if size != serialization.Length(len(buffer)) {
		t.Fatalf("StrideLength() did not return the correct size. Got %d, expected %d", size, len(buffer))
	}
}
