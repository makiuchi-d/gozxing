package encoder

import (
	"math"
	"strings"

	"golang.org/x/text/encoding"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
)

// MinimalEncoderMode represents the encoding mode
type MinimalEncoderMode int

const (
	MinimalEncoderModeASCII MinimalEncoderMode = iota
	MinimalEncoderModeC40
	MinimalEncoderModeTEXT
	MinimalEncoderModeX12
	MinimalEncoderModeEDF
	MinimalEncoderModeB256
)

var MinimalEncoder_c40Shift2Chars = []byte{'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
	':', ';', '<', '=', '>', '?', '@', '[', '\\', ']', '^', '_'}

func MinimalEncoder_isExtendedASCII(ch byte, fnc1 int) bool {
	return int(ch) != fnc1 && ch >= 128 && ch <= 255
}

func MinimalEncoder_isInC40Shift1Set(ch byte) bool {
	return ch <= 31
}

func MinimalEncoder_isInC40Shift2Set(ch byte, fnc1 int) bool {
	for _, c40Shift2Char := range MinimalEncoder_c40Shift2Chars {
		if c40Shift2Char == ch {
			return true
		}
	}
	return int(ch) == fnc1
}

func MinimalEncoder_isInTextShift1Set(ch byte) bool {
	return MinimalEncoder_isInC40Shift1Set(ch)
}

func MinimalEncoder_isInTextShift2Set(ch byte, fnc1 int) bool {
	return MinimalEncoder_isInC40Shift2Set(ch, fnc1)
}

// MinimalEncoder_EncodeHighLevel Performs message encoding of a DataMatrix message
//
// @param msg the message
// @param priorityCharset The preferred encoding.Encoding. When the value of the argument is nil, the algorithm
// chooses charsets that leads to a minimal representation. Otherwise the algorithm will use the priority
// charset to encode any character in the input that can be encoded by it if the charset is among the
// supported charsets.
// @param fnc1 denotes the character in the input that represents the FNC1 character or -1 if this is not a GS1
// bar code. If the value is not -1 then a FNC1 is also prepended.
// @param shape requested shape.
// @return the encoded message (the char values range from 0 to 255)
func MinimalEncoder_EncodeHighLevel(msg string, priorityCharset encoding.Encoding, fnc1 int, shape SymbolShapeHint) ([]byte, error) {
	macroId := 0
	if strings.HasPrefix(msg, HighLevelEncoder_MACRO_05_HEADER) && strings.HasSuffix(msg, HighLevelEncoder_MACRO_TRAILER) {
		macroId = 5
		msg = msg[len(HighLevelEncoder_MACRO_05_HEADER) : len(msg)-2]
	} else if strings.HasPrefix(msg, HighLevelEncoder_MACRO_06_HEADER) && strings.HasSuffix(msg, HighLevelEncoder_MACRO_TRAILER) {
		macroId = 6
		msg = msg[len(HighLevelEncoder_MACRO_06_HEADER) : len(msg)-2]
	}
	return MinimalEncoder_encode([]rune(msg), priorityCharset, fnc1, shape, macroId)
}

// MinimalEncoder_encode Encodes input minimally and returns an array of the codewords
//
// @param input The string to encode
// @param priorityCharset The preferred encoding.Encoding. When the value of the argument is nil, the algorithm
// chooses charsets that leads to a minimal representation. Otherwise the algorithm will use the priority
// charset to encode any character in the input that can be encoded by it if the charset is among the
// supported charsets.
// @param fnc1 denotes the character in the input that represents the FNC1 character or -1 if this is not a GS1
// bar code. If the value is not -1 then a FNC1 is also prepended.
// @param shape requested shape.
// @param macroId Prepends the specified macro function in case that a value of 5 or 6 is specified.
// @return An array of bytes representing the codewords of a minimal encoding.
func MinimalEncoder_encode(input []rune, priorityCharset encoding.Encoding, fnc1 int, shape SymbolShapeHint, macroId int) ([]byte, error) {
	in, err := newMinimalEncoderInput(input, priorityCharset, fnc1, shape, macroId)
	if err != nil {
		return nil, err
	}
	re, err := MinimalEncoder_encodeMinimally(in)
	if err != nil {
		return nil, err
	}
	return re.getBytes(), nil
}

func MinimalEncoder_addEdge(edges [][]*MinimalEncoderEdge, edge *MinimalEncoderEdge) {
	vertexIndex := edge.fromPosition + edge.characterLength
	endModeOrdinal := int(edge.getEndMode())
	if edges[vertexIndex][endModeOrdinal] == nil ||
		edges[vertexIndex][endModeOrdinal].cachedTotalSize > edge.cachedTotalSize {
		edges[vertexIndex][endModeOrdinal] = edge
	}
}

