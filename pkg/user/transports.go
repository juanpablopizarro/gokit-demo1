package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"
)

func RegisterHandlers(logger log.Logger, router *mux.Router) {

	srv := NewServiceInstance()
	ctx := context.Background()
	srv = proxyingAdvMiddleware(ctx, "http://localhost:8081/hash", logger)(srv)
	//srv = proxyingMiddleware("http://localhost:8081/hash", logger, ctx)(srv)
	srv = loggingMiddleware{logger, srv}

	validateHandler := httptransport.NewServer(
		makeValidateEndpoint(srv),
		decodeValidateRequest,
		encodeResponse)

	hashHandler := httptransport.NewServer(
		makeHashEndpoint(srv),
		decodeHashRequest,
		encodeResponse)

	router.Handle("/validate", validateHandler).Methods("POST")
	router.Handle("/hash", hashHandler).Methods("POST")
}

func decodeValidateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req validateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHashRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req hashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
