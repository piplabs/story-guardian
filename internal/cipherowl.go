package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/piplabs/story-guardian/internal/pkg/httpclient"
	"github.com/piplabs/story-guardian/utils/ctxutil"
)

const (
	baseAPIURL      = "https://svc.cipherowl.ai/"
	oAuthTokenPath  = "oauth/token"
	bloomFilterPath = "api/bloom-filter/file/1"
	uploadFilePath  = "api/upload/report/v1"
)

var (
	accessTokenURL     = baseAPIURL + oAuthTokenPath
	BloomFilterFileURL = baseAPIURL + bloomFilterPath
	UploadFileURL      = baseAPIURL + uploadFilePath
)

// oAuthTokenRequest is the payload structure for obtaining an access token.
type oAuthTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

// oAuthTokenResponse represents the response structure for the access token request.
type oAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// presignedURLResponse represents the response structure for the bloom filter file's presigned URL request.
type presignedURLResponse struct {
	PresignedURL string `json:"presignedUrl"`
}

// FetchAccessToken retrieves an OAuth access token using client credentials.
func FetchAccessToken(ctx context.Context, clientID, clientSecret string) (string, error) {
	requestPayload := oAuthTokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     "svc.cipherowl.ai",
		GrantType:    "client_credentials",
	}

	// Serialize the request payload
	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return "", err
	}

	// Use the default HTTP client from the httpclient package
	client := httpclient.DefaultClient()

	header := map[string]string{
		httpclient.ContentTypeHeader: httpclient.ContentTypeJSON,
	}
	// Send POST request
	resp, err := client.Do(ctx, http.MethodPost, accessTokenURL, bytes.NewReader(jsonData), header)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode response JSON
	var tokenResponse oAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

// fetchBloomFilterPresignedURL retrieves the presigned URL for the bloom filter file.
func fetchBloomFilterPresignedURL(ctx context.Context) (string, error) {
	// Use the default HTTP client from the httpclient package
	client := httpclient.DefaultClient()

	header := map[string]string{
		httpclient.ContentTypeHeader:   httpclient.ContentTypeJSON,
		httpclient.AuthorizationHeader: "Bearer " + ctxutil.GetAccessToken(ctx),
	}

	// Perform the HTTP request
	resp, err := client.Do(ctx, http.MethodGet, BloomFilterFileURL, nil, header)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode response JSON
	var urlResp presignedURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&urlResp); err != nil {
		return "", err
	}

	return urlResp.PresignedURL, nil
}

// uploadReportFile uploads the filtered report file to the CipherOwl server.
func uploadReportFile(ctx context.Context, buf *bytes.Buffer, contentType string) error {
	// Use the default HTTP client from the httpclient package
	client := httpclient.DefaultClient()

	header := map[string]string{
		httpclient.ContentTypeHeader:   contentType,
		httpclient.AuthorizationHeader: "Bearer " + ctxutil.GetAccessToken(ctx),
	}

	// Perform the HTTP request
	resp, err := client.Do(ctx, http.MethodPost, UploadFileURL, buf, header)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
