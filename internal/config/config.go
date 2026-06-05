package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
