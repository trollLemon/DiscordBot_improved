package util_test

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"goManip/util"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestContext(params map[string]string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	query := req.URL.Query()
	for key, val := range params {
		query.Set(key, val)
	}
	req.URL.RawQuery = query.Encode()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c
}

func TestParseSaturation(t *testing.T) {
	tests := []struct {
		name        string
		param       map[string]string
		wantErr     bool
		expectedVal float32
	}{
		{
			name: "valid saturation",
			param: map[string]string{
				"saturation": "1.0",
			},
			wantErr:     false,
			expectedVal: 1.0,
		},
		{
			name: "invalid saturation",
			param: map[string]string{
				"saturation": "one",
			},
			wantErr:     true,
			expectedVal: 0.0,
		},
		{
			name: "missing saturation",
			param: map[string]string{
				"notSat": "one",
			},
			wantErr:     true,
			expectedVal: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(tt.param)
			sat, err := util.ParseSaturation(ctx)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expectedVal, sat)
		})
	}
}

func TestParseEdgeDetection(t *testing.T) {
	tests := []struct {
		name         string
		param        map[string]string
		wantErr      bool
		expectedVals [2]float32
	}{
		{
			name: "valid edge detection parameters",
			param: map[string]string{
				"lower":  "120",
				"higher": "200",
			},
			wantErr:      false,
			expectedVals: [2]float32{120, 200},
		},
		{
			name: "invalid lower parameter",
			param: map[string]string{
				"lower":  "abc",
				"higher": "200",
			},
			wantErr:      true,
			expectedVals: [2]float32{0, 0},
		},
		{
			name: "invalid upper parameter",
			param: map[string]string{
				"lower":  "120",
				"higher": "abc",
			},
			wantErr:      true,
			expectedVals: [2]float32{0, 0},
		},
		{
			name: "missing lower parameter",
			param: map[string]string{
				"higher": "200",
			},
			wantErr:      true,
			expectedVals: [2]float32{0, 0},
		},
		{
			name: "missing upper parameter",
			param: map[string]string{
				"lower": "120",
			},
			wantErr:      true,
			expectedVals: [2]float32{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(tt.param)
			tLower, tHigher, err := util.ParseEdgeDetection(ctx)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expectedVals[0], tLower)
			assert.Equal(t, tt.expectedVals[1], tHigher)
		})
	}
}

func TestParseReduce(t *testing.T) {
	tests := []struct {
		name        string
		param       map[string]string
		wantErr     bool
		expectedVal float32
	}{
		{
			name: "valid reduce parameter",
			param: map[string]string{
				"quality": "0.4",
			},
			wantErr:     false,
			expectedVal: 0.4,
		},
		{
			name: "invalid reduce parameter",
			param: map[string]string{
				"quality": "abc",
			},
			wantErr:     true,
			expectedVal: 0.0,
		},
		{
			name:        "missing reduce parameter",
			param:       map[string]string{},
			wantErr:     true,
			expectedVal: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(tt.param)
			reduce, err := util.ParseReduce(ctx)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expectedVal, reduce)
		})
	}
}

func TestParseMorphology(t *testing.T) {
	tests := []struct {
		name         string
		param        map[string]string
		wantErr      bool
		expectedVals [2]int
		expectedOp   string
	}{
		{
			name: "valid morphology: dilation",
			param: map[string]string{
				"type":       "Dilate",
				"kernelSize": "3",
				"iterations": "4",
			},
			wantErr:      false,
			expectedVals: [2]int{3, 4},
			expectedOp:   "Dilate",
		},
		{
			name: "valid morphology: erosion",
			param: map[string]string{
				"type":       "Erosion",
				"kernelSize": "3",
				"iterations": "4",
			},
			wantErr:      false,
			expectedVals: [2]int{3, 4},
			expectedOp:   "Erosion",
		},
		{
			name: "missing type",
			param: map[string]string{
				"kernelSize": "3",
				"iterations": "4",
			},
			wantErr:      true,
			expectedVals: [2]int{0, 0},
			expectedOp:   "",
		},
		{
			name: "missing kernelSize",
			param: map[string]string{
				"type":       "Erosion",
				"iterations": "4",
			},
			wantErr:      true,
			expectedVals: [2]int{0, 0},
			expectedOp:   "",
		},

		{
			name: "missing iterations",
			param: map[string]string{
				"type":       "Erosion",
				"kernelSize": "3",
			},
			wantErr:      true,
			expectedVals: [2]int{0, 0},
			expectedOp:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(tt.param)

			op, kernelSize, iterations, err := util.ParseMorphology(ctx)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expectedOp, op)
			assert.Equal(t, tt.expectedVals[0], kernelSize)
			assert.Equal(t, tt.expectedVals[1], iterations)
		})
	}
}

