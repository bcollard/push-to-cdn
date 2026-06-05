package gcs

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// Client wraps cloud.google.com/go/storage with the small subset we need.
type Client struct {
	c      *storage.Client
	bucket string
}

func New(ctx context.Context, bucket string) (*Client, error) {
	if bucket == "" {
		return nil, fmt.Errorf("bucket name is required")
	}
	c, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating GCS client (is ADC configured? run: gcloud auth application-default login): %w", err)
	}
	return &Client{c: c, bucket: bucket}, nil
}

func (c *Client) Close() error { return c.c.Close() }

// UploadOpts controls object metadata on upload.
type UploadOpts struct {
	ContentType  string // overrides extension-based detection
	CacheControl string // e.g. "public, max-age=3600"
}

// Upload streams src to objectName in the bucket, returning the final object name.
func (c *Client) Upload(ctx context.Context, src io.Reader, objectName string, opts UploadOpts) (string, error) {
	objectName = strings.TrimLeft(objectName, "/")
	obj := c.c.Bucket(c.bucket).Object(objectName)
	w := obj.NewWriter(ctx)
	if opts.ContentType != "" {
		w.ContentType = opts.ContentType
	}
	if opts.CacheControl != "" {
		w.CacheControl = opts.CacheControl
	}
	if _, err := io.Copy(w, src); err != nil {
		_ = w.Close()
		return "", fmt.Errorf("writing object: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("finalizing object: %w", err)
	}
	return objectName, nil
}

// UploadFile is a convenience wrapper that opens path and uploads it,
// inferring content-type from the extension if opts.ContentType is empty.
func (c *Client) UploadFile(ctx context.Context, path, objectName string, opts UploadOpts) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if opts.ContentType == "" {
		if ct := mime.TypeByExtension(filepath.Ext(path)); ct != "" {
			opts.ContentType = ct
		}
	}
	return c.Upload(ctx, f, objectName, opts)
}

// Object describes a single object as listed by List.
type Object struct {
	Name    string
	Size    int64
	Updated string
}

// List returns objects whose name starts with prefix.
func (c *Client) List(ctx context.Context, prefix string) ([]Object, error) {
	it := c.c.Bucket(c.bucket).Objects(ctx, &storage.Query{Prefix: prefix})
	var out []Object
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		out = append(out, Object{
			Name:    attrs.Name,
			Size:    attrs.Size,
			Updated: attrs.Updated.UTC().Format("2006-01-02T15:04:05Z"),
		})
	}
	return out, nil
}

// Delete removes the named object.
func (c *Client) Delete(ctx context.Context, objectName string) error {
	objectName = strings.TrimLeft(objectName, "/")
	return c.c.Bucket(c.bucket).Object(objectName).Delete(ctx)
}
