package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func AddToList(key string, tickers []string) error {
	svc, err := startSession()
	if err != nil {
		return err
	}

	// Read existing object if exists
	obj, err := svc.GetObject(&s3.GetObjectInput{}) // TODO object params

	// Doesn't exist
	if err != nil {
		// write new object
	} else {
		// modify wrie
	}

	return nil
}

func RmFromList(key string, tickers []string) error {
	// check for object
	// if exists read modify wrie
	// otherwise nothing
	return nil
}

func GetList(key string) (string, error) {
	// check for object
	// if exists read and return
	return "", nil
}

// FIXME this should not use the hard coded access keys
func startSession() (*s3.S3, error) {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ""),
	})

	if err != nil {
		return nil, error
	}

	return s3.New(sess)
}

func getObject(name string) {
}

func putObject(name string) {

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
