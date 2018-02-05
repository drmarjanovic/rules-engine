// Package rules contains domain definitions for specifying alarming rules
// and triggers in Mainflux.
package rules

import (
	"encoding/json"
	"net/http"
)

const version string = "0.1.0-rc.1"

type response struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Health exposes an HTTP handler function for retrieving service name and version.
func Health() http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		res := response{Name: "rules-engine", Version: version}

		data, _ := json.Marshal(res)

		rw.Write(data)
	})
}
