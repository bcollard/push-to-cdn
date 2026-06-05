# Claude context for push-to-cdn

## What this project is

A small Go CLI (`pushcdn`) that uploads public assets to a single public GCS
bucket fronted by a Google Cloud Load Balancer at `cdn.runlocal.dev`. The
bucket only ever holds publicly readable content (logos, resume, files shared
by URL) — there is no notion of authenticated reads or per-object ACLs.

The repo also contains the Terraform that provisions the bucket. The
load balancer + DNS are managed outside this repo.

## Architecture at a glance

```
main.go
└── cmd/                  cobra commands — root, version, config, upload, ls, rm, url
    ├── root.go           rootCmd + Execute(version, commit, date)
    ├── config.go         "pushcdn config {show,list,set}"
    ├── upload.go         "pushcdn upload" + destName() helper
    ├── ls.go             "pushcdn ls [prefix]"
    ├── rm.go             "pushcdn rm <obj>..."
    ├── url.go            "pushcdn url <obj>"
    └── version.go
internal/
├── config/               persistent CLI config at ~/.config/pushcdn/config.json
│                         keys: bucket, project, base-url
│                         Resolve() layers env vars over stored config
└── gcs/                  thin wrapper around cloud.google.com/go/storage
                          (Upload/UploadFile/List/Delete)
terraform/                public GCS bucket + IAM allUsers:objectViewer
                          versions.tf · variables.tf · main.tf · outputs.tf
.github/workflows/        ci.yml (build/test/vet/goreleaser-check/terraform-validate)
                          release.yml (goreleaser on v*.*.* tag)
.goreleaser.yaml          binary + homebrew cask -> bcollard/homebrew-push-to-cdn
```

## Design decisions worth knowing

- **Auth model.** The CLI uses Application Default Credentials (`gcloud auth
  application-default login`). No token caching, no service account JSON in the
  CLI itself. This is why `internal/config/` is much smaller than the keycloak-cli
  sibling — there are no credentials to store.

- **Bucket name = domain name.** Default is `cdn.runlocal.dev` (dotted). GCS
  requires the bucket creator to have verified ownership of `runlocal.dev` in
  Google Search Console. If you hit that wall, the `bucket_name` variable can be
  overridden to `cdn-runlocal-dev` (hyphens) without breaking anything — the
  load balancer fronts whatever bucket name you pick.

- **No website-vs-CDN split.** The bucket has a `website {}` block but the LB +
  Cloud CDN are what actually serve traffic. The website config is harmless and
  useful for direct `storage.googleapis.com` access during debugging.

- **Public read is bucket-wide.** `roles/storage.objectViewer` on `allUsers` at
  the bucket level, with `uniform_bucket_level_access = true`. Do not add
  per-object ACLs — they are intentionally disabled.

- **Cache-Control defaults.** Upload sets `public, max-age=3600` unless
  `--cache-control` is passed. Cloud CDN honors this.

- **Reference project.** The CLI structure, goreleaser config, and GHA workflows
  are deliberately parallel to [github.com/bcollard/keycloak-cli](https://github.com/bcollard/keycloak-cli)
  so changes to conventions can be propagated between the two.

## Common tasks

- **Add a new subcommand:** drop a new file in `cmd/`, register it in `init()`
  with `rootCmd.AddCommand(...)`.
- **Add a new config key:** append to `Keys` in `internal/config/config.go`, add
  the corresponding field + Get/Set case, and wire an env-var override in
  `Resolve()`.
- **Cut a release:** tag `vX.Y.Z` on `main` and push. `release.yml` runs goreleaser,
  which publishes to GitHub Releases and bumps the homebrew tap.

## Repos that must exist before a release works

- `github.com/bcollard/push-to-cdn` — this repo (binary releases)
- `github.com/bcollard/homebrew-push-to-cdn` — the tap (cask updates)

## Secrets required in CI

- `GITHUB_TOKEN` — provided automatically
- `HOMEBREW_TAP_GITHUB_TOKEN` — PAT with `contents:write` on the tap repo
