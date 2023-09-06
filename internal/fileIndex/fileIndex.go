package fileIndex

import "awesomeProject/internal/fileArray"

type fileIndex struct {
	fileArray fileArray.FileArray
}

func NewFileIndex(filename string) *fileIndex {
	fIndex := &fileIndex{}

	fa, err := fileArray.NewFileArray[Index](Index{}, filename)

	if err != nil {
		panic("Failed to create file fIndex")
	}

	fIndex.fileArray = *fa

	return fIndex
}
