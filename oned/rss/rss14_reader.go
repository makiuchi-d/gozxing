package rss

import (
	"strconv"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common/util"
	"github.com/makiuchi-d/gozxing/oned"
)

// Decodes RSS-14, including truncated and stacked variants. See ISO/IEC 24724:2006.

var (
	rss14_OUTSIDE_EVEN_TOTAL_SUBSET = []int{1, 10, 34, 70, 126}
	rss14_INSIDE_ODD_TOTAL_SUBSET   = []int{4, 20, 48, 81}
	rss14_OUTSIDE_GSUM              = []int{0, 161, 961, 2015, 2715}
	rss14_INSIDE_GSUM               = []int{0, 336, 1036, 1516}
	rss14_OUTSIDE_ODD_WIDEST        = []int{8, 6, 4, 3, 1}
	rss14_INSIDE_ODD_WIDEST         = []int{2, 4, 6, 8}

	rss14_FINDER_PATTERNS = [][]int{
		{3, 8, 2, 1},
		{3, 5, 5, 1},
		{3, 3, 7, 1},
		{3, 1, 9, 1},
		{2, 7, 4, 1},
		{2, 5, 6, 1},
		{2, 3, 8, 1},
		{1, 5, 7, 1},
		{1, 3, 9, 1},
	}
)

type rss14Reader struct {
	*oned.OneDReader
	*AbstractRSSReader
	possibleLeftPairs  []*Pair
	possibleRightPairs []*Pair
}

func NewRSS14Reader() gozxing.Reader {
	reader := &rss14Reader{
		AbstractRSSReader:  NewAbstractRSSReader(),
		possibleLeftPairs:  make([]*Pair, 0),
		possibleRightPairs: make([]*Pair, 0),
	}
	reader.OneDReader = oned.NewOneDReader(reader)
	return reader
}

func (this *rss14Reader) DecodeRow(rowNumber int, row *gozxing.BitArray, hints map[gozxing.DecodeHintType]interface{}) (*gozxing.Result, error) {
	leftPair := this.decodePair(row, false, rowNumber, hints)
	this.possibleLeftPairs = this.addOrTally(this.possibleLeftPairs, leftPair)
	row.Reverse()
	rightPair := this.decodePair(row, true, rowNumber, hints)
	this.possibleRightPairs = this.addOrTally(this.possibleRightPairs, rightPair)
	row.Reverse()
	for _, left := range this.possibleLeftPairs {
		if left.GetCount() > 1 {
			for _, right := range this.possibleRightPairs {
				if right.GetCount() > 1 && checkChecksum(left, right) {
					return constructResult(left, right), nil
				}
			}
		}
	}
	return nil, gozxing.NewNotFoundException("no possible left/right pairs")
}

func (this *rss14Reader) addOrTally(possiblePairs []*Pair, pair *Pair) []*Pair {
	if pair == nil {
		return possiblePairs
	}
	found := false
	for _, other := range possiblePairs {
		if other.GetValue() == pair.GetValue() {
			other.IncrementCount()
			found = true
			break
		}
	}
	if !found {
		possiblePairs = append(possiblePairs, pair)
	}
	return possiblePairs
}

func (this *rss14Reader) Reset() {
	this.possibleLeftPairs = this.possibleLeftPairs[:0]
	this.possibleRightPairs = this.possibleRightPairs[:0]
}

func constructResult(leftPair, rightPair *Pair) *gozxing.Result {
	symbolValue := 4537077*leftPair.GetValue() + rightPair.GetValue()
	text := strconv.Itoa(symbolValue)

	buffer := make([]byte, 0, 14)
	for i := 13 - len(text); i > 0; i-- {
		buffer = append(buffer, '0')
	}
	buffer = append(buffer, []byte(text)...)

	checkDigit := byte(0)
	for i := 0; i < 13; i++ {
		digit := buffer[i] - '0'
		if (i & 0x01) == 0 {
			checkDigit += 3 * digit
		} else {
			checkDigit += digit
		}
	}
	checkDigit = 10 - (checkDigit % 10)
	if checkDigit == 10 {
		checkDigit = 0
	}
	buffer = append(buffer, checkDigit+'0')

	leftPoints := leftPair.GetFinderPattern().GetResultPoints()
	rightPoints := rightPair.GetFinderPattern().GetResultPoints()
	result := gozxing.NewResult(
		string(buffer),
		nil,
		[]gozxing.ResultPoint{leftPoints[0], leftPoints[1], rightPoints[0], rightPoints[1]},
		gozxing.BarcodeFormat_RSS_14)
	result.PutMetadata(gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER, "]e0")
	return result
}

