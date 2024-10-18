package helpers

import "strings"

// compareSessionNames compares two session names and returns true if a should come before b
func CompareSessionNames(a, b string) bool {
	aParts := strings.Split(a, ", ")
	bParts := strings.Split(b, ", ")

	if len(aParts) != 2 || len(bParts) != 2 {
		return a < b // fallback to lexicographic comparison if format is unexpected
	}

	aSem := aParts[0]
	aYear := aParts[1]
	bSem := bParts[0]
	bYear := bParts[1]

	if aYear != bYear {
		return aYear > bYear // more recent years come first
	}

	// If years are the same, compare semesters
	return aSem > bSem // "Sem 2" comes before "Sem 1" for the same year
}
