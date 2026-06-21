package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"context"
)

func NewR2Client(endpoint, accessKeyID, secretAccessKey string) (*s3.Client, error) {
	// r2Resolver := aws.EndpointResolverWithOptionsFunc(
	// 	func(service, region string, options ...interface{}) (aws.Endpoint, error) {
	// 		return aws.Endpoint{URL: endpoint}, nil
	// 	},
	// )
	cfg, err := config.LoadDefaultConfig(context.Background(),
	config.WithCredentialsProvider(
		credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
	),
	config.WithRegion("auto"), 
	)
	if err != nil {
		return nil, err
	}
 
	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	}), nil
}