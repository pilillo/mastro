package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pilillo/mastro/abstract"
)

type S3ConnDetails struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

type s3Crawler struct {
	Client *minio.Client
}

func (crawler *s3Crawler) InitConnection(details *S3ConnDetails) (*s3Crawler, error) {

	minioClient, err := minio.New(details.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(details.AccessKeyID, details.SecretAccessKey, ""),
		Secure: details.UseSSL,
	})

	if err != nil {
		return nil, err
	}

	return &s3Crawler{
		Client: minioClient,
	}, nil

}

func (crawler *s3Crawler) Walk(bucket string) ([]minio.ObjectInfo, error) {

	exists, errBucketExists := crawler.Client.BucketExists(context.Background(), bucket)
	if errBucketExists != nil {
		return nil, errBucketExists
	}

	if !exists {
		return nil, fmt.Errorf("bucket %s does not exist", bucket)
	}

	return crawler.ListObjects(bucket, "", true, abstract.DefaultManifestFilename)
}

func (crawler *s3Crawler) ListBuckets() ([]minio.BucketInfo, error) {
	return crawler.Client.ListBuckets(context.Background())
}

func (crawler *s3Crawler) ListObjects(bucket string, prefix string, recursive bool, manifest string) ([]minio.ObjectInfo, error) {
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	ctx := context.Background()

	opts := minio.ListObjectsOptions{
		Recursive: recursive,
		Prefix:    prefix,
	}

	objectCh := crawler.Client.ListObjects(ctx, bucket, opts)
	var slice []minio.ObjectInfo

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		slice = append(slice, object)
	}

	return slice, nil
}

func (crawler *s3Crawler) WalkWithFilter(root string, filter string) ([]abstract.Asset, error) {
	return nil, nil
}
