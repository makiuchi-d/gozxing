package oned

import (
	"github.com/makiuchi-d/gozxing"
)

const (
	// These two values are critical for determining how permissive the decoding will be.
	// We've arrived at these values through a lot of trial and error. Setting them any higher
	// lets false positives creep in quickly.
	UPCEANReader_MAX_AVG_VARIANCE        = 0.48
	UPCEANReader_MAX_INDIVIDUAL_VARIANCE = 0.7
)

var (
	// Start/end guard pattern.
	UPCEANReader_START_END_PATTERN = []int{1, 1, 1}

	// Pattern marking the middle of a UPC/EAN pattern, separating the two halves.
	UPCEANReader_MIDDLE_PATTERN = []int{1, 1, 1, 1, 1}

	// end guard pattern.
	UPCEANReader_END_PATTERN = []int{1, 1, 1, 1, 1, 1}

	// "Odd", or "L" patterns used to encode UPC/EAN digits.
	UPCEANReader_L_PATTERNS = [][]int{
		{3, 2, 1, 1}, // 0
		{2, 2, 2, 1}, // 1
		{2, 1, 2, 2}, // 2
		{1, 4, 1, 1}, // 3
		{1, 1, 3, 2}, // 4
		{1, 2, 3, 1}, // 5
		{1, 1, 1, 4}, // 6
		{1, 3, 1, 2}, // 7
		{1, 2, 1, 3}, // 8
		{3, 1, 1, 2}, // 9
	}

	// As above but also including the "even", or "G" patterns used to encode UPC/EAN digits.
	UPCEANReader_L_AND_G_PATTERNS [][]int
)

func init() {
	UPCEANReader_L_AND_G_PATTERNS = make([][]int, 20)
	copy(UPCEANReader_L_AND_G_PATTERNS, UPCEANReader_L_PATTERNS)
	for i := 10; i < 20; i++ {
		widths := UPCEANReader_L_PATTERNS[i-10]
		reversedWidths := make([]int, len(widths))
		for j := 0; j < len(widths); j++ {
			reversedWidths[j] = widths[len(widths)-j-1]
		}
		UPCEANReader_L_AND_G_PATTERNS[i] = reversedWidths
	}
}

func UPCEANReader_findGuardPattern(row *gozxing.BitArray, rowOffset int, whiteFirst bool, pattern []int) ([]int, error) {
	counters := make([]int, len(pattern))
	return UPCEANReader_findGuardPatternWithCounters(row, rowOffset, whiteFirst, pattern, counters)
}

// UPCEANReader_findGuardPatternWithCounters Find guard pattern
// @param row row of black/white values to search
// @param rowOffset position to start search
// @param whiteFirst if true, indicates that the pattern specifies white/black/white/...
// pixel counts, otherwise, it is interpreted as black/white/black/...
// @param pattern pattern of counts of number of black and white pixels that are being
// searched for as a pattern
// @param counters array of counters, as long as pattern, to re-use
// @return start/end horizontal offset of guard pattern, as an array of two ints
// @throws NotFoundException if pattern is not found
func UPCEANReader_findGuardPatternWithCounters(
	row *gozxing.BitArray, rowOffset int, whiteFirst bool, pattern, counters []int) ([]int, error) {

	width := row.GetSize()
	if whiteFirst {
		rowOffset = row.GetNextUnset(rowOffset)
	} else {
		rowOffset = row.GetNextSet(rowOffset)
	}
	counterPosition := 0
	patternStart := rowOffset
	patternLength := len(pattern)
	isWhite := whiteFirst
	for x := rowOffset; x < width; x++ {
		if row.Get(x) != isWhite {
			counters[counterPosition]++
		} else {
			if counterPosition == patternLength-1 {
				if patternMatchVariance(counters, pattern, UPCEANReader_MAX_INDIVIDUAL_VARIANCE) < UPCEANReader_MAX_AVG_VARIANCE {
					return []int{patternStart, x}, nil
				}
				patternStart += counters[0] + counters[1]
				copy(counters[:counterPosition-1], counters[2:counterPosition+1])
				counters[counterPosition-1] = 0
				counters[counterPosition] = 0
				counterPosition--
			} else {
				counterPosition++
			}
			counters[counterPosition] = 1
			isWhite = !isWhite
		}
	}
	return nil, gozxing.GetNotFoundExceptionInstance()
}

// UPCEANReader_decodeDigit Attempts to decode a single UPC/EAN-encoded digit.
//
// @param row row of black/white values to decode
// @param counters the counts of runs of observed black/white/black/... values
// @param rowOffset horizontal offset to start decoding from
// @param patterns the set of patterns to use to decode -- sometimes different encodings
// for the digits 0-9 are used, and this indicates the encodings for 0 to 9 that should
// be used
// @return horizontal offset of first pixel beyond the decoded digit
// @throws NotFoundException if digit cannot be decoded
func UPCEANReader_decodeDigit(row *gozxing.BitArray, counters []int, rowOffset int, patterns [][]int) (int, error) {
	e := recordPattern(row, rowOffset, counters)
	if e != nil {
		return 0, e
	}
	bestVariance := UPCEANReader_MAX_AVG_VARIANCE // worst variance we'll accept
	bestMatch := -1
	max := len(patterns)
	for i := 0; i < max; i++ {
		pattern := patterns[i]
		variance := patternMatchVariance(counters, pattern, UPCEANReader_MAX_INDIVIDUAL_VARIANCE)
		if variance < bestVariance {
			bestVariance = variance
			bestMatch = i
		}
	}
	if bestMatch >= 0 {
		return bestMatch, nil
	} else {
		return 0, gozxing.GetNotFoundExceptionInstance()
	}
}