func TestParseAddText(t *testing.T) {
	tests := []struct {
		name         string
		param        map[string]string
		wantErr      bool
		expectedVals [3]float64
		expectedText string
	}{
		{
			name: "valid text and floats",
			param: map[string]string{
				"text":      "I love golang",
				"fontScale": "1.0",
				"xPerc":     "0.5",
				"yPerc":     "0.5",
			},
			wantErr:      false,
			expectedVals: [3]float64{1.0, 0.5, 0.5},
			expectedText: "I love golang",
		},
		{
			name: "invalid font scale ",
			param: map[string]string{
				"text":      "I love golang",
				"fontScale": "aaa",
				"xPerc":     "0.5",
				"yPerc":     "0.5",
			},
			wantErr:      true,
			expectedVals: [3]float64{0.0, 0.0, 0.0},
			expectedText: "",
		},
		{
			name: "invalid xPerc ",
			param: map[string]string{
				"text":      "I love golang",
				"fontScale": "1.0",
				"xPerc":     "hdjhsk",
				"yPerc":     "0.5",
			},
			wantErr:      true,
			expectedVals: [3]float64{0.0, 0.0, 0.0},
			expectedText: "",
		},
		{
			name: "invalid yPerc ",
			param: map[string]string{
				"text":      "I love golang",
				"fontScale": "1.0",
				"xPerc":     "0.5",
				"yPerc":     "hdjhsk",
			},
			wantErr:      true,
			expectedVals: [3]float64{0.0, 0.0, 0.0},
			expectedText: "",
		},
		{
			name: "missing font scale ",
			param: map[string]string{
				"text":  "I love golang",
				"xPerc": "0.5",
				"yPerc": "0.5",
			},
			wantErr:      true,
			expectedVals: [3]float64{0.0, 0.0, 0.0},
			expectedText: "",
		},
		{
			name: "missing xPerc ",
			param: map[string]string{
				"text":      "I love golang",
				"fontScale": "1.0",
				"yPerc":     "0.5",
			},
			wantErr:      true,
			expectedVals: [3]float64{0.0, 0.0, 0.0},
			expectedText: "",
		},
		{
			name: "missing yPerc ",
			param: map[string]string{
				"text":      "I love golang",
				"fontScale": "1.0",
				"xPerc":     "0.5",
			},
			wantErr:      true,
			expectedVals: [3]float64{0.0, 0.0, 0.0},
			expectedText: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(tt.param)

			text, fontScale, xPerc, yPerc, err := util.ParseAddText(ctx)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expectedText, text)
			assert.Equal(t, tt.expectedVals[0], fontScale)
			assert.Equal(t, tt.expectedVals[1], xPerc)
			assert.Equal(t, tt.expectedVals[2], yPerc)
		})
	}
}

func TestParseRandomFilter(t *testing.T) {
	tests := []struct {
		name         string
		params       map[string]string
		wantErr      bool
		expectedVals [3]int
		normalize    bool
	}{
		{
			name: "valid input with normalize true",
			params: map[string]string{
				"minVal":     "1",
				"maxVal":     "10",
				"kernelSize": "3",
				"normalize":  "true",
			},
			wantErr:      false,
			expectedVals: [3]int{1, 10, 3},
			normalize:    true,
		},
		{
			name: "invalid minVal",
			params: map[string]string{
				"minVal":     "abc",
				"maxVal":     "10",
				"kernelSize": "3",
			},
			wantErr: true,
		},
		{
			name: "invalid maxVal",
			params: map[string]string{
				"minVal":     "10",
				"maxVal":     "abs",
				"kernelSize": "3",
			},
			wantErr: true,
		},
		{
			name: "missing maxVal",
			params: map[string]string{
				"minVal":     "1",
				"kernelSize": "3",
			},
			wantErr: true,
		},
		{
			name: "missing minVal",
			params: map[string]string{
				"maxVal":     "1",
				"kernelSize": "3",
			},
			wantErr: true,
		},
		{
			name: "missing kernelSize",
			params: map[string]string{
				"maxVal": "1",
				"minVal": "-1",
			},
			wantErr: true,
		},

		{
			name: "normalize false",
			params: map[string]string{
				"minVal":     "1",
				"maxVal":     "10",
				"kernelSize": "3",
				"normalize":  "false",
			},
			wantErr:      false,
			expectedVals: [3]int{1, 10, 3},
			normalize:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(tt.params)
			minVal, maxVal, kernelSize, normalize, err := util.ParseRandomFilter(ctx)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expectedVals[0], minVal)
			assert.Equal(t, tt.expectedVals[1], maxVal)
			assert.Equal(t, tt.expectedVals[2], kernelSize)
			assert.Equal(t, tt.normalize, normalize)

		})
	}
}

func TestParseShuffle(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]string
		wantErr     bool
		expectedVal int
	}{
		{
			name: "valid partition param",
			params: map[string]string{
				"partitions": "3",
			},
			wantErr:     false,
			expectedVal: 3,
		},
		{
			name: "invalid partition param",
			params: map[string]string{
				"partitions": "abc",
			},
			wantErr:     true,
			expectedVal: 0,
		},
		{
			name:        "missing partition param",
			params:      map[string]string{},
			wantErr:     true,
			expectedVal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(tt.params)
			partitions, err := util.ParseShuffle(ctx)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expectedVal, partitions)

		})
	}
}
