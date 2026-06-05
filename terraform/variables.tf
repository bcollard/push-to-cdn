variable "project" {
  description = "GCP project ID that owns the bucket."
  type        = string
}

variable "region" {
  description = "Default region for the provider. The bucket itself is multi-region (see var.location)."
  type        = string
  default     = "europe-west1"
}

variable "bucket_name" {
  description = <<EOT
Name of the public GCS bucket. Derived from the domain (cdn.runlocal.dev).

Note: dotted bucket names matching a DNS domain require ownership verification in
Google Search Console (https://search.google.com/search-console) before the
bucket can be created. If you have not verified runlocal.dev, either:
  - verify it first, then `terraform apply`; or
  - override to a hyphenated name (e.g. "cdn-runlocal-dev") — your load balancer
    can still front it at cdn.runlocal.dev regardless of the bucket name.
EOT
  type        = string
  default     = "cdn.runlocal.dev"
}

variable "location" {
  description = "Bucket location. Use a multi-region (EU, US, ASIA) for global CDN distribution."
  type        = string
  default     = "EU"
}

variable "storage_class" {
  description = "Storage class for objects in the bucket."
  type        = string
  default     = "STANDARD"
}
