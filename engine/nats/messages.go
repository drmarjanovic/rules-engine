package nats

import (
	"github.com/gocql/gocql"
	"github.com/MainfluxLabs/rules-engine/engine"
	"github.com/asaskevich/govalidator"
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

type rulesMsg struct {
	Data []rule `json:"rules"`
}

type rule struct {
	UserId     string      `json:"userId"`
	Name       string      `json:"name"`
	Conditions []condition `json:"conditions"`
	Actions    []action    `json:"actions"`
}

type condition struct {
	DeviceID string         `json:"deviceId"`
	Property string         `json:"property"`
	Operator engine.Operator `json:"operator"`
	Value    interface{}    `json:"value"`
}

type bounds struct {
	from float64
	to   float64
}

type action map[string]interface{}

func (msg rulesMsg) validate() error {
	for _, r := range msg.Data {
		if err := r.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (msg rulesMsg) toDomain() ([]engine.Rule, error) {
	var rls []engine.Rule

	if err := msg.validate(); err != nil {
		return rls, err
	}

	for _, r := range msg.Data {
		rls = append(rls, r.toDomain())
	}

	return rls, nil
}

func (r rule) validate() error {
	if !govalidator.IsUUID(r.UserId) || len(r.Actions) == 0 || len(r.Conditions) == 0 {
		return engine.ErrMalformedEntity
	}

	for _, c := range r.Conditions {
		if err := c.validate(); err != nil {
			return err
		}
	}

	for _, a := range r.Actions {
		if err := a.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (r rule) toDomain() (engine.Rule) {
	var (
		actions    []engine.Action
		conditions []engine.Condition
	)

	rule := engine.Rule{
		ID:     gocql.TimeUUID().String(),
		Name:   r.Name,
		UserId: r.UserId,
	}

	for _, a := range r.Actions {
		actions = append(actions, a.toDomain())
	}

	for _, c := range r.Conditions {
		conditions = append(conditions, c.toDomain())
	}

	rule.Actions = actions
	rule.Conditions = conditions

	return rule
}

func (c condition) validate() error {
	if !govalidator.IsUUID(c.DeviceID) || c.Property == "" || c.Operator == engine.Undefined || validateValue(c.Value, c.Operator) != nil {
		return engine.ErrMalformedEntity
	}

	return nil
}

func validateValue(v interface{}, op engine.Operator) error {
	switch v.(type) {
	case bool, string:
		if op != engine.Eq && op != engine.Neq {
			return engine.ErrMalformedEntity
		}
	case float64:
		if op == engine.Btw {
			return engine.ErrMalformedEntity
		}
	case map[string]interface{}:
		bounds, err := convertBounds(v.(map[string]interface{}))
		if err != nil {
			return err
		}
		if bounds.from >= bounds.to {
			return engine.ErrMalformedEntity
		}
	default:
		return engine.ErrMalformedEntity
	}

	return nil
}

func (c condition) toDomain() engine.Condition {
	cnd := engine.Condition{
		DeviceID: c.DeviceID,
		Property: c.Property,
		Operator: c.Operator,
		Value:    c.Value,
	}

	switch c.Value.(type) {
	case bool:
		cnd.Type = engine.Bool
	case float64:
		cnd.Type = engine.Numeric
	case string:
		cnd.Type = engine.String
	case map[string]interface{}:
		cnd.Type = engine.Between
		bounds, _ := convertBounds(c.Value.(map[string]interface{}))
		cnd.Value = engine.Range{
			From: bounds.from,
			To:   bounds.to,
		}
	}

	return cnd
}

func (a action) validate() error {
	name, ok := a[name]

	if !ok {
		return engine.ErrMalformedEntity
	}

	switch name {
	case sendEmail:
		if _, err := requireStrProp(a, content); err != nil {
			return err
		}
		if _, err := requireStrProp(a, recipient); err != nil {
			return err
		}
	case turnOff:
		if id, err := requireStrProp(a, deviceId); err != nil || !govalidator.IsUUID(*id) {
			return engine.ErrMalformedEntity
		}
	default:
		return engine.ErrMalformedEntity
	}

	return nil
}

func (a action) toDomain() engine.Action {
	switch a[name] {
	case sendEmail:
		return engine.SendEmailAction{
			Name:      sendEmail,
			Recipient: a[recipient].(string),
			Content:   a[content].(string),
		}
	case turnOff:
		return engine.TurnOffAction{
			Name:     turnOff,
			DeviceId: a[deviceId].(string),
		}
	}

	return nil
}

func convertBounds(object map[string]interface{}) (*bounds, error) {
	if fp, ok := object[from]; ok {
		if from, ok := fp.(float64); ok {
			if tp, ok := object[to]; ok {
				if to, ok := tp.(float64); ok {
					return &bounds{from, to}, nil
				}
			}
		}
	}

	return nil, engine.ErrMalformedEntity
}

func requireStrProp(object map[string]interface{}, prop string) (*string, error) {
	if p, ok := object[prop]; ok {
		if sp, ok := p.(string); ok && sp != "" {
			return &sp, nil
		}
	}

	return nil, engine.ErrMalformedEntity
}
