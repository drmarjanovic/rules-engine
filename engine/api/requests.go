package api

import (
	"github.com/asaskevich/govalidator"
	"github.com/MainfluxLabs/rules-engine/engine"
)

type apiReq interface {
	validate() error
}

type viewRuleReq struct {
	userId string
	ruleId string
}

func (req viewRuleReq) validate() error {
	if !govalidator.IsUUID(req.userId) || !govalidator.IsUUID(req.ruleId) {
		return engine.ErrMalformedUrl
	}

	return nil
}

type listRulesReq struct {
	userId string
}

func (req listRulesReq) validate() error {
	if !govalidator.IsUUID(req.userId) {
		return engine.ErrMalformedUrl
	}

	return nil
}
