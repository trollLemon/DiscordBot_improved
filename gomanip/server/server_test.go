package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"gocv.io/x/gocv"

	"goManip/jobdispatch"
	"goManip/jobs"
	"goManip/server"
	"goManip/worker"
)

func generateTestImage(width, height int) *gocv.Mat {
	image := gocv.NewMatWithSize(width, height, gocv.MatTypeCV8UC3)
	rng := gocv.TheRNG()
	rng.Fill(&image, gocv.RNGDistNormal, 1.0, 0.0, false)

	return &image
}

var tests = []struct {
	name           string
	URI            string
	fileType       string
	wantStatusCode int
	wantErr        *server.GomanipError
	image          *gocv.Mat
}{

	{
		name:           "test invert endpoint png",
		URI:            "/invert/",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test invert endpoint jpeg",
		URI:            "/invert/",
		fileType:       "jpeg",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test invert endpoint bad file type",
		URI:            "/invert/",
		fileType:       "webp",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Given filetype is not supported, please use jpeg or png images.",
		},
	},

	{
		name:           "test saturation endpoint png",
		URI:            "/saturate/?saturation=2.0",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test saturation endpoint jpeg",
		URI:            "/saturate/?saturation=2.0",
		fileType:       "jpeg",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test saturation endpoint bad file type",
		URI:            "/saturate/?saturation=2.0",
		fileType:       "webp",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Given filetype is not supported, please use jpeg or png images.",
		},
	},
	{
		name:           "test saturation endpoint bad params",
		URI:            "/saturate/?saturation=-2.0",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "expected saturation value to be greater than 0, got -2.0",
		},
	},
	{
		name:           "test saturation endpoint parse failure",
		URI:            "/saturate/?saturateeion=-2.0",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Failed to parse query parameters. Check the request URI.",
		},
	},
	{
		name:           "test edge detection endpoint png",
		URI:            "/edgeDetection/?lower=120&higher=200",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test edge detection endpoint jpeg",
		URI:            "/edgeDetection/?lower=120&higher=200",
		fileType:       "jpeg",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test edge detection endpoint bad file type",
		URI:            "/edgeDetection/?lower=120&higher=200",
		fileType:       "webp",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Given filetype is not supported, please use jpeg or png images.",
		},
	},
	{
		name:           "test edge detection endpoint bad params",
		URI:            "/edgeDetection/?lower=-120&higher=200",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "expected t_lower and t_higher to be greater than or equal to 0, got -120.0 and 200.0",
		},
	},
	{
		name:           "test edge detection endpoint parse failure",
		URI:            "/edgeDetection/?lower=-120&higheeeer=200",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Failed to parse query parameters. Check the request URI.",
		},
	},

	{
		name:           "test morpholoy endpoint png",
		URI:            "/morphology/?type=Erode&kernelSize=3&iterations=5",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test morpholoy endpoint jpeg",
		URI:            "/morphology/?type=Dilate&kernelSize=3&iterations=5",
		fileType:       "jpeg",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test morpholoy endpoint bad file type",
		URI:            "/morphology/?type=Erode&kernelSize=3&iterations=5",
		fileType:       "webp",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Given filetype is not supported, please use jpeg or png images.",
		},
	},
	{
		name:           "test morpholoy endpoint bad params",
		URI:            "/morphology/?type=Erode&kernelSize=-3&iterations=5",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "expected kernel size and iterations to be greater than 0, got -3 and 5",
		},
	},
	{
		name:           "test morpholoy endpoint parse failure",
		URI:            "/morphology/?tyype=Erode",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Failed to parse query parameters. Check the request URI.",
		},
	},
	{
		name:           "test reduce endpoint png",
		URI:            "/reduction/?quality=0.4",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test reduce endpoint jpeg",
		URI:            "/reduction/?quality=0.4",
		fileType:       "jpeg",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test reduce endpoint bad file type",
		URI:            "/reduction/?quality=0.4",
		fileType:       "webp",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Given filetype is not supported, please use jpeg or png images.",
		},
	},
	{
		name:           "test reduce endpoint bad params",
		URI:            "/reduction/?quality=-0.4",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "expected quality to be greater than 0.0, got -0.4",
		},
	},
	{
		name:           "test reduce endpoint parse failure",
		URI:            "/reduction/?qualeity=0.4",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Failed to parse query parameters. Check the request URI.",
		},
	},
	{
		name:           "test text endpoint png",
		URI:            "/text/?text=foo&fontScale=1.0&xPerc=0.5&yPerc=0.5",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test text endpoint jpeg",
		URI:            "/text/?text=foo&fontScale=1.0&xPerc=0.5&yPerc=0.5",
		fileType:       "jpeg",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test text endpoint bad file type",
		URI:            "/text/?text=foo&fontScale=1.0&xPerc=0.5&yPerc=0.5",
		fileType:       "webp",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Given filetype is not supported, please use jpeg or png images.",
		},
	},
	{
		name:           "test text endpoint bad params",
		URI:            "/text/?text=foo&fontScale=0.0&xPerc=0.5&yPerc=0.5",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "expected font scale to be greater than 0, got 0.0",
		},
	},
	{
		name:           "test text endpoint bad params again",
		URI:            "/text/?text=foo&fontScale=1.0&xPerc=-0.5&yPerc=0.5",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "expected x and y percentages to be between 0 and 1, got -0.5 and 0.5",
		},
	},
	{
		name:           "test text endpoint parse failure",
		URI:            "/text/?text=foo&fontScalfe=1.0&xPerc=0.5&yPerc=0.5",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Failed to parse query parameters. Check the request URI.",
		},
	},
	{
		name:           "test random filter endpoint png",
		URI:            "/randomFilter/?minVal=-1&maxVal=1&kernelSize=5&normalize=true",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test random filter endpoint jpeg",
		URI:            "/randomFilter/?minVal=-1&maxVal=1&kernelSize=5&normalize=true",
		fileType:       "jpeg",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test random filter endpoint bad file type",
		URI:            "/randomFilter/?minVal=-1&maxVal=1&kernelSize=5&normalize=true",
		fileType:       "webp",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Given filetype is not supported, please use jpeg or png images.",
		},
	},
	{
		name:           "test random filter endpoint bad params",
		URI:            "/randomFilter/?minVal=-1&maxVal=1&kernelSize=-5&normalize=true",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "expected kernel size to be greater than 0, got -5",
		},
	},
	{
		name:           "test random filter endpoint parse failure",
		URI:            "/randomFilter/?minVal=-0.1&maxVal=1.0&normalize=true",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Failed to parse query parameters. Check the request URI.",
		},
	},

	{
		name:           "test shuffle endpoint png",
		URI:            "/shuffle/?partitions=4",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test shuffle endpoint jpeg",
		URI:            "/shuffle/?partitions=4",
		fileType:       "jpeg",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusOK,
	},

	{
		name:           "test shuffle endpoint bad file type",
		URI:            "/shuffle/?partitions=4",
		fileType:       "webp",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Given filetype is not supported, please use jpeg or png images.",
		},
	},
	{
		name:           "test shuffle endpoint bad params",
		URI:            "/shuffle/?partitions=-4",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "expected partitions to be greater than 1, got -4",
		},
	},
	{
		name:           "test shuffle endpoint too many partitions",
		URI:            "/shuffle/?partitions=99999999",
		fileType:       "png",
		image:          generateTestImage(192, 80),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "cannot fit 99999999 partitions in a 192 by 80 image",
		},
	},
	{
		name:           "test suffle endpoint parse failure",
		URI:            "/shuffle/?paretsts=4",
		fileType:       "png",
		image:          generateTestImage(1920, 1080),
		wantStatusCode: http.StatusBadRequest,
		wantErr: &server.GomanipError{
			Status: "400",
			Detail: "Failed to parse query parameters. Check the request URI.",
		},
	},
}

