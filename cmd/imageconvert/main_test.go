package main

import (
	"testing"

	"github.com/kmulvey/humantime"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name            string
	compress        bool
	force           bool
	watch           bool
	threads         int
	timerange       humantime.TimeRange
	resizeThreshold string
	resizeSize      string
	expectedError   bool
}

// nolint: gochecknoglobals
var testCases = []testCase{
	{
		name:            "valid params",
		compress:        true,
		force:           false,
		watch:           false,
		threads:         2,
		timerange:       humantime.TimeRange{},
		resizeThreshold: "2560x1440",
		resizeSize:      "5120x2880",
		expectedError:   false,
	},
	{
		name:            "invalid threads",
		compress:        true,
		force:           false,
		watch:           false,
		threads:         -1,
		timerange:       humantime.TimeRange{},
		resizeThreshold: "2560x1440",
		resizeSize:      "5120x2880",
		expectedError:   true,
	},
	{
		name:            "invalid resize threshold format",
		compress:        true,
		force:           false,
		watch:           false,
		threads:         1,
		timerange:       humantime.TimeRange{},
		resizeThreshold: "2560-1440",
		resizeSize:      "5120x2880",
		expectedError:   true,
	},
	{
		name:            "invalid resize size format",
		compress:        true,
		force:           false,
		watch:           false,
		threads:         1,
		timerange:       humantime.TimeRange{},
		resizeThreshold: "2560x1440",
		resizeSize:      "5120-2880",
		expectedError:   true,
	},
}

func TestGetResizeValue(t *testing.T) {
	t.Parallel()

	assert.Equal(t, uint16(5120), getResizeValue("5120"))
}

func TestParseParams(t *testing.T) {
	t.Parallel()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := parseParams(tt.compress, tt.force, tt.watch, tt.threads, tt.timerange, tt.resizeThreshold, tt.resizeSize)
			if tt.expectedError {
				assert.Error(t, err, tt.name)
			} else {
				assert.NoError(t, err, tt.name)
			}
		})
	}
}
