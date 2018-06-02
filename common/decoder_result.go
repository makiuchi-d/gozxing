package common

type DecoderResult struct {
	rawBytes                       []byte
	numBits                        int
	text                           string
	byteSegments                   [][]byte
	ecLevel                        string
	errorsCorrected                int
	erasures                       int
	other                          interface{}
	structuredAppendParity         int
	structuredAppendSequenceNumber int
}

func NewDecoderResult(rawBytes []byte, text string, byteSegments [][]byte, ecLevel string) *DecoderResult {
	return NewDecoderResultWithSA(rawBytes, text, byteSegments, ecLevel, -1, -1)
}

func NewDecoderResultWithSA(rawBytes []byte, text string, byteSegments [][]byte, ecLevel string, saSequence, saParity int) *DecoderResult {
	return &DecoderResult{
		rawBytes:                       rawBytes,
		numBits:                        8 * len(rawBytes),
		text:                           text,
		byteSegments:                   byteSegments,
		ecLevel:                        ecLevel,
		structuredAppendParity:         saParity,
		structuredAppendSequenceNumber: saSequence,
	}
}

func (this *DecoderResult) GetRawBytes() []byte {
	return this.rawBytes
}

func (this *DecoderResult) GetNumBits() int {
	return this.numBits
}

func (this *DecoderResult) SetNumBits(numBits int) {
	this.numBits = numBits
}

func (this *DecoderResult) GetText() string {
	return this.text
}

func (this *DecoderResult) GetByteSegments() [][]byte {
	return this.byteSegments
}

func (this *DecoderResult) GetECLevel() string {
	return this.ecLevel
}

func (this *DecoderResult) GetErrorsCorrected() int {
	return this.errorsCorrected
}

func (this *DecoderResult) SetErrorsCorrected(errorsCorrected int) {
	this.errorsCorrected = errorsCorrected
}

func (this *DecoderResult) GetErasures() int {
	return this.erasures
}

func (this *DecoderResult) SetErasures(erasures int) {
	this.erasures = erasures
}

func (this *DecoderResult) GetOther() interface{} {
	return this.other
}

func (this *DecoderResult) SetOther(other interface{}) {
	this.other = other
}

func (this *DecoderResult) HasStructuredAppend() bool {
	return this.structuredAppendParity >= 0 && this.structuredAppendSequenceNumber >= 0
}

func (this *DecoderResult) GetStructuredAppendParity() int {
	return this.structuredAppendParity
}

func (this *DecoderResult) GetStructuredAppendSequenceNumber() int {
	return this.structuredAppendSequenceNumber
}
