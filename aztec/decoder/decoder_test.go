package decoder

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/aztec/detector"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestDecoder_Decode(t *testing.T) {
	dec := NewDecoder()

	bmp, _ := gozxing.ParseStringToBitMatrix(""+
		"    ##    ##  ####        ##  \n"+
		"  ######    ##  ######      ##\n"+
		"    ####        ##  ##  ##    \n"+
		"##########################    \n"+
		"####  ##              ##      \n"+
		"    ####  ##########  ##  ##  \n"+
		"  ##  ##  ##      ##  ##      \n"+
		"  ######  ##  ##  ##  ########\n"+
		"  ######  ##      ##  ##      \n"+
		"  ######  ##########  ####    \n"+
		"    ####              ######  \n"+
		"##    ####################  ##\n"+
		"##        ##    ##  ##        \n"+
		"####      ######  ##  ##    ##\n"+
		"########    ####  ####  ##  ##\n",
		"##", "  ")
	bmp = testutil.ExpandBitMatrix(bmp, 3)
	ddata, _ := detector.NewDetector(bmp).Detect(false)

	r, e := dec.Decode(ddata)
	if e != nil {
		t.Fatalf("Decode error: %v", e)
	}
	wants := "Histórico"
	if r.GetText() != wants {
		t.Fatalf("Decode = %q, wants %q", r.GetText(), wants)
	}

	// reedsolomon error
	bmp = ddata.GetBits()
	bmp.SetRegion(0, 0, 15, 3)
	_, e = dec.Decode(ddata)
	if e == nil {
		t.Fatalf("Decode must be error")
	}

	// invalid bits with correct ecwords
	bmp, _ = gozxing.ParseStringToBitMatrix(""+
		"    ####      ##  ##  ##  ####\n"+
		"            ##            ##  \n"+
		"    ####            ##  ######\n"+
		"  ########################    \n"+
		"      ##              ##    ##\n"+
		"    ####  ##########  ##  ####\n"+
		"  ##  ##  ##      ##  ####    \n"+
		"####  ##  ##  ##  ##  ##      \n"+
		"####  ##  ##      ##  ########\n"+
		"##    ##  ##########  ######  \n"+
		"  ######              ##    ##\n"+
		"      ######################  \n"+
		"  ##          ####        ##  \n"+
		"      ####        ##  ##  ####\n"+
		"##  ##      ######  ####    ##\n",
		"##", "  ")
	bmp = testutil.ExpandBitMatrix(bmp, 3)
	ddata, _ = detector.NewDetector(bmp).Detect(false)
	_, e = dec.Decode(ddata)
	if e == nil {
		t.Fatalf("Decode must be error")
	}

	// full, layer=5
	bmp, _ = gozxing.ParseStringToBitMatrix(""+
		"                                                                              \n"+
		"  ##############  ######    ##  ####  ####  ##    ####  ########    ######    \n"+
		"  ##    ####    ####      ######        ##        ######  ##    ##  ##    ##  \n"+
		"  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  \n"+
		"        ##    ##  ##  ####        ##        ##    ##        ##  ##      ##    \n"+
		"      ##    ####    ####    ######    ##  ######  ##  ####      ##    ##  ##  \n"+
		"  ##    ####  ##          ##  ####  ##  ########    ##        ######      ##  \n"+
		"      ##    ##########  ##  ##  ##    ##  ##        ##    ####  ##    ##  ##  \n"+
		"    ##  ####      ##  ######    ####    ######      ######  ########          \n"+
		"  ##  ####  ########    ######      ######  ##  ######  ########  ########    \n"+
		"    ##  ##  ########  ##    ##          ######  ####  ##    ####  ####        \n"+
		"      ########  ####    ##########    ##########  ##  ####  ##      ####  ##  \n"+
		"          ##          ######    ##              ##  ######        ####    ##  \n"+
		"  ##  ######  ##  ####################################    ##  ##  ##########  \n"+
		"  ##    ####      ##      ##                      ##  ####  ####    ##        \n"+
		"    ####      ####  ##    ##  ##################  ##  ########  ####  ##  ##  \n"+
		"  ##        ##      ##    ##  ##              ##  ##  ##        ####    ##    \n"+
		"      ##  ##  ########  ####  ##  ##########  ##  ####  ####  ##  ######  ##  \n"+
		"  ####        ######  ##  ##  ##  ##      ##  ##  ##      ##      ##      ##  \n"+
		"  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  \n"+
		"  ##                  ######  ##  ##      ##  ##  ##  ##  ##        ##        \n"+
		"  ##  ##########  ##      ##  ##  ##########  ##  ####  ####  ##  ##  ##      \n"+
		"            ##  ##      ####  ##              ##  ######                ##    \n"+
		"    ######  ##    ####  ####  ##################  ##      ####    ######  ##  \n"+
		"              ####  ##    ##                      ##        ##########  ##    \n"+
		"  ##  ##  ##  ####    ##  ############################  ##      ####  ##      \n"+
		"  ##      ##  ########        ####      ####  ##      ######      ##      ##  \n"+
		"  ##  ##    ##  ##    ####  ####  ############      ####################      \n"+
		"    ##  ##########  ######        ##    ####    ####  ####  ##  ##  ##    ##  \n"+
		"    ####      ##  ##  ######  ####    ##      ####    ##  ####  ####  ##  ##  \n"+
		"  ##      ##  ##  ##  ####  ##  ######    ##    ####    ##  ####  ##    ####  \n"+
		"  ##############  ##  ########  ####  ####  ##########    ##      ##  ####    \n"+
		"          ######    ##  ##    ######    ##    ##    ##        ##    ##  ####  \n"+
		"    ####      ######    ####    ##########  ####    ##  ##  ####    ########  \n"+
		"  ##        ##        ####  ##  ##      ##  ######        ##      ####  ##    \n"+
		"  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  \n"+
		"    ##    ######        ##      ######        ####    ####    ##  ##    ####  \n"+
		"  ########  ####    ##  ##  ####  ##  ##  ######  ##########  ##  ######  ##  \n"+
		"                                                                              \n",
		"##", "  ")
	bmp = testutil.ExpandBitMatrix(bmp, 3)
	ddata, e = detector.NewDetector(bmp).Detect(false)
	r, e = dec.Decode(ddata)
	if e != nil {
		t.Fatalf("Decode error: %v", e)
	}
	wants = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if txt := r.GetText(); txt != wants {
		t.Fatalf("Decode = %q, wants %q", txt, wants)
	}
}

