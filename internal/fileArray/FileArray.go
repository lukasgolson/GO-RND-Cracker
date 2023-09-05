package fileArray

import (
	"awesomeProject/internal/serialization"
	"encoding/binary"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"io"
	"os"
)

type FileArray struct {
	header      Header
	memoryMap   mmap.MMap
	backingFile *os.File
}

// NewFileArray initializes a new FileArray instance.
//
// Parameters:
//   - serializer: The serializer used to serialize data.
//   - filename: The path to the backing file.
//
// Returns:
//   - *FileArray: A pointer to the FileArray instance.
//   - error: An error if initialization fails.
func NewFileArray[T serialization.Serializer[T]](serializer T, filename string) (*FileArray, error) {
	fileArray := &FileArray{}

	file, err := openAndInitializeFile(serializer, filename)
	if err != nil {
		return nil, err
	}

	memoryMap, err := openMmap(file)
	if err != nil {
		return nil, err
	}

	fileArray.backingFile = file
	fileArray.memoryMap = memoryMap

	header, err := readHeader(fileArray.getHeaderSlice())
	err = verifyHeader(serializer, header)

	fileArray.header = header

	if err != nil {
		return nil, err
	}

	return fileArray, nil
}

// openAndInitializeFile opens or creates a file with the given filename and initializes it with a header if it's a new file.
//
// Parameters:
//   - serializer: The serializer used to serialize data.
//   - filename: The path to the backing file.
//
// Returns:
//   - *os.File: A pointer to the opened file.
//   - error: An error if opening or initialization fails.
func openAndInitializeFile[T serialization.Serializer[T]](serializer T, filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	fileSize, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	if fileSize == 0 {
		_, err := file.Write(generateHeader(serializer))
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

// openMmap maps the opened file into memory for read and write access using mmap.
//
// Parameters:
//   - file: The opened file to be memory-mapped.
//
// Returns:
//   - mmap.MMap: The memory-mapped region.
//   - error: An error if memory mapping fails.
func openMmap(file *os.File) (mmap.MMap, error) {
	memoryMap, err := mmap.Map(file, mmap.RDWR, 0)
	if err != nil {
		return nil, err
	}

	return memoryMap, nil
}

// Count returns the current count of elements stored in the FileArray instance.
func (fileArray *FileArray) Count() uint64 {
	counterSlice := fileArray.getCounterSlice()
	count := binary.BigEndian.Uint64(counterSlice)
	return count
}

// setCount sets the count of elements in the FileArray to the specified value.
func (fileArray *FileArray) setCount(value uint64) {

	counterSlice := fileArray.getCounterSlice()
	binary.BigEndian.PutUint64(counterSlice, value)
}

// incrementCount increments the count of elements in the FileArray by one.
func (fileArray *FileArray) incrementCount() {
	fileArray.setCount(fileArray.Count() + 1)
}

// getDataSlice returns a slice containing the data stored in the FileArray, excluding the header.
func (fileArray *FileArray) getDataSlice() []byte {
	return fileArray.memoryMap[headerLength:]
}

// getHeaderSlice returns a slice containing the data stored in the FileArray header.
func (fileArray *FileArray) getHeaderSlice() []byte {
	return fileArray.memoryMap[:headerLength]
}

// getHeaderSlice returns a slice containing the header data stored in the FileArray.
func (fileArray *FileArray) getCounterSlice() []byte {
	return fileArray.memoryMap[headerLength-8 : headerLength]
}

// expandMemoryMapSize increases the size of the memory-mapped region by the specified expansionSize.
//
// Parameters:
//   - expansionSize: The size by which to expand the memory-mapped region.
//
// Returns:
//   - error: An error if the expansion fails.
func (fileArray *FileArray) expandMemoryMapSize(expansionSize int64) error {
	currentSize, err := fileArray.backingFile.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	if err := fileArray.memoryMap.Unmap(); err != nil {
		return err
	}

	if err := fileArray.backingFile.Truncate(currentSize + expansionSize); err != nil {
		return err
	}

	memoryMap, err := mmap.Map(fileArray.backingFile, mmap.RDWR, 0)
	if err != nil {
		return err
	}

	fileArray.memoryMap = memoryMap

	return nil
}

// multiplyMemoryMapSize increases the size of the memory-mapped region by multiplying the current size with a multiplier.
//
// Parameters:
//   - multiplier: The multiplier for increasing the memory-mapped region size (should be greater than 1.0).
//
// Returns:
//   - error: An error if the operation fails.
func (fileArray *FileArray) multiplyMemoryMapSize(multiplier float64) error {
	if multiplier <= 1.0 {
		return fmt.Errorf("multiplier should be greater than 1.0")
	}

	currentSize, err := fileArray.backingFile.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	newSize := int64(float64(currentSize)*multiplier) - currentSize

	if err := fileArray.expandMemoryMapSize(newSize); err != nil {
		return err
	}

	return nil
}

// shrinkFileSizeToDataSize reduces the size of the backing file to match the actual data size, excluding the header.
//
// Parameters:
//   - itemSize: The size of each item stored in the FileArray.
//
// Returns:
//   - error: An error if the operation fails.
func (fileArray *FileArray) shrinkFileSizeToDataSize(itemSize uint64) error {

	dataSize := int64(itemSize*fileArray.Count()) + headerLength

	err := (*fileArray).memoryMap.Unmap()
	if err != nil {
		return err
	}

	if err := (*fileArray).backingFile.Truncate(dataSize); err != nil {
		return err
	}

	memoryMap, err := mmap.Map((*fileArray).backingFile, mmap.RDWR, 0)
	if err != nil {
		return err
	}

	(*fileArray).memoryMap = memoryMap

	return nil
}

// shrinkFileSizeToDataSize reduces the size of the backing file to match the actual data size.
//
// Parameters:
//   - itemSize: The size of each item stored in the FileArray.
//
// Returns:
//   - error: An error if the operation fails.
func (fileArray *FileArray) hasSpace(dataSize uint64) bool {
	return uint64(len(fileArray.getDataSlice())) > (dataSize)
}

// Close unmaps the memory-mapped region and closes the backing file.
//
// Returns:
//   - error: An error if unmap or file close operations fail.
func (fileArray *FileArray) Close() error {
	var err error

	if fileArray.memoryMap != nil {
		err = fileArray.memoryMap.Unmap()
	}

	if fileArray.backingFile != nil {
		err = fileArray.backingFile.Close()
	}

	return err
}

// GetFileName returns the name of the backing file.
func (fileArray *FileArray) GetFileName() string {
	return fileArray.backingFile.Name()
}
