package storage

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type FakeWriterAt struct {
    w io.Writer
}

type BucketStore struct {
	uploader 	*s3manager.Uploader
	downloader	*s3manager.Downloader
	Bucket   	string
	Prefix   	string
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
	downloader := s3manager.NewDownloader(sess)
	downloader.Concurrency = 1

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

func (b *BucketStore) Retrieve(data io.Writer, filename string, done chan bool) error {
	key := fmt.Sprintf("%v/%v", b.Prefix, filename)

	go func() error {
		_, err := b.downloader.Download(
			FakeWriterAt{
				w: data,
			},
			&s3.GetObjectInput{
				Bucket: aws.String(b.Bucket),
				Key:    aws.String(key),
			},
		)
		if err != nil {
			return fmt.Errorf("Unable to download: %w", err)
		}

		done <- true

		return nil
	}()

	return nil
}

func (fw FakeWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
    // ignore 'offset' because we forced sequential downloads
    return fw.w.Write(p)
}
