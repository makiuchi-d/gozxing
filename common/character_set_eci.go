package common

import (
	"github.com/makiuchi-d/gozxing"
)

type CharacterSetECI struct {
	values             []int
	name               string
	otherEncodingNames []string
}

var (
	valueToECI = map[int]*CharacterSetECI{}
	nameToECI  = map[string]*CharacterSetECI{}

	CharacterSetECI_Cp437      = newCharsetECI([]int{0, 2}, "Cp437")
	CharacterSetECI_ISO8859_1  = newCharsetECI([]int{1, 3}, "ISO-8859-1", "ISO8859_1")
	CharacterSetECI_ISO8859_2  = newCharsetECI([]int{4}, "ISO-8859-2", "ISO8859_2")
	CharacterSetECI_ISO8859_3  = newCharsetECI([]int{5}, "ISO-8859-3", "ISO8859_3")
	CharacterSetECI_ISO8859_4  = newCharsetECI([]int{6}, "ISO-8859-4", "ISO8859_4")
	CharacterSetECI_ISO8859_5  = newCharsetECI([]int{7}, "ISO-8859-5", "ISO8859_5")
	CharacterSetECI_ISO8859_6  = newCharsetECI([]int{8}, "ISO-8859-6", "ISO8859_6")
	CharacterSetECI_ISO8859_7  = newCharsetECI([]int{9}, "ISO-8859-7", "ISO8859_7")
	CharacterSetECI_ISO8859_8  = newCharsetECI([]int{10}, "ISO-8859-8", "ISO8859_8")
	CharacterSetECI_ISO8859_9  = newCharsetECI([]int{11}, "ISO-8859-9", "ISO8859_9")
	CharacterSetECI_ISO8859_10 = newCharsetECI([]int{12}, "ISO-8859-10", "ISO8859_10")
	//CharacterSetECI_ISO8859_11 = newCharsetECI([]int{13}, "TIS-620", "ISO-8859-11", "ISO8859_11") // golang does not support
	CharacterSetECI_ISO8859_13         = newCharsetECI([]int{15}, "ISO-8859-13", "ISO8859_13")
	CharacterSetECI_ISO8859_14         = newCharsetECI([]int{16}, "ISO-8859-14", "ISO8859_14")
	CharacterSetECI_ISO8859_15         = newCharsetECI([]int{17}, "ISO-8859-15", "ISO8859_15")
	CharacterSetECI_ISO8859_16         = newCharsetECI([]int{18}, "ISO-8859-16", "ISO8859_16")
	CharacterSetECI_SJIS               = newCharsetECI([]int{20}, "Shift_JIS", "SJIS")
	CharacterSetECI_Cp1250             = newCharsetECI([]int{21}, "windows-1250", "Cp1250")
	CharacterSetECI_Cp1251             = newCharsetECI([]int{22}, "windows-1251", "Cp1251")
	CharacterSetECI_Cp1252             = newCharsetECI([]int{23}, "windows-1252", "Cp1252")
	CharacterSetECI_Cp1256             = newCharsetECI([]int{24}, "windows-1256", "Cp1256")
	CharacterSetECI_UnicodeBigUnmarked = newCharsetECI([]int{25}, "UTF-16BE", "UnicodeBig", "UnicodeBigUnmarked")
	CharacterSetECI_UTF8               = newCharsetECI([]int{26}, "UTF-8", "UTF8")
	CharacterSetECI_ASCII              = newCharsetECI([]int{27, 170}, "ASCII", "US-ASCII")
	CharacterSetECI_Big5               = newCharsetECI([]int{28}, "Big5")
	CharacterSetECI_GB18030            = newCharsetECI([]int{29}, "GB18030", "GB2312", "EUC_CN", "GBK") // BG18030 is upward compatible with others
	CharacterSetECI_EUC_KR             = newCharsetECI([]int{30}, "EUC-KR", "EUC_KR")
)

func newCharsetECI(values []int, encodingNames ...string) *CharacterSetECI {
	c := &CharacterSetECI{
		values:             values,
		name:               encodingNames[0],
		otherEncodingNames: encodingNames[1:],
	}
	for _, val := range values {
		valueToECI[val] = c
	}
	for _, name := range encodingNames {
		nameToECI[name] = c
	}
	return c
}

func (this *CharacterSetECI) GetValue() int {
	return this.values[0]
}

func (this *CharacterSetECI) Name() string {
	return this.name
}

func GetCharacterSetECIByValue(value int) (*CharacterSetECI, error) {
	if value < 0 || value >= 900 {
		return nil, gozxing.NewFormatException()
	}
	return valueToECI[value], nil
}

func GetCharacterSetECIByName(name string) *CharacterSetECI {
	return nameToECI[name]
}