// MinimalEncoder_getNumberOfC40Words returns the number of words in which the string starting at from can be encoded in c40 or text mode.
// The number of characters encoded is returned in characterLength.
// The number of characters encoded is also minimal in the sense that the algorithm stops as soon
// as a character encoding fills a C40 word competely (three C40 values). An exception is at the
// end of the string where two C40 values are allowed (according to the spec the third c40 value
// is filled  with 0 (Shift 1) in this case).
func MinimalEncoder_getNumberOfC40Words(input *MinimalEncoderInput, from int, c40 bool, characterLength *int) int {
	thirdsCount := 0
	for i := from; i < input.Length(); i++ {
		if isECI, _ := input.IsECI(i); isECI {
			*characterLength = 0
			return 0
		}
		ci := input.CharAt(i)
		if (c40 && HighLevelEncoder_isNativeC40(ci)) || (!c40 && HighLevelEncoder_isNativeText(ci)) {
			thirdsCount++ // native
		} else if !MinimalEncoder_isExtendedASCII(ci, input.GetFNC1Character()) {
			thirdsCount += 2 // shift
		} else {
			asciiValue := int(ci) & 0xff
			if asciiValue >= 128 && ((c40 && HighLevelEncoder_isNativeC40(byte(asciiValue-128))) ||
				(!c40 && HighLevelEncoder_isNativeText(byte(asciiValue-128)))) {
				thirdsCount += 3 // shift, Upper shift
			} else {
				thirdsCount += 4 // shift, Upper shift, shift
			}
		}

		if thirdsCount%3 == 0 || ((thirdsCount-2)%3 == 0 && i+1 == input.Length()) {
			*characterLength = i - from + 1
			return int(math.Ceil(float64(thirdsCount) / 3.0))
		}
	}
	*characterLength = 0
	return 0
}

func MinimalEncoder_addEdges(input *MinimalEncoderInput, edges [][]*MinimalEncoderEdge, from int, previous *MinimalEncoderEdge) {
	if isECI, _ := input.IsECI(from); isECI {
		MinimalEncoder_addEdge(edges, newMinimalEncoderEdge(input, MinimalEncoderModeASCII, from, 1, previous))
		return
	}

	ch := input.CharAt(from)
	if previous == nil || previous.getEndMode() != MinimalEncoderModeEDF { // not possible to unlatch a full EDF edge to something else
		if HighLevelEncoder_isDigit(ch) && input.HaveNCharacters(from, 2) &&
			HighLevelEncoder_isDigit(input.CharAt(from+1)) {
			// two digits ASCII encoded
			MinimalEncoder_addEdge(edges, newMinimalEncoderEdge(input, MinimalEncoderModeASCII, from, 2, previous))
		} else {
			// one ASCII encoded character or an extended character via Upper Shift
			MinimalEncoder_addEdge(edges, newMinimalEncoderEdge(input, MinimalEncoderModeASCII, from, 1, previous))
		}

		modes := []MinimalEncoderMode{MinimalEncoderModeC40, MinimalEncoderModeTEXT}
		for _, mode := range modes {
			var characterLength int
			if MinimalEncoder_getNumberOfC40Words(input, from, mode == MinimalEncoderModeC40, &characterLength) > 0 {
				MinimalEncoder_addEdge(edges, newMinimalEncoderEdge(input, mode, from, characterLength, previous))
			}
		}

		if input.HaveNCharacters(from, 3) &&
			HighLevelEncoder_isNativeX12(input.CharAt(from)) &&
			HighLevelEncoder_isNativeX12(input.CharAt(from+1)) &&
			HighLevelEncoder_isNativeX12(input.CharAt(from+2)) {
			MinimalEncoder_addEdge(edges, newMinimalEncoderEdge(input, MinimalEncoderModeX12, from, 3, previous))
		}

		MinimalEncoder_addEdge(edges, newMinimalEncoderEdge(input, MinimalEncoderModeB256, from, 1, previous))
	}

	// We create 4 EDF edges, with 1, 2 3 or 4 characters length. The fourth normally doesn't have a latch to ASCII
	// unless it is 2 characters away from the end of the input.
	var i int
	for i = 0; i < 3; i++ {
		pos := from + i
		if input.HaveNCharacters(pos, 1) && HighLevelEncoder_isNativeEDIFACT(input.CharAt(pos)) {
			MinimalEncoder_addEdge(edges, newMinimalEncoderEdge(input, MinimalEncoderModeEDF, from, i+1, previous))
		} else {
			break
		}
	}
	if i == 3 && input.HaveNCharacters(from, 4) && HighLevelEncoder_isNativeEDIFACT(input.CharAt(from+3)) {
		MinimalEncoder_addEdge(edges, newMinimalEncoderEdge(input, MinimalEncoderModeEDF, from, 4, previous))
	}
}

