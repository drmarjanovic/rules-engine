package tests

import (
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/MainfluxLabs/rules-engine/engine/mocks"
	"github.com/MainfluxLabs/rules-engine/engine"
)

var (
	rulesRepo engine.RuleRepository = mocks.NewRuleRepository()
	svc       engine.Service        = engine.NewService(rulesRepo)
)

func TestViewRule(t *testing.T) {
	existingRule := engine.Rule{"1", "1", "test-rule-1", make([]engine.Condition, 0), make([]engine.Action, 0)}
	rulesRepo.Save(existingRule)

	cases := []struct {
		userId string
		ruleId string
		rule   *engine.Rule
		err    error
	}{
		{"1", "1", &existingRule, nil},
		{"1", "2", &engine.Rule{}, engine.ErrNotFound},
	}

	for i, tc := range cases {
		r, err := svc.ViewRule(tc.userId, tc.ruleId)
		assert.Equal(t, tc.err, err, fmt.Sprintf("failed at %d\n", i))
		assert.Equal(t, tc.rule, r, fmt.Sprintf("failed at %d\n", i))
	}
}

func TestListRules(t *testing.T) {
	r1 := engine.Rule{"1", "2", "test-rule-1", make([]engine.Condition, 0), make([]engine.Action, 0)}
	r2 := engine.Rule{"2", "2", "test-rule-2", make([]engine.Condition, 0), make([]engine.Action, 0)}
	rulesRepo.Save(r1)
	rulesRepo.Save(r2)

	cases := []struct {
		userId string
		rules  []engine.Rule
		err    error
	}{
		{"2", []engine.Rule{r1, r2}, nil},
		{"3", make([]engine.Rule, 0), nil},
	}

	for i, tc := range cases {
		r, err := svc.ListRules(tc.userId)
		assert.Equal(t, tc.err, err, fmt.Sprintf("failed at %d\n", i))
		assert.ElementsMatch(t, tc.rules, r, fmt.Sprintf("failed at %d\n", i))
	}
}

func TestRemoveRule(t *testing.T) {
	cases := []struct {
		userId string
		ruleId string
		err    error
	}{
		{"1", "1", nil},
		{"3", "2", nil},
	}

	for i, tc := range cases {
		err := svc.RemoveRule(tc.userId, tc.ruleId)
		assert.Equal(t, tc.err, err, fmt.Sprintf("failed at %d\n", i))
	}
}
