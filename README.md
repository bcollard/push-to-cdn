# push-to-cdn

A tiny Go CLI (`pushcdn`) to publish public assets — a logo for an email signature,
a resume PDF, files you share by URL — to a Google Cloud Storage bucket fronted by
a Google Cloud Load Balancer at [cdn.runlocal.dev](https://cdn.runlocal.dev).

No auth on reads. No auth on the CLI side beyond your normal `gcloud` Application
Default Credentials.


## Install

### Homebrew (macOS / Linux)

```bash
brew tap bcollard/push-to-cdn
brew install --cask push-to-cdn
```

### go install

```bash
go install github.com/bcollard/push-to-cdn@latest
```

### Download binary

Grab the latest release from [GitHub Releases](https://github.com/bcollard/push-to-cdn/releases),
extract the archive and move `pushcdn` to somewhere on your `PATH`.

### Build from source

```bash
git clone https://github.com/bcollard/push-to-cdn
cd push-to-cdn
make install
```


## Prerequisites

- `pushcdn` binary (see Install above)
- `gcloud auth application-default login` (the CLI uses ADC — no service account JSON files)
- A public GCS bucket — provisioned by the [Terraform](./terraform/) in this repo


## Quick start

```bash
# one-time setup
gcloud auth application-default login
pushcdn config set bucket   cdn.runlocal.dev
pushcdn config set base-url https://cdn.runlocal.dev

# upload
pushcdn upload ~/Pictures/logo.png
# logo.png -> cdn.runlocal.dev/logo.png
#   https://cdn.runlocal.dev/logo.png

# list
pushcdn ls

# delete
pushcdn rm old-thing.png

# print the public URL for an object
pushcdn url logo.png
```


## Commands

### `pushcdn upload`

```bash
pushcdn upload <file> [more-files...]
  -d, --dest string            destination object name or folder (folders end with /)
      --content-type string    override Content-Type (default: inferred from extension)
      --cache-control string   Cache-Control header (default: "public, max-age=3600")
      --no-progress            disable the per-file progress bar
```

Each upload shows a byte-level progress bar on stderr when stdout is a terminal;
piping the command (or passing `--no-progress`) suppresses it cleanly.

Examples:

```bash
pushcdn upload logo.png                        # → /logo.png
pushcdn upload logo.png -d brand/              # → /brand/logo.png
pushcdn upload logo.png -d brand/icon.png      # → /brand/icon.png (rename)
pushcdn upload *.png -d gallery/               # multiple files into a prefix
pushcdn upload resume.pdf --cache-control "public, max-age=300"
```

### `pushcdn ls`

```bash
pushcdn ls                       # everything in the bucket
pushcdn ls brand/                # prefix match — objects under "brand/"
pushcdn ls '*gorilla*'           # glob — names with "gorilla", no '/'
pushcdn ls '**gorilla**'         # glob — same, but spans folder boundaries
pushcdn ls 'logo.???'            # logo.png, logo.svg, …
```

Quote glob patterns so your shell doesn't expand them first. The arg is sent to
GCS as a server-side `matchGlob`, so listings stay efficient even for large
buckets.

### `pushcdn rm`

```bash
pushcdn rm old.png
pushcdn rm a.png b.png c.png
```

### `pushcdn url`

```bash
pushcdn url logo.png
# https://cdn.runlocal.dev/logo.png
```

### `pushcdn config`

```bash
pushcdn config list                                  # keys + descriptions + examples
pushcdn config show                                  # current values
pushcdn config set bucket   cdn.runlocal.dev        # bare name; gs:// stripped if present
pushcdn config set project  my-gcp-project           # optional, informational
pushcdn config set base-url https://cdn.runlocal.dev # scheme required; trailing / stripped
```

### `pushcdn version`

```bash
pushcdn version
pushcdn --version
```


## Configuration

Resolution order: **env var → stored config → built-in default** (only `base-url`
has a default, derived from `bucket`).

| Key        | Env var               | Description |
|------------|-----------------------|---|
| `bucket`   | `PUSHCDN_BUCKET`      | GCS bucket name |
| `project`  | `PUSHCDN_PROJECT`     | GCP project ID (informational; ADC handles auth) |
| `base-url` | `PUSHCDN_BASE_URL`    | Public URL prefix used by `pushcdn url` and post-upload output |

Stored at `~/.config/pushcdn/config.json` (mode `0600`).


## Infrastructure

[`./terraform/`](./terraform/) provisions the public GCS bucket. The Google Cloud
Load Balancer (URL map, target proxy, forwarding rule, SSL cert) and DNS A/AAAA
record for `cdn.runlocal.dev` are managed outside this repo.

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
# edit terraform.tfvars to set project = "..."
terraform init
terraform apply
```

See [`terraform/README.md`](./terraform/README.md) for the domain-verification
gotcha (the default `cdn.runlocal.dev` bucket name is dotted, which GCS only
allows after verifying domain ownership in Google Search Console).


## Distribution

Tagged releases (`v*.*.*`) are built with [goreleaser](https://goreleaser.com/)
via `.github/workflows/release.yml` and published to:

- GitHub Releases (tarballs for macOS/Linux × amd64/arm64)
- The `bcollard/homebrew-push-to-cdn` tap as a cask
