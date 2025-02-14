package imageconvert

import (
	"testing"

	"github.com/kmulvey/imageconvert/v2/testimages"
)

func TestMain(_ *testing.M) {
	compressTestCases = make([]compressTestCase, len(testimages.TestCases))
	for i, tc := range testimages.TestCases {
		compressTestCases[i] = compressTestCase{
			TestCase: tc,
			ShouldCompress: func() bool {
				return tc.ImageType == imageExt
			}(),
			PartialErrString: func() string {
				if tc.ImageType != "jpeg" {
					return "Not a JPEG file:"
				}
				return ""
			}(),
		}
	}

	convertTestCases = make([]convertTestCase, len(testimages.TestCases))
	for i, tc := range testimages.TestCases {
		convertTestCases[i] = convertTestCase{
			TestCase: tc,
			ShouldConvert: func() bool {
				return tc.ImageType != imageExt
			}(),
			PartialErrString: func() string {
				if tc.ImageType != "jpeg" {
					return "Not a JPEG file:"
				}
				return ""
			}(),
		}
	}
}
