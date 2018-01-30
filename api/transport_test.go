package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	rulesEngine "github.com/MainfluxLabs/rules-engine"
)

func TestHealth(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(rulesEngine.Health())
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK, "bad status code")
}
