package services

import (
	"go-template/internal/common"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectMinio() (*minio.Client, error) {
	secretKey := common.GetEnv("MINIO_SECRET_KEY")
	accessKey := common.GetEnv("MINIO_ACCESS_KEY")
	// endpoint := "minio-test.alieexe.tech"
	endpoint := "s3.amazonaws.com"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
		Region: "ap-southeast-1",
	})

	if err != nil {
		return nil, err
	}

	// buckets, err := minioClient.ListBuckets(context.Background())
	// if err != nil {
	// 	return nil, err
	// }

	// for _, bucket := range buckets {
	// 	fmt.Println("bucket name:", bucket.Name)
	// }

	// if err != nil {
	// 	return nil, err
	// }

	return minioClient, nil
}
