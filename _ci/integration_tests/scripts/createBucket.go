package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"log"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getConfig(bucketAddress string) (*aws.Config, error) {
	var cfg aws.Config
	var err error

	if strings.HasPrefix(bucketAddress, "s3://") {
		cfg, err = config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	} else if strings.HasPrefix(bucketAddress, "https://") || strings.HasPrefix(bucketAddress, "http://") {
		// e.g. http://127.0.0.1:9000/test becomes http://127.0.0.1:9000
		server, _ := path.Split(bucketAddress)
		// https://stackoverflow.com/questions/67575681/is-aws-go-sdk-v2-integrated-with-local-minio-server
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               server,
				SigningRegion:     "us-east-1",
				HostnameImmutable: true,
			}, nil
		})

		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion("us-east-1"),
			config.WithEndpointResolverWithOptions(resolver),
		)
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return nil, errors.New("bucketAddress did not contain s3://, http://, or https:// prefix")
}

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

func upload(ctx context.Context, serverName, bucketName, key string, b []byte) error {
	cfg, err := getConfig(serverName)
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(s3.NewFromConfig(*cfg))
	uploader.Concurrency = 3

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
	cfg, err := getConfig(bucketAddress)
	if err != nil {
		return nil, err
	}
	return &s3Client{
		Client: s3.NewFromConfig(*cfg),
	}, nil
}

func main() {
	ctx := context.Background()

	serverAddress := flag.String("bucket_address", "http://127.0.0.1:9000/test", "address to bucket")
	s3c, err := NewClient(*serverAddress)
	if err != nil {
		log.Fatal(err)
	}
	_, bucket := path.Split(*serverAddress)
	if err := s3c.createBucket(ctx, bucket); err != nil {
		if !strings.Contains(err.Error(), "BucketAlreadyOwnedByYou") {
			log.Fatal(err)
		}
	}
	b := []byte("blah")
	if err := upload(ctx, *serverAddress, bucket, "thing", b); err != nil {
		log.Fatal(err)
	}

}
