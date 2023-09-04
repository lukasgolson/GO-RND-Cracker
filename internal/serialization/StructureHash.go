package serialization

import (
	"hash/crc32"
	"reflect"
	"sort"
	"strings"
)

func GenerateStructStructureHash(data interface{}) uint32 {
	val := reflect.TypeOf(data)

	if val.Kind() != reflect.Struct {
		panic("Input is not a struct")
	}

	var fieldInfo []string

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Hash the field's name and type
		fieldName := field.Name
		fieldType := field.Type.String()

		// Concatenate the field name and type and store in the slice
		fieldInfo = append(fieldInfo, fieldName+fieldType)
	}

	sort.Strings(fieldInfo)

	combinedData := strings.Join(fieldInfo, "")

	crcHash := crc32.ChecksumIEEE([]byte(combinedData))

	return crcHash
}
