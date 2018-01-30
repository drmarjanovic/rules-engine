package api

import (
	"net/http"

	rulesEngine "github.com/MainfluxLabs/rules-engine"
	"github.com/go-zoo/bone"
)

func MakeHandler() http.Handler {
	r := bone.New()

	r.GetFunc("/health", rulesEngine.Health())

	return r
}