func MinimalEncoder_encodeMinimally(input *MinimalEncoderInput) (*MinimalEncoderResult, error) {
	inputLength := input.Length()

	// Array that represents vertices. There is a vertex for every character and mode.
	// The last dimension in the array below encodes the 6 modes ASCII, C40, TEXT, X12, EDF and B256
	edges := make([][]*MinimalEncoderEdge, inputLength+1)
	for i := range edges {
		edges[i] = make([]*MinimalEncoderEdge, 6)
	}
	MinimalEncoder_addEdges(input, edges, 0, nil)

	for i := 1; i <= inputLength; i++ {
		for j := 0; j < 6; j++ {
			if edges[i][j] != nil && i < inputLength {
				MinimalEncoder_addEdges(input, edges, i, edges[i][j])
			}
		}
		// optimize memory by removing edges that have been passed.
		for j := 0; j < 6; j++ {
			edges[i-1][j] = nil
		}
	}

	minimalJ := -1
	minimalSize := math.MaxInt
	for j := 0; j < 6; j++ {
		if edges[inputLength][j] != nil {
			edge := edges[inputLength][j]
			size := edge.cachedTotalSize
			if j >= 1 && j <= 3 { // C40, TEXT and X12 need an extra unlatch at the end
				size++
			}
			if size < minimalSize {
				minimalSize = size
				minimalJ = j
			}
		}
	}

	if minimalJ < 0 {
		return nil, gozxing.NewWriterException("IllegalStateException: Failed to encode \"%s\"", input)
	}
	return newMinimalEncoderResult(edges[inputLength][minimalJ]), nil
}

type MinimalEncoderEdge struct {
	input           *MinimalEncoderInput
	mode            MinimalEncoderMode // the mode at the start of this edge.
	fromPosition    int
	characterLength int
	previous        *MinimalEncoderEdge
	cachedTotalSize int
}

var (
	minimalEncoderEdge_allCodewordCapacities = []int{
		3, 5, 8, 10, 12, 16, 18, 22, 30, 32, 36, 44, 49, 62, 86, 114,
		144, 174, 204, 280, 368, 456, 576, 696, 816, 1050, 1304, 1558,
	}
	minimalEncoderEdge_squareCodewordCapacities = []int{
		3, 5, 8, 12, 18, 22, 30, 36, 44, 62, 86, 114, 144, 174, 204,
		280, 368, 456, 576, 696, 816, 1050, 1304, 1558,
	}
	minimalEncoderEdge_rectangularCodewordCapacities = []int{5, 10, 16, 33, 32, 49}
)

