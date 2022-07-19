package storage

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type BucketStore struct {
	uploader *s3manager.Uploader
	Bucket   string
	Prefix   string
}

type BucketOptions struct {
	Data     io.Reader
	Filename string
}

func NewBucketStorer(region, bucket, prefix string) (*BucketStore, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create session: %w", err)
	}

	uploader := s3manager.NewUploader(sess)

	return &BucketStore{
		uploader: uploader,
		Bucket:   bucket,
		Prefix:   prefix,
	}, nil
}

func (b *BucketStore) Store(data io.Reader, filename string, done chan bool, result chan string) error {
	filename = fmt.Sprintf(FilenameFormat, filename)
	filename = time.Now().Format(filename)
	key := fmt.Sprintf("%v/%v", b.Prefix, filename)

	go func() error {
		out, err := b.uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(b.Bucket),
			Key:    aws.String(key),
			Body:   data,
		})
		if err != nil {
			return fmt.Errorf("Unable to upload: %w", err)
		}

		result <- out.Location
		done <- true

		return nil
	}()

	return nil
}
