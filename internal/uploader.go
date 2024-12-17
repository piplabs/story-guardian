package internal

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// UploadReportFile uploads the filtered report file to the CipherOwl server.
func UploadReportFile(ctx context.Context, filePath string) error {
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	srcFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	dstFile, err := w.CreateFormFile("file", filepath.Base(srcFile.Name()))
	if err != nil {
		return err
	}
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	// Upload the report file
	if err := uploadReportFile(ctx, buf, w.FormDataContentType()); err != nil {
		return err
	}

	// Remove the file after uploading
	return os.Remove(filePath)
}
