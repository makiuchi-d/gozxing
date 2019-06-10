package rss

import (
	"fmt"
)

// Encapsulates a since character value in an RSS barcode, including its checksum information.

type DataCharacter struct {
	value           int
	checksumPortion int
}

func NewDataCharacter(value, checksumPortion int) *DataCharacter {
	return &DataCharacter{value, checksumPortion}
}

func (this *DataCharacter) GetValue() int {
	return this.value
}

func (this *DataCharacter) GetChecksumPortion() int {
	return this.checksumPortion
}

func (this *DataCharacter) String() string {
	return fmt.Sprintf("%v(%v)", this.value, this.checksumPortion)
}

func (this *DataCharacter) Equals(o interface{}) bool {
	that, ok := o.(*DataCharacter)
	if !ok {
		return false
	}
	return this.value == that.value &&
		this.checksumPortion == that.checksumPortion
}

func (this *DataCharacter) HashCode() int {
	return this.value ^ this.checksumPortion
}
