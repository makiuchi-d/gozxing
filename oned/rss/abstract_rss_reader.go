package rss

import (
	"math"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
)

// Superclass of {@link OneDReader} implementations that read barcodes in the RSS family
// of formats.

const (
	rssReader_MAX_AVG_VARIANCE        = 0.2
	rssReader_MAX_INDIVIDUAL_VARIANCE = 0.45

	rssReader_MIN_FINDER_PATTERN_RATIO = 9.5 / 12.0
	rssReader_MAX_FINDER_PATTERN_RATIO = 12.5 / 14.0
)

type AbstractRSSReader struct {
	decodeFinderCounters  []int
	dataCharacterCounters []int
	oddRoundingErrors     []float64
	evenRoundingErrors    []float64
	oddCounts             []int
	evenCounts            []int
}

func NewAbstractRSSReader() *AbstractRSSReader {
	dcCountersLen := 8
	reader := &AbstractRSSReader{
		decodeFinderCounters:  make([]int, 4),
		dataCharacterCounters: make([]int, dcCountersLen),
		oddRoundingErrors:     make([]float64, 4),
		evenRoundingErrors:    make([]float64, 4),
		oddCounts:             make([]int, dcCountersLen/2),
		evenCounts:            make([]int, dcCountersLen/2),
	}
	return reader
}

func (this *AbstractRSSReader) GetDecodeFinderCounters() []int {
	return this.decodeFinderCounters
}

func (this *AbstractRSSReader) GetDataCharacterCounters() []int {
	return this.dataCharacterCounters
}

func (this *AbstractRSSReader) GetOddRoundingErrors() []float64 {
	return this.oddRoundingErrors
}

func (this *AbstractRSSReader) GetEvenRoundingErrors() []float64 {
	return this.evenRoundingErrors
}

func (this *AbstractRSSReader) GetOddCounts() []int {
	return this.oddCounts
}

func (this *AbstractRSSReader) GetEvenCounts() []int {
	return this.evenCounts
}

func RSSReader_parseFinderValue(counters []int, finderPatterns [][]int) (int, error) {
	for value := 0; value < len(finderPatterns); value++ {
		v := oned.PatternMatchVariance(counters, finderPatterns[value], rssReader_MAX_INDIVIDUAL_VARIANCE)
		if v < rssReader_MAX_AVG_VARIANCE {
			return value, nil
		}
	}
	return 0, gozxing.NewNotFoundException("No matched finder pattern")
}

// /**
//  * @param array values to sum
//  * @return sum of values
//  * @deprecated call {@link MathUtils#sum(int[])}
//  */
// @Deprecated
// protected static int count(int[] array) {
//   return MathUtils.sum(array);
// }

func RSSReader_increment(array []int, errors []float64) {
	index := 0
	biggestError := errors[0]
	for i := 1; i < len(array); i++ {
		if errors[i] > biggestError {
			biggestError = errors[i]
			index = i
		}
	}
	array[index]++
}

func RSSReader_decrement(array []int, errors []float64) {
	index := 0
	biggestError := errors[0]
	for i := 1; i < len(array); i++ {
		if errors[i] < biggestError {
			biggestError = errors[i]
			index = i
		}
	}
	array[index]--
}

func RSSReader_isFinderPattern(counters []int) bool {
	firstTwoSum := counters[0] + counters[1]
	sum := firstTwoSum + counters[2] + counters[3]
	ratio := float64(firstTwoSum) / float64(sum)
	if ratio >= rssReader_MIN_FINDER_PATTERN_RATIO && ratio <= rssReader_MAX_FINDER_PATTERN_RATIO {
		// passes ratio test in spec, but see if the counts are unreasonable
		minCounter := math.MaxInt32
		maxCounter := math.MinInt32
		for _, counter := range counters {
			if counter > maxCounter {
				maxCounter = counter
			}
			if counter < minCounter {
				minCounter = counter
			}
		}
		return maxCounter < 10*minCounter
	}
	return false
}
