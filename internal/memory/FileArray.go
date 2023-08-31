package memory

import (
	"encoding/binary"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"io"
	"os"
)

type FileArray struct {
	mMap mmap.MMap
	file *os.File
}

func NewFileArray(filename string) (*FileArray, error) {
	fileSlice := &FileArray{}

	file, err := openAndInitializeFile(filename, 8)
	if err != nil {
		return nil, err
	}

	memoryMap, err := openMmap(file)
	if err != nil {
		return nil, err
	}

	fileSlice.file = file
	fileSlice.mMap = memoryMap

	return fileSlice, nil
}

func openAndInitializeFile(filename string, size int64) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644) // Use octal notation for file mode.
	if err != nil {
		return nil, err
	}

	fileSize, err := file.Seek(0, io.SeekEnd) // Check error here.
	if err != nil {
		return nil, err
	}

	if fileSize == 0 {
		_, err := file.Write(make([]byte, size)) // Use _ to ignore the return value if not needed.
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

func openMmap(file *os.File) (mmap.MMap, error) {
	memoryMap, err := mmap.Map(file, mmap.RDWR, 0)
	if err != nil {
		return nil, err
	}

	return memoryMap, nil
}

func (fileArray *FileArray) GetLength() uint64 {
	bytes := fileArray.mMap

	if len(bytes) < 8 {
		return 0
	}
	counterSlice := bytes[:8]
	count := binary.BigEndian.Uint64(counterSlice)
	return count
}

func (fileArray *FileArray) setLength(value uint64) {
	if len(fileArray.mMap) < 8 {
		// Increase the length of bytes to 8.
		// Do this using mmap.Resize if available.
	}

	counterSlice := fileArray.mMap[:8]
	binary.BigEndian.PutUint64(counterSlice, value)
}

func (fileArray *FileArray) incrementLength() {
	currentSize := fileArray.GetLength()
	fileArray.setLength(currentSize + 1)
}

func (fileArray *FileArray) getSlice() []byte {
	return fileArray.mMap[8:]
}

func (fileArray *FileArray) increaseMemoryMapSize(newSize int64) error {
	currentSize, err := fileArray.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	if newSize <= currentSize {
		return nil // No need to truncate, the new size is not smaller
	}

	// Unmap the existing memory-mapped file
	if err := fileArray.mMap.Unmap(); err != nil {
		return err
	}

	// Expand the file size
	if err := fileArray.file.Truncate(newSize); err != nil {
		return err
	}

	// Remap the file with the new size
	memoryMap, err := mmap.Map(fileArray.file, mmap.RDWR, 0)
	if err != nil {
		return err
	}

	fileArray.mMap = memoryMap

	return nil
}

func (fileArray *FileArray) AdjustFileSize(multiplier float64) error {
	if multiplier <= 1.0 {
		return fmt.Errorf("multiplier should be greater than 1.0")
	}

	currentSize, err := fileArray.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	newSize := int64(float64(currentSize) * multiplier)

	if err := fileArray.increaseMemoryMapSize(newSize); err != nil {
		return err
	}

	return nil
}

func (fileArray *FileArray) shrinkFileSizeToDataSize(itemSize uint64) error {

	dataSize := int64(itemSize*fileArray.GetLength()) + 8

	if err := (*fileArray).file.Truncate(dataSize); err != nil {
		return err
	}

	memoryMap, err := mmap.Map((*fileArray).file, mmap.RDWR, 0)
	if err != nil {
		return err
	}

	(*fileArray).mMap = memoryMap

	return nil
}

func (fileArray *FileArray) Close() error {
	var err error

	if fileArray.mMap != nil {
		err = fileArray.mMap.Unmap()
	}

	if fileArray.file != nil {
		err = fileArray.file.Close()
	}

	return err
}
