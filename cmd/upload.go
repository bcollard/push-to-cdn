package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bcollard/push-to-cdn/internal/config"
	"github.com/bcollard/push-to-cdn/internal/gcs"
	"github.com/spf13/cobra"
)

var (
	uploadDest         string
	uploadContentType  string
	uploadCacheControl string
)

var uploadCmd = &cobra.Command{
	Use:     "upload <file> [more-files...]",
	Aliases: []string{"up", "put"},
	Short:   "Upload one or more files to the CDN bucket",
	Long: `Upload one or more files to the configured bucket. Each file is stored at the
object name derived from its basename, unless --dest is given.

When --dest ends in '/', it is treated as a folder and each file keeps its basename
under that prefix. When --dest does not end in '/' and a single file is uploaded,
--dest is used as the exact object name. Content-Type is inferred from the file
extension unless --content-type is set.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Resolve()
		if err != nil {
			return err
		}
		bucket, err := cfg.RequireBucket()
		if err != nil {
			return err
		}

		ctx := context.Background()
		client, err := gcs.New(ctx, bucket)
		if err != nil {
			return err
		}
		defer client.Close()

		opts := gcs.UploadOpts{
			ContentType:  uploadContentType,
			CacheControl: uploadCacheControl,
		}

		for _, path := range args {
			objectName := destName(uploadDest, path, len(args))
			name, err := client.UploadFile(ctx, path, objectName, opts)
			if err != nil {
				return fmt.Errorf("uploading %s: %w", path, err)
			}
			fmt.Printf("%s -> %s/%s\n", path, bucket, name)
			if cfg.BaseURL != "" {
				fmt.Printf("  %s/%s\n", strings.TrimRight(cfg.BaseURL, "/"), name)
			}
		}
		return nil
	},
}

// destName resolves a local file path + --dest flag into the target object name.
func destName(dest, localPath string, count int) string {
	base := filepath.Base(localPath)
	switch {
	case dest == "":
		return base
	case strings.HasSuffix(dest, "/"):
		return strings.TrimLeft(dest, "/") + base
	case count > 1:
		// multiple files with a non-/ dest — treat as folder anyway
		return strings.TrimLeft(dest, "/") + "/" + base
	default:
		return strings.TrimLeft(dest, "/")
	}
}

func init() {
	uploadCmd.Flags().StringVarP(&uploadDest, "dest", "d", "", "destination object name or folder (folders end with /)")
	uploadCmd.Flags().StringVar(&uploadContentType, "content-type", "", "override Content-Type (default: inferred from extension)")
	uploadCmd.Flags().StringVar(&uploadCacheControl, "cache-control", "public, max-age=3600", "Cache-Control header to set on the object")
	rootCmd.AddCommand(uploadCmd)
}
