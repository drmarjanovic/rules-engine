package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/MainfluxLabs/rules-engine/engine"
)

func TestHealth(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(engine.Health())
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK, "bad status code")
}
