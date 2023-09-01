package fileArray

import (
	"encoding/binary"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"io"
	"os"
)

type FileArray struct {
	memoryMap   mmap.MMap
	backingFile *os.File
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

	fileSlice.backingFile = file
	fileSlice.memoryMap = memoryMap

	return fileSlice, nil
}

func openAndInitializeFile(filename string, size int64) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	fileSize, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	if fileSize == 0 {
		_, err := file.Write(make([]byte, size))
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

func (fileArray *FileArray) Count() uint64 {
	bytes := fileArray.memoryMap

	if len(bytes) < 8 {
		return 0
	}
	counterSlice := bytes[:8]
	count := binary.BigEndian.Uint64(counterSlice)
	return count
}

func (fileArray *FileArray) setCount(value uint64) {

	counterSlice := fileArray.memoryMap[:8]
	binary.BigEndian.PutUint64(counterSlice, value)
}

func (fileArray *FileArray) incrementCount() {
	fileArray.setCount(fileArray.Count() + 1)
}

func (fileArray *FileArray) getSlice() []byte {
	return fileArray.memoryMap[8:]
}

func (fileArray *FileArray) increaseMemoryMapSize(newSize int64) error {
	currentSize, err := fileArray.backingFile.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	if newSize <= currentSize {
		return nil // No need to truncate, the new size is not smaller
	}

	if err := fileArray.memoryMap.Unmap(); err != nil {
		return err
	}

	if err := fileArray.backingFile.Truncate(newSize); err != nil {
		return err
	}

	memoryMap, err := mmap.Map(fileArray.backingFile, mmap.RDWR, 0)
	if err != nil {
		return err
	}

	fileArray.memoryMap = memoryMap

	return nil
}

func (fileArray *FileArray) multiplyFileSize(multiplier float64) error {
	if multiplier <= 1.0 {
		return fmt.Errorf("multiplier should be greater than 1.0")
	}

	currentSize, err := fileArray.backingFile.Seek(0, io.SeekEnd)
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

	dataSize := int64(itemSize*fileArray.Count()) + 8

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

func (fileArray *FileArray) hasSpace(dataSize uint64) bool {
	return uint64(len(fileArray.memoryMap)) >= (dataSize + 8)
}

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
