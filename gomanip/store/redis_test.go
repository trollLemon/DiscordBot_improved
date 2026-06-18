package store

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func TestRedisStore_SaveAndGetImage(t *testing.T) {
	tests := []struct {
		name      string
		jobID     string
		image     []byte
		setup     func(context.Context, *RedisStore, string) error
		wantImage []byte
		wantTrace string
		wantErr   error
	}{
		{
			name:      "Save and get image success",
			jobID:     "job-123",
			image:     []byte{0x01, 0x02, 0x03, 0x04},
			wantImage: []byte{0x01, 0x02, 0x03, 0x04},
			wantTrace: "",
		},
		{
			name:      "Get missing image returns not found",
			jobID:     "does-not-exist",
			wantImage: nil,
			wantTrace: "",
			wantErr:   ErrImageNotFound,
		},
		{
			name:  "Get returns job failed and trace when error is stored",
			jobID: "job-err-456",
			setup: func(ctx context.Context, s *RedisStore, id string) error {
				return s.WriteError(ctx, id, errors.New("processing failed"), "trace-abc-123")
			},
			wantImage: nil,
			wantTrace: "trace-abc-123",
			wantErr:   ErrJobFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mr, err := miniredis.Run()
			if err != nil {
				t.Fatalf("failed to start miniredis: %v", err)
			}
			defer mr.Close()

			store := NewRedisStore(mr.Addr(), "", "", 0, time.Hour)
			defer store.client.Close()

			if tt.image != nil {
				if err := store.SaveImage(ctx, tt.jobID, tt.image); err != nil {
					t.Fatalf("SaveImage() unexpected error: %v", err)
				}
			}

			if tt.setup != nil {
				if err := tt.setup(ctx, store, tt.jobID); err != nil {
					t.Fatalf("setup() unexpected error: %v", err)
				}
			}

			got, gotTrace, err := store.GetImage(ctx, tt.jobID)

			if tt.wantErr == nil && err != nil {
				t.Fatalf("GetImage() unexpected error: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("GetImage() error = %v, want errors.Is(..., %v)", err, tt.wantErr)
			}

			if !bytes.Equal(got, tt.wantImage) {
				t.Fatalf("GetImage() = %v, want %v", got, tt.wantImage)
			}

			if gotTrace != tt.wantTrace {
				t.Fatalf("GetImage() trace = %q, want %q", gotTrace, tt.wantTrace)
			}
		})
	}
}