func newMinimalEncoderEdge(input *MinimalEncoderInput, mode MinimalEncoderMode, fromPosition int, characterLength int, previous *MinimalEncoderEdge) *MinimalEncoderEdge {
	this := &MinimalEncoderEdge{
		input:           input,
		mode:            mode,
		fromPosition:    fromPosition,
		characterLength: characterLength,
		previous:        previous,
	}

	size := 0
	if previous != nil {
		size = previous.cachedTotalSize
	}

	previousMode := this.getPreviousMode()

	// Switching modes
	// ASCII -> C40: latch 230
	// ASCII -> TEXT: latch 239
	// ASCII -> X12: latch 238
	// ASCII -> EDF: latch 240
	// ASCII -> B256: latch 231
	// C40 -> ASCII: word(c1,c2,c3), 254
	// TEXT -> ASCII: word(c1,c2,c3), 254
	// X12 -> ASCII: word(c1,c2,c3), 254
	// EDIFACT -> ASCII: Unlatch character,0,0,0 or c1,Unlatch character,0,0 or c1,c2,Unlatch character,0 or
	// c1,c2,c3,Unlatch character
	// B256 -> ASCII: without latch after n bytes
	switch mode {
	case MinimalEncoderModeASCII:
		size++
		if isECI, _ := input.IsECI(fromPosition); isECI || MinimalEncoder_isExtendedASCII(input.CharAt(fromPosition), input.GetFNC1Character()) {
			size++
		}
		if previousMode == MinimalEncoderModeC40 ||
			previousMode == MinimalEncoderModeTEXT ||
			previousMode == MinimalEncoderModeX12 {
			size++ // unlatch 254 to ASCII
		}
	case MinimalEncoderModeB256:
		size++
		if previousMode != MinimalEncoderModeB256 {
			size++ // byte count
		} else if this.getB256Size() == 250 {
			size++ // extra byte count
		}
		if previousMode == MinimalEncoderModeASCII {
			size++ // latch to B256
		} else if previousMode == MinimalEncoderModeC40 ||
			previousMode == MinimalEncoderModeTEXT ||
			previousMode == MinimalEncoderModeX12 {
			size += 2 // unlatch to ASCII, latch to B256
		}
	case MinimalEncoderModeC40, MinimalEncoderModeTEXT, MinimalEncoderModeX12:
		if mode == MinimalEncoderModeX12 {
			size += 2
		} else {
			var charLen int
			size += MinimalEncoder_getNumberOfC40Words(input, fromPosition, mode == MinimalEncoderModeC40, &charLen) * 2
		}

		if previousMode == MinimalEncoderModeASCII || previousMode == MinimalEncoderModeB256 {
			size++ // additional byte for latch from ASCII to this mode
		} else if previousMode != mode && (previousMode == MinimalEncoderModeC40 ||
			previousMode == MinimalEncoderModeTEXT ||
			previousMode == MinimalEncoderModeX12) {
			size += 2 // unlatch 254 to ASCII followed by latch to this mode
		}
	case MinimalEncoderModeEDF:
		size += 3
		if previousMode == MinimalEncoderModeASCII || previousMode == MinimalEncoderModeB256 {
			size++ // additional byte for latch from ASCII to this mode
		} else if previousMode == MinimalEncoderModeC40 ||
			previousMode == MinimalEncoderModeASCII ||
			previousMode == MinimalEncoderModeX12 {
			size += 2 // unlatch 254 to ASCII followed by latch to this mode
		}
	}
	this.cachedTotalSize = size
	return this
}

// getB256Size does not count beyond 250
func (e *MinimalEncoderEdge) getB256Size() int {
	cnt := 0
	current := e
	for current != nil && current.mode == MinimalEncoderModeB256 && cnt <= 250 {
		cnt++
		current = current.previous
	}
	return cnt
}

func (e *MinimalEncoderEdge) getPreviousStartMode() MinimalEncoderMode {
	if e.previous == nil {
		return MinimalEncoderModeASCII
	}
	return e.previous.mode
}

func (e *MinimalEncoderEdge) getPreviousMode() MinimalEncoderMode {
	if e.previous == nil {
		return MinimalEncoderModeASCII
	}
	return e.previous.getEndMode()
}

// getEndMode Returns ModeASCII in case that:
//   - Mode is EDIFACT and characterLength is less than 4 or the remaining characters can be encoded in at most 2
//     ASCII bytes.
//   - Mode is C40, TEXT or X12 and the remaining characters can be encoded in at most 1 ASCII byte.
//     Returns mode in all other cases.
func (e *MinimalEncoderEdge) getEndMode() MinimalEncoderMode {
	if e.mode == MinimalEncoderModeEDF {
		if e.characterLength < 4 {
			return MinimalEncoderModeASCII
		}
		lastASCII := e.getLastASCII() // see 5.2.8.2 EDIFACT encodation Rules
		if lastASCII > 0 && e.getCodewordsRemaining(e.cachedTotalSize+lastASCII) <= 2-lastASCII {
			return MinimalEncoderModeASCII
		}
	}
	if e.mode == MinimalEncoderModeC40 ||
		e.mode == MinimalEncoderModeTEXT ||
		e.mode == MinimalEncoderModeX12 {

		// see 5.2.5.2 C40 encodation rules and 5.2.7.2 ANSI X12 encodation rules
		if e.fromPosition+e.characterLength >= e.input.Length() && e.getCodewordsRemaining(e.cachedTotalSize) == 0 {
			return MinimalEncoderModeASCII
		}
		lastASCII := e.getLastASCII()
		if lastASCII == 1 && e.getCodewordsRemaining(e.cachedTotalSize+1) == 0 {
			return MinimalEncoderModeASCII
		}
	}
	return e.mode
}

func (e *MinimalEncoderEdge) getMode() MinimalEncoderMode {
	return e.mode
}

