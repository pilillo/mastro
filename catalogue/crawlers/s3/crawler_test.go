package s3

import (
	"context"
	"os"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
)

var (
	details *S3ConnDetails
	crawler *S3Crawler
	connErr error
)

func TestMain(m *testing.M) {
	details = &S3ConnDetails{
		Endpoint:        "play.min.io",
		AccessKeyID:     "Q3AM3UQ867SPQQA43P2F",
		SecretAccessKey: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
		UseSSL:          true,
	}

	crawler, connErr = InitConnection(details)

	// run module tests
	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestConnection(t *testing.T) {
	t.Log("TestConnection running")
	assert := assert.New(t)
	assert.Equal(connErr, nil)
	assert.NotEqual(crawler.Client, nil)
}

func TestWalk(t *testing.T) {
	t.Log("TestWalk running")
	assert := assert.New(t)

	bucketName := "mastrobucket"
	// create bucket if not existing
	crawler.Client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})

	// Name of the object
	objectName := "exampleObject"
	// Path to file to be uploaded
	filePath := "file.csv"

	f, err := os.Create(filePath)
	_, _ = f.WriteString("hello world\n")
	f.Close()

	size, err := crawler.Client.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: "application/csv"})

	assert.Equal(err, nil)
	assert.NotEqual(size, 0)

	// test walk function
	fs, err := crawler.Walk(bucketName)

	assert.Equal(err, nil)
	assert.NotEqual(fs, nil)

	// remove remote object and bucket and ciao
	err = crawler.Client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	assert.Equal(err, nil)
	err = crawler.Client.RemoveBucket(context.Background(), bucketName)
	assert.Equal(err, nil)

	err = os.Remove(filePath)
	assert.Equal(err, nil)
}
