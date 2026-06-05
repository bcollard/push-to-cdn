output "bucket_name" {
  description = "Name of the created bucket — pass this to `pushcdn config set bucket`."
  value       = google_storage_bucket.cdn.name
}

output "bucket_self_link" {
  description = "GCS API self-link for the bucket."
  value       = google_storage_bucket.cdn.self_link
}

output "bucket_gs_uri" {
  description = "gs:// URI of the bucket."
  value       = "gs://${google_storage_bucket.cdn.name}"
}

output "next_steps" {
  description = "Wire-up steps the user owns (domain + load balancer)."
  value       = <<EOT
Next steps (manual — outside of this Terraform):

1. Create a backend bucket pointing at ${google_storage_bucket.cdn.name}:
     gcloud compute backend-buckets create cdn-runlocal-dev \
       --gcs-bucket-name=${google_storage_bucket.cdn.name} \
       --enable-cdn

2. Add a URL map + target HTTPS proxy + global forwarding rule that
   routes cdn.runlocal.dev to the backend bucket above.

3. Point your DNS A/AAAA record for cdn.runlocal.dev at the LB's IP.

4. Configure your CLI:
     pushcdn config set bucket   ${google_storage_bucket.cdn.name}
     pushcdn config set base-url https://cdn.runlocal.dev
EOT
}
