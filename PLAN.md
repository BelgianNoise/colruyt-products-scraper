# Planned changes

Status: ⬜ todo · 🟨 in progress · ✅ done

---

## 1. Move blob storage from GCS to Versity S3 Gateway on the NAS ⬜

**Goal:** Store the JSON dumps on the NAS as plain files, so the existing NAS syncs back them up. Use Versity S3 Gateway (`versitygw`, POSIX backend) for S3 access, with its WebGUI to manage users and keys. Move all existing data off GCS and stop writing to it.

**Why Versity:** It keeps objects as plain files (`<bucket>/<key>` on disk), has built-in user and key management with a WebGUI, and runs in a container. Garage, RustFS and SeaweedFS use their own storage formats. MinIO's community edition was archived in Feb 2026.

### 1a. Install versitygw on the NAS ⬜
- [ ] Create a dataset/folder for the data (e.g. `/data/s3`), plus separate folders for IAM (`/data/iam`) and versioning (`/data/versioning`)
- [ ] Check that the data filesystem supports extended attributes (xattrs); object metadata such as Content-Type is stored there. On ZFS, consider `xattr=sa`
- [ ] Add a `docker-compose.yml` with `versity/versitygw:latest`:
  - `VGW_BACKEND=posix` pointed at the data folder
  - `VGW_IAM_DIR` pointed at the IAM folder (users are stored as a JSON file)
  - `VGW_PORT` (S3 API), `VGW_ADMIN_PORT` (admin API), `VGW_WEBUI_PORT` (WebGUI; off unless set)
  - Root credentials: the docs use both `ROOT_ACCESS_KEY`/`ROOT_SECRET_KEY` and `ROOT_ACCESS_KEY_ID`/`ROOT_SECRET_ACCESS_KEY`; check which ones the image reads
  - Run as a non-root user (the image has no `USER` directive)
- [ ] Add Traefik routes (Docker labels on the versitygw container); the existing Cloudflare Tunnel → Traefik setup handles the rest:
  - `s3.<domain>` → S3 API port (`VGW_PORT`); this is the public one, so GitHub Actions can reach it
  - WebGUI and admin API: separate hosts, kept local-only (e.g. an `ipAllowList` middleware)
  - Don't rewrite the `Host` header (S3 request signatures include it)
  - Set `--cors-allow-origin` on versitygw to the WebGUI host; the WebGUI calls the S3 and admin APIs from the browser
- [ ] Check that a ~50 MB upload goes through (product dumps are ~30–50 MB uncompressed). Don't add a Traefik `buffering` limit below that; Cloudflare's proxy caps request bodies at 100 MB on the Free plan
- [ ] In the WebGUI, create the bucket `colruyt-products` and the users:
  - `writer`: owns the bucket (read/write), for scraper, comparer and pp (all three write files)
  - `reader`: read-only access, for the inserter (bucket policy)
- [ ] Public read access for the website (it loads `drastische-dalers/`, `sterkste-stijgers/`, `prettige-prijzen/` from the public bucket). The GCS bucket is fully public today. Either:
  - set a `public-read` ACL on the bucket (same as now), or
  - add a bucket policy allowing anonymous `GetObject` on just those three prefixes (check that versitygw supports an anonymous principal)
  - and add the website origin to `--cors-allow-origin` if the browser fetches the files directly
- [ ] Backups: add the data folder to the existing syncs, **excluding `.sgwtmp/`** (temporary upload folders). Check the sync keeps xattrs, and also back up the IAM folder

### 1b. Copy the existing data from GCS (~22 GB) ⬜
- [ ] Set up `rclone` with a GCS remote and an S3 remote (`provider = Other`, path-style URLs, endpoint `https://s3.<domain>`)
- [ ] Bulk copy, keeping the key layout unchanged (`colruyt-products/*.json`, `promotions/*.json`):
  `rclone copy gcs:colruyt-products nas:colruyt-products --progress`
  - One-time GCS egress cost: up to ~$2.50–3 for 22 GB
- [ ] Verify: `rclone check gcs:colruyt-products nas:colruyt-products` (counts and sizes match)

