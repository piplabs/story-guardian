package cmd

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/piplabs/story-guardian/internal"
	"github.com/piplabs/story-guardian/utils"
	"github.com/piplabs/story-guardian/utils/ctxutil"
)

func Test_downloadAndRetry(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	outputDir = utils.GetDefaultPath()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "successful download and save",
			args: args{
				ctx: func() context.Context {
					ctx := context.Background()
					ctxutil.WithAccessToken(ctx, "test_access_token")
					return ctx
				}(),
			},
			mock: func() {
				httpmock.RegisterResponder(http.MethodGet, internal.BloomFilterFileURL,
					httpmock.NewStringResponder(http.StatusOK, `{"presignedUrl": "test_presigned_url"}`))

				httpmock.RegisterResponder(http.MethodGet, "test_presigned_url",
					httpmock.NewStringResponder(http.StatusOK, "bloom_filter_data"))
			},
		}, {
			name: "ctx canceled",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					ctxutil.WithAccessToken(ctx, "test_access_token")
					cancel()
					return ctx
				}(),
			},
			mock: func() {
				httpmock.RegisterResponder(http.MethodGet, internal.BloomFilterFileURL,
					func(request *http.Request) (response *http.Response, err error) {
						return nil, context.Canceled
					})
			},
		},
	}
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		t.Run(tt.name, func(t *testing.T) {
			downloadAndRetry(tt.args.ctx)
		})
	}
}

func Test_uploadAndRetry(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	filteredReportFilePath = filepath.Join(os.TempDir(), "filtered_report.log")

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name: "successful upload",
			args: args{
				ctx: func() context.Context {
					ctx := context.Background()
					ctxutil.WithAccessToken(ctx, "test_access_token")
					return ctx
				}(),
			},
			mock: func() {
				httpmock.RegisterResponder(http.MethodPost, internal.UploadFileURL,
					httpmock.NewStringResponder(http.StatusOK, `{"status": "success"}`))
			},
			wantErr: false,
		}, {
			name: "ctx canceled",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					ctxutil.WithAccessToken(ctx, "test_access_token")
					cancel()
					return ctx
				}(),
			},
			mock: func() {
				httpmock.RegisterResponder(http.MethodPost, internal.UploadFileURL,
					func(request *http.Request) (response *http.Response, err error) {
						return nil, context.Canceled
					})
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.OpenFile(filteredReportFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			_, err = file.WriteString("timestamp: 2024-11-14T17:14:05+08:00, filtered_address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266, tx_hash: 0xe3bcd00a87ca32a507c30864511e1469badbed066d719e48c43e4b2fbe2e8b85, type: 0, from: 0x32E89fEAd3b7E77dD8B26206c0607ecC6FAFBa58, to: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266, value: 0, nonce: 0, gas: 0, gas_price: 0")
			if err != nil {
				t.Fatal(err)
			}
			if tt.wantErr {
				defer os.Remove(filteredReportFilePath)
			}

			uploadAndRetry(tt.args.ctx)
		})
	}
}
