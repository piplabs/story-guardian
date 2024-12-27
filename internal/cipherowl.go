package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
	bloomFilterFileURL = baseAPIURL + bloomFilterPath
	uploadFileURL      = baseAPIURL + uploadFilePath
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
func FetchAccessToken(clientID, clientSecret string) (string, error) {
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

	// Send POST request
	resp, err := http.Post(accessTokenURL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Verify HTTP response code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("access token request failed with status %d", resp.StatusCode)
	}

	// Decode response JSON
	var tokenResponse oAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

// fetchBloomFilterPresignedURL retrieves the presigned URL for the bloom filter file.
func fetchBloomFilterPresignedURL(ctx context.Context) (string, error) {
	// Create HTTP GET request
	req, err := http.NewRequest(http.MethodGet, bloomFilterFileURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+ctxutil.GetAccessToken(ctx))

	// Perform the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Verify HTTP response code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("presigned URL request failed with status %d", resp.StatusCode)
	}

	// Decode response JSON
	var urlResp presignedURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&urlResp); err != nil {
		return "", err
	}

	return urlResp.PresignedURL, nil
}

func uploadReportFile(ctx context.Context, buf *bytes.Buffer, contentType string) error {
	req, err := http.NewRequest(http.MethodPost, uploadFileURL, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+ctxutil.GetAccessToken(ctx))

	// Perform the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verify HTTP response code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload request failed with status %d", resp.StatusCode)
	}

	return nil
}
