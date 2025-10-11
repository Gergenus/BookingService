package s3

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// юзер и пароль
func InitS3Storage(endpoint, accessKeyId, secretAccessKey, bucketName string) *minio.Client {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		panic("s3 storage does not work: " + err.Error())
	}

	ok, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		panic("bucket checking error: " + err.Error())
	}
	if !ok {
		err := client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{
			Region: "moscow",
		})
		if err != nil {
			panic("s3 bucket was not created: " + err.Error())
		}
	}
	return client
}
