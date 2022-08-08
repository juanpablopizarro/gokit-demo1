package password

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type hashRequest struct {
	Pass string `json:"pass"`
}

type hashResponse struct {
	Hash string `json:"hash"`
	Err  string `json:"err"`
}

type checkRequest struct {
	Pass string `json:"pass"`
	Hash string `json:"hash"`
}

type checkResponse struct {
	Check string `json:"check"`
	Err   string `json:"err"`
}

func makeHashEndpoint(ps Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(hashRequest)
		h, err := ps.Hash(req.Pass)
		if err != nil {
			return hashResponse{"", err.Error()}, err
		}
		return hashResponse{h, ""}, nil
	}
}

func makeCheckEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(checkRequest)
		ok := s.Check(req.Pass, req.Hash)
		if ok {
			return checkResponse{"ok", ""}, nil
		}
		return checkResponse{"fail", ""}, nil
	}
}
