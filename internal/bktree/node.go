package bktree

import (
	"encoding/binary"
	"io"
)

const nodeWordLength = 32

const nodeByteSize = 8 + nodeWordLength

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
	binary.LittleEndian.PutUint32(idBytes, node.ID)
	_, err := writer.Write(idBytes)
	if err != nil {
		return err
	}

	_, err = writer.Write(node.Word[:])
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

func DeserializeNodeFromBinaryStream(reader io.Reader) (*Node, error) {
	data := make([]byte, 8+nodeWordLength)

	_, err := io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}

	nodeID := binary.LittleEndian.Uint32(data[:4])

	seed := binary.LittleEndian.Uint32(data[4:8])

	word := make([]byte, nodeWordLength)
	copy(word, data[8:8+nodeWordLength])

	return &Node{
		ID:   nodeID,
		Word: word,
		Seed: int32(seed),
	}, nil
}
