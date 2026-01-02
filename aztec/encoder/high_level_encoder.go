package encoder

import (
	"fmt"

	"golang.org/x/text/encoding"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
)

var (
	MODE_NAMES = []string{"UPPER", "LOWER", "DIGIT", "MIXED", "PUNCT"}

	MODE_UPPER = 0 // 5 bits
	MODE_LOWER = 1 // 5 bits
	MODE_DIGIT = 2 // 4 bits
	MODE_MIXED = 3 // 5 bits
	MODE_PUNCT = 4 // 5 bits
)

// The Latch Table shows, for each pair of Modes, the optimal method for
// getting from one mode to another.  In the worst possible case, this can
// be up to 14 bits.  In the best possible case, we are already there!
// The high half-word of each entry gives the number of bits.
// The low half-word of each entry are the actual bits necessary to change
var LATCH_TABLE = [][]int{
	{
		0,
		(5 << 16) + 28,              // UPPER -> LOWER
		(5 << 16) + 30,              // UPPER -> DIGIT
		(5 << 16) + 29,              // UPPER -> MIXED
		(10 << 16) + (29 << 5) + 30, // UPPER -> MIXED -> PUNCT
	},
	{
		(9 << 16) + (30 << 4) + 14, // LOWER -> DIGIT -> UPPER
		0,
		(5 << 16) + 30,              // LOWER -> DIGIT
		(5 << 16) + 29,              // LOWER -> MIXED
		(10 << 16) + (29 << 5) + 30, // LOWER -> MIXED -> PUNCT
	},
	{
		(4 << 16) + 14,             // DIGIT -> UPPER
		(9 << 16) + (14 << 5) + 28, // DIGIT -> UPPER -> LOWER
		0,
		(9 << 16) + (14 << 5) + 29, // DIGIT -> UPPER -> MIXED
		(14 << 16) + (14 << 10) + (29 << 5) + 30,
		// DIGIT -> UPPER -> MIXED -> PUNCT
	},
	{
		(5 << 16) + 29,              // MIXED -> UPPER
		(5 << 16) + 28,              // MIXED -> LOWER
		(10 << 16) + (29 << 5) + 30, // MIXED -> UPPER -> DIGIT
		0,
		(5 << 16) + 30, // MIXED -> PUNCT
	},
	{
		(5 << 16) + 31,              // PUNCT -> UPPER
		(10 << 16) + (31 << 5) + 28, // PUNCT -> UPPER -> LOWER
		(10 << 16) + (31 << 5) + 30, // PUNCT -> UPPER -> DIGIT
		(10 << 16) + (31 << 5) + 29, // PUNCT -> UPPER -> MIXED
		0,
	},
}

// A reverse mapping from [mode][char] to the encoding for that character
// in that mode.  An entry of 0 indicates no mapping exists.
var CHAR_MAP [5][256]int

// GetCharMap is exported for testing
func GetCharMap(mode, ch int) int {
	return CHAR_MAP[mode][ch]
}

func init() {
	// Initialize CHAR_MAP
	CHAR_MAP[MODE_UPPER][' '] = 1
	for c := 'A'; c <= 'Z'; c++ {
		CHAR_MAP[MODE_UPPER][c] = int(c) - 'A' + 2
	}
	CHAR_MAP[MODE_LOWER][' '] = 1
	for c := 'a'; c <= 'z'; c++ {
		CHAR_MAP[MODE_LOWER][c] = int(c) - 'a' + 2
	}
	CHAR_MAP[MODE_DIGIT][' '] = 1
	for c := '0'; c <= '9'; c++ {
		CHAR_MAP[MODE_DIGIT][c] = int(c) - '0' + 2
	}
	CHAR_MAP[MODE_DIGIT][','] = 12
	CHAR_MAP[MODE_DIGIT]['.'] = 13
	mixedTable := []byte{
		0, ' ', 1, 2, 3, 4, 5, 6, 7, '\b', '\t', '\n',
		11, '\f', '\r', 27, 28, 29, 30, 31, '@', '\\', '^',
		'_', '`', '|', '~', 127,
	}
	for i := 0; i < len(mixedTable); i++ {
		CHAR_MAP[MODE_MIXED][mixedTable[i]] = i
	}
	punctTable := []byte{
		0, '\r', 0, 0, 0, 0, '!', '"', '#', '$', '%', '&', '\'',
		'(', ')', '*', '+', ',', '-', '.', '/', ':', ';', '<', '=', '>', '?',
		'[', ']', '{', '}',
	}
	for i := 0; i < len(punctTable); i++ {
		if punctTable[i] > 0 {
			CHAR_MAP[MODE_PUNCT][punctTable[i]] = i
		}
	}
}