func checkChecksum(leftPair, rightPair *Pair) bool {
	checkValue := (leftPair.GetChecksumPortion() + 16*rightPair.GetChecksumPortion()) % 79
	targetCheckValue :=
		9*leftPair.GetFinderPattern().GetValue() + rightPair.GetFinderPattern().GetValue()
	if targetCheckValue > 72 {
		targetCheckValue--
	}
	if targetCheckValue > 8 {
		targetCheckValue--
	}
	return checkValue == targetCheckValue
}

func (this *rss14Reader) decodePair(row *gozxing.BitArray, right bool, rowNumber int, hints map[gozxing.DecodeHintType]interface{}) *Pair {

	startEnd, e := this.findFinderPattern(row, right)
	if e != nil {
		return nil // ignore NotFoundException
	}

	pattern, e := this.parseFoundFinderPattern(row, rowNumber, right, startEnd)
	if e != nil {
		return nil // ignore NotFoundException
	}

	if resultPointCallback, ok := hints[gozxing.DecodeHintType_NEED_RESULT_POINT_CALLBACK].(gozxing.ResultPointCallback); ok {
		startEnd = pattern.GetStartEnd()
		center := float64(startEnd[0]+startEnd[1]-1) / 2.0
		if right {
			// row is actually reversed
			center = float64(row.GetSize()) - 1 - center
		}
		resultPointCallback(gozxing.NewResultPoint(center, float64(rowNumber)))
	}

	outside, e := this.decodeDataCharacter(row, pattern, true)
	if e != nil {
		return nil // ignore NotFoundException
	}

	inside, e := this.decodeDataCharacter(row, pattern, false)
	if e != nil {
		return nil // ignore NotFoundException
	}

	return NewPair(1597*outside.GetValue()+inside.GetValue(),
		outside.GetChecksumPortion()+4*inside.GetChecksumPortion(),
		pattern)
}

