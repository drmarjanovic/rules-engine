package cassandra

import (
	"encoding/json"

	"github.com/gocql/gocql"
	"github.com/MainfluxLabs/rules-engine/engine"
)

var _ engine.RuleRepository = (*ruleRepository)(nil)

type ruleRepository struct {
	session *gocql.Session
}

// NewRuleRepository instantiates Cassandra rule repository.
func NewRuleRepository(session *gocql.Session) engine.RuleRepository {
	return &ruleRepository{session}
}

func (repo *ruleRepository) Save(rule engine.Rule) error {
	cql := `INSERT INTO rules (id, user_id, name, conditions, actions) VALUES (?, ?, ?, ?, ?)`

	actions, _ := json.Marshal(fromDomain(rule.Actions))
	conditions, _ := json.Marshal(rule.Conditions)

	if err := repo.session.Query(cql, rule.ID, rule.UserId, rule.Name, conditions, actions).Exec(); err != nil {
		return err
	}

	return nil
}

func (repo *ruleRepository) One(userId string, ruleId string) (*engine.Rule, error) {
	cql := `SELECT name, conditions, actions FROM rules WHERE user_id = ? AND id = ? LIMIT 1`
	var (
		conditions, actions []byte
		dbActions           []dbAction
	)

	r := &engine.Rule{
		ID:     ruleId,
		UserId: userId,
	}

	if err := repo.session.Query(cql, userId, ruleId).Scan(&r.Name, &conditions, &actions); err != nil {
		return nil, engine.ErrNotFound
	}

	if err := json.Unmarshal(conditions, &r.Conditions); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(actions, &dbActions); err != nil {
		return nil, err
	}

	acts, err := toDomain(dbActions)
	if err != nil {
		return r, err
	}
	r.Actions = acts

	return r, nil
}

func (repo *ruleRepository) All(userId string) []engine.Rule {
	cql := `SELECT id, name, conditions, actions FROM rules WHERE user_id = ?`
	var (
		id, name            string
		conditions, actions []byte
		dbActions           []dbAction
	)

	iter := repo.session.Query(cql, userId).Iter()
	defer iter.Close()

	rulesList := make([]engine.Rule, 0)

	for iter.Scan(&id, &name, &conditions, &actions) {
		r := engine.Rule{
			ID:     id,
			UserId: userId,
			Name:   name,
		}

		if json.Unmarshal(conditions, &r.Conditions) != nil {
			return rulesList
		}

		if json.Unmarshal(actions, &dbActions) != nil {
			return rulesList
		}

		acts, err := toDomain(dbActions)
		if err != nil {
			return rulesList
		}
		r.Actions = acts

		rulesList = append(rulesList, r)
	}

	return rulesList
}

func (repo *ruleRepository) Remove(userId string, ruleId string) error {
	cql := `DELETE FROM rules WHERE user_id = ? AND id = ?`
	return repo.session.Query(cql, userId, ruleId).Exec()
}
