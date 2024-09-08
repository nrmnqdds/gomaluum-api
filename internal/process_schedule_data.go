package internal

import (
	"strings"
)

func ProcessScrapedData(rawText string) [][]string {
	// Split the raw text into lines
	lines := strings.Split(rawText, "\n")

	// Process each line
	var result [][]string
	for _, line := range lines {
		// Trim spaces and split the line into fields
		fields := strings.Fields(strings.TrimSpace(line))

		// Only add non-empty lines to the result
		if len(fields) > 0 {
			result = append(result, fields)
		}
	}

	return result
}
