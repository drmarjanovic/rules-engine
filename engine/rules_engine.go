package engine

import (
	"github.com/mainflux/mainflux/writer"
)

var _ Service = (*ruleService)(nil)

type ruleService struct {
	rules RuleRepository
}

// NewService instantiates the domain service implementation.
func NewService(rules RuleRepository) Service {
	return &ruleService{
		rules: rules,
	}
}

func (rs *ruleService) SaveRule(rule Rule) error {
	return rs.rules.Save(rule)
}

func (rs *ruleService) ViewRule(userId string, ruleId string) (*Rule, error) {
	return rs.rules.One(userId, ruleId)
}

func (rs *ruleService) ListRules(userId string) ([]Rule, error) {
	return rs.rules.All(userId), nil
}

func (rs *ruleService) RemoveRule(userId string, ruleId string) error {
	return rs.rules.Remove(userId, ruleId)
}

func (rs *ruleService) ApplyRules(userId string, events []writer.Message) error {
	rls, err := rs.ListRules(userId)
	if err != nil {
		return err
	}

	for _, event := range events {
		for _, rule := range rls {
			if rule.IsMatchedBy(event) {
				for _, action := range rule.Actions {
					action.Execute()
				}
			}
		}
	}
	return nil
}
