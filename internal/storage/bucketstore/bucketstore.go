package bucketstore

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FakeWriterAt struct {
	w io.WriteCloser
}

type BucketStore struct {
	s3client   *s3.Client
	downloader *manager.Downloader
	uploader   *manager.Uploader
	Bucket     string
	Prefix     string
}

type BucketOptions struct {
	Data     io.Reader
	Filename string
}

const (
	FilenameFormat = "%v-20060102-150405.sql.gz"
)

func New(region, bucket, prefix string) (*BucketStore, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	s3client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(s3client)
	downloader := manager.NewDownloader(s3client)
	downloader.Concurrency = 1

	return &BucketStore{
		s3client:   s3client,
		uploader:   uploader,
		downloader: downloader,
		Bucket:     bucket,
		Prefix:     prefix,
	}, nil
}

func (b *BucketStore) Store(data io.Reader, filename string) (string, error) {
	filename = fmt.Sprintf(FilenameFormat, filename)
	filename = time.Now().Format(filename)
	key := fmt.Sprintf("%v/%v", b.Prefix, filename)

	out, err := b.uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(b.Bucket),
		Key:    aws.String(key),
		Body:   data,
	})
	if err != nil {
		return "", fmt.Errorf("Unable to upload: %w", err)
	}

	return out.Location, nil
}

func (b *BucketStore) Retrieve(data io.WriteCloser, filename string) error {
	key := fmt.Sprintf("%v/%v", b.Prefix, filename)

	w := FakeWriterAt{
		w: data,
	}

	_, err := b.downloader.Download(context.Background(), w,
		&s3.GetObjectInput{
			Bucket: aws.String(b.Bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return fmt.Errorf("Unable to download: %w", err)
	}

	data.Close()

	return nil
}

func (b *BucketStore) List() ([]string, error) {
	var res []string

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(b.Bucket),
		Prefix: aws.String(b.Prefix),
	}

	out, err := b.s3client.ListObjectsV2(context.Background(), input)
	if err != nil {
		return res, fmt.Errorf("unable to list objects: %w", err)
	}

	for _, v := range out.Contents {
		v := strings.TrimPrefix(*v.Key, b.Prefix)
		v = strings.TrimPrefix(v, "/")
		res = append(res, v)
	}

	return res, nil
}

func (fw FakeWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	return fw.w.Write(p)
}
