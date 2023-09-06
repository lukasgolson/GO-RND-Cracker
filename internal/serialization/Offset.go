package serialization

// Offset represents an offset within a file array.
type Offset uint64

// Length is an alias for the Offset type, indicating a length or size within the file array.
type Length = Offset

// MaxOffset returns the maximum value representable by the Offset type.
// It does so by performing a bitwise NOT operation on an Offset initialized with 0,
// effectively flipping all its bits to set them to 1, resulting in the maximum possible Offset value.
func MaxOffset() Offset {
	return ^Offset(0)
}
