package api

import (
	"fmt"
	"net/http"

	"github.com/MainfluxLabs/rules-engine"
)

const contentType = "application/json; charset=utf-8"

type apiRes interface {
	code() int
	headers() map[string]string
	empty() bool
}

type viewRuleRes struct {
	rules.Rule
}

func (res viewRuleRes) code() int {
	return http.StatusOK
}

func (res viewRuleRes) headers() map[string]string {
	return map[string]string{}
}

func (res viewRuleRes) empty() bool {
	return false
}

type listRulesRes struct {
	Rules []rules.Rule `json:"rules"`
	count int
}

func (res listRulesRes) code() int {
	return http.StatusOK
}

func (res listRulesRes) headers() map[string]string {
	return map[string]string{
		"X-Count": fmt.Sprintf("%d", res.count),
	}
}

func (res listRulesRes) empty() bool {
	return false
}

type removeRes struct{}

func (res removeRes) code() int {
	return http.StatusNoContent
}

func (res removeRes) headers() map[string]string {
	return map[string]string{}
}

func (res removeRes) empty() bool {
	return false
}
