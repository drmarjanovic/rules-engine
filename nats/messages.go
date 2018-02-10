package nats

import (
	"github.com/gocql/gocql"
	"github.com/MainfluxLabs/rules-engine"
)

const (
	sendEmail string = "SEND EMAIL"
	turnOff   string = "TURN OFF"

	name      string = "name"
	content   string = "content"
	recipient string = "recipient"
	from      string = "from"
	to        string = "to"
	deviceId  string = "deviceId"
)

type RawRules struct {
	Data []RawRule `json:"rules"`
}

type RawRule struct {
	UserId     *string                  `json:"userId"`
	Name       string                   `json:"name"`
	Conditions []RawCondition           `json:"conditions"`
	Actions    []map[string]interface{} `json:"actions"`
}

type RawCondition struct {
	DeviceID *string     `json:"deviceId"`
	Property *string     `json:"property"`
	Operator *string     `json:"operator"`
	Value    interface{} `json:"value"`
}

func (msg *RawRules) toRules() ([]rules.Rule, error) {
	var rls []rules.Rule

	for _, r := range msg.Data {
		rule, err := parseRule(r)
		if err != nil {
			return nil, err
		}
		rls = append(rls, *rule)
	}
	return rls, nil
}

func parseRule(rawRule RawRule) (*rules.Rule, error) {
	var (
		actions    []rules.Action
		conditions []rules.Condition
	)
	rule := &rules.Rule{
		ID:   gocql.TimeUUID().String(),
		Name: rawRule.Name,
	}

	if rawRule.UserId == nil {
		return nil, rules.ErrMalformedEntity
	}
	rule.UserId = *rawRule.UserId

	for _, act := range rawRule.Actions {
		action, err := parseAction(act)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}

	for _, cnd := range rawRule.Conditions {
		condition, err := parseCondition(cnd)
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, condition)
	}

	rule.Actions = actions
	rule.Conditions = conditions
	return rule, nil
}

func parseAction(action map[string]interface{}) (rules.Action, error) {
	if name, ok := action[name]; ok {
		switch name {
		case sendEmail:
			return parseSendEmail(action)
		case turnOff:
			return parseTurnOff(action)
		}
	}

	return nil, rules.ErrMalformedEntity
}

func parseSendEmail(action map[string]interface{}) (*rules.SendEmailAction, error) {
	se := &rules.SendEmailAction{Name: sendEmail}

	c, err := parseStringProperty(action, content)
	if err != nil {
		return nil, err
	}
	se.Content = *c

	r, err := parseStringProperty(action, recipient)
	if err != nil {
		return nil, err
	}
	se.Recipient = *r

	return se, nil
}

func parseTurnOff(action map[string]interface{}) (*rules.TurnOffAction, error) {
	t := &rules.TurnOffAction{Name: turnOff}

	id, err := parseStringProperty(action, deviceId)
	if err != nil {
		return nil, err
	}
	t.DeviceId = *id

	return t, nil
}

func parseStringProperty(object map[string]interface{}, property string) (*string, error) {
	prop, ok := object[property]
	if !ok {
		return nil, rules.ErrMalformedEntity
	}

	strProp, ok := prop.(string)
	if !ok {
		return nil, rules.ErrMalformedEntity
	}

	return &strProp, nil
}

func parseCondition(rawCnd RawCondition) (rules.Condition, error) {
	cndData, err := parseConditionData(rawCnd)
	if err != nil {
		return nil, err
	}

	switch rawCnd.Value.(type) {
	case bool:
		return rules.BooleanCondition{ConditionData: cndData, Value: rawCnd.Value.(bool)}, nil
	case float64:
		return rules.NumericCondition{ConditionData: cndData, Value: rawCnd.Value.(float64)}, nil
	case string:
		return rules.StringCondition{ConditionData: cndData, Value: rawCnd.Value.(string)}, nil
	case map[string]interface{}:
		return parseBetweenCondition(rawCnd.Value, cndData)
	}

	return nil, rules.ErrMalformedEntity
}

func parseConditionData(rawCnd RawCondition) (*rules.ConditionData, error) {
	cndData := &rules.ConditionData{
		DeviceID: rawCnd.DeviceID,
		Property: rawCnd.Property,
		Operator: rawCnd.Operator,
	}

	if err := cndData.Validate(); err != nil {
		return nil, err
	}
	return cndData, nil
}

func parseBetweenCondition(value interface{}, data *rules.ConditionData) (*rules.BetweenCondition, error) {
	bounds, err := convertBounds(value)
	if err != nil {
		return nil, rules.ErrMalformedEntity
	}

	cnd := &rules.BetweenCondition{
		ConditionData: data,
		From:          bounds[from],
		To:            bounds[to],
	}

	if err := cnd.Validate(); err != nil {
		return nil, err
	}
	return cnd, nil
}

func convertBounds(originalMap interface{}) (map[string]*float64, error) {
	convertedMap := map[string]*float64{}
	for key, value := range originalMap.(map[string]interface{}) {
		v, ok := value.(float64)
		if !ok {
			return nil, rules.ErrMalformedEntity
		}
		convertedMap[key] = &v
	}
	return convertedMap, nil
}