func (this *rss14Reader) decodeDataCharacter(row *gozxing.BitArray, pattern *FinderPattern, outsideChar bool) (*DataCharacter, error) {

	counters := this.GetDataCharacterCounters()
	for x := 0; x < len(counters); x++ {
		counters[x] = 0
	}

	if outsideChar {
		oned.RecordPatternInReverse(row, pattern.GetStartEnd()[0], counters)
	} else {
		oned.RecordPattern(row, pattern.GetStartEnd()[1], counters)
		// reverse it
		for i, j := 0, len(counters)-1; i < j; i, j = i+1, j-1 {
			counters[i], counters[j] = counters[j], counters[i]
		}
	}

	numModules := 15
	if outsideChar {
		numModules = 16
	}
	elementWidth := float64(util.MathUtils_Sum(counters)) / float64(numModules)

	oddCounts := this.GetOddCounts()
	evenCounts := this.GetEvenCounts()
	oddRoundingErrors := this.GetOddRoundingErrors()
	evenRoundingErrors := this.GetEvenRoundingErrors()

	for i := 0; i < len(counters); i++ {
		value := float64(counters[i]) / elementWidth
		count := int(value + 0.5) // Round
		if count < 1 {
			count = 1
		} else if count > 8 {
			count = 8
		}
		offset := i / 2
		if (i & 0x01) == 0 {
			oddCounts[offset] = count
			oddRoundingErrors[offset] = value - float64(count)
		} else {
			evenCounts[offset] = count
			evenRoundingErrors[offset] = value - float64(count)
		}
	}

	e := this.adjustOddEvenCounts(outsideChar, numModules)
	if e != nil {
		return nil, e
	}

	oddSum := 0
	oddChecksumPortion := 0
	for i := len(oddCounts) - 1; i >= 0; i-- {
		oddChecksumPortion *= 9
		oddChecksumPortion += oddCounts[i]
		oddSum += oddCounts[i]
	}
	evenChecksumPortion := 0
	evenSum := 0
	for i := len(evenCounts) - 1; i >= 0; i-- {
		evenChecksumPortion *= 9
		evenChecksumPortion += evenCounts[i]
		evenSum += evenCounts[i]
	}
	checksumPortion := oddChecksumPortion + 3*evenChecksumPortion

	if outsideChar {
		if (oddSum&0x01) != 0 || oddSum > 12 || oddSum < 4 {
			return nil, gozxing.NewNotFoundException("oddSum = %v", oddSum)
		}
		group := (12 - oddSum) / 2
		oddWidest := rss14_OUTSIDE_ODD_WIDEST[group]
		evenWidest := 9 - oddWidest
		vOdd := RSSUtils_getRSSvalue(oddCounts, oddWidest, false)
		vEven := RSSUtils_getRSSvalue(evenCounts, evenWidest, true)
		tEven := rss14_OUTSIDE_EVEN_TOTAL_SUBSET[group]
		gSum := rss14_OUTSIDE_GSUM[group]
		return NewDataCharacter(vOdd*tEven+vEven+gSum, checksumPortion), nil
	} else {
		if (evenSum&0x01) != 0 || evenSum > 10 || evenSum < 4 {
			return nil, gozxing.NewNotFoundException("evenSum = %v", oddSum)
		}
		group := (10 - evenSum) / 2
		oddWidest := rss14_INSIDE_ODD_WIDEST[group]
		evenWidest := 9 - oddWidest
		vOdd := RSSUtils_getRSSvalue(oddCounts, oddWidest, true)
		vEven := RSSUtils_getRSSvalue(evenCounts, evenWidest, false)
		tOdd := rss14_INSIDE_ODD_TOTAL_SUBSET[group]
		gSum := rss14_INSIDE_GSUM[group]
		return NewDataCharacter(vEven*tOdd+vOdd+gSum, checksumPortion), nil
	}
}

func (this *rss14Reader) findFinderPattern(row *gozxing.BitArray, rightFinderPattern bool) ([]int, error) {

	counters := this.GetDecodeFinderCounters()
	counters[0] = 0
	counters[1] = 0
	counters[2] = 0
	counters[3] = 0

	width := row.GetSize()
	isWhite := false
	rowOffset := 0
	for rowOffset < width {
		isWhite = !row.Get(rowOffset)
		if rightFinderPattern == isWhite {
			// Will encounter white first when searching for right finder pattern
			break
		}
		rowOffset++
	}

	counterPosition := 0
	patternStart := rowOffset
	for x := rowOffset; x < width; x++ {
		if row.Get(x) != isWhite {
			counters[counterPosition]++
		} else {
			if counterPosition == 3 {
				if RSSReader_isFinderPattern(counters) {
					return []int{patternStart, x}, nil
				}
				patternStart += counters[0] + counters[1]
				counters[0] = counters[2]
				counters[1] = counters[3]
				counters[2] = 0
				counters[3] = 0
				counterPosition--
			} else {
				counterPosition++
			}
			counters[counterPosition] = 1
			isWhite = !isWhite
		}
	}
	return nil, gozxing.NewNotFoundException("finder pattern not found")
}

func (this *rss14Reader) parseFoundFinderPattern(row *gozxing.BitArray, rowNumber int, right bool, startEnd []int) (*FinderPattern, error) {
	// Actually we found elements 2-5
	firstIsBlack := row.Get(startEnd[0])
	firstElementStart := startEnd[0] - 1
	// Locate element 1
	for firstElementStart >= 0 && firstIsBlack != row.Get(firstElementStart) {
		firstElementStart--
	}
	firstElementStart++
	firstCounter := startEnd[0] - firstElementStart
	// Make 'counters' hold 1-4
	counters := this.GetDecodeFinderCounters()
	copy(counters[1:], counters[:len(counters)-1])
	counters[0] = firstCounter
	value, e := RSSReader_parseFinderValue(counters, rss14_FINDER_PATTERNS)
	if e != nil {
		return nil, e
	}
	start := firstElementStart
	end := startEnd[1]
	if right {
		// row is actually reversed
		start = row.GetSize() - 1 - start
		end = row.GetSize() - 1 - end
	}
	return NewFinderPattern(value, []int{firstElementStart, startEnd[1]}, start, end, rowNumber), nil
}

