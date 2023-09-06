package tree

import (
	"awesomeProject/internal/serialization"
	"encoding/binary"
	"io"
)

const NodeWordSize = 32

type Node struct {
	ID   serialization.Offset
	Word [NodeWordSize]byte
	Seed int32
}

func NewNode(ID serialization.Offset, word [NodeWordSize]byte, seed int32) *Node {
	return &Node{
		ID:   ID,
		Word: word,
		Seed: seed,
	}
}

func (n Node) SerializeToBinaryStream(writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, n.ID)
	if err != nil {
		return err
	}

	err = binary.Write(writer, binary.LittleEndian, n.Word[:])
	if err != nil {
		return err
	}

	err = binary.Write(writer, binary.LittleEndian, n.Seed)
	if err != nil {
		return err
	}

	return nil
}

func (n Node) DeserializeFromBinaryStream(reader io.Reader) (Node, error) {
	var ID serialization.Offset
	err := binary.Read(reader, binary.LittleEndian, &ID)
	if err != nil {
		return n, err
	}

	var word [NodeWordSize]byte
	err = binary.Read(reader, binary.LittleEndian, &word)
	if err != nil {
		return n, err
	}

	var seed int32
	err = binary.Read(reader, binary.LittleEndian, &seed)
	if err != nil {
		return n, err
	}

	n.ID = ID
	n.Word = word
	n.Seed = seed

	return n, nil
}

func (n Node) StrideLength() serialization.Length {
	return serialization.Length(binary.Size(n.ID) + len(n.Word) + binary.Size(n.Seed))
}

func (n Node) IDByte() byte {
	return 'N'
}