// A map showing the available shift codes.  (The shifts to BINARY are not
// shown
var SHIFT_TABLE [6][6]int

func init() {
	// Initialize SHIFT_TABLE
	for i := range SHIFT_TABLE {
		for j := range SHIFT_TABLE[i] {
			SHIFT_TABLE[i][j] = -1
		}
	}
	SHIFT_TABLE[MODE_UPPER][MODE_PUNCT] = 0
	SHIFT_TABLE[MODE_LOWER][MODE_PUNCT] = 0
	SHIFT_TABLE[MODE_LOWER][MODE_UPPER] = 28
	SHIFT_TABLE[MODE_MIXED][MODE_PUNCT] = 0
	SHIFT_TABLE[MODE_DIGIT][MODE_PUNCT] = 0
	SHIFT_TABLE[MODE_DIGIT][MODE_UPPER] = 15
}

type HighLevelEncoder struct {
	text    []byte
	charset encoding.Encoding
}

func NewHighLevelEncoder(text []byte) *HighLevelEncoder {
	return &HighLevelEncoder{
		text:    text,
		charset: nil,
	}
}

func NewHighLevelEncoderWithCharset(text []byte, charset encoding.Encoding) *HighLevelEncoder {
	return &HighLevelEncoder{
		text:    text,
		charset: charset,
	}
}

// Encode returns text represented by this encoder encoded as a BitArray
func (this *HighLevelEncoder) Encode() *gozxing.BitArray {
	initialState := InitialState
	if this.charset != nil {
		eci, ok := common.GetCharacterSetECI(this.charset)
		if !ok {
			panic("No ECI code for character set")
		}
		initialState = initialState.AppendFLGn(eci.GetValue())
	}
	states := []*State{initialState}
	for index := 0; index < len(this.text); index++ {
		var pairCode int
		var nextChar byte
		if index+1 < len(this.text) {
			nextChar = this.text[index+1]
		}
		switch this.text[index] {
		case '\r':
			if nextChar == '\n' {
				pairCode = 2
			}
		case '.':
			if nextChar == ' ' {
				pairCode = 3
			}
		case ',':
			if nextChar == ' ' {
				pairCode = 4
			}
		case ':':
			if nextChar == ' ' {
				pairCode = 5
			}
		}
		if pairCode > 0 {
			// We have one of the four special PUNCT pairs.  Treat them specially.
			// Get a new set of states for the two new characters.
			states = this.updateStateListForPair(states, index, pairCode)
			index++
		} else {
			// Get a new set of states for the new character.
			states = this.updateStateListForChar(states, index)
		}
	}
	// We are left with a set of states.  Find the shortest one.
	if len(states) == 0 {
		panic("No states after encoding")
	}
	minState := states[0]
	for _, state := range states[1:] {
		if state.GetBitCount() < minState.GetBitCount() {
			minState = state
		}
	}
	// Convert it to a bit array, and return.
	return minState.ToBitArray(this.text)
}

// We update a set of states for a new character by updating each state
// for the new character, merging the results, and then removing the
// non-optimal states.
func (this *HighLevelEncoder) updateStateListForChar(states []*State, index int) []*State {
	result := make([]*State, 0)
	for _, state := range states {
		this.updateStateForChar(state, index, &result)
	}
	if len(result) == 0 {
		// This should never happen, but if it does, we have a bug
		panic(fmt.Sprintf("updateStateListForChar: no states created for char at index %d (char=%c, states input=%d)",
			index, this.text[index], len(states)))
	}
	simplified := this.simplifyStates(result)
	if len(simplified) == 0 {
		panic(fmt.Sprintf("simplifyStates removed all states for char at index %d", index))
	}
	return simplified
}

