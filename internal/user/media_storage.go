package user

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MediaObjectInfo struct {
	Key         string
	Size        int64
	ContentType string
	ETag        string
}

type UserMediaStorage interface {
	CreatePresignedUploadURL(ctx context.Context, key string, contentType string, expires time.Duration) (string, error)
	GetObjectInfo(ctx context.Context, key string) (*MediaObjectInfo, error)
	Delete(ctx context.Context, key string) error
	PublicURL(key string) string
}

type r2UserMediaStorage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucketName    string
	publicBaseURL string
}

func NewR2UserMediaStorage(client *s3.Client, bucketName string, publicBaseURL string) UserMediaStorage {
	return &r2UserMediaStorage{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucketName:    bucketName,
		publicBaseURL: strings.TrimRight(publicBaseURL, "/"),
	}
}

func (s *r2UserMediaStorage) CreatePresignedUploadURL(ctx context.Context, key string, contentType string, expires time.Duration) (string, error) {
	req, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(options *s3.PresignOptions) {
		options.Expires = expires
	})
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (s *r2UserMediaStorage) GetObjectInfo(ctx context.Context, key string) (*MediaObjectInfo, error) {
	out, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return &MediaObjectInfo{
		Key:         key,
		Size:        aws.ToInt64(out.ContentLength),
		ContentType: aws.ToString(out.ContentType),
		ETag:        strings.Trim(aws.ToString(out.ETag), "\""),
	}, nil
}

func (s *r2UserMediaStorage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	return err
}

func (s *r2UserMediaStorage) PublicURL(key string) string {
	if s.publicBaseURL == "" {
		return key
	}
	escapedKey := strings.ReplaceAll(url.PathEscape(key), "%2F", "/")
	return s.publicBaseURL + "/" + escapedKey
}
