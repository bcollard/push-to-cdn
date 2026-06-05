package gcs

import (
	"context"
	"fmt"
	"io"
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

// Bucket returns the bucket name this client was constructed with.
func (c *Client) Bucket() string { return c.bucket }

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

// Object describes a single object as listed by List.
type Object struct {
	Name    string
	Size    int64
	Updated string
}

// ListQuery selects which objects List returns. Exactly one of Prefix or
// MatchGlob is typically set; both may be empty to list everything.
//
// MatchGlob uses Cloud Storage server-side glob syntax:
//   - "*"  matches any character except "/"
//   - "**" matches any character including "/"
//   - "?"  matches exactly one character except "/"
//   - "[abc]" / "[a-z]" character classes
//
// See https://cloud.google.com/storage/docs/json_api/v1/objects/list#parameters
type ListQuery struct {
	Prefix    string
	MatchGlob string
}

// List returns objects matching q.
func (c *Client) List(ctx context.Context, q ListQuery) ([]Object, error) {
	it := c.c.Bucket(c.bucket).Objects(ctx, &storage.Query{
		Prefix:    q.Prefix,
		MatchGlob: q.MatchGlob,
	})
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
