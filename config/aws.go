package config

import (
	"context"
	"log"
	"os"

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

	log.Println("AWS S3 initialized successfully")
}

func GetBucketName() string {
	return os.Getenv("AWS_BUCKET_NAME")
}

func GetCloudFrontURL() string {
	return os.Getenv("CLOUDFRONT_URL")
}
