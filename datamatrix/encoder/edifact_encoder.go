package encoder

import (
	"errors"
)

type EdifactEncoder struct{}

func NewEdifactEncoder() Encoder {
	return EdifactEncoder{}
}

func (this EdifactEncoder) getEncodingMode() int {
	return HighLevelEncoder_EDIFACT_ENCODATION
}

func (this EdifactEncoder) encode(context *EncoderContext) error {
	//step F
	buffer := make([]byte, 0)
	for context.HasMoreCharacters() {
		c := context.GetCurrentChar()
		var e error
		buffer, e = edifactEncodeChar(c, buffer)
		if e != nil {
			return e
		}
		context.pos++

		count := len(buffer)
		if count >= 4 {
			codewords, _ := edifactEncodeToCodewords(buffer, 0)
			context.WriteCodewords(codewords)
			buffer = buffer[4:]

			newMode := HighLevelEncoder_lookAheadTest(context.GetMessage(), context.pos, this.getEncodingMode())
			if newMode != this.getEncodingMode() {
				// Return to ASCII encodation, which will actually handle latch to new mode
				context.SignalEncoderChange(HighLevelEncoder_ASCII_ENCODATION)
				break
			}
		}
	}
	buffer = append(buffer, 31) //Unlatch
	return edifactHandleEOD(context, buffer)
}

// edifactHandleEOD Handle "end of data" situations
//
// @param context the encoder context
// @param buffer  the buffer with the remaining encoded characters
//
func edifactHandleEOD(context *EncoderContext, buffer []byte) error {
	defer context.SignalEncoderChange(HighLevelEncoder_ASCII_ENCODATION)

	count := len(buffer)
	if count == 0 {
		return nil //Already finished
	}
	if count == 1 {
		//Only an unlatch at the end
		e := context.UpdateSymbolInfo()
		if e != nil {
			return e
		}

		available := context.GetSymbolInfo().GetDataCapacity() - context.GetCodewordCount()
		remaining := context.GetRemainingCharacters()
		// The following two lines are a hack inspired by the 'fix' from https://sourceforge.net/p/barcode4j/svn/221/
		if remaining > available {
			e := context.UpdateSymbolInfoByLength(context.GetCodewordCount() + 1)
			if e != nil {
				return e
			}
			available = context.GetSymbolInfo().GetDataCapacity() - context.GetCodewordCount()
		}
		if remaining <= available && available <= 2 {
			return nil //No unlatch
		}
	}

	if count > 4 {
		return errors.New("IllegalStateException: Count must not exceed 4")
	}
	restChars := count - 1
	encoded, _ := edifactEncodeToCodewords(buffer, 0)
	endOfSymbolReached := !context.HasMoreCharacters()
	restInAscii := endOfSymbolReached && restChars <= 2

	if restChars <= 2 {
		e := context.UpdateSymbolInfoByLength(context.GetCodewordCount() + restChars)
		if e != nil {
			return e
		}
		available := context.GetSymbolInfo().GetDataCapacity() - context.GetCodewordCount()
		if available >= 3 {
			restInAscii = false
			e := context.UpdateSymbolInfoByLength(context.GetCodewordCount() + len(encoded))
			if e != nil {
				return e
			}
			//available = context.symbolInfo.dataCapacity - context.getCodewordCount();
		}
	}

	if restInAscii {
		context.ResetSymbolInfo()
		context.pos -= restChars
	} else {
		context.WriteCodewords(encoded)
	}

	return nil
}

func edifactEncodeChar(c byte, sb []byte) ([]byte, error) {
	if c >= ' ' && c <= '?' {
		sb = append(sb, c)
	} else if c >= '@' && c <= '^' {
		sb = append(sb, c-64)
	} else {
		return sb, illegalCharacter(c)
	}
	return sb, nil
}

func edifactEncodeToCodewords(sb []byte, startPos int) ([]byte, error) {
	len := len(sb) - startPos
	if len == 0 {
		return sb, errors.New("IllegalStateException: StringBuilder must not be empty")
	}
	c1 := int(sb[startPos])
	c2 := 0
	if len >= 2 {
		c2 = int(sb[startPos+1])
	}
	c3 := 0
	if len >= 3 {
		c3 = int(sb[startPos+2])
	}
	c4 := 0
	if len >= 4 {
		c4 = int(sb[startPos+3])
	}

	v := (c1 << 18) + (c2 << 12) + (c3 << 6) + c4
	cw1 := byte((v >> 16) & 255)
	cw2 := byte((v >> 8) & 255)
	cw3 := byte(v & 255)
	res := make([]byte, 0, 3)
	res = append(res, cw1)
	if len >= 2 {
		res = append(res, cw2)
	}
	if len >= 3 {
		res = append(res, cw3)
	}
	return res, nil
}
