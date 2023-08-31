package memory

import (
	"encoding/binary"
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

	fileSlice.file = file
	fileSlice.mMap = memoryMap

	return fileSlice, nil
}

func openAndInitializeFile(filename string, size int64) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		return nil, err
	}

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return nil, err
	}
	fileSize, _ := file.Seek(0, io.SeekCurrent)
	if fileSize == 0 {
		if _, err := file.Write(make([]byte, size)); err != nil {
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

func (fileArray *FileArray) GetSize() uint64 {
	bytes := fileArray.mMap

	if len(bytes) < 8 {
		return 0
	}
	counterSlice := bytes[:8]
	count := binary.BigEndian.Uint64(counterSlice)
	return count
}

func (fileArray *FileArray) setSize(value uint64) {
	bytes := fileArray.mMap

	if len(bytes) < 8 {
		//increase the length of bytes to 8
	}

	counterSlice := bytes[:8]
	binary.BigEndian.PutUint64(counterSlice, value)
}

func (fileArray *FileArray) incrementSize() {
	currentSize := fileArray.GetSize()
	fileArray.setSize(currentSize + 1)
}

func (fileArray *FileArray) getSlice() []byte {
	return fileArray.mMap[8:]
}

func (fileArray *FileArray) Close() error {

	var err error

	if fileArray.mMap != nil {
		err = fileArray.mMap.Unmap()
	}

	if fileArray.file != nil {
		err = fileArray.file.Close()
	}

	if err != nil {
		return err
	}

	return nil
}
