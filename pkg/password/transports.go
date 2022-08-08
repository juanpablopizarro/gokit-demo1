package password

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/log"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func RegisterHandlers(logger log.Logger, router *mux.Router) {
	s := NewPasswordService()
	s = loggingMiddleware{logger, s}

	hashHandler := httptransport.NewServer(
		makeHashEndpoint(s),
		decodeHashRequest,
		encodeResponse)

	checkHandler := httptransport.NewServer(
		makeCheckEndpoint(s),
		decodeCheckRequest,
		encodeResponse)

	router.Handle("/hash", hashHandler).Methods("POST")
	router.Handle("/check", checkHandler).Methods("POST")
}

func decodeHashRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req hashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req checkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
