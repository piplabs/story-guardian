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

var (
	testReportFilePath = filepath.Join(os.TempDir(), "filtered_report.log")
)

func TestUploadReportFile(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	ctx := context.Background()
	ctxutil.WithAccessToken(ctx, "test_access_token")

	type args struct {
		ctx      context.Context
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "successful upload",
			args: args{
				ctx:      ctx,
				filePath: testReportFilePath,
			},
			wantErr: false,
			mock: func() {
				httpmock.RegisterResponder(http.MethodPost, UploadFileURL,
					httpmock.NewStringResponder(http.StatusOK, `{"status": "success"}`))
			},
		},
		{
			name: "failed upload",
			args: args{
				ctx:      ctx,
				filePath: testReportFilePath,
			},
			wantErr: true,
			mock: func() {
				httpmock.RegisterResponder(http.MethodPost, UploadFileURL,
					httpmock.NewStringResponder(http.StatusBadRequest, `{"error": "invalid_file"}`))
			},
		},
	}
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.OpenFile(testReportFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			_, err = file.WriteString("timestamp: 2024-11-14T17:14:05+08:00, filtered_address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266, tx_hash: 0xe3bcd00a87ca32a507c30864511e1469badbed066d719e48c43e4b2fbe2e8b85, type: 0, from: 0x32E89fEAd3b7E77dD8B26206c0607ecC6FAFBa58, to: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266, value: 0, nonce: 0, gas: 0, gas_price: 0")
			if err != nil {
				t.Fatal(err)
			}
			if !tt.wantErr {
				defer os.Remove(testReportFilePath)
			}

			if err := UploadReportFile(tt.args.ctx, tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("UploadReportFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
