package common

import (
	"fmt"
	"math"

	"golang.org/x/text/encoding"
	"golang.org/x/xerrors"
)

// MinimalECIInput Class that converts a character string into a sequence of ECIs and bytes
//
// The implementation uses the Dijkstra algorithm to produce minimal encodings
type MinimalECIInput struct {
	bytes []int
	fnc1  int
}

var _ ECIInput = (*MinimalECIInput)(nil)

const (
	costPerECI = 3 // approximated (latch + 2 codewords)
	fnc1Value  = 1000
)

// NewMinimalECIInput Constructs a minimal input
//
// @param stringToEncode the character string to encode
// @param priorityCharset The preferred encoding.Encoding. When the value of the argument is nil, the algorithm
// chooses charsets that leads to a minimal representation. Otherwise the algorithm will use the priority
// charset to encode any character in the input that can be encoded by it if the charset is among the
// supported charsets.
//
// @param fnc1 denotes the character in the input that represents the FNC1 character or -1 if this is not GS1
// input.
func NewMinimalECIInput(stringToEncode []rune, priorityCharset encoding.Encoding, fnc1 int) (*MinimalECIInput, error) {
	encoderSet := NewECIEncoderSet(stringToEncode, priorityCharset, fnc1)
	var bytes []int
	if encoderSet.Length() == 1 {
		// optimization for the case when all can be encoded without ECI in ISO-8859-1
		bytes = make([]int, len(stringToEncode))
		for i := 0; i < len(bytes); i++ {
			c := stringToEncode[i]
			if c == rune(fnc1) {
				bytes[i] = fnc1Value
			} else {
				bytes[i] = int(c)
			}
		}
	} else {
		var err error
		bytes, err = minimalECIInput_encodeMinimally(stringToEncode, encoderSet, fnc1)
		if err != nil {
			return nil, err
		}
	}
	return &MinimalECIInput{
		bytes: bytes,
		fnc1:  fnc1,
	}, nil
}

func (this *MinimalECIInput) GetFNC1Character() int {
	return this.fnc1
}

// Length Returns the length of this input. The length is the number
// of bytes, FNC1 characters or ECIs in the sequence.
//
// @return  the number of chars in this sequence
func (this *MinimalECIInput) Length() int {
	return len(this.bytes)
}

func (this *MinimalECIInput) HaveNCharacters(index, n int) bool {
	if index+n-1 >= len(this.bytes) {
		return false
	}
	for i := 0; i < n; i++ {
		if isECI, _ := this.IsECI(index + i); isECI {
			return false
		}
	}
	return true
}

// CharAt Returns the byte value at the specified index. An index ranges from zero
// to length() - 1. The first byte value of the sequence is at
// index zero, the next at index one, and so on, as for array
// indexing.
//
// @param   index the index of the byte value to be returned
//
// @return  the specified byte value as character or the FNC1 character
func (this *MinimalECIInput) CharAt(index int) (int, error) {
	if index < 0 || index >= this.Length() {
		return 0, xerrors.Errorf("index out of bounds: %d", index)
	}
	if isECI, _ := this.IsECI(index); isECI {
		return 0, xerrors.Errorf("value at %d is not a character but an ECI", index)
	}
	if isFNC1, _ := this.IsFNC1(index); isFNC1 {
		return this.fnc1, nil
	}
	return this.bytes[index], nil
}

