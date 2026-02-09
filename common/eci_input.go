package common

// ECIInput Interface to navigate a sequence of ECIs and bytes.
type ECIInput interface {
	// Length Returns the length of this input. The length is the number
	// of bytes in or ECIs in the sequence.
	Length() int

	// CharAt Returns the byte value at the specified index. An index ranges from zero
	// to length() - 1. The first byte value of the sequence is at
	// index zero, the next at index one, and so on, as for array
	// indexing.
	//
	// @param   index the index of the byte value to be returned
	//
	// @return  the specified byte value as character or the FNC1 character
	CharAt(index int) (int, error)

	// SubSequence Returns a subsequence of this sequence.
	// The subsequence starts with the char value at the specified index and
	// ends with the char value at index end - 1. The length
	// (in chars) of the
	// returned sequence is end - start, so if start == end
	// then an empty sequence is returned.
	SubSequence(start, end int) ([]int, error)

	// IsECI Determines if a value is an ECI
	//
	// @param   index the index of the value
	//
	// @return  true if the value at position index is an ECI
	IsECI(index int) (bool, error)

	// GetECIValue Returns the int ECI value at the specified index. An index ranges from zero
	// to length() - 1. The first byte value of the sequence is at
	// index zero, the next at index one, and so on, as for array
	// indexing.
	//
	// @param   index the index of the int value to be returned
	//
	// @return  the specified int ECI value.
	//          The ECI specified the encoding of all bytes with a higher index until the
	//          next ECI or until the end of the input if no other ECI follows.
	GetECIValue(index int) (int, error)

	// HaveNCharacters returns true if there are n characters starting at index
	HaveNCharacters(index, n int) bool
}
