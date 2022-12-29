package aws

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (awsSession *AwsSession) UploadToS3Bucket(bucket string, key string, jsonData string) error {
	data := strings.NewReader(jsonData)
	uploader := s3manager.NewUploader(awsSession.Session)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   data,
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Stolen from: https://stackoverflow.com/a/65360986
func (awsSession *AwsSession) S3KeyExists(bucket string, key string) (bool, error) {
	s3Session := s3.New(awsSession.Session)
	_, err := s3Session.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}

func (awsSession *AwsSession) GetFromS3Bucket(bucket string, key string) ([]byte, error) {
	s3Client := s3.New(awsSession.Session)

	requestInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := s3Client.GetObject(requestInput)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer result.Body.Close()
	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return data, nil
}

func (awsSession *AwsSession) GetPresignedUploadUrl(bucket string, key string, expire time.Duration) (string, error) {
	s3Client := s3.New(awsSession.Session)
	putObjReq, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	uploadUrl, err := putObjReq.Presign(expire)
	return uploadUrl, err
}

func (awsSession *AwsSession) GetPresignedDownloadUrl(bucket string, key string, expire time.Duration) (string, error) {
	s3Client := s3.New(awsSession.Session)
	getObjReq, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	downloadUrl, err := getObjReq.Presign(15 * time.Minute)
	return downloadUrl, err
}
