package ts3

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestRunTestingS3(t *testing.T) {
	uploadedContent := new(bytes.Buffer)
	sess, cleanup := TestS3(t, uploadedContent)
	defer cleanup()

	s := s3.New(sess)
	_, err := s.PutObject(&s3.PutObjectInput{
		Bucket:               TestingBucket,
		Key:                  aws.String("key"),
		Body:                 aws.ReadSeekCloser(strings.NewReader("payload")),
		ACL:                  aws.String("authenticated-read"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		t.Errorf("%v", err)
	}
	if uploadedContent.String() != "payload" {
		t.Errorf(`should be "payload" != %s`, uploadedContent.String())
	}
}
