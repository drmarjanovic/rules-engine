package mocks

import (
	"sync"
	"strings"
	"fmt"

	"github.com/MainfluxLabs/rules-engine"
)

var _ rules.RuleRepository = (*ruleRepositoryMock)(nil)

type ruleRepositoryMock struct {
	mu    sync.Mutex
	rules map[string]rules.Rule
}

// NewRuleRepository instantiates in-memory rule repository.
func NewRuleRepository() rules.RuleRepository {
	return &ruleRepositoryMock{
		rules: make(map[string]rules.Rule),
	}
}

func (repo *ruleRepositoryMock) Save(rule rules.Rule) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.rules[key(rule.UserId, rule.ID)] = rule

	return nil
}

func (repo *ruleRepositoryMock) One(userId string, ruleId string) (rules.Rule, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if r, ok := repo.rules[key(userId, ruleId)]; ok {
		return r, nil
	}

	return rules.Rule{}, rules.ErrNotFound
}

func (repo *ruleRepositoryMock) All(userId string) []rules.Rule {
	prefix := fmt.Sprintf("%s-", userId)

	rulesList := make([]rules.Rule, 0)

	for k, v := range repo.rules {
		if strings.HasPrefix(k, prefix) {
			rulesList = append(rulesList, v)
		}
	}

	return rulesList
}

func (repo *ruleRepositoryMock) Remove(userId string, ruleId string) error {
	delete(repo.rules, key(userId, ruleId))
	return nil
}

func key(userId, ruleId string) string {
	return fmt.Sprintf("%s-%s", userId, ruleId)
}
