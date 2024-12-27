package internal

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/piplabs/story-guardian/utils/ctxutil"
)

func Test_FetchAccessToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	type args struct {
		clientID     string
		clientSecret string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		mock    func()
	}{
		{
			name: "successful fetch access token",
			args: args{
				clientID:     "test_client_id",
				clientSecret: "test_client_secret",
			},
			want:    "test_access_token",
			wantErr: false,
			mock: func() {
				httpmock.RegisterResponder(http.MethodPost, accessTokenURL,
					httpmock.NewStringResponder(http.StatusOK, `{"access_token": "test_access_token"}`))
			},
		},
		{
			name: "failed fetch access token",
			args: args{
				clientID:     "test_client_id",
				clientSecret: "test_client_secret",
			},
			want:    "",
			wantErr: true,
			mock: func() {
				httpmock.RegisterResponder(http.MethodPost, accessTokenURL,
					httpmock.NewStringResponder(http.StatusBadRequest, `{"error": "invalid_client"}`))
			},
		},
	}
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchAccessToken(tt.args.clientID, tt.args.clientSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fetchAccessToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fetchBloomFilterPresignedURL(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	ctx := context.Background()
	ctxutil.WithAccessToken(ctx, "test_access_token")

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		mock    func()
	}{
		{
			name: "successful fetch bloom filter presigned URL",
			args: args{
				ctx: ctx,
			},
			want:    "test_presigned_url",
			wantErr: false,
			mock: func() {
				httpmock.RegisterResponder(http.MethodGet, bloomFilterFileURL,
					httpmock.NewStringResponder(http.StatusOK, `{"presignedUrl": "test_presigned_url"}`))
			},
		},
		{
			name: "failed fetch bloom filter presigned URL",
			args: args{
				ctx: ctx,
			},
			want:    "",
			wantErr: true,
			mock: func() {
				httpmock.RegisterResponder(http.MethodGet, bloomFilterFileURL,
					httpmock.NewStringResponder(http.StatusUnauthorized, `{"error": "invalid_token"}`))
			},
		},
	}
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchBloomFilterPresignedURL(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchBloomFilterPresignedURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fetchBloomFilterPresignedURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uploadReportFile(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	ctx := context.Background()
	ctxutil.WithAccessToken(ctx, "test_access_token")

	type args struct {
		ctx         context.Context
		buf         *bytes.Buffer
		contentType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "successful upload report file",
			args: args{
				ctx:         ctx,
				buf:         bytes.NewBuffer([]byte("test_report_file")),
				contentType: "application/json",
			},
			wantErr: false,
			mock: func() {
				httpmock.RegisterResponder(http.MethodPost, uploadFileURL,
					httpmock.NewStringResponder(http.StatusOK, `{"status": "success"}`))
			},
		},
		{
			name: "failed upload report file",
			args: args{
				ctx:         ctx,
				buf:         bytes.NewBuffer([]byte("test_report_file")),
				contentType: "application/json",
			},
			wantErr: true,
			mock: func() {
				httpmock.RegisterResponder(http.MethodPost, uploadFileURL,
					httpmock.NewStringResponder(http.StatusBadRequest, `{"error": "invalid_request"}`))
			},
		},
	}
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := uploadReportFile(tt.args.ctx, tt.args.buf, tt.args.contentType); (err != nil) != tt.wantErr {
				t.Errorf("uploadReportFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
