package encoder

import (
	"fmt"

	"github.com/makiuchi-d/gozxing"
)

// State represents all information about a sequence necessary to generate the current output.
// Note that a state is immutable.
type State struct {
	mode                int
	token               Token
	binaryShiftByteCount int
	bitCount            int
	binaryShiftCost     int
}

var InitialState = NewState(GetEmptyToken(), MODE_UPPER, 0, 0)

func NewState(token Token, mode, binaryBytes, bitCount int) *State {
	return &State{
		token:               token,
		mode:                mode,
		binaryShiftByteCount: binaryBytes,
		bitCount:            bitCount,
		binaryShiftCost:     calculateBinaryShiftCost(binaryBytes),
	}
}

func (this *State) GetMode() int {
	return this.mode
}

func (this *State) GetToken() Token {
	return this.token
}

func (this *State) GetBinaryShiftByteCount() int {
	return this.binaryShiftByteCount
}

func (this *State) GetBitCount() int {
	return this.bitCount
}

func (this *State) AppendFLGn(eci int) *State {
	result := this.ShiftAndAppend(MODE_PUNCT, 0) // 0: FLG(n)
	token := result.token
	bitsAdded := 3
	if eci < 0 {
		token = token.Add(0, 3) // 0: FNC1
	} else if eci > 999999 {
		panic("ECI code must be between 0 and 999999")
	} else {
		eciStr := fmt.Sprintf("%d", eci)
		eciDigits := []byte(eciStr)
		token = token.Add(len(eciDigits), 3) // 1-6: number of ECI digits
		for _, eciDigit := range eciDigits {
			token = token.Add(int(eciDigit)-'0'+2, 4)
		}
		bitsAdded += len(eciDigits) * 4
	}
	return NewState(token, this.mode, 0, this.bitCount+bitsAdded)
}

// LatchAndAppend Create a new state representing this state with a latch to a (not
// necessary different) mode, and then a code.
func (this *State) LatchAndAppend(mode, value int) *State {
	bitCount := this.bitCount
	token := this.token
	if mode != this.mode {
		latch := LATCH_TABLE[this.mode][mode]
		token = token.Add(latch&0xFFFF, latch>>16)
		bitCount += latch >> 16
	}
	latchModeBitCount := 5
	if mode == MODE_DIGIT {
		latchModeBitCount = 4
	}
	token = token.Add(value, latchModeBitCount)
	return NewState(token, mode, 0, bitCount+latchModeBitCount)
}

// ShiftAndAppend Create a new state representing this state, with a temporary shift
// to a different mode to output a single value.
func (this *State) ShiftAndAppend(mode, value int) *State {
	token := this.token
	thisModeBitCount := 5
	if this.mode == MODE_DIGIT {
		thisModeBitCount = 4
	}
	// Shifts exist only to UPPER and PUNCT, both with tokens size 5.
	token = token.Add(SHIFT_TABLE[this.mode][mode], thisModeBitCount)
	token = token.Add(value, 5)
	return NewState(token, this.mode, 0, this.bitCount+thisModeBitCount+5)
}

// AddBinaryShiftChar Create a new state representing this state, but an additional character
// output in Binary Shift mode.
func (this *State) AddBinaryShiftChar(index int) *State {
	token := this.token
	mode := this.mode
	bitCount := this.bitCount
	if this.mode == MODE_PUNCT || this.mode == MODE_DIGIT {
		latch := LATCH_TABLE[mode][MODE_UPPER]
		token = token.Add(latch&0xFFFF, latch>>16)
		bitCount += latch >> 16
		mode = MODE_UPPER
	}
	deltaBitCount := 8
	if this.binaryShiftByteCount == 0 || this.binaryShiftByteCount == 31 {
		deltaBitCount = 18
	} else if this.binaryShiftByteCount == 62 {
		deltaBitCount = 9
	}
	result := NewState(token, mode, this.binaryShiftByteCount+1, bitCount+deltaBitCount)
	if result.binaryShiftByteCount == 2047+31 {
		// The string is as long as it's allowed to be.  We should end it.
		result = result.EndBinaryShift(index + 1)
	}
	return result
}

// EndBinaryShift Create the state identical to this one, but we are no longer in
// Binary Shift mode.
func (this *State) EndBinaryShift(index int) *State {
	if this.binaryShiftByteCount == 0 {
		return this
	}
	token := this.token
	token = token.AddBinaryShift(index-this.binaryShiftByteCount, this.binaryShiftByteCount)
	return NewState(token, this.mode, 0, this.bitCount)
}

// IsBetterThanOrEqualTo Returns true if "this" state is better (or equal) to be in than "that"
// state under all possible circumstances.
func (this *State) IsBetterThanOrEqualTo(other *State) bool {
	newModeBitCount := this.bitCount + (LATCH_TABLE[this.mode][other.mode] >> 16)
	if this.binaryShiftByteCount < other.binaryShiftByteCount {
		// add additional B/S encoding cost of other, if any
		newModeBitCount += other.binaryShiftCost - this.binaryShiftCost
	} else if this.binaryShiftByteCount > other.binaryShiftByteCount && other.binaryShiftByteCount > 0 {
		// maximum possible additional cost (we end up exceeding the 31 byte boundary and other state can stay beneath it)
		newModeBitCount += 10
	}
	return newModeBitCount <= other.bitCount
}

func (this *State) ToBitArray(text []byte) *gozxing.BitArray {
	symbols := make([]Token, 0)
	endState := this.EndBinaryShift(len(text))
	// Traverse the token chain, stopping when we hit the empty token (previous is nil)
	for token := endState.token; token != nil && token.GetPrevious() != nil; token = token.GetPrevious() {
		symbols = append(symbols, token)
	}
	bitArray := gozxing.NewEmptyBitArray()
	// Add each token to the result in forward order
	for i := len(symbols) - 1; i >= 0; i-- {
		symbols[i].AppendTo(bitArray, text)
	}
	return bitArray
}

func (this *State) String() string {
	return fmt.Sprintf("%s bits=%d bytes=%d", MODE_NAMES[this.mode], this.bitCount, this.binaryShiftByteCount)
}

func calculateBinaryShiftCost(binaryShiftByteCount int) int {
	if binaryShiftByteCount > 62 {
		return 21 // B/S with extended length
	}
	if binaryShiftByteCount > 31 {
		return 20 // two B/S
	}
	if binaryShiftByteCount > 0 {
		return 10 // one B/S
	}
	return 0
}

