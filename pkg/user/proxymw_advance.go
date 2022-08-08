package user

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"

	"github.com/sony/gobreaker"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

func proxyingAdvMiddleware(ctx context.Context, instances string, logger log.Logger) ServiceMiddleware {
	// If instances is empty, don't proxy.
	if instances == "" {
		logger.Log("proxy_to", "none")
		return func(next Service) Service { return next }
	}

	// Set some parameters for our client.
	var (
		qps         = 100                     // beyond which we will return an error
		maxAttempts = 3                       // per request, before giving up
		maxTime     = 2000 * time.Millisecond // wallclock time, before giving up
	)

	// Otherwise, construct an endpoint for each instance in the list, and add
	// it to a fixed set of endpoints. In a real service, rather than doing this
	// by hand, you'd probably use package sd's support for your service
	// discovery system.
	var (
		instanceList = split(instances)
		endpointer   sd.FixedEndpointer
	)

	logger.Log("proxy_to", fmt.Sprint(instanceList))
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeHashProxyAdv(ctx, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}

	// Now, build a single, retrying, load-balancing endpoint out of all of
	// those individual endpoints.
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	// And finally, return the ServiceMiddleware, implemented by proxymw.
	return func(next Service) Service {
		return proxymwAdv{ctx, next, retry}
	}
}

type proxymwAdv struct {
	ctx   context.Context
	next  Service
	proxy endpoint.Endpoint
}

func (mw proxymwAdv) Validate(email, pass string) (*User, error) {
	return mw.next.Validate(email, pass)
}

func (mw proxymwAdv) Hash(s string) (string, error) {
	response, err := mw.proxy(mw.ctx, hashRequest{s})
	if err != nil {
		return "", err
	}

	resp := response.(hashResponse)
	if resp.Err != "" {
		return resp.Hash, errors.New(resp.Err)
	}
	return resp.Hash, nil
}

func makeHashProxyAdv(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}

	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}

	if u.Path == "" {
		u.Path = "/hash"
	}

	return httptransport.NewClient(
		"POST",
		u,
		encodeHashProxyRequest,
		decodeHashProxyResponse,
	).Endpoint()
}

func split(s string) []string {
	a := strings.Split(s, ",")
	for i := range a {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}
