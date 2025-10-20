package bucket

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/tiredkangaroo/ajiteshcc/env"
)

var S3Client *s3.Client

func Init() error {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				env.DefaultEnv.R2_ACCESS_KEY_ID,
				env.DefaultEnv.R2_SECRET_ACCESS_KEY,
				"",
			),
		),
		config.WithRegion("auto"),
	)
	if err != nil {
		return err
	}
	S3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", env.DefaultEnv.R2_ACCOUNT_ID))
	})
	return nil
}

type Object struct {
	Name      string            `json:"name"`       // object key
	Size      int64             `json:"size"`       // size in bytes
	PublicURL string            `json:"public_url"` // public URL of the object
	Metadata  map[string]string `json:"metadata"`   // metadata associated with the object
}

func ListAllObjectsInBucket(ctx context.Context, bucketName string) ([]Object, error) {
	output, err := S3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		return nil, fmt.Errorf("list objects in bucket: %w", err)
	}
	objects := make([]Object, len(output.Contents))
	pubURL := *env.DefaultEnv.R2_PHOTOS_BUCKET_PUBLIC_URL
	for i, obj := range output.Contents {
		md, err := GetObjectMetadata(ctx, bucketName, *obj.Key)
		if err != nil {
			return nil, fmt.Errorf("get object metadata for %s: %w", *obj.Key, err)
		}
		pubURL.Path = "/" + *obj.Key
		objects[i] = Object{
			Name:      *obj.Key,
			Size:      *obj.Size,
			PublicURL: pubURL.String(),
			Metadata:  md,
		}
	}
	return objects, nil
}

func PutObjectInBucket(ctx context.Context, bucketName, objectKey string, metadata map[string]string, body io.Reader) error {
	_, err := S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:   &bucketName,
		Key:      &objectKey,
		Body:     body,
		Metadata: metadata,
	})
	return err
}

func GetObjectMetadata(ctx context.Context, bucketName, objectKey string) (map[string]string, error) {
	headOutput, err := S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		return nil, fmt.Errorf("head object %s: %w", objectKey, err)
	}
	return headOutput.Metadata, nil
}

func UpdateObjectMetadata(ctx context.Context, bucketName, objectKey string, newMetadata map[string]string) error {
	_, err := S3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:            &bucketName,
		Key:               &objectKey,
		CopySource:        aws.String(bucketName + "/" + objectKey),
		Metadata:          newMetadata,
		MetadataDirective: types.MetadataDirectiveReplace,
	})
	return err
}