// getLastASCII Peeks ahead and returns 1 if the postfix consists of exactly two digits, 2 if the postfix consists of exactly
// two consecutive digits and a non extended character or of 4 digits.
// Returns 0 in any other case
func (e *MinimalEncoderEdge) getLastASCII() int {
	length := e.input.Length()
	from := e.fromPosition + e.characterLength
	if length-from > 4 || from >= length {
		return 0
	}
	if length-from == 1 {
		if MinimalEncoder_isExtendedASCII(e.input.CharAt(from), e.input.GetFNC1Character()) {
			return 0
		}
		return 1
	}
	if length-from == 2 {
		if MinimalEncoder_isExtendedASCII(e.input.CharAt(from), e.input.GetFNC1Character()) ||
			MinimalEncoder_isExtendedASCII(e.input.CharAt(from+1), e.input.GetFNC1Character()) {
			return 0
		}
		if HighLevelEncoder_isDigit(e.input.CharAt(from)) && HighLevelEncoder_isDigit(e.input.CharAt(from+1)) {
			return 1
		}
		return 2
	}
	if length-from == 3 {
		if HighLevelEncoder_isDigit(e.input.CharAt(from)) && HighLevelEncoder_isDigit(e.input.CharAt(from+1)) &&
			!MinimalEncoder_isExtendedASCII(e.input.CharAt(from+2), e.input.GetFNC1Character()) {
			return 2
		}
		if HighLevelEncoder_isDigit(e.input.CharAt(from+1)) && HighLevelEncoder_isDigit(e.input.CharAt(from+2)) &&
			!MinimalEncoder_isExtendedASCII(e.input.CharAt(from), e.input.GetFNC1Character()) {
			return 2
		}
		return 0
	}
	if HighLevelEncoder_isDigit(e.input.CharAt(from)) && HighLevelEncoder_isDigit(e.input.CharAt(from+1)) &&
		HighLevelEncoder_isDigit(e.input.CharAt(from+2)) && HighLevelEncoder_isDigit(e.input.CharAt(from+3)) {
		return 2
	}
	return 0
}

// getMinSymbolSize Returns the capacity in codewords of the smallest symbol that has enough capacity to fit the given minimal
// number of codewords.
func (e *MinimalEncoderEdge) getMinSymbolSize(minimum int) int {
	switch e.input.getShapeHint() {
	case SymbolShapeHint_FORCE_SQUARE:
		for _, capacity := range minimalEncoderEdge_squareCodewordCapacities {
			if capacity >= minimum {
				return capacity
			}
		}
	case SymbolShapeHint_FORCE_RECTANGLE:
		for _, capacity := range minimalEncoderEdge_rectangularCodewordCapacities {
			if capacity >= minimum {
				return capacity
			}
		}
	}
	for _, capacity := range minimalEncoderEdge_allCodewordCapacities {
		if capacity >= minimum {
			return capacity
		}
	}
	return minimalEncoderEdge_allCodewordCapacities[len(minimalEncoderEdge_allCodewordCapacities)-1]
}

// getCodewordsRemaining Returns the remaining capacity in codewords of the smallest symbol that has enough capacity to fit the given
// minimal number of codewords.
func (e *MinimalEncoderEdge) getCodewordsRemaining(minimum int) int {
	return e.getMinSymbolSize(minimum) - minimum
}

func MinimalEncoderEdge_getBytes(c int) []byte {
	return []byte{byte(c)}
}

func MinimalEncoderEdge_getBytes2(c1, c2 int) []byte {
	return []byte{byte(c1), byte(c2)}
}

func MinimalEncoderEdge_setC40Word(bytes []byte, offset int, c1, c2, c3 int) {
	val16 := (1600 * (c1 & 0xff)) + (40 * (c2 & 0xff)) + (c3 & 0xff) + 1
	bytes[offset] = byte(val16 / 256)
	bytes[offset+1] = byte(val16 % 256)
}

func MinimalEncoderEdge_getX12Value(c byte) int {
	if c == 13 {
		return 0
	}
	if c == 42 {
		return 1
	}
	if c == 62 {
		return 2
	}
	if c == 32 {
		return 3
	}
	if c >= 48 && c <= 57 {
		return int(c) - 44
	}
	if c >= 65 && c <= 90 {
		return int(c) - 51
	}
	return int(c)
}

func (e *MinimalEncoderEdge) getX12Words() []byte {
	result := make([]byte, e.characterLength/3*2)
	for i := 0; i < len(result); i += 2 {
		MinimalEncoderEdge_setC40Word(result, i,
			MinimalEncoderEdge_getX12Value(e.input.CharAt(e.fromPosition+i/2*3)),
			MinimalEncoderEdge_getX12Value(e.input.CharAt(e.fromPosition+i/2*3+1)),
			MinimalEncoderEdge_getX12Value(e.input.CharAt(e.fromPosition+i/2*3+2)))
	}
	return result
}

