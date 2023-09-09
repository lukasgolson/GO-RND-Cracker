package tree

import (
	"awesomeProject/internal/serialization"
	"bytes"
	"testing"
)

func TestEdgeSerializeDeserialize(t *testing.T) {
	edge1 := edge{
		ChildIndex: 43,
		Distance:   12,
	}

	var buffer bytes.Buffer
	err := edge1.SerializeToBinaryStream(&buffer)
	if err != nil {
		t.Fatalf("Failed to serialize edge: %v", err)
	}

	var edge2 edge
	edge2, err = edge2.DeserializeFromBinaryStream(&buffer)

	if err != nil {
		t.Fatalf("Failed to deserialize edge: %v", err)
	}

	if edge2.ChildIndex != edge1.ChildIndex {
		t.Fatalf("ChildIndex did not match. Got %d, expected %d", edge2.ChildIndex, edge1.ChildIndex)
	}

	if edge2.Distance != edge1.Distance {
		t.Fatalf("Distance did not match. Got %d, expected %d", edge2.Distance, edge1.Distance)
	}

}

func TestEdgeSerializedSize(t *testing.T) {
	edge := edge{
		ChildIndex: 43,
		Distance:   12,
	}

	size := edge.StrideLength()

	var buffer bytes.Buffer
	err := edge.SerializeToBinaryStream(&buffer)
	if err != nil {
		t.Fatalf("Failed to serialize edge: %v", err)
	}

	if size != serialization.Length(len(buffer.Bytes())) {
		t.Fatalf("StrideLength() did not return the correct size. Got %d, expected %d", size, len(buffer.Bytes()))
	}
}

func TestNewEdge(t *testing.T) {
	edge := NewEdge(42, 43, 12)

	if edge.ChildIndex != 43 {
		t.Fatalf("ChildIndex did not match. Got %d, expected %d", edge.ChildIndex, 43)
	}

	if edge.Distance != 12 {
		t.Fatalf("Distance did not match. Got %d, expected %d", edge.Distance, 12)
	}
}
