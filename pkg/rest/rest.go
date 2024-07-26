package rest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/wazadio/realtime-weather/pkg/logger"
)

const ERROR_THIRD_PARTY = "third party error"

type RestRequest struct {
	Scheme  string
	BaseUrl string
	Enpoint string
	Headers map[string][]string
	Method  string
	Params  map[string]string
	Body    *[]byte
}

type RestResponse struct {
	Status  int
	Headers map[string][]string
	Body    []byte
}

type rest struct {
}

type Rest interface {
	Call(ctx context.Context, req RestRequest) (res RestResponse, err error)
}

func NewRest() Rest {
	return &rest{}
}

func (r rest) Call(ctx context.Context, req RestRequest) (res RestResponse, err error) {
	reqContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	scheme := "https"
	if req.Scheme != "" {
		scheme = req.Scheme
	}

	params := url.Values{}
	for key, val := range req.Params {
		params.Add(key, val)
	}

	fullUrl := fmt.Sprintf("%s://%s/%s?%s", scheme, req.BaseUrl, req.Enpoint, params.Encode())
	client := &http.Client{}
	payload := bytes.NewBuffer([]byte{})
	if req.Body != nil {
		payload = bytes.NewBuffer(*req.Body)
	}

	newRequest, err := http.NewRequestWithContext(reqContext, req.Method, fullUrl, payload)
	if err != nil {
		return
	}

	for key, values := range req.Headers {
		for _, val := range values {
			newRequest.Header.Add(key, val)
		}
	}

	defer func() {
		logger.Print(ctx, logger.INFO, map[string]any{
			"request":       fmt.Sprintf("%+v", newRequest),
			"response":      fmt.Sprintf("%+v", res),
			"response body": string(res.Body),
		})
	}()

	resp, err := client.Do(newRequest)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	res.Status = resp.StatusCode
	res.Headers = make(map[string][]string)
	for key, values := range resp.Header {
		res.Headers[key] = values
	}

	res.Body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}
