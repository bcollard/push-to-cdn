terraform {
  required_version = ">= 1.5.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }

  # Backend is configured locally — copy backend.tf.example to backend.tf
  # (gitignored) and fill it in. With no backend.tf present, Terraform
  # falls back to local state, which is fine for `validate` in CI.
}

provider "google" {
  project = var.project
  region  = var.region
}
