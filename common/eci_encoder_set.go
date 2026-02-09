package common

import (
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"

	"github.com/makiuchi-d/gozxing"
)

// ECIEncoderSet Set of CharsetEncoders for a given input string
//
// Invariants:
//   - The list contains only encoders from CharacterSetECI (list is shorter then the list of encoders available on
//     the platform for which ECI values are defined).
//   - The list contains encoders at least one encoder for every character in the input.
//   - The first encoder in the list is always the ISO-8859-1 encoder even of no character in the input can be encoded
//     by it.
//   - If the input contains a character that is not in ISO-8859-1 then the last two entries in the list will be the
//     UTF-8 encoder and the UTF-16BE encoder.
type ECIEncoderSet struct {
	encoders             []encoding.Encoding
	priorityEncoderIndex int
}

var encodersList = []encoding.Encoding{
	charmap.CodePage437, // "IBM437"
	charmap.ISO8859_2,
	charmap.ISO8859_3,
	charmap.ISO8859_4,
	charmap.ISO8859_5,
	charmap.ISO8859_6,
	charmap.ISO8859_7,
	charmap.ISO8859_8,
	charmap.ISO8859_9,
	charmap.ISO8859_10,
	//charmap.ISO8859_11, // golang does not support
	charmap.ISO8859_13,
	charmap.ISO8859_14,
	charmap.ISO8859_15,
	charmap.ISO8859_16,
	charmap.Windows1250,
	charmap.Windows1251,
	charmap.Windows1252,
	charmap.Windows1256,
	japanese.ShiftJIS,
}

// NewECIEncoderSet Constructs an encoder set
//
// @param stringToEncode the string that needs to be encoded
// @param priorityCharset The preferred encoding.Encoding or nil.
// @param fnc1 fnc1 denotes the character in the input that represents the FNC1 character or -1 for a non-GS1 bar
// code. When specified, it is considered an error to pass it as argument to the methods canEncode() or encode().
func NewECIEncoderSet(stringToEncode []rune, priorityCharset encoding.Encoding, fnc1 int) *ECIEncoderSet {
	neededEncoders := make([]encoding.Encoding, 0)

	// we always need the ISO-8859-1 encoder. It is the default encoding
	neededEncoders = append(neededEncoders, charmap.ISO8859_1)

	needUnicodeEncoder := false
	if priorityCharset != nil {
		if name, err := ianaindex.IANA.Name(priorityCharset); err == nil {
			needUnicodeEncoder = strings.HasPrefix(name, "UTF")
		}
	}

	// Walk over the input string and see if all characters can be encoded with the list of encoders
	for i := 0; i < len(stringToEncode); i++ {
		canEncode := false
		c := stringToEncode[i]
		if c == rune(fnc1) {
			canEncode = true
		} else {
			for _, encoder := range neededEncoders {
				if canEncodeChar(encoder, c) {
					canEncode = true
					break
				}
			}
		}

		if !canEncode {
			// for the character at position i we don't yet have an encoder in the list
			for _, encoder := range encodersList {
				if canEncodeChar(encoder, c) {
					// Good, we found an encoder that can encode the character. We add him to the list and continue scanning
					// the input
					neededEncoders = append(neededEncoders, encoder)
					canEncode = true
					break
				}
			}
		}

		if !canEncode {
			// The character is not encodeable by any of the single byte encoders so we remember that we will need a
			// Unicode encoder.
			needUnicodeEncoder = true
		}
	}

	var encoders []encoding.Encoding
	if len(neededEncoders) == 1 && !needUnicodeEncoder {
		// the entire input can be encoded by the ISO-8859-1 encoder
		encoders = []encoding.Encoding{neededEncoders[0]}
	} else {
		// we need more than one single byte encoder or we need a Unicode encoder.
		// In this case we append a UTF-8 and UTF-16 encoder to the list
		encoders = make([]encoding.Encoding, len(neededEncoders)+2)
		copy(encoders, neededEncoders)
		encoders[len(neededEncoders)] = unicode.UTF8
		utf16be, _ := ianaindex.IANA.Encoding("UTF-16BE")
		encoders[len(neededEncoders)+1] = utf16be
	}

	// Compute priorityEncoderIndex by looking up priorityCharset in encoders
	priorityEncoderIndexValue := -1
	if priorityCharset != nil {
		priorityName, _ := ianaindex.IANA.Name(priorityCharset)
		for i := 0; i < len(encoders); i++ {
			if encoders[i] != nil {
				if name, err := ianaindex.IANA.Name(encoders[i]); err == nil {
					if name == priorityName {
						priorityEncoderIndexValue = i
						break
					}
				}
			}
		}
	}

	return &ECIEncoderSet{
		encoders:             encoders,
		priorityEncoderIndex: priorityEncoderIndexValue,
	}
}

func canEncodeChar(enc encoding.Encoding, c rune) bool {
	encoder := enc.NewEncoder()
	_, err := encoder.Bytes([]byte(string(c)))
	return err == nil
}

func (this *ECIEncoderSet) Length() int {
	return len(this.encoders)
}

func (this *ECIEncoderSet) GetCharsetName(index int) string {
	if index >= len(this.encoders) {
		return ""
	}
	name, _ := ianaindex.IANA.Name(this.encoders[index])
	return name
}

func (this *ECIEncoderSet) GetCharset(index int) encoding.Encoding {
	if index >= len(this.encoders) {
		return nil
	}
	return this.encoders[index]
}

func (this *ECIEncoderSet) GetECIValue(encoderIndex int) int {
	if encoderIndex >= len(this.encoders) {
		return -1
	}
	if eci, ok := GetCharacterSetECI(this.encoders[encoderIndex]); ok {
		return eci.GetValue()
	}
	return -1
}

// GetPriorityEncoderIndex returns -1 if no priority charset was defined
func (this *ECIEncoderSet) GetPriorityEncoderIndex() int {
	return this.priorityEncoderIndex
}

func (this *ECIEncoderSet) CanEncode(c rune, encoderIndex int) bool {
	if encoderIndex >= len(this.encoders) {
		return false
	}
	return canEncodeChar(this.encoders[encoderIndex], c)
}

func (this *ECIEncoderSet) EncodeChar(c rune, encoderIndex int) ([]byte, error) {
	if encoderIndex >= len(this.encoders) {
		return nil, gozxing.NewWriterException("encoderIndex out of bounds")
	}
	encoder := this.encoders[encoderIndex].NewEncoder()
	return encoder.Bytes([]byte(string(c)))
}

func (this *ECIEncoderSet) EncodeString(s string, encoderIndex int) ([]byte, error) {
	if encoderIndex >= len(this.encoders) {
		return nil, gozxing.NewWriterException("encoderIndex out of bounds")
	}
	encoder := this.encoders[encoderIndex].NewEncoder()
	return encoder.Bytes([]byte(s))
}
