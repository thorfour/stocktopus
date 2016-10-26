package aws

import (
	"bytes"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	bucketName = "stocktopus"
)

func AddToList(key string, tickers []string) error {

	// Read existing object if exists
	obj, sess, err := getObject(key, nil)

	if err != nil { // Doesn't exist
		obj = []byte(strings.Join(tickers, " "))
	} else { // Exists, add tickers to list
		obj = append(obj, []byte(strings.Join(tickers, " "))...)
	}

	// Write object
	_, err = putObject(key, obj, sess)

	return err
}

func RmFromList(key string, tickers []string) error {
	//--------
	// FIXME
	//--------

	// check for object
	_, _, err := getObject(key, nil)

	if err == nil {
		// if exists read modify wrie
		// FIXME
	}

	return err
}

func GetList(key string) (string, error) {

	obj, _, err := getObject(key, nil)
	if err != nil {
		return "", nil
	}

	return string(obj), nil
}

// FIXME this should not use the hard coded access keys
func startSession() (*s3.S3, error) {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ""),
	})

	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}

func getObject(name string, svc *s3.S3) ([]byte, *s3.S3, error) {

	var err error

	// If there is no active session start one
	if svc == nil {
		svc, err = startSession()
		if err != nil {
			return nil, nil, err
		}
	}

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(name),
	}

	resp, err := svc.GetObject(params)
	if err != nil {
		return nil, svc, err
	}

	// Read the payload into slice
	obj := make([]byte, *resp.ContentLength)
	_, err = resp.Body.Read(obj)
	return obj, svc, err
}

func putObject(name string, data []byte, svc *s3.S3) (*s3.S3, error) {

	var err error
	if svc == nil {
		svc, err = startSession()
		if err != nil {
			return nil, err
		}
	}

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(name),
		Body:   bytes.NewReader(data),
		Metadata: map[string]*string{
			"Key": aws.String(name),
		},
	}
	_, err = svc.PutObject(params)

	return svc, err
}

//-----------------------------------
// DEBUG FUNCTIONS
//-----------------------------------

// List of all buckets
func DebugListBuckets() (string, error) {

	svc, err := startSession()
	if err != nil {
		return "", err
	}

	var params *s3.ListBucketsInput
	resp, err := svc.ListBuckets(params)
	if err != nil {
		return "", err
	}

	return *resp.Buckets[0].Name, err
}
