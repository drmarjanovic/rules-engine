package engine

import (
	"encoding/json"
	"fmt"
)

// Operator represents possible operators that should be used to compare
// values using condition
type Operator int

const (
	Undefined Operator = iota
	Eq
	Neq
	Lt
	Lte
	Gt
	Gte
	Btw
)

func (op Operator) String() string {
	return operatorIDs[op]
}

var operatorIDs = map[Operator]string{
	Eq:  "=",
	Neq: "!=",
	Lt:  "<",
	Lte: "<=",
	Gt:  ">",
	Gte: ">=",
	Btw: "BETWEEN",
}

var operatorNames = idsToNames(operatorIDs)

func idsToNames(operatorIDs map[Operator]string) map[string]Operator {
	var names = make(map[string]Operator)
	for k, v := range operatorIDs {
		names[v] = k
	}

	return names
}

func (op *Operator) MarshalJSON() ([]byte, error) {
	if op == nil {
		return []byte("null"), nil
	}

	name, ok := operatorIDs[*op]
	if !ok {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf(`"%s"`, name)), nil
}

func (op *Operator) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	opId, ok := operatorNames[s]
	if !ok {
		return ErrMalformedEntity
	}

	*op = opId
	return nil
}

// Compare compares two values using specified operator
func (op Operator) Compare(expected, actual interface{}) bool {
	switch op {
	case Eq:
		return expected == actual
	case Neq:
		return expected != actual
	case Lt:
		return expected.(float64) < actual.(float64)
	case Lte:
		return expected.(float64) <= actual.(float64)
	case Gt:
		return expected.(float64) > actual.(float64)
	case Gte:
		return expected.(float64) >= actual.(float64)
	case Btw:
		value := actual.(float64)
		ex := expected.(Range)
		return ex.From <= value && value <= ex.To
	}

	return false
}
