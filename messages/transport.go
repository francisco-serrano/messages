package messages

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"net/http"
)

func MakeInsertEndpoint(svc MessageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(insertRequest)
		if err := svc.Create(req.Content); err != nil {
			return insertResponse{Result: "ERR_NOT_CREATED", Error: err.Error()}, nil
		}

		return insertResponse{Result: "MSG_CREATED", Error: ""}, nil
	}
}

func DecodeInsertRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request insertRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

type insertRequest struct {
	Content string `json:"content"`
}

type insertResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

func MakeGetAllEndpoint(svc MessageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		cars, err := svc.GetAll()
		if err != nil {
			return nil, err
		}

		return getAllResponse{Messages: cars}, nil
	}
}

func DecodeGetAllRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return &getAllRequest{}, nil
}

type getAllRequest struct {
}

type getAllResponse struct {
	Messages []Message `json:"messages"`
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err, ok := response.(error); ok && err.Error() != "" {
		encodeError(ctx, err, w)
		return nil
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
		break
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