func MinimalEncoderEdge_getShiftValue(c byte, c40 bool, fnc1 int) int {
	if (c40 && MinimalEncoder_isInC40Shift1Set(c)) || (!c40 && MinimalEncoder_isInTextShift1Set(c)) {
		return 0
	}
	if (c40 && MinimalEncoder_isInC40Shift2Set(c, fnc1)) || (!c40 && MinimalEncoder_isInTextShift2Set(c, fnc1)) {
		return 1
	}
	return 2
}

func MinimalEncoderEdge_getC40Value(c40 bool, setIndex int, c byte, fnc1 int) int {
	if int(c) == fnc1 {
		return 27
	}
	if c40 {
		if c <= 31 {
			return int(c)
		}
		if c == 32 {
			return 3
		}
		if c <= 47 {
			return int(c) - 33
		}
		if c <= 57 {
			return int(c) - 44
		}
		if c <= 64 {
			return int(c) - 43
		}
		if c <= 90 {
			return int(c) - 51
		}
		if c <= 95 {
			return int(c) - 69
		}
		if c <= 127 {
			return int(c) - 96
		}
		return int(c)
	} else {
		if c == 0 {
			return 0
		}
		if setIndex == 0 && c <= 3 {
			return int(c) - 1 // is this a bug in the spec?
		}
		if setIndex == 1 && c <= 31 {
			return int(c)
		}
		if c == 32 {
			return 3
		}
		if c >= 33 && c <= 47 {
			return int(c) - 33
		}
		if c >= 48 && c <= 57 {
			return int(c) - 44
		}
		if c >= 58 && c <= 64 {
			return int(c) - 43
		}
		if c >= 65 && c <= 90 {
			return int(c) - 64
		}
		if c >= 91 && c <= 95 {
			return int(c) - 69
		}
		if c == 96 {
			return 0
		}
		if c >= 97 && c <= 122 {
			return int(c) - 83
		}
		if c >= 123 && c <= 127 {
			return int(c) - 96
		}
		return int(c)
	}
}

func (e *MinimalEncoderEdge) getC40Words(c40 bool, fnc1 int) []byte {
	c40Values := make([]byte, 0)
	for i := 0; i < e.characterLength; i++ {
		ci := e.input.CharAt(e.fromPosition + i)
		if (c40 && HighLevelEncoder_isNativeC40(ci)) || (!c40 && HighLevelEncoder_isNativeText(ci)) {
			c40Values = append(c40Values, byte(MinimalEncoderEdge_getC40Value(c40, 0, ci, fnc1)))
		} else if !MinimalEncoder_isExtendedASCII(ci, fnc1) {
			shiftValue := MinimalEncoderEdge_getShiftValue(ci, c40, fnc1)
			c40Values = append(c40Values, byte(shiftValue)) // Shift[123]
			c40Values = append(c40Values, byte(MinimalEncoderEdge_getC40Value(c40, shiftValue, ci, fnc1)))
		} else {
			asciiValue := byte((int(ci) & 0xff) - 128)
			if (c40 && HighLevelEncoder_isNativeC40(asciiValue)) ||
				!c40 && HighLevelEncoder_isNativeText(asciiValue) {
				c40Values = append(c40Values, 1)  // Shift 2
				c40Values = append(c40Values, 30) // Upper Shift
				c40Values = append(c40Values, byte(MinimalEncoderEdge_getC40Value(c40, 0, asciiValue, fnc1)))
			} else {
				c40Values = append(c40Values, 1)  // Shift 2
				c40Values = append(c40Values, 30) // Upper Shift
				shiftValue := MinimalEncoderEdge_getShiftValue(asciiValue, c40, fnc1)
				c40Values = append(c40Values, byte(shiftValue)) // Shift[123]
				c40Values = append(c40Values, byte(MinimalEncoderEdge_getC40Value(c40, shiftValue, asciiValue, fnc1)))
			}
		}
	}

	if len(c40Values)%3 != 0 {
		// pad with 0 (Shift 1)
		c40Values = append(c40Values, 0)
	}

	result := make([]byte, len(c40Values)/3*2)
	byteIndex := 0
	for i := 0; i < len(c40Values); i += 3 {
		MinimalEncoderEdge_setC40Word(result, byteIndex, int(c40Values[i]&0xff), int(c40Values[i+1]&0xff), int(c40Values[i+2]&0xff))
		byteIndex += 2
	}
	return result
}

