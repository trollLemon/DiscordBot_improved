package util_test

import (
	"bot/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaturateQuery(t *testing.T) {

	tests := []struct {
		name     string
		param    float32
		expected string
	}{
		{
			name:     "Test single digit",
			param:    1.0,
			expected: "?saturation=1.00",
		},
		{
			name:     "Test several digits",
			param:    25.0,
			expected: "?saturation=25.00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			queryStr := util.SaturateQuery(tt.param)
			assert.Equal(t, tt.expected, queryStr)

		})
	}

}

func TestEdgeDetectQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   [2]int64
		expected string
	}{
		{
			name:     "Test positive numbers",
			params:   [2]int64{100, 200},
			expected: "?lower=100&higher=200",
		},
		{
			name:     "Test negative numbers",
			params:   [2]int64{-100, -200},
			expected: "?lower=-100&higher=-200",
		},
		{
			name:     "Test negative and positive",
			params:   [2]int64{-100, 200},
			expected: "?lower=-100&higher=200",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			queryStr := util.EdgeDetectQuery(tt.params[0], tt.params[1])
			assert.Equal(t, tt.expected, queryStr)

		})
	}
}

func TestDilateQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   [2]int64
		expected string
	}{
		{
			name:     "Test kernel size and iterations",
			params:   [2]int64{3, 4},
			expected: "?kernelSize=3&iterations=4&type=Dilate",
		},
		{
			name:     "Test negative values (invalid for manip server but ensure query is formatted correctly)",
			params:   [2]int64{-3, -4},
			expected: "?kernelSize=-3&iterations=-4&type=Dilate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			queryStr := util.DilateQuery(tt.params[0], tt.params[1])
			assert.Equal(t, tt.expected, queryStr)
		})
	}
}

func TestErodeQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   [2]int64
		expected string
	}{
		{
			name:     "Test kernel size and iterations",
			params:   [2]int64{3, 4},
			expected: "?kernelSize=3&iterations=4&type=Erode",
		},
		{
			name:     "Test negative values (invalid for manip server but ensure query is formatted correctly)",
			params:   [2]int64{-3, -4},
			expected: "?kernelSize=-3&iterations=-4&type=Erode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			queryStr := util.ErodeQuery(tt.params[0], tt.params[1])
			assert.Equal(t, tt.expected, queryStr)
		})
	}
}

func TestReduceQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   float32
		expected string
	}{
		{
			name:     "Test query",
			params:   0.2,
			expected: "?quality=0.20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryStr := util.ReduceQuery(tt.params)
			assert.Equal(t, tt.expected, queryStr)
		})
	}
}

func TestAddTextQuery(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		params   [3]float32
		expected string
	}{
		{
			name:     "Test adding text query",
			text:     "text",
			params:   [3]float32{1.0, 0.5, 0.5},
			expected: "?text=text&fontScale=1.00&xPerc=0.50&yPerc=0.50",
		},
		{
			name:     "Test adding text query with spaces",
			text:     "I love golang",
			params:   [3]float32{1.0, 0.5, 0.5},
			expected: "?text=I+love+golang&fontScale=1.00&xPerc=0.50&yPerc=0.50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryStr := util.AddTextQuery(tt.text, tt.params[0], tt.params[1], tt.params[2])
			assert.Equal(t, tt.expected, queryStr)
		})
	}
}

func TestRandomFilterQuery(t *testing.T) {
	tests := []struct {
		name      string
		params    [3]int64
		normalize bool
		expected  string
	}{
		{
			name:      "Test random filter query normalize=false",
			params:    [3]int64{5, -1, 1},
			normalize: false,
			expected:  "?minVal=-1&maxVal=1&kernelSize=5&normalize=false",
		},
		{
			name:      "Test random filter query normalize=true",
			params:    [3]int64{5, -1, 1},
			normalize: true,
			expected:  "?minVal=-1&maxVal=1&kernelSize=5&normalize=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryStr := util.RandomFilterQuery(tt.params[0], tt.params[1], tt.params[2], tt.normalize)
			assert.Equal(t, tt.expected, queryStr)
		})
	}
}

func TestShuffleQuery(t *testing.T) {
	tests := []struct {
		name     string
		param    int64
		expected string
	}{
		{
			name:     "Test shuffle query small number",
			param:    8,
			expected: "?partitions=8",
		},
		{
			name:     "Test shuffle query large number",
			param:    99999,
			expected: "?partitions=99999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryStr := util.ShuffleQuery(tt.param)
			assert.Equal(t, tt.expected, queryStr)
		})
	}
}
