package rules

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

func (rs *ruleService) ViewRule(userId string, ruleId string) (Rule, error) {
	return rs.rules.One(userId, ruleId)
}

func (rs *ruleService) ListRules(userId string) ([]Rule, error) {
	return rs.rules.All(userId), nil
}

func (rs *ruleService) RemoveRule(userId string, ruleId string) error {
	return rs.rules.Remove(userId, ruleId)
}
