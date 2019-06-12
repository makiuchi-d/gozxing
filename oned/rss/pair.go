package rss

type Pair struct {
	*DataCharacter
	finderPattern *FinderPattern
	count         int
}

func NewPair(value, checksumPortion int, finderPattern *FinderPattern) *Pair {
	return &Pair{
		DataCharacter: NewDataCharacter(value, checksumPortion),
		finderPattern: finderPattern,
	}
}

func (this *Pair) GetFinderPattern() *FinderPattern {
	return this.finderPattern
}

func (this *Pair) GetCount() int {
	return this.count
}

func (this *Pair) IncrementCount() {
	this.count++
}
