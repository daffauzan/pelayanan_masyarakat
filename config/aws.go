package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client

func InitAWS() {
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatal("Unable to load AWS config: ", err)
	}

	S3Client = s3.NewFromConfig(cfg)

	if endpoint := strings.TrimSpace(os.Getenv("S3_ENDPOINT")); endpoint != "" {
		S3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		})
		log.Printf("AWS S3 initialized with custom endpoint: %s", endpoint)
		return
	}

	log.Println("AWS S3 initialized successfully")
}

func GetBucketName() string {
	return os.Getenv("AWS_BUCKET_NAME")
}

func GetCloudFrontURL() string {
	return os.Getenv("CLOUDFRONT_URL")
}

func GetObjectBaseURL() string {
	if cloudfront := strings.TrimSpace(GetCloudFrontURL()); cloudfront != "" {
		return strings.TrimRight(cloudfront, "/")
	}

	bucket := strings.TrimSpace(GetBucketName())
	if endpoint := strings.TrimSpace(os.Getenv("S3_ENDPOINT")); endpoint != "" {
		return fmt.Sprintf("%s/%s", strings.TrimRight(endpoint, "/"), bucket)
	}

	region := strings.TrimSpace(os.Getenv("AWS_REGION"))
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucket, region)
}
