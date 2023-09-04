package fileArray

import (
	"awesomeProject/internal/serialization"
	"bytes"
	"testing"
)

func TestGenerateHeader(t *testing.T) {

	Num := serialization.NewNumber(64)

	headerBytes := generateHeader(Num)

	if len(headerBytes) != headerLength {
		t.Errorf("Header length is incorrect, got: %d, want: %d", len(headerBytes), 16)
	}

	if bytes.Equal(headerBytes[0:6], []byte(signature[0:6])) != true {
		t.Errorf("Header signature is incorrect, got: %s, want: %s", headerBytes[0:6], signature[0:6])
	}

	if headerBytes[6] != version {
		t.Errorf("Header version is incorrect, got: %d, want: %d", headerBytes[6], version)
	}

	println(headerBytes)
}

func TestReadHeader(t *testing.T) {

	Num := serialization.NewNumber(64)

	headerBytes := generateHeader(serialization.Number{})

	header, err := readHeader(headerBytes)
	if err != nil {
		t.Errorf("Failed to read header: %v", err)
	}

	if bytes.Equal(header.signature, []byte(signature[0:6])) != true {
		t.Errorf("Header signature is incorrect, got: %s, want: %s", header.signature, signature[0:6])
	}

	if header.version != version {
		t.Errorf("Header version is incorrect, got: %d, want: %d", header.version, version)
	}

	if header.serializerID != Num.IDByte() {
		t.Errorf("Header serializer ID is incorrect, got: %d, want: %d", header.serializerID, Num.IDByte())
	}

	if header.structHash != serialization.GenerateStructStructureHash(Num) {
		t.Errorf("Header struct hash is incorrect, got: %d, want: %d", header.structHash, serialization.GenerateStructStructureHash(Num))
	}

	if header.strideLength != Num.StrideLength() {
		t.Errorf("Header stride length is incorrect, got: %d, want: %d", header.strideLength, Num.StrideLength())
	}

}

func TestReadHeaderTooShort(t *testing.T) {

	headerBytes := []byte{0x00}

	_, err := readHeader(headerBytes)
	if err == nil {
		t.Errorf("Failed to read header: %v", err)
	}
}

func TestVerifyHeader(t *testing.T) {

	Num := serialization.NewNumber(64)

	headerBytes := generateHeader(Num)

	header, err := readHeader(headerBytes)
	if err != nil {
		t.Errorf("Failed to read header: %v", err)
	}

	err = verifyHeader(Num, header)
	if err != nil {
		t.Errorf("Failed to verify header: %v", err)
	}
}
