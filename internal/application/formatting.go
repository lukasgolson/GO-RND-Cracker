package application

import (
	"encoding/csv"
	"fmt"
	"strings"
)

// FormatByteArrayAsNumbers formats a byte slice as decimal numbers with newlines after a configurable number of items.

func FormatByteArrayAsNumbers(byteArray []byte, itemsPerLine int) string {
	var formattedNumbers []string

	for i, b := range byteArray {
		// Append the decimal number with zero padding to the hundreds place
		formattedNumbers = append(formattedNumbers, fmt.Sprintf("%03d", b))

		// Check if we need to add a newline
		if (i+1)%itemsPerLine == 0 && i != len(byteArray)-1 {
			formattedNumbers = append(formattedNumbers, "\n")
		} else if i != len(byteArray)-1 {
			// Add a comma if it's not the last item
			formattedNumbers = append(formattedNumbers, ", ")
		}
	}

	return strings.Join(formattedNumbers, "")
}

func FormatByteArrayAsCSV(byteArray []byte) (string, error) {
	// Create a buffer to store the CSV data
	var csvBuffer strings.Builder

	// Initialize a CSV writer
	writer := csv.NewWriter(&csvBuffer)

	// Convert the byte array to a slice of strings
	var stringSlice []string
	for _, b := range byteArray {
		stringSlice = append(stringSlice, fmt.Sprintf("%d", b))
	}

	// Write the stringSlice as a CSV record
	err := writer.Write(stringSlice)
	if err != nil {
		return "", err
	}

	// Flush any buffered data to the writer
	writer.Flush()

	// Return the CSV string
	return csvBuffer.String(), nil
}