// Return a set of states that represent the possible ways of updating this
// state for the next character.  The resulting set of states are added to
// the "result" list.
func (this *HighLevelEncoder) updateStateForChar(state *State, index int, result *[]*State) {
	ch := this.text[index] & 0xFF
	charInCurrentTable := CHAR_MAP[state.GetMode()][ch] > 0
	var stateNoBinary *State
	for mode := 0; mode <= MODE_PUNCT; mode++ {
		charInMode := CHAR_MAP[mode][ch]
		if charInMode > 0 {
			if stateNoBinary == nil {
				// Only create stateNoBinary the first time it's required.
				stateNoBinary = state.EndBinaryShift(index)
			}
			// Try generating the character by latching to its mode
			if !charInCurrentTable || mode == state.GetMode() || mode == MODE_DIGIT {
				// If the character is in the current table, we don't want to latch to
				// any other mode except possibly digit (which uses only 4 bits).  Any
				// other latch would be equally successful *after* this character, and
				// so wouldn't save any bits.
				latchState := stateNoBinary.LatchAndAppend(mode, charInMode)
				*result = append(*result, latchState)
			}
			// Try generating the character by switching to its mode.
			if !charInCurrentTable && SHIFT_TABLE[state.GetMode()][mode] >= 0 {
				// It never makes sense to temporarily shift to another mode if the
				// character exists in the current mode.  That can never save bits.
				shiftState := stateNoBinary.ShiftAndAppend(mode, charInMode)
				*result = append(*result, shiftState)
			}
		}
	}
	if state.GetBinaryShiftByteCount() > 0 || CHAR_MAP[state.GetMode()][ch] == 0 {
		// It's never worthwhile to go into binary shift mode if you're not already
		// in binary shift mode, and the character exists in your current mode.
		// That can never save bits over just outputting the char in the current mode.
		binaryState := state.AddBinaryShiftChar(index)
		*result = append(*result, binaryState)
	}
}

func (this *HighLevelEncoder) updateStateListForPair(states []*State, index int, pairCode int) []*State {
	result := make([]*State, 0)
	for _, state := range states {
		this.updateStateForPair(state, index, pairCode, &result)
	}
	return this.simplifyStates(result)
}

func (this *HighLevelEncoder) updateStateForPair(state *State, index int, pairCode int, result *[]*State) {
	stateNoBinary := state.EndBinaryShift(index)
	// Possibility 1.  Latch to MODE_PUNCT, and then append this code
	*result = append(*result, stateNoBinary.LatchAndAppend(MODE_PUNCT, pairCode))
	if state.GetMode() != MODE_PUNCT {
		// Possibility 2.  Shift to MODE_PUNCT, and then append this code.
		// Every state except MODE_PUNCT (handled above) can shift
		*result = append(*result, stateNoBinary.ShiftAndAppend(MODE_PUNCT, pairCode))
	}
	if pairCode == 3 || pairCode == 4 {
		// both characters are in DIGITS.  Sometimes better to just add two digits
		digitState := stateNoBinary.
			LatchAndAppend(MODE_DIGIT, 16-pairCode). // period or comma in DIGIT
			LatchAndAppend(MODE_DIGIT, 1)            // space in DIGIT
		*result = append(*result, digitState)
	}
	if state.GetBinaryShiftByteCount() > 0 {
		// It only makes sense to do the characters as binary if we're already
		// in binary mode.
		binaryState := state.AddBinaryShiftChar(index).AddBinaryShiftChar(index + 1)
		*result = append(*result, binaryState)
	}
}

func (this *HighLevelEncoder) simplifyStates(states []*State) []*State {
	// Match Java implementation exactly: use a deque-like structure
	// Java uses LinkedList with addFirst, which we simulate by prepending
	result := make([]*State, 0)
	for _, newState := range states {
		add := true
		// Iterate through existing states and compare
		for i := 0; i < len(result); {
			oldState := result[i]
			if oldState.IsBetterThanOrEqualTo(newState) {
				add = false
				break
			}
			if newState.IsBetterThanOrEqualTo(oldState) {
				// Remove oldState - it's worse than newState
				result = append(result[:i], result[i+1:]...)
				// Don't increment i, since we removed an element
				continue
			}
			i++
		}
		if add {
			// Add to front (like Java's addFirst) to match Java behavior
			result = append([]*State{newState}, result...)
		}
	}
	return result
}
