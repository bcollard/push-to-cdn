resource "google_storage_bucket" "cdn" {
  name     = var.bucket_name
  location = var.location

  storage_class               = var.storage_class
  uniform_bucket_level_access = true
  public_access_prevention    = "inherited"

  versioning {
    enabled = false
  }

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }

  cors {
    origin          = ["*"]
    method          = ["GET", "HEAD"]
    response_header = ["*"]
    max_age_seconds = 3600
  }
}

# Public read for every object. The bucket only ever holds public assets
# (logos, resume, files shared by URL) — this is the whole point.
resource "google_storage_bucket_iam_member" "public_read" {
  bucket = google_storage_bucket.cdn.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}
