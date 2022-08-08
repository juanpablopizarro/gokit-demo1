package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

type hashProxyMiddleware struct {
	next   Service
	proxy  endpoint.Endpoint
	logger log.Logger
	ctx    context.Context
}

func (mw hashProxyMiddleware) Validate(email, pass string) (*User, error) {
	return mw.next.Validate(email, pass)
}

func (mw hashProxyMiddleware) Hash(pass string) (string, error) {
	hash, err := mw.proxy(mw.ctx, hashRequest{pass})
	if err != nil {
		return "", err
	}

	response := hash.(hashResponse)
	if response.Err != "" {
		return "", errors.New(response.Err)
	}

	// here we can change the response if we want !

	return response.Hash, nil
}

type ServiceMiddleware func(Service) Service

func proxyingMiddleware(proxyURL string, logger log.Logger, ctx context.Context) ServiceMiddleware {
	return func(next Service) Service {
		return hashProxyMiddleware{next, makeHashProxy(proxyURL), logger, ctx}
	}
}

func makeHashProxy(proxyURL string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		parseURL(proxyURL),
		encodeHashProxyRequest,
		decodeHashProxyResponse).Endpoint()
}

func parseURL(proxyURL string) *url.URL {
	u, err := url.Parse(proxyURL)
	if err != nil {
		panic(err)
	}
	return u
}

func encodeHashProxyRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeHashProxyResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response hashResponse

	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
