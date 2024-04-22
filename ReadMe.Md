# Upload File to A S3 bucket using PresingedURLS in Go 

This repository contains a simple rest api that uploads files to an s3 bucket and delete files files from an s3 bucket 

Read More about this project [here](https://dev.to/sha254/upload-and-delete-file-from-amazon-s3-bucket-in-go-using-presigned-urls-17fo)

## Prerequisites

Must have an access key and access secret key from AWS assinged to a user with the appropriate permission for the s3 bucket.

## Usage 
Clone the repository 
```bash
git clone https://github.com/Sean-Miningah/s3BucketWithGo.git

```

Create .env file in the root and populate it with configurations in this format 

```json
AWS_BUCKET_NAME="<bucket-name-here>"
AWS_REGION="<aws-region-here>"
AWS_S3_BUCKET_ACCESS_KEY="<aws-bucket-access-key-here>"
AWS_S3_BUCKET_SECRET_ACCESS_KEY="<aws-bucket-secret-key-here>"
```
## Run the Server
To run the server, execute the go run command 

```bash
go run main.go
```
