package user

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type validateRequest struct {
	Email    string `json:"email"`
	Password string `json:"pass"`
}

type validateResponse struct {
	Response string `json:"result"`
	Error    string `json:"error,omitempty"`
}

type hashRequest struct {
	Pass string `json:"pass"`
}

type hashResponse struct {
	Hash string    `json:"hash"`
	Time time.Time `json:"time"`
	Err  string    `json:"err,omitempty"`
}

func makeValidateEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(validateRequest)
		_, err := s.Validate(req.Email, req.Password) // we ignore the Validate's return value and just watch if there is an error or no
		if err != nil {
			return validateResponse{"", err.Error()}, err
		}
		return validateResponse{"OK", ""}, nil
	}
}

func makeHashEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(hashRequest)
		h, err := s.Hash(req.Pass)
		if err != nil {
			return hashResponse{"", time.Now(), err.Error()}, err
		}
		return hashResponse{h, time.Now(), ""}, nil
	}
}
