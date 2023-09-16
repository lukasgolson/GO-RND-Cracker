package fileArray

import (
	"awesomeProject/internal/serialization"
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	headerLength       = 28
	signature          = "LGO-FA"
	version      uint8 = 2
)

type Header struct {
	signature    []byte
	version      uint8
	serializerID uint8
	structHash   uint32
	strideLength serialization.Length
}

func generateHeader[T serialization.Serializer[T]](serializer T) []byte {
	header := make([]byte, headerLength)

	// Layout:
	// 6 bytes signature,
	// 1 byte version, 1 byte serializer ID
	// 4 bytes data struct hash, 8 bytes stride length,
	// 8 bytes array count

	copy(header[0:6], signature[0:6])

	header[6] = version

	// Convert serializer.IDByte() to a byte slice and copy it
	serializerID := []byte{serializer.IDByte()}
	copy(header[7:8], serializerID)

	binary.LittleEndian.PutUint32(header[8:12], serialization.GenerateStructStructureHash(serializer))
	binary.LittleEndian.PutUint64(header[12:20], uint64(serializer.StrideLength()))

	binary.LittleEndian.PutUint64(header[20:28], 0)

	return header
}

func readHeader(header []byte) (Header, error) {
	h := Header{}

	if len(header) < headerLength {
		return h, fmt.Errorf("header is too short")
	}

	h.signature = make([]byte, 6)
	copy(h.signature[:], header[0:6])

	h.version = header[6]

	h.serializerID = header[7]

	h.structHash = binary.LittleEndian.Uint32(header[8:12])

	h.strideLength = serialization.Length(binary.LittleEndian.Uint64(header[12:20]))

	return h, nil
}

func verifyHeader[T serialization.Serializer[T]](serializer T, header Header) error {

	if bytes.Equal(header.signature, []byte(signature[0:6])) != true {
		return fmt.Errorf("invalid signature. Got: %s, want: %s", header.signature, signature[0:6])
	}

	if header.version != version {
		return fmt.Errorf("invalid version. Got: %d, compatible: %d", header.version, version)
	}

	if header.serializerID != serializer.IDByte() {
		return fmt.Errorf("invalid serializer ID")
	}

	expectedStructHash := serialization.GenerateStructStructureHash(serializer)
	if header.structHash != expectedStructHash {
		return fmt.Errorf("invalid struct hash. Got: %d, want: %d", header.structHash, expectedStructHash)
	}

	if header.strideLength != serializer.StrideLength() {
		fmt.Println("Warning: Stride length does not match the serializer's stride length")
	}

	return nil
}
