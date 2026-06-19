package worker_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/goleak"
	"gocv.io/x/gocv"

	"goManip/jobs"
	"goManip/worker"
)

type NoOpStore struct{}

func (s NoOpStore) SaveImage(ctx context.Context, jobId string, image []byte) error {
	return nil
}

func (s NoOpStore) WriteError(ctx context.Context, jobId string, err error, traceID string) error {
	return nil
}

func (s NoOpStore) GetImage(ctx context.Context, jobId string) ([]byte, string, error) {
	return nil, "", nil
}

type NoopAttribute struct{}

func (NoopAttribute) ToAttributes() []attribute.KeyValue {
	return nil
}

type MockOperationSuccess struct {
	NoopAttribute
}
type MockOperationErr struct {
	NoopAttribute
}

type MockOperationTimeOut struct {
	NoopAttribute
}

func (m MockOperationErr) Run(input *gocv.Mat) (*gocv.Mat, error) {
	return nil, errors.New("error processing job")
}

func (m MockOperationSuccess) Run(input *gocv.Mat) (*gocv.Mat, error) {
	clone := input.Clone()
	return &clone, nil
}

func (m MockOperationTimeOut) Run(input *gocv.Mat) (*gocv.Mat, error) {
	time.Sleep(1 * time.Second)
	clone := input.Clone()

	return &clone, nil
}

func TestWorker(t *testing.T) {

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	var tests = []struct {
		name      string
		operation jobs.Operation
	}{
		{
			name:      "Test Success",
			operation: MockOperationSuccess{},
		},
		{
			name:      "Test Error",
			operation: MockOperationErr{},
		},
	}

	defer goleak.VerifyNone(t)

	timeLimit := 5 * time.Second

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testImage := gocv.NewMatWithSize(64, 64, gocv.MatTypeCV8UC3)
			defer testImage.Close()

			shutDownCtx, shutDownCancel := context.WithCancel(context.Background())

			workerWg := &sync.WaitGroup{}

			jobReqs := make(chan *jobs.JobRequest)

			worker := worker.NewWorker(shutDownCtx, 0, jobReqs, workerWg, NoOpStore{}, logger)

			workerWg.Add(1)
			go worker.Work()

			job := jobs.NewJob("0", tt.operation, &testImage)

			jobReqs <- &jobs.JobRequest{
				Job: job,
				Ctx: context.Background(),
			}

			shutDownCancel()
			close(jobReqs)
			workerWg.Wait()

			deadline := time.Now().Add(timeLimit)
			for job.GetEndTime().IsZero() && time.Now().Before(deadline) {
				time.Sleep(10 * time.Millisecond)
			}

			if job.GetEndTime().IsZero() {
				t.Fatalf("TestWorker() %s, job was not processed before timeout", tt.name)
			}
		})
	}
}

func TestWorkerShutdown(t *testing.T) {
	wg := &sync.WaitGroup{}
	jobReqs := make(chan *jobs.JobRequest)
	ctx, cancel := context.WithCancel(context.Background())
	worker := worker.NewWorker(ctx, 0, jobReqs, wg, NoOpStore{}, slog.New(slog.NewJSONHandler(io.Discard, nil)))

	wg.Add(1)
	go worker.Work()
	cancel()
	close(jobReqs)
	wg.Wait()
}

func TestWorkerCancel(t *testing.T) {
	wg := &sync.WaitGroup{}
	jobReqs := make(chan *jobs.JobRequest)
	testImage := gocv.NewMat()
	defer testImage.Close()
	ctx, cancel := context.WithCancel(context.Background())
	worker := worker.NewWorker(ctx, 0, jobReqs, wg, NoOpStore{}, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	wg.Add(1)
	go worker.Work()

	jobReqs <- &jobs.JobRequest{
		Job: jobs.NewJob("0", MockOperationTimeOut{}, &testImage),
		Ctx: context.Background(),
	}
	time.Sleep(50 * time.Millisecond)
	cancel()

	close(jobReqs)
	wg.Wait()
}