func strToBools(s string) []bool {
	r := make([]bool, len(s))
	for i, c := range s {
		r[i] = c == '1'
	}
	return r
}

func TestDecoder_HighLevelDecode(t *testing.T) {
	dec := NewDecoder()
	bits := strToBools("00010" + "11100" + "11011" + "11110" + "1101") // 'A', CTRL_LL, 'z', CTLR_DL, '.'
	r, e := dec.HighLevelDecode(bits)
	if e != nil {
		t.Fatalf("HighLevelDecode error: %v", e)
	}
	wants := "Az."
	if r != wants {
		t.Fatalf("HighLevelDecode = %q, wants %q", r, wants)
	}

	// GS1 data
	// P/S FLG(n) 0  D/L 1 0 1 2 3 P/S FLG(n) 0 3 7 4 2
	bits = strToBools("" +
		"00000" + "00000" + "000" + "11110" + "0011" + "0010" + "0011" + "0100" + "0101" +
		"0000" + "00000" + "000" + "0101" + "1001" + "0110" + "0100")
	r, e = dec.HighLevelDecode(bits)
	if e != nil {
		t.Fatalf("HighLevelDecode error: %v", e)
	}
	wants = "\x1d10123\x1d3742"
	if r != wants {
		t.Fatalf("HighLevelDecode = %q, wants %q", r, wants)
	}
}

