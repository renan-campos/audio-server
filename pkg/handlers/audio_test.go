package handlers

import (
	"fmt"
	"regexp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Audio handler endpoint behavior", func() {
	type rangeTestParams struct {
		RangeStr      string
		FileSize      int
		ExpectedStart int64
		ExpectedEnd   int64
		ExpectedErr   error
	}
	DescribeTable("parseRange helper", func(params rangeTestParams) {
		subject := &audioHandler{
			rangeRegex: regexp.MustCompile(rangeRegex),
		}
		start, end, err := subject.parseRange(params.RangeStr, params.FileSize)
		if params.ExpectedErr != nil {
			Expect(err).ToNot(BeNil())
		} else {
			Expect(err).To(BeNil())
			Expect(start).To(Equal(params.ExpectedStart))
			Expect(end).To(Equal(params.ExpectedEnd))
		}
	},
		Entry("no range specified results in full range", rangeTestParams{
			RangeStr:      "",
			FileSize:      100,
			ExpectedStart: 0,
			ExpectedEnd:   100,
			ExpectedErr:   nil,
		}),
		Entry("start and end range specified returns expected values", rangeTestParams{
			RangeStr:      "bytes=25-75",
			FileSize:      100,
			ExpectedStart: 25,
			ExpectedEnd:   75,
			ExpectedErr:   nil,
		}),
		Entry("no end range specified results in filesize as end range", rangeTestParams{
			RangeStr:      "bytes=50-",
			FileSize:      100,
			ExpectedStart: 50,
			ExpectedEnd:   100,
			ExpectedErr:   nil,
		}),
		// This one needs fixing
		Entry("only end range specified results in last n bytes of file", rangeTestParams{
			RangeStr:      "bytes=-20",
			FileSize:      100,
			ExpectedStart: 80,
			ExpectedEnd:   100,
			ExpectedErr:   nil,
		}),
		Entry("rangeStart larger than file size results in error", rangeTestParams{
			RangeStr:    "bytes=101-",
			FileSize:    100,
			ExpectedErr: fmt.Errorf("error occurred"),
		}),
		Entry("rangeEnd larger than file size results in error", rangeTestParams{
			RangeStr:    "bytes=0-200",
			FileSize:    100,
			ExpectedErr: fmt.Errorf("error occurred"),
		}),
		Entry("rangeStart larger than rangeEnd results in error", rangeTestParams{
			RangeStr:    "bytes=75-50",
			FileSize:    100,
			ExpectedErr: fmt.Errorf("error occurred"),
		}),
	)
})
