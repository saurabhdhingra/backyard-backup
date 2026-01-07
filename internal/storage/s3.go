package storage

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	Config Config
	sess   *session.Session
}

func NewS3(cfg Config) (*S3, error) {
	awsConfig := &aws.Config{
		Region: aws.String(cfg.Region),
	}

	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, "")
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	return &S3{Config: cfg, sess: sess}, nil
}

func (s *S3) Upload(localPath string, remotePath string) error {
	f, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", localPath, err)
	}
	defer f.Close()

	uploader := s3manager.NewUploader(s.sess)

	// Combine BasePath and remotePath if BasePath is set (as prefix)
	key := remotePath
	if s.Config.BasePath != "" {
		key = fmt.Sprintf("%s/%s", s.Config.BasePath, remotePath)
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	return nil
}

func (s *S3) Download(remotePath string, localPath string) error {
	f, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", localPath, err)
	}
	defer f.Close()

	downloader := s3manager.NewDownloader(s.sess)

	key := remotePath
	if s.Config.BasePath != "" {
		key = fmt.Sprintf("%s/%s", s.Config.BasePath, remotePath)
	}

	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to download file, %v", err)
	}

	return nil
}

func (s *S3) StreamUpload(reader io.Reader, remotePath string) error {
	uploader := s3manager.NewUploader(s.sess)

	key := remotePath
	if s.Config.BasePath != "" {
		key = fmt.Sprintf("%s/%s", s.Config.BasePath, remotePath)
	}

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	if err != nil {
		return fmt.Errorf("failed to upload stream, %v", err)
	}

	return nil
}
