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

// Clear Delete the file entirely
func Clear(key string) error {

	svc, err := startSession()
	if err != nil {
		return err
	}

	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	_, err = svc.DeleteObject(params)

	return err
}

// AddToList adds an item to a watchlist file
func AddToList(key string, tickers []string) error {

	// Read existing object if exists
	obj, sess, err := getObject(key, nil)

	if err != nil { // Doesn't exist
		obj = []byte(strings.Join(tickers, " "))
	} else { // Exists, add tickers to list
		if len(obj) != 0 {
			obj = append(obj, []byte(" ")...)
		}
		obj = append(obj, []byte(strings.Join(tickers, " "))...)
	}

	// Write object
	_, err = putObject(key, obj, sess)

	return err
}

// RmFromList removes an item from a watchlist file
func RmFromList(key string, tickers []string) error {

	// check for object
	obj, sess, err := getObject(key, nil)
	if err == nil { // if exists read modify wrie
		list := strings.Split(string(obj), " ")
		for i := range list {
			if list[i] == tickers[0] {
				list = append(list[:i], list[i+1:]...) // Remove ticker from list
				_, err := putObject(key, []byte(strings.Join(list, " ")), sess)
				return err
			}
		}
	}

	// Object didn't exist or ticker didn't exist
	return nil
}

// GetList Returns a watchlist in a file
func GetList(key string) (string, error) {

	obj, _, err := getObject(key, nil)
	if err != nil {
		return "", err
	}

	return string(obj), nil
}

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
	resp.Body.Read(obj)
	return obj, svc, nil
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
func debugListBuckets() (string, error) {

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
