package tree

import (
	"encoding/binary"
	"io"
)

const FixedWordSize = 32

type Node struct {
	ID   uint32
	Word [FixedWordSize]byte
	Seed int32
}

func NewNode(ID uint32, word [FixedWordSize]byte, seed int32) *Node {
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
	var ID uint32
	err := binary.Read(reader, binary.LittleEndian, &ID)
	if err != nil {
		return n, err
	}

	var word [FixedWordSize]byte
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

func (n Node) StrideLength() uint64 {
	return uint64(binary.Size(n.ID) + len(n.Word) + binary.Size(n.Seed))
}

func (n Node) IDByte() []byte {
	return []byte{'N'}
}
