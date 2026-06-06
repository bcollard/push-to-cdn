# Terraform — public CDN bucket

Provisions a single public GCS bucket (`cdn.runlocal.dev` by default) that holds
publicly readable assets fronted by a Google Cloud Load Balancer.

## What this manages

- `google_storage_bucket.cdn` — uniform bucket-level access, CORS open for GET/HEAD,
  optional website config for `index.html` / `404.html`.
- `google_storage_bucket_iam_member.public_read` — `roles/storage.objectViewer`
  for `allUsers`.

## What this does NOT manage

The domain and load balancer (URL map, target proxy, forwarding rule, SSL cert,
DNS A/AAAA record) are handled outside this stack. See `outputs.tf` for the
gcloud commands to wire them up.

## Prerequisites

- `gcloud auth application-default login` (or a SA key via `GOOGLE_APPLICATION_CREDENTIALS`)
- A GCP project with billing enabled
- **If using the dotted default `cdn.runlocal.dev`:** ownership of `runlocal.dev`
  must be verified in [Google Search Console](https://search.google.com/search-console)
  by the same identity running `terraform apply`. Otherwise GCS refuses to create
  the bucket. If you don't want to verify, override `bucket_name` to
  `cdn-runlocal-dev` — the LB can still serve it at `cdn.runlocal.dev`.

## State backend

By default the tracked Terraform has no backend block, so state is stored
locally in `terraform.tfstate` — fine for trying things out. For a real
deployment you'll want a remote backend; copy the example and fill it in:

```bash
cp backend.tf.example backend.tf
# edit backend.tf to set bucket and prefix (and optionally credentials)
```

`backend.tf` is gitignored. Anything other than `backend.tf.example` is yours
to keep private.

## Usage

From the repo root, the Makefile wraps the common flow:

```bash
make tf-init        # terraform init (uses backend.tf if present)
make tf-plan
make tf-apply
make tf-destroy     # if you ever need to tear down
```

Or invoke terraform directly:

```bash
cp terraform.tfvars.example terraform.tfvars
# edit terraform.tfvars to set project = "..."
terraform init
terraform plan
terraform apply
```