func (e *MinimalEncoderEdge) getEDFBytes() []byte {
	numberOfThirds := int(math.Ceil(float64(e.characterLength) / 4.0))
	result := make([]byte, numberOfThirds*3)
	pos := e.fromPosition
	endPos := e.fromPosition + e.characterLength - 1
	if endPos > e.input.Length()-1 {
		endPos = e.input.Length() - 1
	}
	for i := 0; i < numberOfThirds; i += 3 {
		edfValues := make([]int, 4)
		for j := 0; j < 4; j++ {
			if pos <= endPos {
				edfValues[j] = int(e.input.CharAt(pos)) & 0x3f
				pos++
			} else {
				if pos == endPos+1 {
					edfValues[j] = 0x1f
				} else {
					edfValues[j] = 0
				}
			}
		}
		val24 := edfValues[0]<<18 | edfValues[1]<<12 | edfValues[2]<<6 | edfValues[3]
		result[i] = byte((val24 >> 16) & 0xff)
		result[i+1] = byte((val24 >> 8) & 0xff)
		result[i+2] = byte(val24 & 0xff)
	}
	return result
}

func (e *MinimalEncoderEdge) getLatchBytes() []byte {
	switch e.getPreviousMode() {
	case MinimalEncoderModeASCII, MinimalEncoderModeB256: // after B256 ends (via length) we are back to ASCII
		switch e.mode {
		case MinimalEncoderModeB256:
			return MinimalEncoderEdge_getBytes(231)
		case MinimalEncoderModeC40:
			return MinimalEncoderEdge_getBytes(230)
		case MinimalEncoderModeTEXT:
			return MinimalEncoderEdge_getBytes(239)
		case MinimalEncoderModeX12:
			return MinimalEncoderEdge_getBytes(238)
		case MinimalEncoderModeEDF:
			return MinimalEncoderEdge_getBytes(240)
		}
	case MinimalEncoderModeC40, MinimalEncoderModeTEXT, MinimalEncoderModeX12:
		if e.mode != e.getPreviousMode() {
			switch e.mode {
			case MinimalEncoderModeASCII:
				return MinimalEncoderEdge_getBytes(254)
			case MinimalEncoderModeB256:
				return MinimalEncoderEdge_getBytes2(254, 231)
			case MinimalEncoderModeC40:
				return MinimalEncoderEdge_getBytes2(254, 230)
			case MinimalEncoderModeTEXT:
				return MinimalEncoderEdge_getBytes2(254, 239)
			case MinimalEncoderModeX12:
				return MinimalEncoderEdge_getBytes2(254, 238)
			case MinimalEncoderModeEDF:
				return MinimalEncoderEdge_getBytes2(254, 240)
			}
		}
	case MinimalEncoderModeEDF:
		// The rightmost EDIFACT edge always contains an unlatch character
		break
	}
	return []byte{}
}

// getDataBytes Important: The function does not return the length bytes (one or two) in case of B256 encoding
func (e *MinimalEncoderEdge) getDataBytes() []byte {
	switch e.mode {
	case MinimalEncoderModeASCII:
		if isECI, _ := e.input.IsECI(e.fromPosition); isECI {
			eciValue, _ := e.input.GetECIValue(e.fromPosition)
			return MinimalEncoderEdge_getBytes2(241, eciValue+1)
		}
		if MinimalEncoder_isExtendedASCII(e.input.CharAt(e.fromPosition), e.input.GetFNC1Character()) {
			return MinimalEncoderEdge_getBytes2(235, int(e.input.CharAt(e.fromPosition))-127)
		}
		if e.characterLength == 2 {
			return MinimalEncoderEdge_getBytes((int(e.input.CharAt(e.fromPosition))-'0')*10 + int(e.input.CharAt(e.fromPosition+1)) - '0' + 130)
		}
		if isFNC1, _ := e.input.IsFNC1(e.fromPosition); isFNC1 {
			return MinimalEncoderEdge_getBytes(232)
		}
		return MinimalEncoderEdge_getBytes(int(e.input.CharAt(e.fromPosition)) + 1)
	case MinimalEncoderModeB256:
		return MinimalEncoderEdge_getBytes(int(e.input.CharAt(e.fromPosition)))
	case MinimalEncoderModeC40:
		return e.getC40Words(true, e.input.GetFNC1Character())
	case MinimalEncoderModeTEXT:
		return e.getC40Words(false, e.input.GetFNC1Character())
	case MinimalEncoderModeX12:
		return e.getX12Words()
	case MinimalEncoderModeEDF:
		return e.getEDFBytes()
	}
	return []byte{}
}

type MinimalEncoderResult struct {
	bytes []byte
}

