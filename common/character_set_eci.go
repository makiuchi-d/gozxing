package common

import (
	"github.com/makiuchi-d/gozxing"
)

type CharacterSetECI struct {
	values             []int
	otherEncodingNames []string
}

var (
	valueToECI = map[int]*CharacterSetECI{}
	nameToECI  = map[string]*CharacterSetECI{}
)

func init() {
	mapping(&CharacterSetECI{[]int{0, 2}, []string{}}) //Cp437
	mapping(&CharacterSetECI{[]int{1, 3}, []string{"ISO-8859-1"}})
	mapping(&CharacterSetECI{[]int{4}, []string{"ISO-8859-2"}})
	mapping(&CharacterSetECI{[]int{5}, []string{"ISO-8859-3"}})
	mapping(&CharacterSetECI{[]int{6}, []string{"ISO-8859-4"}})
	mapping(&CharacterSetECI{[]int{7}, []string{"ISO-8859-5"}})
	mapping(&CharacterSetECI{[]int{8}, []string{"ISO-8859-6"}})
	mapping(&CharacterSetECI{[]int{9}, []string{"ISO-8859-7"}})
	mapping(&CharacterSetECI{[]int{10}, []string{"ISO-8859-8"}})
	mapping(&CharacterSetECI{[]int{11}, []string{"ISO-8859-9"}})
	mapping(&CharacterSetECI{[]int{12}, []string{"ISO-8859-10"}})
	mapping(&CharacterSetECI{[]int{13}, []string{"ISO-8859-11"}})
	mapping(&CharacterSetECI{[]int{15}, []string{"ISO-8859-13"}})
	mapping(&CharacterSetECI{[]int{16}, []string{"ISO-8859-14"}})
	mapping(&CharacterSetECI{[]int{17}, []string{"ISO-8859-15"}})
	mapping(&CharacterSetECI{[]int{18}, []string{"ISO-8859-16"}})
	mapping(&CharacterSetECI{[]int{20}, []string{"Shift_JIS"}})
	mapping(&CharacterSetECI{[]int{21}, []string{"windows-1250"}})
	mapping(&CharacterSetECI{[]int{22}, []string{"windows-1251"}})
	mapping(&CharacterSetECI{[]int{23}, []string{"windows-1252"}})
	mapping(&CharacterSetECI{[]int{24}, []string{"windows-1256"}})
	mapping(&CharacterSetECI{[]int{25}, []string{"UTF-16BE", "UnicodeBig"}})
	mapping(&CharacterSetECI{[]int{26}, []string{"UTF-8"}})
	mapping(&CharacterSetECI{[]int{27, 170}, []string{"US-ASCII"}})
	mapping(&CharacterSetECI{[]int{28}, []string{}})
	mapping(&CharacterSetECI{[]int{29}, []string{"GB2312", "EUC_CN", "GBK"}})
	mapping(&CharacterSetECI{[]int{30}, []string{"EUC-KR"}})
}

func mapping(c *CharacterSetECI) {
	for _, val := range c.values {
		valueToECI[val] = c
	}
	for _, name := range c.otherEncodingNames {
		nameToECI[name] = c
	}
}

func (this *CharacterSetECI) GetValue() int {
	return this.values[0]
}

func GetCharacterSetECIByValue(value int) (*CharacterSetECI, error) {
	if value < 0 || value >= 900 {
		return nil, gozxing.GetFormatExceptionInstance()
	}
	return valueToECI[value], nil
}

func GetCharacterSetECIByName(name string) *CharacterSetECI {
	return nameToECI[name]
}
