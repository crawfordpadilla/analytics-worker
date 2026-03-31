package analytics_worker

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type HttpRequest struct {
	HTTPMethod string
	Endpoint   string
	Headers    map[string]string
	Body       string
}

type SNSMessage struct {
	DeviceID string
	Event    string
	Product  string
	Timestamp int64
}

func getSignatureKey(key, datestamp, region, service string) []byte {
	kDate := hmac.NewSHA256([]byte(datestamp))
	kDate.Write([]byte(region))
	kDate.Write([]byte(service))
	kSigning := hmac.NewSHA256([]byte(key))
	kSigning.Write(kDate.Sum(nil))
	return kSigning.Sum(nil)
}

func getHeadersSignature(method, path, datestamp, region, service, signingKey string) string {
	signature := base64.StdEncoding.EncodeToString(hmac.NewSHA256(getSignatureKey(signingKey, datestamp, region, service)).Sum(nil))
	return fmt.Sprintf("%s %s:%s,%s date:%s,host:%s,signature:%s,signaturekey:%s", method, path, datestamp, region, service, signature, signingKey)
}

func getAwsSnsClient() *sns.SNS {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))}, nil)
	if err != nil {
		log.Fatal(err)
	}
	return sns.New(sess)
}

func sendAnalyticsEvent(deviceID, event, product string, timestamp int64) error {
	snsClient := getAwsSnsClient()
	snsData := SNSMessage{
		DeviceID: deviceID,
		Event:    event,
		Product:  product,
		Timestamp: timestamp,
	}

	jsonData, err := json.Marshal(snsData)
	if err != nil {
		return err
	}

	params := &sns.PublishInput{
		TopicArn: aws.String(os.Getenv("SNS_TOPIC_ARN")),
		Message:  aws.String(string(jsonData)),
	}

	_, err = snsClient.Publish(params)
	return err
}

func sendHttpRequest(req *HttpRequest) (*http.Response, error) {
	// Create a new HTTP request
	reqBody := &bytes.Buffer{}
	if req.Body != "" {
		_, _ = reqBody.Write([]byte(req.Body))
	}

	reqHeaders := make(http.Header)
	for k, v := range req.Headers {
		reqHeaders.Set(k, v)
	}

	httpReq, err := http.NewRequest(req.HTTPMethod, req.Endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	for k, v := range reqHeaders {
		httpReq.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	return resp, err
}

func sendHttpRequestWithRetry(req *HttpRequest, maxAttempts int) (*http.Response, error) {
	for attempt := 0; attempt < maxAttempts; attempt++ {
		resp, err := sendHttpRequest(req)
		if err == nil {
			return resp, nil
		}
		log.Printf("Attempt %d failed: %v", attempt, err)
		time.Sleep(time.Second * 2)
	}
	return nil, fmt.Errorf("max attempts reached")
}

func parseQueryParams(req *HttpRequest) map[string]string {
	u, err := url.Parse(req.Endpoint)
	if err != nil {
		log.Fatal(err)
	}
	query := u.Query()
	params := make(map[string]string)
	for k, v := range query {
		params[k] = strings.Join(v, ",")
	}
	return params
}