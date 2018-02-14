package engine

import (
	"testing"
	"github.com/mainflux/mainflux/writer"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestIsSatisfied(t *testing.T) {
	cases := []struct {
		cnd       Condition
		event     writer.Message
		satisfied bool
	}{
		{Condition{"id", "active", Eq, Bool, true}, writer.Message{Publisher: "id", Name: "active", BoolValue: true}, true},
		{Condition{"mismatchedId", "active", Eq, Bool, true}, writer.Message{Publisher: "id", Name: "active", BoolValue: true}, false},
		{Condition{"id", "mismatchedProperty", Eq, Bool, true}, writer.Message{Publisher: "id", Name: "active", BoolValue: true}, false},
		{Condition{"id", "active", Eq, Bool, false}, writer.Message{Publisher: "id", Name: "active", BoolValue: true}, false},
		{Condition{"id",  "temp", Btw, Between, Range{15, 20}}, writer.Message{Publisher: "id", Name: "temp", Value: 18}, true},
		{Condition{"id",  "temp", Btw, Between, Range{15, 20}}, writer.Message{Publisher: "id", Name: "temp", Value: 15}, true},
		{Condition{"id",  "temp", Btw, Between, Range{15, 20}}, writer.Message{Publisher: "id", Name: "temp", Value: 20}, true},
		{Condition{"id",  "temp", Btw, Between, Range{15, 20}}, writer.Message{Publisher: "id", Name: "temp", Value: 30}, false},
		{Condition{"id", "temp", Eq, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 15}, true},
		{Condition{"id", "temp", Eq, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 20}, false},
		{Condition{"id", "temp", Lt, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 13}, false},
		{Condition{"id", "temp", Lt, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 20}, true},
		{Condition{"id", "temp", Lt, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 15}, false},
		{Condition{"id", "temp", Lte, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 10}, false},
		{Condition{"id", "temp",  Lte, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 15}, true},
		{Condition{"id", "temp", Lte, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 20}, true},
		{Condition{"id", "temp", Gt, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 20}, false},
		{Condition{"id", "temp", Gt, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 15}, false},
		{Condition{"id", "temp", Gt, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 10}, true},
		{Condition{"id", "temp", Gte, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 20}, false},
		{Condition{"id", "temp", Gte, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 15}, true},
		{Condition{"id", "temp", Gte, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 13}, true},
		{Condition{"id", "temp", Neq, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 20}, true},
		{Condition{"id", "temp", Neq, Numeric, float64(15)}, writer.Message{Publisher: "id", Name: "temp", Value: 15}, false},
		{Condition{"id",  "name", Eq, String,  "a"}, writer.Message{Publisher: "id", Name: "name", StringValue: "a"}, true},
		{Condition{"id",  "name", Eq, String,  "a"}, writer.Message{Publisher: "id", Name: "name", StringValue: "b"}, false},
		{Condition{"id",  "name", Neq, String,  "a"}, writer.Message{Publisher: "id", Name: "name", StringValue: "b"}, true},
		{Condition{"id",  "name", Neq, String,  "a"}, writer.Message{Publisher: "id", Name: "name", StringValue: "a"}, false},
	}
	for i, tc := range cases {
		satisfied := tc.cnd.isSatisfied(tc.event)
		assert.Equal(t, tc.satisfied, satisfied, fmt.Sprintf("failed at %d\n", i))
	}
}
