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

const headerLength = 24
const signature = "LOFA"
const version uint8 = 1

func NewFileArray(filename string) (*FileArray, error) {
	fileSlice := &FileArray{}

	file, err := openAndInitializeFile(filename)
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

func openAndInitializeFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	fileSize, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	if fileSize == 0 {
		_, err := file.Write(generateHeader())
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
	counterSlice := fileArray.getCounterSlice()
	count := binary.BigEndian.Uint64(counterSlice)
	return count
}

func (fileArray *FileArray) setCount(value uint64) {

	counterSlice := fileArray.getCounterSlice()
	binary.BigEndian.PutUint64(counterSlice, value)
}

func (fileArray *FileArray) incrementCount() {
	fileArray.setCount(fileArray.Count() + 1)
}

func (fileArray *FileArray) getDataSlice() []byte {
	return fileArray.memoryMap[headerLength:]
}

func (fileArray *FileArray) getHeaderSlice() []byte {
	return fileArray.memoryMap[:headerLength]
}

func (fileArray *FileArray) getCounterSlice() []byte {
	return fileArray.memoryMap[headerLength-8:]
}

func generateHeader() []byte {
	header := make([]byte, headerLength)

	// Layout: 4 bytes signature, 1 byte version, 2 byte reserve,
	// 4 bytes data type, 4 bytes data size,
	// 8 bytes array count

	copy(header[0:4], signature[0:4])
	header[4] = version
	binary.LittleEndian.PutUint16(header[5:7], 0)   // For 2-byte reserve
	binary.LittleEndian.PutUint32(header[7:11], 0)  // For 4-byte data type and data size
	binary.LittleEndian.PutUint64(header[11:19], 0) // For 8-byte array count

	return header
}

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

func (fileArray *FileArray) hasSpace(dataSize uint64) bool {
	return uint64(len(fileArray.getDataSlice())) > (dataSize)
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
