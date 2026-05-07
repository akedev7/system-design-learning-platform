package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	cfg "go-course-service/internal/config"
)

type S3Client struct {
	client    *s3.Client
	presign   *s3.PresignClient
	bucket    string
	publicURL string
}

func NewS3Client(s3Cfg *cfg.S3Config) (*S3Client, error) {
	if s3Cfg.Endpoint == "" {
		return nil, errors.New("S3 endpoint is required")
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(s3Cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s3Cfg.AccessKey,
			s3Cfg.SecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	var opts []func(*s3.Options)
	if s3Cfg.Endpoint != "" {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(s3Cfg.Endpoint)
		})
	}
	if s3Cfg.ForcePathStyle {
		opts = append(opts, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}

	s3Client := s3.NewFromConfig(awsCfg, opts...)
	presignClient := s3.NewPresignClient(s3Client)

	return &S3Client{
		client:    s3Client,
		presign:   presignClient,
		bucket:    s3Cfg.Bucket,
		publicURL: s3Cfg.PublicURL,
	}, nil
}

func (c *S3Client) GeneratePresignedUpload(ctx context.Context, key string, contentType string, expiresIn time.Duration) (string, error) {
	presignedURL, err := c.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return presignedURL.URL, nil
}

func (c *S3Client) GeneratePresignedDownload(ctx context.Context, key string, expiresIn time.Duration) (string, error) {
	presignedURL, err := c.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return presignedURL.URL, nil
}

func (c *S3Client) Upload(ctx context.Context, key string, body []byte, contentType string) error {
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}
	return nil
}

func (c *S3Client) Delete(ctx context.Context, key string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}
	return nil
}

func (c *S3Client) GetPublicURL(key string) string {
	if c.publicURL != "" {
		return fmt.Sprintf("%s/%s", c.publicURL, key)
	}
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", c.bucket, key)
}

func (c *S3Client) GetObject(ctx context.Context, key string) ([]byte, error) {
	resp, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}