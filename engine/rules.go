package engine

import "github.com/mainflux/mainflux/writer"

// Rule represents base model for Mainflux rule.
type Rule struct {
	ID         string      `json:"id"`
	UserId     string      `json:"-"`
	Name       string      `json:"name,omitempty"`
	Conditions []Condition `json:"conditions"`
	Actions    []Action    `json:"actions"`
}

// Action represents base action specification.
type Action interface {
	Execute()
}

// IsMatchedBy checks that all event satisfies all conditions
// specified by rule.
func (rule Rule) IsMatchedBy(event writer.Message) bool {
	for _, cnd := range rule.Conditions {
		if !cnd.isSatisfied(event) {
			return false
		}
	}
	return true
}

// SendEmailAction represents model for triggering an email sending.
type SendEmailAction struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	Recipient string `json:"recipient"`
}

var _ Action = (*SendEmailAction)(nil)

func (action SendEmailAction) Execute() {}

// TurnOffAction represents model for triggering action to turn off
// the device.
type TurnOffAction struct {
	Name     string `json:"name"`
	DeviceId string `json:"deviceId"`
}

var _ Action = (*TurnOffAction)(nil)

func (action TurnOffAction) Execute() {}

// RuleRepository specifies API for rules managing.
type RuleRepository interface {
	// Save persists the rule. A non-nil error is returned to indicate
	// operation failure.
	Save(Rule) error

	// One retrieves specific rule by its owner and unique identifier.
	// A non-nil error is returned to indicate operation failure.
	One(string, string) (*Rule, error)

	// All retrieves list of rules for specific user.
	All(string) []Rule

	// Remove removes specific rule from database. A non-nil error is
	// returned to indicate operation failure.
	Remove(string, string) error
}
