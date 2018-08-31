package ts3

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
)

// TestingBucket defines the bucket name
var TestingBucket = "127.0.0.1" // The bucket name must be 127.0.0.1

// TestS3 returns session and cleanup function.
func TestS3(t *testing.T, dst io.Writer) (*session.Session, func()) {
	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		for k, v := range r.Header {
			t.Logf("%s %s", k, v)
		}
		io.Copy(dst, r.Body)
	}))
	cleanup := func() {
		ts3.Close()
	}

	s3EndpointFn := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		u, err := url.Parse(ts3.URL)
		if err != nil {
			return endpoints.ResolvedEndpoint{}, err
		}
		return endpoints.ResolvedEndpoint{
			URL: "http://:" /* 127.0.0.1 */ + u.Port(),
		}, nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(endpoints.ApNortheast1RegionID),
		EndpointResolver: endpoints.ResolverFunc(s3EndpointFn),
	})
	if err != nil {
		t.Fatal(err)
	}

	return sess, cleanup
}