func TestDecoder_getEncodedData(t *testing.T) {
	dec := NewDecoder()

	bits := strToBools("" +
		"00000" + "00000" + "000" + // CTRL_PS, FLG(n), 0
		"00010" + "11100" + "11011" + "11110" + "1101" + "1110" + // 'A', CTRL_LL, 'z', CTLR_DL, '.' CTRL_UL
		"00000" + "00000" + "010" + "0100" + "1000" + // CTRL_PS, FLG(n), 2, '2', '6' (charset = UTF-8)
		"11111" + "00000" + "00000001001" + // CTRL_BS, length=40 (0, 9)
		"11100101" + "10101111" + "10111111" + // "寿"
		"11100101" + "10001111" + "10111000" + // "司"
		"11110000" + "10011111" + "10001101" + "10100011" + // "🍣"
		"11101001" + "10000101" + "10010010" + // "酒"
		"11110000" + "10011111" + "10001101" + "10110110" + // "🍶"
		"11100011" + "10000011" + "10010100" + // "ピ"
		"11100011" + "10000010" + "10110110" + // "ザ"
		"11110000" + "10011111" + "10001101" + "10010101" + // "🍕"
		"11100011" + "10000011" + "10010011" + // "ビ"
		"11100011" + "10000011" + "10111100" + // "ー"
		"11100011" + "10000011" + "10101011" + // "ル"
		"11110000" + "10011111" + "10001101" + "10111010" + // "🍺"
		"")
	r, e := dec.getEncodedData(bits)
	wants := "\035Az.寿司🍣酒🍶ピザ🍕ビール🍺"
	if e != nil {
		t.Fatalf("getEncodedData error: %v", e)
	}
	if r != wants {
		t.Fatalf("result = %q, wants %q", r, wants)
	}

	// FLG(7)
	bits = strToBools("00000" + "00000" + "111") // CTRL_PS, FLG(n), 7
	_, e = dec.getEncodedData(bits)
	if e == nil {
		t.Fatalf("getEncodedData(FLG(7)) must be error")
	}

	// encoding digit error
	bits = strToBools("00000" + "00000" + "001" + "0001") // CTRL_PS, FLG(n), 1, ' '
	_, e = dec.getEncodedData(bits)
	if e == nil {
		t.Fatalf("getEncodeData(FLG(1),1) must be error")
	}

	// eci value error
	bits = strToBools("00000" + "00000" + "011" + "1011" + "0010" + "0010") // CTRL_PS, FLG(n), 3, "900"
	_, e = dec.getEncodedData(bits)
	if e == nil {
		t.Fatalf("getEncodeData(FLG(3),100) must be error")
	}

	// break with incomplete bits

	bits = strToBools("00010" + "11111" + "00") // 'A', CTRL_BS, ...
	r, e = dec.getEncodedData(bits)
	if e != nil {
		t.Fatalf("getEncodedData error: %v", e)
	}
	if r != "A" {
		t.Fatalf("result = %q, wants %q", r, "A")
	}

	bits = strToBools("00010" + "11111" + "00000" + "00") // 'A', CTRL_BS, 0, ...
	r, e = dec.getEncodedData(bits)
	if e != nil {
		t.Fatalf("getEncodedData error: %v", e)
	}
	if r != "A" {
		t.Fatalf("result = %q, wants %q", r, "A")
	}

	bits = strToBools("00010" + "11111" + "00001" + "00") // 'A', CTRL_BS, length=1, ...
	r, e = dec.getEncodedData(bits)
	if e != nil {
		t.Fatalf("getEncodedData error: %v", e)
	}
	if r != "A" {
		t.Fatalf("result = %q, wants %q", r, "A")
	}

	bits = strToBools("00010" + "00") // 'A', ...
	r, e = dec.getEncodedData(bits)
	if e != nil {
		t.Fatalf("getEncodedData error: %v", e)
	}
	if r != "A" {
		t.Fatalf("result = %q, wants %q", r, "A")
	}

	bits = strToBools("00010" + "00000" + "00000" + "00") // 'A', CTRL_PS, FLG(n), ...
	r, e = dec.getEncodedData(bits)
	if e != nil {
		t.Fatalf("getEncodedData error: %v", e)
	}
	if r != "A" {
		t.Fatalf("result = %q, wants %q", r, "A")
	}

	bits = strToBools("00010" + "00000" + "00000" + "100" + "00") // 'A', CTRL_PS, FLG(n), 3, ...
	r, e = dec.getEncodedData(bits)
	if e != nil {
		t.Fatalf("getEncodedData error: %v", e)
	}
	if r != "A" {
		t.Fatalf("result = %q, wants %q", r, "A")
	}
}

func TestGetTable(t *testing.T) {
	tests := []struct {
		t     byte
		wants Table
	}{
		{'L', TableLOWER},
		{'P', TablePUNCT},
		{'M', TableMIXED},
		{'D', TableDIGIT},
		{'B', TableBINARY},
		{'U', TableUPPER},
		{'-', TableUPPER},
	}
	for _, test := range tests {
		r := getTable(test.t)
		if r != test.wants {
			t.Fatalf("getTable(%v) = %v, wants %v", test.t, r, test.wants)
		}
	}
}

