package shared

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

// ErrNotFound is returned by FileStore.Get when the key does not exist.
var ErrNotFound = errors.New("file not found")

// FileStore is the single entry point for file storage. It talks to any
// S3-compatible service (versitygw, R2, AWS S3, ...) and is configured
// through environment variables:
//
//	S3_ENDPOINT           e.g. https://s3.example.com (required)
//	S3_BUCKET             bucket name (required)
//	S3_ACCESS_KEY_ID      optional; without keys requests are anonymous
//	S3_SECRET_ACCESS_KEY  optional
//	S3_REGION             optional, defaults to us-east-1
type FileStore struct {
	client *s3.Client
	bucket string
}

type FileInfo struct {
	Key          string
	LastModified time.Time
	Size         int64
}

func NewFileStore() (*FileStore, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	bucket := os.Getenv("S3_BUCKET")
	if endpoint == "" || bucket == "" {
		return nil, errors.New("S3_ENDPOINT and S3_BUCKET must be set")
	}
	region := os.Getenv("S3_REGION")
	if region == "" {
		region = "us-east-1"
	}

	var creds aws.CredentialsProvider = aws.AnonymousCredentials{}
	if key := os.Getenv("S3_ACCESS_KEY_ID"); key != "" {
		creds = credentials.NewStaticCredentialsProvider(key, os.Getenv("S3_SECRET_ACCESS_KEY"), "")
	}

	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(endpoint),
		Region:       region,
		Credentials:  creds,
		// Self-hosted S3 servers expect https://host/bucket/key URLs.
		UsePathStyle: true,
		// Newer SDK versions add checksums to every request by default, which
		// several S3-compatible servers reject. Only send them when required.
		RequestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
		ResponseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
	})

	return &FileStore{client: client, bucket: bucket}, nil
}

func (f *FileStore) Put(ctx context.Context, key string, data []byte, contentType string) error {
	_, err := f.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(f.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("put %q: %w", key, err)
	}
	return nil
}

func (f *FileStore) PutJSON(ctx context.Context, key string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("serialize %q: %w", key, err)
	}
	return f.Put(ctx, key, data, "application/json")
}

// Get returns the contents of key, or ErrNotFound if it does not exist.
func (f *FileStore) Get(ctx context.Context, key string) ([]byte, error) {
	out, err := f.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("get %q: %w", key, ErrNotFound)
		}
		return nil, fmt.Errorf("get %q: %w", key, err)
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", key, err)
	}
	return data, nil
}

// List returns every file whose key starts with prefix.
func (f *FileStore) List(ctx context.Context, prefix string) ([]FileInfo, error) {
	var files []FileInfo
	pages := s3.NewListObjectsV2Paginator(f.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(f.bucket),
		Prefix: aws.String(prefix),
	})
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list %q: %w", prefix, err)
		}
		for _, obj := range page.Contents {
			files = append(files, FileInfo{
				Key:          aws.ToString(obj.Key),
				LastModified: aws.ToTime(obj.LastModified),
				Size:         aws.ToInt64(obj.Size),
			})
		}
	}
	return files, nil
}

// Copy copies srcKey to dstKey within the bucket, keeping its metadata.
func (f *FileStore) Copy(ctx context.Context, srcKey, dstKey string) error {
	_, err := f.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(f.bucket),
		Key:        aws.String(dstKey),
		CopySource: aws.String(escapePath(f.bucket + "/" + srcKey)),
	})
	if err != nil {
		return fmt.Errorf("copy %q to %q: %w", srcKey, dstKey, err)
	}
	return nil
}

func (f *FileStore) Delete(ctx context.Context, key string) error {
	_, err := f.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("delete %q: %w", key, err)
	}
	return nil
}

func isNotFound(err error) bool {
	var noSuchKey *types.NoSuchKey
	if errors.As(err, &noSuchKey) {
		return true
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "NoSuchKey", "NotFound":
			return true
		}
	}
	return false
}

// escapePath URL-encodes each segment of a path but keeps the slashes,
// as required for the CopySource header. Everything except unreserved
// characters is encoded; a literal "+" must be %2B or servers read a space.
func escapePath(p string) string {
	segments := strings.Split(p, "/")
	for i, s := range segments {
		segments[i] = strings.ReplaceAll(url.QueryEscape(s), "+", "%20")
	}
	return strings.Join(segments, "/")
}
