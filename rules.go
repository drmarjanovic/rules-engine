package rules

// Rule represents base model for Mainflux rule.
type Rule struct {
	ID         string      `json:"id"`
	UserId     string      `json:"-"`
	Name       string      `json:"name,omitempty"`
	Conditions []Condition `json:"conditions"`
	Actions    []Action    `json:"actions"`
}

// Condition represents definition what needs to be satisfied in order to trigger
// an action.
type Condition interface{}

// ConditionData represents common condition properties.
type ConditionData struct {
	DeviceID *string `json:"deviceId"`
	Property *string `json:"property"`
	Operator *string `json:"operator"`
}

// Validate validates condition data properties.
func (data *ConditionData) Validate() error {
	if data.DeviceID == nil || data.Property == nil || data.Operator == nil {
		return ErrMalformedEntity
	}
	return nil
}

// BooleanCondition represents condition with boolean value.
type BooleanCondition struct {
	*ConditionData
	Value bool `json:"value"`
}

// BetweenCondition represents condition with value bounded in range.
type BetweenCondition struct {
	*ConditionData
	From *float64 `json:"from"`
	To   *float64 `json:"to"`
}

// Validate validates between condition properties.
func (cnd *BetweenCondition) Validate() error {
	if cnd.From == nil || cnd.To == nil || *cnd.From > *cnd.To {
		return ErrMalformedEntity
	}
	return nil
}

// BetweenCondition represents condition with numeric value.
type NumericCondition struct {
	*ConditionData
	Value float64 `json:"value"`
}

// StringCondition represents condition with string value.
type StringCondition struct {
	*ConditionData
	Value string `json:"value"`
}

// Action represents base action specification.
type Action interface{}

// SendEmailAction represents model for triggering an email sending.
type SendEmailAction struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	Recipient string `json:"recipient"`
}

// TurnOffAction represents model for triggering action to turn off the device.
type TurnOffAction struct {
	Name     string `json:"name"`
	DeviceId string `json:"deviceId"`
}

// RuleRepository specifies API for rules managing.
type RuleRepository interface {
	// Save persists the rule. A non-nil error is returned to indicate
	// operation failure.
	Save(Rule) error

	// One retrieves specific rule by its owner and unique identifier.
	// A non-nil error is returned to indicate operation failure.
	One(string, string) (Rule, error)

	// All retrieves list of rules for specific user.
	All(string) []Rule

	// Remove removes specific rule from database. A non-nil error is
	// returned to indicate operation failure.
	Remove(string, string) error
}
