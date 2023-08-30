package bktree

import (
	"encoding/binary"
	"io"
)

type Node struct {
	ID   uint32
	Word []byte
	Seed int32
}

func NewNode(ID uint32, word []byte, seed int32) *Node {
	return &Node{
		ID:   ID,
		Word: word,
		Seed: seed,
	}
}

func (node *Node) SerializeToBinaryStream(writer io.Writer) error {
	idBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(idBytes, uint32(node.ID))
	_, err := writer.Write(idBytes)
	if err != nil {
		return err
	}

	wordLenBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(wordLenBytes, uint16(len(node.Word)))
	_, err = writer.Write(wordLenBytes)
	if err != nil {
		return err
	}
	_, err = writer.Write(node.Word)
	if err != nil {
		return err
	}

	seedBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(seedBytes, uint32(node.Seed))
	_, err = writer.Write(seedBytes)
	if err != nil {
		return err
	}

	return nil
}
