package engine

import "github.com/mainflux/mainflux/writer"

// Condition represents definition what needs to be satisfied in order to trigger
// an action.
type Condition struct {
	DeviceID string        `json:"deviceId"`
	Property string        `json:"property"`
	Operator Operator      `json:"operator"`
	Type     ConditionType `json:"-"`
	Value    interface{}   `json:"value"`
}

// ConditionType represent possible condition types based on value type
type ConditionType int

const (
	Bool    ConditionType = iota
	String
	Numeric
	Between
)

// Range represents upper and lower bound in between condition.
type Range struct {
	From float64 `json:"from"`
	To   float64 `json:"to"`
}

func (cnd Condition) isSatisfied(event writer.Message) bool {
	var actual interface{}

	satisfied := cnd.DeviceID == event.Publisher && cnd.Property == event.Name

	switch cnd.Type {
	case Bool:
		actual = event.BoolValue
	case String:
		actual = event.StringValue
	case Numeric, Between:
		actual = event.Value
	}

	return satisfied && cnd.Operator.Compare(cnd.Value, actual)
}
