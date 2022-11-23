package httpclient

import "github.com/imroc/req/v3"

func NewHttpClient() *req.Client {
	client := req.C().
		SetBaseURL("http://localhost:8080").
		SetCommonRetryCount(3)

	return client
}