### 1c. One shared S3 file service ✅
All storage access goes through one S3-compatible file service in `shared/pkg`. No other package imports a storage SDK, and switching S3 hosts only means changing env vars.

**API** (`shared/pkg/filestore.go`, on top of `aws-sdk-go-v2`):
```go
type FileStore struct { client *s3.Client; bucket string }
type FileInfo struct { Key string; LastModified time.Time; Size int64 }
var ErrNotFound = errors.New("file not found")

func NewFileStore() (*FileStore, error) // from S3_ENDPOINT, S3_BUCKET, S3_ACCESS_KEY_ID, S3_SECRET_ACCESS_KEY, S3_REGION (optional)
func (f *FileStore) Put(ctx context.Context, key string, data []byte, contentType string) error
func (f *FileStore) PutJSON(ctx context.Context, key string, v any) error
func (f *FileStore) Get(ctx context.Context, key string) ([]byte, error) // ErrNotFound if missing
func (f *FileStore) List(ctx context.Context, prefix string) ([]FileInfo, error)
func (f *FileStore) Copy(ctx context.Context, srcKey, dstKey string) error
func (f *FileStore) Delete(ctx context.Context, key string) error
```
- Without `S3_ACCESS_KEY_ID`, requests are anonymous (works against a public-read bucket)
- Client settings: `UsePathStyle: true`, `RequestChecksumCalculation: WhenRequired` (some S3-compatible servers reject the SDK's default checksums), region `us-east-1`
- `List` uses the `ListObjectsV2` paginator and returns errors (the current loop spins forever on a persistent error)
- Each `main()` creates one `FileStore` and passes it down. The bucket name comes from env, so callers only pass keys

**Migration of call sites:**
- [x] Add `filestore.go`; delete `save-json-to-bucket.go`, `get-object-from-bucket.go`, `list-bucket-objects.go`, `gcs-bucket-name.go`
- [x] `shared/pkg/get-products.go`, `compare.go`: take a `*FileStore`. `compare.go` checks `ErrNotFound` instead of an empty result
- [x] `scraper`: `save-to-gcs.go` → `save-to-storage.go` using `PutJSON`; update the log lines
- [x] `inserter`: the promotions time filter uses `FileInfo.LastModified` (was GCS `attrs.Updated`)
- [x] `comparer`: `compare-today.go:21` uses the literal prefix `"colruyt-products/"` (it relied on the bucket name equalling the folder name)
- [x] `pp`: replace the direct GCS client with `Copy` / `Delete` / `Get`
- [x] Remove the `cloud.google.com/go/storage` / `google.golang.org/api` dependencies
- [x] Tested against a local versitygw container: put/get, missing key → `ErrNotFound`, copy/delete (incl. keys with spaces and `+`), listing past 1000 keys, compare + cached compare, anonymous access, plain files on disk with Content-Type in xattrs

### 1d. Cut over ⬜
- [ ] Add GitHub secrets: `S3_ENDPOINT`, `S3_BUCKET`, plus the `writer` keys (scraper, comparer, pp) and `reader` keys (inserter)
- [ ] Update `.github/workflows/scraper.yaml`, `comparer.yaml`, `pp.yaml`, `inserter.yaml`: write the S3 vars to `.env`, remove `GCP_SERVICE_ACCOUNT_KEY` / `GOOGLE_APPLICATION_CREDENTIALS`
- [ ] Add the S3 vars to the local `.env`
- [ ] End-to-end run against the NAS: scrape → compare → pp → insert
- [ ] Update the README (it still describes the public GCS bucket)
- [ ] Point the website at `https://s3.<domain>/colruyt-products/...` instead of the GCS URLs
- [ ] Merge, then run a final `rclone copy` from GCS to the NAS for anything written between the bulk copy and the switch
- [ ] Watch the first few daily runs (the files appear on the NAS and the inserter reads them)
- [ ] Keep the GCS bucket for ~2–4 weeks as a fallback, then delete it and the GCP service account / `GCP_SERVICE_ACCOUNT_KEY` secret