// ここまで書き換え済み
func newMinimalEncoderResult(solution *MinimalEncoderEdge) *MinimalEncoderResult {
	input := solution.input
	size := 0
	bytesAL := make([]byte, 0)
	randomizePostfixLength := make([]int, 0)
	randomizeLengths := make([]int, 0)
	if (solution.mode == MinimalEncoderModeC40 ||
		solution.mode == MinimalEncoderModeTEXT ||
		solution.mode == MinimalEncoderModeX12) &&
		solution.getEndMode() != MinimalEncoderModeASCII {
		size += MinimalEncoderResult_prepend(MinimalEncoderEdge_getBytes(254), &bytesAL)
	}
	current := solution
	for current != nil {
		size += MinimalEncoderResult_prepend(current.getDataBytes(), &bytesAL)

		if current.previous == nil || current.getPreviousStartMode() != current.getMode() {
			if current.getMode() == MinimalEncoderModeB256 {
				if size <= 249 {
					bytesAL = append([]byte{byte(size)}, bytesAL...)
					size++
				} else {
					bytesAL = append([]byte{byte(size % 250), byte(size/250 + 249)}, bytesAL...)
					size += 2
				}
				randomizePostfixLength = append(randomizePostfixLength, len(bytesAL))
				randomizeLengths = append(randomizeLengths, size)
			}
			MinimalEncoderResult_prepend(current.getLatchBytes(), &bytesAL)
			size = 0
		}

		current = current.previous
	}
	if input.getMacroId() == 5 {
		size += MinimalEncoderResult_prepend(MinimalEncoderEdge_getBytes(236), &bytesAL)
	} else if input.getMacroId() == 6 {
		size += MinimalEncoderResult_prepend(MinimalEncoderEdge_getBytes(237), &bytesAL)
	}

	if input.GetFNC1Character() > 0 {
		size += MinimalEncoderResult_prepend(MinimalEncoderEdge_getBytes(232), &bytesAL)
	}
	for i := 0; i < len(randomizePostfixLength); i++ {
		MinimalEncoderResult_applyRandomPattern(&bytesAL, len(bytesAL)-randomizePostfixLength[i], randomizeLengths[i])
	}
	// add padding
	capacity := solution.getMinSymbolSize(len(bytesAL))
	if len(bytesAL) < capacity {
		bytesAL = append(bytesAL, 129)
	}
	for len(bytesAL) < capacity {
		bytesAL = append(bytesAL, byte(MinimalEncoderResult_randomize253State(len(bytesAL)+1)))
	}

	return &MinimalEncoderResult{bytes: bytesAL}
}

func MinimalEncoderResult_prepend(bytes []byte, into *[]byte) int {
	for i := len(bytes) - 1; i >= 0; i-- {
		*into = append([]byte{bytes[i]}, *into...)
	}
	return len(bytes)
}

func MinimalEncoderResult_randomize253State(codewordPosition int) int {
	pseudoRandom := ((149 * codewordPosition) % 253) + 1
	tempVariable := 129 + pseudoRandom
	if tempVariable <= 254 {
		return tempVariable
	}
	return tempVariable - 254
}

func MinimalEncoderResult_applyRandomPattern(bytesAL *[]byte, startPosition int, length int) {
	for i := 0; i < length; i++ {
		// See "B.1 253-state algorithm
		padCodewordPosition := startPosition + i
		padCodewordValue := int((*bytesAL)[padCodewordPosition] & 0xff)
		pseudoRandomNumber := ((149 * (padCodewordPosition + 1)) % 255) + 1
		tempVariable := padCodewordValue + pseudoRandomNumber
		if tempVariable <= 255 {
			(*bytesAL)[padCodewordPosition] = byte(tempVariable)
		} else {
			(*bytesAL)[padCodewordPosition] = byte(tempVariable - 256)
		}
	}
}

func (r *MinimalEncoderResult) getBytes() []byte {
	return r.bytes
}

type MinimalEncoderInput struct {
	*common.MinimalECIInput
	shape   SymbolShapeHint
	macroId int
}

func newMinimalEncoderInput(stringToEncode []rune, priorityCharset encoding.Encoding, fnc1 int, shape SymbolShapeHint, macroId int) (*MinimalEncoderInput, error) {
	in, err := common.NewMinimalECIInput(stringToEncode, priorityCharset, fnc1)
	if err != nil {
		return nil, gozxing.WrapWriterException(err)
	}
	return &MinimalEncoderInput{
		MinimalECIInput: in,
		shape:           shape,
		macroId:         macroId,
	}, nil
}

func (in *MinimalEncoderInput) CharAt(index int) byte {
	// no error occurs: the index range and IsECI have been validated.
	c, _ := in.MinimalECIInput.CharAt(index)
	return byte(c)
}

func (in *MinimalEncoderInput) getMacroId() int {
	return in.macroId
}

func (in *MinimalEncoderInput) getShapeHint() SymbolShapeHint {
	return in.shape
}
