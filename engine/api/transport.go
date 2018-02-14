package api

import (
	"net/http"
	"context"
	"encoding/json"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/MainfluxLabs/rules-engine/engine"
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc engine.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Get("/users/:userId/rules", kithttp.NewServer(
		retrieveRulesEndpoint(svc),
		decodeList,
		encodeResponse,
		opts...,
	))

	r.Get("/users/:userId/rules/:ruleId", kithttp.NewServer(
		retrieveRuleEndpoint(svc),
		decodeView,
		encodeResponse,
		opts...,
	))

	r.Delete("/users/:userId/rules/:ruleId", kithttp.NewServer(
		removeRuleEndpoint(svc),
		decodeView,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/health", engine.Health())

	return r
}

func decodeView(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewRuleReq{
		userId: bone.GetValue(r, "userId"),
		ruleId: bone.GetValue(r, "ruleId"),
	}

	return req, nil
}

func decodeList(_ context.Context, r *http.Request) (interface{}, error) {
	req := listRulesReq{
		userId: bone.GetValue(r, "userId"),
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(apiRes); ok {
		for k, v := range ar.headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(ar.code())

		if ar.empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case engine.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case engine.ErrMalformedUrl:
		w.WriteHeader(http.StatusBadRequest)
	case engine.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		if _, ok := err.(*json.SyntaxError); ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	}
}