// SubSequence Returns a subsequence of this sequence.
// The subsequence starts with the char value at the specified index and
// ends with the char value at index end - 1. The length
// (in chars) of the
// returned sequence is end - start, so if start == end
// then an empty sequence is returned.
func (this *MinimalECIInput) SubSequence(start, end int) ([]int, error) {
	if start < 0 || start > end || end > this.Length() {
		return nil, xerrors.Errorf("index out of bounds: start=%d, end=%d", start, end)
	}
	result := make([]int, 0, end-start)
	for i := start; i < end; i++ {
		if isECI, _ := this.IsECI(i); isECI {
			return nil, xerrors.Errorf("value at %d is not a character but an ECI", i)
		}
		c, err := this.CharAt(i)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

// IsECI Determines if a value is an ECI
//
// @param   index the index of the value
//
// @return  true if the value at position index is an ECI
func (this *MinimalECIInput) IsECI(index int) (bool, error) {
	if index < 0 || index >= this.Length() {
		return false, xerrors.Errorf("index out of bounds: %d", index)
	}
	return this.bytes[index] > 255 && this.bytes[index] <= 999, nil
}

// IsFNC1 Determines if a value is the FNC1 character
//
// @param   index the index of the value
//
// @return  true if the value at position index is the FNC1 character
func (this *MinimalECIInput) IsFNC1(index int) (bool, error) {
	if index < 0 || index >= this.Length() {
		return false, xerrors.Errorf("index out of bounds: %d", index)
	}
	return this.bytes[index] == fnc1Value, nil
}

// GetECIValue Returns the int ECI value at the specified index. An index ranges from zero
// to length() - 1. The first byte value of the sequence is at
// index zero, the next at index one, and so on, as for array
// indexing.
//
// @param   index the index of the int value to be returned
//
// @return  the specified int ECI value.
// The ECI specified the encoding of all bytes with a higher index until the
// next ECI or until the end of the input if no other ECI follows.
func (this *MinimalECIInput) GetECIValue(index int) (int, error) {
	if index < 0 || index >= this.Length() {
		return 0, xerrors.Errorf("index out of bounds: %d", index)
	}
	if isECI, _ := this.IsECI(index); !isECI {
		return 0, xerrors.Errorf("value at %d is not an ECI but a character", index)
	}
	return this.bytes[index] - 256, nil
}

func (this *MinimalECIInput) String() string {
	var result []byte
	for i := range this.Length() {
		if i > 0 {
			result = fmt.Append(result, ", ")
		}
		if isECI, _ := this.IsECI(i); isECI {
			result = fmt.Append(result, "ECI(")
			v, err := this.GetECIValue(i)
			if err != nil {
				result = fmt.Appendf(result, "error:%v", err)
			} else {
				result = fmt.Append(result, v)
			}
			result = append(result, ')')
		} else {
			c, err := this.CharAt(i)
			if err != nil {
				result = fmt.Appendf(result, "error:%v", err)
			} else if c < 128 {
				result = fmt.Appendf(result, "'%c'", c)
			} else {
				result = fmt.Appendf(result, "%d", c)
			}
		}
	}
	return string(result)
}

func minimalECIInput_addEdge(edges [][]*inputEdge, to int, edge *inputEdge) {
	if edges[to][edge.encoderIndex] == nil ||
		edges[to][edge.encoderIndex].cachedTotalSize > edge.cachedTotalSize {
		edges[to][edge.encoderIndex] = edge
	}
}

func minimalECIInput_addEdges(
	stringToEncode []rune,
	encoderSet *ECIEncoderSet,
	edges [][]*inputEdge,
	from int,
	previous *inputEdge,
	fnc1 int) {

	ch := stringToEncode[from]

	start := 0
	end := encoderSet.Length()
	if encoderSet.GetPriorityEncoderIndex() >= 0 && (ch == rune(fnc1) || encoderSet.CanEncode(ch, encoderSet.GetPriorityEncoderIndex())) {
		start = encoderSet.GetPriorityEncoderIndex()
		end = start + 1
	}

	for i := start; i < end; i++ {
		if ch == rune(fnc1) || encoderSet.CanEncode(ch, i) {
			minimalECIInput_addEdge(edges, from+1, newInputEdge(ch, encoderSet, i, previous, fnc1))
		}
	}
}

func minimalECIInput_encodeMinimally(stringToEncode []rune, encoderSet *ECIEncoderSet, fnc1 int) ([]int, error) {
	inputLength := len(stringToEncode)

	// Array that represents vertices. There is a vertex for every character and encoding.
	edges := make([][]*inputEdge, inputLength+1)
	for i := range edges {
		edges[i] = make([]*inputEdge, encoderSet.Length())
	}
	minimalECIInput_addEdges(stringToEncode, encoderSet, edges, 0, nil, fnc1)

	for i := 1; i <= inputLength; i++ {
		for j := 0; j < encoderSet.Length(); j++ {
			if edges[i][j] != nil && i < inputLength {
				minimalECIInput_addEdges(stringToEncode, encoderSet, edges, i, edges[i][j], fnc1)
			}
		}
		// optimize memory by removing edges that have been passed.
		for j := 0; j < encoderSet.Length(); j++ {
			edges[i-1][j] = nil
		}
	}
	minimalJ := -1
	minimalSize := math.MaxInt
	for j := 0; j < encoderSet.Length(); j++ {
		if edges[inputLength][j] != nil {
			edge := edges[inputLength][j]
			if edge.cachedTotalSize < minimalSize {
				minimalSize = edge.cachedTotalSize
				minimalJ = j
			}
		}
	}
	if minimalJ < 0 {
		return nil, xerrors.Errorf("Failed to encode \"%s\"", string(stringToEncode))
	}
	intsAL := make([]int, 0)
	current := edges[inputLength][minimalJ]
	for current != nil {
		if current.c == fnc1Value {
			intsAL = append([]int{fnc1Value}, intsAL...)
		} else {
			bytes, _ := encoderSet.EncodeChar(rune(current.c), current.encoderIndex)
			for i := len(bytes) - 1; i >= 0; i-- {
				intsAL = append([]int{int(bytes[i] & 0xFF)}, intsAL...)
			}
		}
		previousEncoderIndex := 0
		if current.previous != nil {
			previousEncoderIndex = current.previous.encoderIndex
		}
		if previousEncoderIndex != current.encoderIndex {
			intsAL = append([]int{256 + encoderSet.GetECIValue(current.encoderIndex)}, intsAL...)
		}
		current = current.previous
	}
	return intsAL, nil
}

type inputEdge struct {
	c               int // character or fnc1Value
	encoderIndex    int // the encoding of this edge
	previous        *inputEdge
	cachedTotalSize int
}

func newInputEdge(c rune, encoderSet *ECIEncoderSet, encoderIndex int, previous *inputEdge, fnc1 int) *inputEdge {
	var cValue int
	var size int
	if c == rune(fnc1) {
		cValue = fnc1Value
		size = 1
	} else {
		cValue = int(c)
		bytes, err := encoderSet.EncodeChar(c, encoderIndex)
		if err != nil {
			// This should not happen if CanEncode was checked
			size = 1 // fallback
		} else {
			size = len(bytes)
		}
	}

	previousEncoderIndex := 0
	if previous != nil {
		previousEncoderIndex = previous.encoderIndex
	}
	if previousEncoderIndex != encoderIndex {
		size += costPerECI
	}
	if previous != nil {
		size += previous.cachedTotalSize
	}

	return &inputEdge{
		c:               cValue,
		encoderIndex:    encoderIndex,
		previous:        previous,
		cachedTotalSize: size,
	}
}

func (e *inputEdge) isFNC1() bool {
	return e.c == fnc1Value
}
