package repo

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	AWS_S3_REGION         = ""
	AWS_S3_BUCKET_NAME    = ""
	AWS_ACCESS_KEY_ID     = ""
	AWS_SECRET_ACCESS_KEY = ""
)

type Repo struct {
	s3Client *s3.Client
}

func NewS3Client() *Repo {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			AWS_ACCESS_KEY_ID,
			AWS_SECRET_ACCESS_KEY,
			"",
		)),
		config.WithRegion(AWS_S3_REGION),
	)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	client := s3.NewFromConfig(cfg)
	return &Repo{
		s3Client: client,
	}
}

func (repo Repo) UploadFile(bucketName *string, objectKey *string, fileName string, file io.Reader) error {

	_, err := repo.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: bucketName,
		Key:    objectKey,
		Body:   file,
	})
	if err != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			fileName, bucketName, objectKey, err)
		return err
	}
	return nil
}
