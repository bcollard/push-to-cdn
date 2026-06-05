package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Bucket  string `json:"bucket,omitempty"`
	Project string `json:"project,omitempty"`
	BaseURL string `json:"base_url,omitempty"`
}

// Keys lists all settable config keys in display order.
var Keys = []string{
	"bucket",
	"project",
	"base-url",
}

// Descriptions explains what each key holds and how it's used.
var Descriptions = map[string]string{
	"bucket":   "GCS bucket name. Bare name only — no gs:// prefix, no scheme.",
	"project":  "GCP project ID that owns the bucket. Informational only — ADC handles auth.",
	"base-url": "Public URL prefix served by your load balancer. Include the scheme; no trailing slash.",
}

// Examples are shown in help / list output to make the expected format obvious.
var Examples = map[string]string{
	"bucket":   "cdn.runlocal.dev",
	"project":  "my-gcp-project",
	"base-url": "https://cdn.runlocal.dev",
}

func (c *Config) Get(key string) string {
	switch key {
	case "bucket":
		return c.Bucket
	case "project":
		return c.Project
	case "base-url":
		return c.BaseURL
	}
	return ""
}

func (c *Config) Set(key, value string) bool {
	switch key {
	case "bucket":
		c.Bucket = value
	case "project":
		c.Project = value
	case "base-url":
		c.BaseURL = value
	default:
		return false
	}
	return true
}

// Normalize cleans up a value before storing it. Returns the normalized value
// plus a human-readable note when something was changed, so the caller can echo it.
func Normalize(key, value string) (string, string) {
	value = strings.TrimSpace(value)
	switch key {
	case "bucket":
		orig := value
		value = strings.TrimPrefix(value, "gs://")
		value = strings.TrimSuffix(value, "/")
		if value != orig {
			return value, "stripped gs:// prefix or trailing /"
		}
	case "base-url":
		orig := value
		value = strings.TrimSuffix(value, "/")
		if value != orig {
			return value, "stripped trailing /"
		}
	}
	return value, ""
}

// Validate returns a non-nil error if the value is unusable for its key.
func Validate(key, value string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", key)
	}
	switch key {
	case "bucket":
		if strings.Contains(value, "/") {
			return fmt.Errorf("bucket %q contains a slash — pass the bare name (e.g. %s)", value, Examples["bucket"])
		}
	case "base-url":
		if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			return fmt.Errorf("base-url %q is missing a scheme — include http:// or https:// (e.g. %s)", value, Examples["base-url"])
		}
	}
	return nil
}

func dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(home, ".config", "pushcdn")
	return d, os.MkdirAll(d, 0700)
}

func path() (string, error) {
	d, err := dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "config.json"), nil
}

func Load() (*Config, error) {
	p, err := path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func Save(c *Config) error {
	p, err := path()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}

// Resolve returns effective values, with priority: env var > stored config.
// Defaults are applied if both are empty.
func Resolve() (*Config, error) {
	c, err := Load()
	if err != nil {
		return nil, err
	}
	if v := os.Getenv("PUSHCDN_BUCKET"); v != "" {
		c.Bucket = v
	}
	if v := os.Getenv("PUSHCDN_PROJECT"); v != "" {
		c.Project = v
	}
	if v := os.Getenv("PUSHCDN_BASE_URL"); v != "" {
		c.BaseURL = v
	}
	if c.BaseURL == "" && c.Bucket != "" {
		c.BaseURL = fmt.Sprintf("https://%s", c.Bucket)
	}
	return c, nil
}

// RequireBucket returns the bucket name or an error suitable for user display.
func (c *Config) RequireBucket() (string, error) {
	if c.Bucket == "" {
		return "", fmt.Errorf("no bucket configured — run: pushcdn config set bucket <name> (or export PUSHCDN_BUCKET)")
	}
	return c.Bucket, nil
}