func TestServer(t *testing.T) {
	defer goleak.VerifyNone(t)
	numWorkers := 1
	jobReqs := make(chan *jobs.JobRequest, numWorkers)
	maxTime := time.Second * 10
	jobDispatcher := jobdispatch.NewJobDispatcher(jobReqs, maxTime)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	for workerId := range numWorkers {
		wg.Add(1)
		go worker.Worker(ctx, workerId+1, jobReqs, wg)
	}

	defer server.GraceFullShutdown(jobDispatcher, wg, cancel)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			server.InitRouting(e, jobDispatcher)

			contentType := "image/" + tt.fileType

			// encode the raw mat as an image format before sending the bytes over
			imgBytes, err := gocv.IMEncode(gocv.FileExt("."+tt.fileType), *tt.image)
			assert.Nil(t, err)
			defer tt.image.Close()

			reader := bytes.NewReader(imgBytes.GetBytes())

			req := httptest.NewRequest(http.MethodPost, tt.URI, reader)
			req.Header.Set("Content-Type", contentType)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatusCode, rec.Code)

			body, err := io.ReadAll(rec.Body)
			assert.Nil(t, err)

			if tt.wantErr != nil {
				var goManipErr *server.GomanipError
				err = json.Unmarshal(body, &goManipErr)
				assert.Nil(t, err)
				assert.Equal(t, tt.wantErr, goManipErr)
			}
		})

	}

}
