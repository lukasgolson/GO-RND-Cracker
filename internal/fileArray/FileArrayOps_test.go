package fileArray

import (
	"awesomeProject/internal/serialization"
	"os"
	"testing"
)

func TestAppendItemSpace(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.NewNumber(42)

	err = fA.expandMemoryMapSize(int64(num.SerializedSize()))
	if err != nil {
		t.Fatalf("Failed to expand memory map size: %v", err)
	}

	err = AppendItem(fA, num)
	if err != nil {
		t.Fatalf("Failed to append item when space is available: %v", err)
	}

}

func TestAppendItemNoSpace(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.NewNumber(42)

	err = AppendItem(fA, num)

	if r := recover(); r != nil {
		t.Fatalf("Failed to append item when space is not available: %v", err)
	}

	if err != nil {
		t.Fatalf("Failed to append item when space is not available: %v", err)
	}
}

func TestSetItemAtIndex(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.NewNumber(42)

	err = SetItemAtIndex(fA, num, 0)
	if err != nil {
		t.Fatalf("Failed to set item at index: %v", err)
	}
}

func TestSetItemAtIndexWithIndexOutOfBounds(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.NewNumber(42)

	err = SetItemAtIndex(fA, num, 1)
	if err == nil {
		t.Fatalf("SetItemAtIndex did not fail with index out of bounds")
	}
}

func TestGetItemFromIndex(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.NewNumber(128)

	err = SetItemAtIndex[serialization.Number](fA, num, 0)
	if err != nil {
		t.Fatalf("Failed to set item at index: %v", err)
	}

	item, err := GetItemFromIndex[serialization.Number](fA, 0)

	if err != nil {
		t.Fatalf("Failed to get item from index: %v", err)
	}

	if item.Value != num.Value {
		t.Fatalf("GetItemFromIndex returned incorrect value. Got %d, expected %d", item.Value, num.Value)
	} else {
		println("GetItemFromIndex returned correct value. Got", item.Value, "expected", num.Value)
	}
}

func TestAppendItemAndCount(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())

	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fileArray, err := NewFileArray(tmpFile.Name())

	if err != nil {
		t.Fatalf("Failed to create file array: %v", err)
	}

	//Repeatedly append items to the file array and check the count
	for i := 0; i < 100; i++ {
		num := serialization.NewNumber(int64(i))
		err = AppendItem(fileArray, num)
		if err != nil {
			t.Fatalf("Failed to append item: %v", err)
		}

		count := fileArray.Count()
		if count != uint64(i+1) {
			t.Fatalf("Count() returned %d, expected %d", count, i+1)
		}

		item, err := GetItemFromIndex[serialization.Number](fileArray, uint64(i))
		if err != nil {
			t.Fatalf("Failed to get item from index: %v", err)
		}

		if item.Value != num.Value {
			t.Fatalf("GetItemFromIndex returned incorrect value. Got %d, expected %d", item.Value, num.Value)
		}

	}
}

func TestGetItemFromIndexWithIndexOutOfBounds(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.NewNumber(128)

	err = SetItemAtIndex[serialization.Number](fA, num, 0)
	if err != nil {
		t.Fatalf("Failed to set item at index: %v", err)
	}

	_, err = GetItemFromIndex[serialization.Number](fA, 1)

	if err == nil {
		t.Fatalf("GetItemFromIndex did not fail with index out of bounds")
	} else {
		println("GetItemFromIndex failed with index out of bounds")
	}
}

func TestGetItemFromIndexWithIndexEqualToCount(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.NewNumber(128)

	err = SetItemAtIndex[serialization.Number](fA, num, 0)
	if err != nil {
		t.Fatalf("Failed to set item at index: %v", err)
	}

	_, err = GetItemFromIndex[serialization.Number](fA, 0)

	if err != nil {
		t.Fatalf("GetItemFromIndex failed with index equal to count")
	} else {
		println("GetItemFromIndex succeeded with index equal to count")
	}
}

func TestAppendAndGetItem(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())

	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fileArray, err := NewFileArray(tmpFile.Name())

	if err != nil {
		t.Fatalf("Failed to create file array: %v", err)
	}

	const count = 100

	//Repeatedly append items to the file array
	for i := 0; i < count; i++ {
		num := serialization.NewNumber(int64(i))
		err = AppendItem(fileArray, num)
		if err != nil {
			t.Fatalf("Failed to append item: %v", err)
		}
	}

	//Repeatedly get items from the file array
	for i := 0; i < count; i++ {
		item, err := GetItemFromIndex[serialization.Number](fileArray, uint64(i))
		if err != nil {
			t.Fatalf("Failed to get item from index: %v", err)
		}

		if item.Value != int64(i) {
			t.Fatalf("GetItemFromIndex returned incorrect value. Got %d, expected %d", item.Value, i)
		}
	}
}

func TestShrinkwrapFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())

	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fileArray, err := NewFileArray(tmpFile.Name())

	if err != nil {
		t.Fatalf("Failed to create file array: %v", err)
	}

	err = AppendItem(fileArray, serialization.NewNumber(42))
	if err != nil {
		t.Fatalf("Failed to append item: %v", err)
	}

	fi, err := tmpFile.Stat()
	if err != nil {
		t.Fatalf("Failed to retrieve file stats...")
	}
	initialSize := fi.Size()

	err = fileArray.expandMemoryMapSize(1024 * 1024)
	if err != nil {
		t.Fatalf("Failed to expand memory map size: %v", err)
	}

	fi, err = tmpFile.Stat()
	if err != nil {
		t.Fatalf("Failed to retrieve file stats...")
	}
	expandedSize := fi.Size()

	err = ShrinkWrapFileArray[serialization.Number](fileArray)

	if err != nil {
		t.Fatalf("Failed to shrinkwrap file: %v", err)
	}

	fi, err = tmpFile.Stat()
	shrunkSize := fi.Size()

	if shrunkSize == expandedSize {
		t.Fatalf("ShrinkWrapFileArray did not shrink the file. Expanded size: %d, Shrunk size: %d", expandedSize, shrunkSize)
	}

	if shrunkSize != initialSize {
		t.Fatalf("ShrinkWrapFileArray shrunk the file to an incorrect size. Initial size: %d, Shrunk size: %d", initialSize, shrunkSize)
	}

	println("ShrinkWrapFileArray succeeded. Initial size:", initialSize, "Expanded size:", expandedSize, "Shrunk size:", shrunkSize)

}
