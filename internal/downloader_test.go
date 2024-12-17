package internal

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/piplabs/story-guardian/utils/ctxutil"
)

func TestDownloader_DownloadAndSaveBloomFilter(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	ctx := context.Background()
	ctxutil.WithAccessToken(ctx, "test_access_token")

	type args struct {
		ctx       context.Context
		outputDir string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		mock    func()
	}{
		{
			name: "successful download and save",
			args: args{
				ctx:       ctx,
				outputDir: os.TempDir(),
			},
			want:    "bloom_filter_data",
			wantErr: false,
			mock: func() {
				httpmock.RegisterResponder(http.MethodGet, bloomFilterFileURL,
					httpmock.NewStringResponder(http.StatusOK, `{"presignedUrl": "test_presigned_url"}`))

				httpmock.RegisterResponder(http.MethodGet, "test_presigned_url",
					httpmock.NewStringResponder(http.StatusOK, "bloom_filter_data"))
			},
		},
		{
			name: "network error in downloading",
			args: args{
				ctx:       ctx,
				outputDir: os.TempDir(),
			},
			wantErr: true,
			mock: func() {
				httpmock.RegisterResponder(http.MethodGet, bloomFilterFileURL,
					httpmock.NewStringResponder(http.StatusOK, `{"presignedUrl": "test_presigned_url"}`))

				httpmock.RegisterResponder(http.MethodGet, "test_presigned_url",
					httpmock.NewErrorResponder(http.ErrServerClosed))
			},
		},
	}
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := DownloadAndSaveBloomFilter(tt.args.ctx, tt.args.outputDir); (err != nil) != tt.wantErr {
				t.Errorf("DownloadAndSaveBloomFilter() error = %v, wantErr %v", err, tt.wantErr)
			}

			// If the test should succeed, check if the file was written correctly.
			if !tt.wantErr {
				filePath := filepath.Join(tt.args.outputDir, bloomFilterFilename)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Expected file does not exist: %v", filePath)
				}
				defer os.Remove(filePath)

				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read downloaded file: %v", err)
				}

				if string(content) != tt.want {
					t.Errorf("Expected file content %s but got %s", tt.want, string(content))
				}
			}
		})
	}
}
