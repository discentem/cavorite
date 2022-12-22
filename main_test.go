package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"path"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	pantris3 "github.com/discentem/pantri_but_go/stores/s3"
)

type s3Client struct {
	*s3.Client
}

func (c s3Client) createBucket(ctx context.Context, bucketAddress string) error {
	_, bucket := path.Split(bucketAddress)
	_, err := c.Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucket,
	})
	if err == nil {
		log.Printf("created bucket %s", bucketAddress)
	}
	return err
}

func upload(ctx context.Context, bucketAddress, key string, b []byte) error {
	cfg, err := pantris3.GetConfig(bucketAddress)
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(s3.NewFromConfig(*cfg))
	uploader.Concurrency = 3

	_, bucketName := path.Split(bucketAddress)
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(b),
	})
	if err != nil {
		log.Print("uploader.Upload returned an error")
		return err
	}
	log.Printf("uploaded %s", key)
	return nil
}

func (c s3Client) listObjects(bucketName string) error {
	listObjsResponse, err := c.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(""),
	})

	if err != nil {
		return err
	}

	for _, object := range listObjsResponse.Contents {
		log.Printf("%s (%d bytes, class %v) \n", *object.Key, object.Size, object.StorageClass)
	}
	return nil
}

func NewClient(bucketAddress string) (*s3Client, error) {
	cfg, err := pantris3.GetConfig(bucketAddress)
	if err != nil {
		return nil, err
	}
	return &s3Client{
		Client: s3.NewFromConfig(*cfg),
	}, nil
}

func TestMain(T *testing.T) {
	ctx := context.Background()

	serverAddress := flag.String("bucket_address", "http://127.0.0.1:9000/test", "address to bucket")
	// sourceRepo := flag.String("source_repo", "./_ci/integration_tests/artifacts/fake_git_repo", "path to source repo")
	flag.Parse()
	s3c, err := NewClient(*serverAddress)
	if err != nil {
		T.Fatal(err)
	}
	_, bucket := path.Split(*serverAddress)
	if err := s3c.createBucket(ctx, bucket); err != nil {
		if !strings.Contains(err.Error(), "BucketAlreadyOwnedByYou") {
			T.Fatal(err)
		}
	}

	main()
	// cmd := exec.Command(
	// 	*binPath,
	// 	"--source_repo",
	// 	*sourceRepo,
	// 	"init",
	// 	"--pantri_address",
	// 	*serverAddress,
	// 	"s3",
	// )
	// cmd.Env = os.Environ()
	// cmd.Env = append(cmd.Env, "AWS_ACCESS_KEY_ID=minioadmin")
	// cmd.Env = append(cmd.Env, "AWS_SECRET_ACCESS_KEY=minioadmin")
	// stdout, err := cmd.Output()
	// if err != nil {
	// 	log.Print(err)
	// }
	// log.Print(string(stdout))

}
