package detector

import (
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
	"github.com/makiuchi-d/gozxing/qrcode/detector"
)

// MultiDetector Encapsulates logic that can detect one or more QR Codes in an image,
// even if the QR Code is rotated or skewed, or partially obscured.
type MultiDetector struct {
	*detector.Detector
}

func NewMultiDetector(image *gozxing.BitMatrix) *MultiDetector {
	return &MultiDetector{
		detector.NewDetector(image),
	}
}

func (this *MultiDetector) DetectMulti(hints map[gozxing.DecodeHintType]interface{}) ([]*common.DetectorResult, error) {
	image := this.GetImage()
	resultPointCallback, _ := hints[gozxing.DecodeHintType_NEED_RESULT_POINT_CALLBACK].(gozxing.ResultPointCallback)

	finder := NewMultiFinderPatternFinder(image, resultPointCallback)
	infos, e := finder.FindMulti(hints)
	if e != nil || len(infos) == 0 {
		return nil, gozxing.WrapNotFoundException(e)
	}

	result := make([]*common.DetectorResult, 0)
	for _, info := range infos {
		r, e := this.ProcessFinderPatternInfo(info)
		if e != nil {
			// ignore
			continue
		}
		result = append(result, r)
	}
	return result, nil
}