func (this *rss14Reader) adjustOddEvenCounts(outsideChar bool, numModules int) error {

	oddSum := util.MathUtils_Sum(this.GetOddCounts())
	evenSum := util.MathUtils_Sum(this.GetEvenCounts())

	incrementOdd := false
	decrementOdd := false
	incrementEven := false
	decrementEven := false

	if outsideChar {
		if oddSum > 12 {
			decrementOdd = true
		} else if oddSum < 4 {
			incrementOdd = true
		}
		if evenSum > 12 {
			decrementEven = true
		} else if evenSum < 4 {
			incrementEven = true
		}
	} else {
		if oddSum > 11 {
			decrementOdd = true
		} else if oddSum < 5 {
			incrementOdd = true
		}
		if evenSum > 10 {
			decrementEven = true
		} else if evenSum < 4 {
			incrementEven = true
		}
	}

	mismatch := oddSum + evenSum - numModules
	oddParityBad := (oddSum & 0x01) == 0
	if outsideChar {
		oddParityBad = (oddSum & 0x01) == 1
	}
	evenParityBad := (evenSum & 0x01) == 1
	/*if (mismatch == 2) {
	    if (!(oddParityBad && evenParityBad)) {
	      throw ReaderException.getInstance();
	    }
	    decrementOdd = true;
	    decrementEven = true;
	  } else if (mismatch == -2) {
	    if (!(oddParityBad && evenParityBad)) {
	      throw ReaderException.getInstance();
	    }
	    incrementOdd = true;
	    incrementEven = true;
	  } else */
	switch mismatch {
	case 1:
		if oddParityBad {
			if evenParityBad {
				return gozxing.NewNotFoundException("adjustOddEvenCounts mismatch=1")
			}
			decrementOdd = true
		} else {
			if !evenParityBad {
				return gozxing.NewNotFoundException("adjustOddEvenCounts mismatch=1")
			}
			decrementEven = true
		}
		break
	case -1:
		if oddParityBad {
			if evenParityBad {
				return gozxing.NewNotFoundException("adjustOddEvenCounts mismatch=-1")
			}
			incrementOdd = true
		} else {
			if !evenParityBad {
				return gozxing.NewNotFoundException("adjustOddEvenCounts mismatch=-1")
			}
			incrementEven = true
		}
		break
	case 0:
		if oddParityBad {
			if !evenParityBad {
				return gozxing.NewNotFoundException("adjustOddEvenCounts mismatch=0")
			}
			// Both bad
			if oddSum < evenSum {
				incrementOdd = true
				decrementEven = true
			} else {
				decrementOdd = true
				incrementEven = true
			}
		} else {
			if evenParityBad {
				return gozxing.NewNotFoundException("adjustOddEvenCounts mismatch=0")
			}
			// Nothing to do!
		}
		break
	default:
		return gozxing.NewNotFoundException("adjustOddEvenCounts mismatch=%v", mismatch)
	}

	if incrementOdd {
		if decrementOdd {
			return gozxing.NewNotFoundException("adjustOddEvenCounts incrementOdd & decrementOdd")
		}
		RSSReader_increment(this.GetOddCounts(), this.GetOddRoundingErrors())
	}
	if decrementOdd {
		RSSReader_decrement(this.GetOddCounts(), this.GetOddRoundingErrors())
	}
	if incrementEven {
		if decrementEven {
			return gozxing.NewNotFoundException("adjustOddEvenCounts incrementEven & decrementEven")
		}
		RSSReader_increment(this.GetEvenCounts(), this.GetEvenRoundingErrors())
	}
	if decrementEven {
		RSSReader_decrement(this.GetEvenCounts(), this.GetEvenRoundingErrors())
	}
	return nil
}
