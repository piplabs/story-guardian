package internal

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/piplabs/story-guardian/internal/pkg/httpclient"
)

const (
	bloomFilterFilename = "bloom_filter.gob"
)

// DownloadAndSaveBloomFilter retrieves and saves the bloom filter file to the specified location.
func DownloadAndSaveBloomFilter(ctx context.Context, outputDir string) error {
	// Ensure the output directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}
	}

	// Retrieve presigned file URL
	presignedURL, err := fetchBloomFilterPresignedURL(ctx)
	if err != nil {
		return err
	}

	// Use the default HTTP client from the httpclient package
	client := httpclient.DefaultClient()

	resp, err := client.Do(ctx, http.MethodGet, presignedURL, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Save file to output directory
	filePath := filepath.Join(outputDir, bloomFilterFilename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}
