package cassandra

import (
	"encoding/json"

	"github.com/MainfluxLabs/rules-engine"
	"github.com/gocql/gocql"
)

var _ rules.RuleRepository = (*ruleRepository)(nil)

type ruleRepository struct {
	session *gocql.Session
}

// NewRuleRepository instantiates Cassandra rule repository.
func NewRuleRepository(session *gocql.Session) rules.RuleRepository {
	return &ruleRepository{session}
}

func (repo *ruleRepository) Save(rule rules.Rule) error {
	cql := `INSERT INTO rules (id, user_id, name, conditions, actions) VALUES (?, ?, ?, ?, ?)`

	actions, _ := json.Marshal(rule.Actions)
	conditions, _ := json.Marshal(rule.Conditions)

	if err := repo.session.Query(cql, rule.ID, rule.UserId, rule.Name, conditions, actions).Exec(); err != nil {
		return err
	}

	return nil
}

func (repo *ruleRepository) One(userId string, ruleId string) (rules.Rule, error) {
	cql := `SELECT name, conditions, actions FROM rules WHERE user_id = ? AND id = ? LIMIT 1`

	rule := rules.Rule{
		ID:     ruleId,
		UserId: userId,
	}

	var conditions []byte
	var actions []byte

	if err := repo.session.Query(cql, userId, ruleId).
		Scan(&rule.Name, &conditions, &actions); err != nil {
		return rule, rules.ErrNotFound
	}

	c := make([]rules.Condition, 0)
	a := make([]rules.Action, 0)

	if err := json.Unmarshal(conditions, &c); err != nil {
		return rule, err
	}

	if err := json.Unmarshal(actions, &a); err != nil {
		return rule, err
	}

	rule.Conditions = c
	rule.Actions = a

	return rule, nil
}

func (repo *ruleRepository) All(userId string) []rules.Rule {
	cql := `SELECT id, name, conditions, actions FROM rules WHERE user_id = ?`
	var (
		id, name            string
		conditions, actions []byte
	)

	iter := repo.session.Query(cql, userId).Iter()
	defer iter.Close()

	rulesList := make([]rules.Rule, 0)

	for iter.Scan(&id, &name, &conditions, &actions) {
		c := make([]rules.Condition, 0)
		a := make([]rules.Action, 0)

		if json.Unmarshal(conditions, &c) != nil {
			return rulesList
		}

		if json.Unmarshal(actions, &a) != nil {
			return rulesList
		}

		r := rules.Rule{
			ID:         id,
			UserId:     userId,
			Name:       name,
			Conditions: c,
			Actions:    a,
		}

		rulesList = append(rulesList, r)
	}

	return rulesList
}

func (repo *ruleRepository) Remove(userId string, ruleId string) error {
	cql := `DELETE FROM rules WHERE user_id = ? AND id = ?`
	return repo.session.Query(cql, userId, ruleId).Exec()
}