func TestGetCharacter(t *testing.T) {
	tests := []struct {
		table Table
		code  int
		wants string
	}{
		{TableUPPER, 0, "CTRL_PS"},
		{TableUPPER, 5, "D"},
		{TableLOWER, 10, "i"},
		{TableMIXED, 3, "\002"},
		{TableMIXED, 20, "@"},
		{TablePUNCT, 0, "FLG(n)"},
		{TablePUNCT, 2, "\r\n"},
		{TableDIGIT, 11, "9"},
	}
	for _, test := range tests {
		r, e := getCharacter(test.table, test.code)
		if e != nil {
			t.Fatalf("getCharacter(%v, %v) error: %v", test.table, test.code, e)
		}
		if r != test.wants {
			t.Fatalf("getCharacter(%v, %v) = %v, wants %v", test.table, test.code, r, test.wants)
		}
	}

	// invalid table
	_, e := getCharacter(Table(100), 0)
	if e == nil {
		t.Fatalf("getCharacter(Table(100)) must be error")
	}

	// invalid code
	_, e = getCharacter(TableDIGIT, len(DIGIT_TABLE))
	if e == nil {
		t.Fatalf("getCharacter(TableDIGIT, %v) must be error", len(DIGIT_TABLE))
	}
}

func TestDecoder_correctBits(t *testing.T) {
	dec := NewDecoder()

	// not enough bits
	dec.ddata = detector.NewAztecDetectorResult(nil, nil, true, 100, 24)
	_, e := dec.correctBits([]bool{})
	if e == nil {
		t.Fatalf("correctBits({}) must be error")
	}

	// reedsolomon collect error
	dec.ddata = detector.NewAztecDetectorResult(nil, nil, true, 11, 1)
	bits := []bool{
		false, false, false, true, false, false, true, true, true, true, false, false, false, true, false, true,
		false, true, false, true, false, false, true, false, true, false, true, true, true, true, true, false,
		true, false, false, false, false, true, true, true, true, true, false, false, true, true, true, false,
		false, true, true, false, true, false, true, false, false, false, true, false, false, true, false, false,
		false, false, true, true, false, false, false, true, false, false, false, false, false, false, false, true,
		true, false, false, false, false, false, false, true, false, true, true, true, true, false, false, true,
		true, false, false, false, false, true, true, true,
	}
	bits[0] = !bits[0]
	bits[16] = !bits[16]
	bits[32] = !bits[32]
	bits[48] = !bits[48]
	bits[64] = !bits[64]
	bits[80] = !bits[80]
	bits[96] = !bits[96]
	_, e = dec.correctBits(bits)
	if e == nil {
		t.Fatalf("correctBits({...}) must be error")
	}

	// 0 filled
	dec.ddata = detector.NewAztecDetectorResult(nil, nil, false, 150, 9)
	_, e = dec.correctBits(make([]bool, 2304))
	if e == nil {
		t.Fatalf("correctBits({false...}) must be error")
	}

	// no error
	dec.ddata = detector.NewAztecDetectorResult(nil, nil, true, 11, 1)
	bits[0] = !bits[0]
	//bits[16] = !bits[16]
	//bits[32] = !bits[32]
	//bits[48] = !bits[48]
	bits[64] = !bits[64]
	bits[80] = !bits[80]
	bits[96] = !bits[96]
	r, e := dec.correctBits(bits)
	if e != nil {
		t.Fatalf("correctBits({...}) error: %v", e)
	}
	if l := len(r.correctBits); l != 65 {
		t.Fatalf("correctBits length = %v, wants 65", l)
	}
	if l := r.ecLevel; l != 35 {
		t.Fatalf("ecLevel = %v, wants 35", l)
	}
}

func TestConvertBoolArrayToByteArray(t *testing.T) {
	bools := []bool{
		true, false, false, true, true, true, false, true,
		false, true, false, false, true, false, false, false,
		true, true, false, true, false,
	}
	wants := []byte{0b10011101, 0b01001000, 0b11010000}
	bytes := convertBoolArrayToByteArray(bools)

	if !reflect.DeepEqual(bytes, wants) {
		t.Fatalf("convertBoolArrayToByteArray = %v, wants %v", bytes, wants)
	}
}

func TestTotalBitsInLayer(t *testing.T) {
	tests := []struct {
		layers  int
		compact bool
		wants   int
	}{
		{2, false, 288},
		{3, false, 480},
		{2, true, 240},
		{3, true, 408},
	}

	for _, test := range tests {
		r := totalBitsInLayer(test.layers, test.compact)
		if r != test.wants {
			t.Fatalf("totalBitsInLayer(%v,%v) = %v, wants %v", test.layers, test.compact, r, test.wants)
		}
	}
}
