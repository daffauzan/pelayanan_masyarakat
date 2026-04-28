package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"pelayanan_publik/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// UploadFile mengunggah file ke S3 dan mengembalikan URL CloudFront-nya.
// folder: subfolder di dalam bucket, misal "surat" atau "pengaduan"
func UploadFile(file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	ext := filepath.Ext(header.Filename)
	key := fmt.Sprintf("%s/%d%s", folder, time.Now().UnixNano(), ext)

	_, err := config.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(config.GetBucketName()),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("gagal upload ke S3: %w", err)
	}

	url := fmt.Sprintf("%s/%s", config.GetCloudFrontURL(), key)
	return url, nil
}

// DeleteFile menghapus file dari S3 berdasarkan key-nya.
func DeleteFile(key string) error {
	_, err := config.S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(config.GetBucketName()),
		Key:    aws.String(key),
	})
	return err
}
