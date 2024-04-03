package repo

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
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
	s3Client          *s3.Client
	s3PresignedClient *s3.PresignClient
}

func NewS3Client(accessKey string, secretKey string, s3BucketRegion string) *Repo {
	fmt.Printf("The keys are: %s", accessKey)
	options := s3.Options{
		Region:      s3BucketRegion,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	}

	client := s3.New(options, func(o *s3.Options) {
		o.Region = s3BucketRegion
		o.UseAccelerate = false
	})

	presignClient := s3.NewPresignClient(client)
	return &Repo{
		s3Client:          client,
		s3PresignedClient: presignClient,
	}
}

func (repo Repo) PutObject(bucketName string, objectKey string, lifetimeSecs int64, size int64) (*v4.PresignedHTTPRequest, error) {
	request, err := repo.s3PresignedClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		// ContentType:   aws.String("jpeg"),
		// ContentLength: aws.Int64(size),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

func (repo Repo) UploadFile(file image.Image, url string) error {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, file, nil)
	if err != nil {
		return nil
	}
	body := io.Reader(&buf)
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "image/jpeg")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close() // Ensure proper closing

	// body, err = httputil.DumpResponse(resp, true)
	// if err != nil {
	// 	log.Println("Error reading response body:", err)
	// 	return err
	// }

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("AWS upload file response body: ", string(bytes))
	// log.Println("AWS ERROR: ", err)
	return err
}
