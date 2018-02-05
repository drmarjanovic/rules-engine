package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/MainfluxLabs/rules-engine"
)

func retrieveRuleEndpoint(svc rules.Service) endpoint.Endpoint {
	return func(_ context.Context, body interface{}) (interface{}, error) {
		b := body.(viewRuleReq)

		if err := b.validate(); err != nil {
			return nil, err
		}

		rule, err := svc.ViewRule(b.userId, b.ruleId)
		if err != nil {
			return nil, err
		}

		return viewRuleRes{rule}, nil
	}
}

func retrieveRulesEndpoint(svc rules.Service) endpoint.Endpoint {
	return func(_ context.Context, body interface{}) (interface{}, error) {
		b := body.(listRulesReq)

		if err := b.validate(); err != nil {
			return nil, err
		}

		rulesList, err := svc.ListRules(b.userId)
		if err != nil {
			return nil, err
		}

		return listRulesRes{rulesList, len(rulesList)}, nil
	}
}

func removeRuleEndpoint(svc rules.Service) endpoint.Endpoint {
	return func(_ context.Context, body interface{}) (interface{}, error) {
		b := body.(viewRuleReq)

		if err := b.validate(); err != nil {
			return nil, err
		}

		if err := svc.RemoveRule(b.userId, b.ruleId); err != nil {
			return nil, err
		}

		return removeRes{}, nil
	}
}
