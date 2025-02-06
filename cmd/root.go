package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/piplabs/story-guardian/internal"
	"github.com/piplabs/story-guardian/internal/config"
	"github.com/piplabs/story-guardian/utils"
	"github.com/piplabs/story-guardian/utils/ctxutil"
)

const (
	retryDelay    = 3 * time.Second
	retryAttempts = 6

	// filteredReportFileName represents the log filename for storing filtered transactions.
	filteredReportFileName = "filtered_report.log"
)

// Global variable to hold the output directory.
var outputDir string
var filteredReportFilePath = filepath.Join(utils.GetDefaultPath(), filteredReportFileName)

// rootCmd is the root command for the Story Guardian.
var rootCmd = &cobra.Command{
	Use:   "story-guardian",
	Short: "A tool that regularly downloads Bloom filter files and uploads filter report files.",
	Run: func(cmd *cobra.Command, args []string) {
		startTask(cmd.Context())
	},
}

// Init Initializes the command-line flags and bind them to Viper configurations.
func Init() {
	// Register the output directory flag and bind it with Viper.
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "o", utils.GetDefaultPath(), "Directory to store the bloom filter file")
	if err := viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output-dir")); err != nil {
		log.Fatalf("failed to bind output flag with viper, err: %v", err)
	}

	conf, err := config.NewAppConfig()
	if err != nil {
		log.Fatalf("failed to initialize configuration: %v", err)
	}
	rootCmd.SetContext(ctxutil.WithAppConfig(context.Background(), conf))

	log.Println("Configuration initialized successfully.")
}

// Execute is the main entry point to start the Cobra CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("command execution failed, err: %v", err)
	}
}

// startTask initializes a periodic task, downloading Bloom filter files and uploading filter report files once a day.
func startTask(ctx context.Context) {
	for {
		// Calculate the time to next midnight.
		now := time.Now()
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		sleepDuration := nextMidnight.Sub(now)

		// Create a timer to wait until next midnight or until the context is done.
		timer := time.NewTimer(sleepDuration)
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Println("startTask: received context cancellation, shutting down.")
			return
		case <-timer.C:
			// Continue with the task execution below.
		}

		conf := ctxutil.GetAppConfig(ctx)
		accessToken, err := internal.FetchAccessToken(ctx, conf.ClientID, conf.ClientSecret)
		if err != nil {
			log.Printf("failed to fetch access token: %v", err)
			continue
		}
		// Add the fetched access token into the context
		ctx = ctxutil.WithAccessToken(ctx, accessToken)

		// Retry and download the file again after the sleep period.
		downloadAndRetry(ctx)

		// Retry and upload the file again after the sleep period.
		// TODO: @stevemilk - Deal with the filtered report file
		// uploadAndRetry(ctx)
	}
}

// downloadAndRetry downloads the bloom filter file with a retry mechanism.
func downloadAndRetry(ctx context.Context) {
	err := retry.Do(
		func() error {
			// Attempt to download and store bloom filter
			if err := internal.DownloadAndSaveBloomFilter(ctx, outputDir); err != nil {
				return fmt.Errorf("download failed: %w", err)
			}
			return nil
		},
		retry.Delay(retryDelay),
		retry.Attempts(retryAttempts),
		retry.RetryIf(func(err error) bool {
			// Check for context-related errors
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				log.Printf("Context-related error occurred: %v, will not retry", err)
				return false
			}
			return true
		}),
	)
	if err != nil {
		log.Printf("Failed to download bloom filter after retries: %v", err)
	} else {
		log.Printf("Successfully downloaded bloom filter to %s", outputDir)
	}
}

// uploadAndRetry uploads the report file with a retry mechanism.
func uploadAndRetry(ctx context.Context) {
	err := retry.Do(
		func() error {
			// Attempt to upload bloom filter
			if err := internal.UploadReportFile(ctx, filteredReportFilePath); err != nil {
				return fmt.Errorf("upload failed: %w", err)
			}
			return nil
		},
		retry.Delay(retryDelay),
		retry.Attempts(retryAttempts),
		retry.RetryIf(func(err error) bool {
			// Check for context-related errors
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				log.Printf("Context-related error occurred: %v, will not retry", err)
				return false
			}
			return true
		}),
	)
	if err != nil {
		log.Printf("Failed to upload report file after retries: %v", err)
	} else {
		log.Printf("Successfully uploaded report file")
	}
}
