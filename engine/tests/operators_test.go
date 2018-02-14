package tests

import (
	"testing"
	"encoding/json"
	"fmt"
	"bytes"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/MainfluxLabs/rules-engine/engine"
)

type testCond struct {
	Operator engine.Operator `json:"op"`
}

func TestOperatorMarshaling(t *testing.T) {
	cases := []struct {
		op   engine.Operator
		json string
		err  error
	}{
		{engine.Eq, `"="`, nil},
		{engine.Neq, `"!="`, nil},
		{engine.Lt, `"<"`, nil},
		{engine.Lte, `"<="`, nil},
		{engine.Gt, `">"`, nil},
		{engine.Gte, `">="`, nil},
		{engine.Btw, `"BETWEEN"`, nil},
	}
	for i, tc := range cases {
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(&tc.op)
		actual := strings.TrimSpace(fmt.Sprintf("%s", buffer.Bytes()))

		assert.Equal(t, tc.err, err, fmt.Sprintf("failed at %d\n", i))
		assert.Equal(t, tc.json, actual, fmt.Sprintf("failed at %d\n", i))
	}
}

func TestOperatorUnmarshaling(t *testing.T) {
	cases := []struct {
		json string
		op   engine.Operator
		err  error
	}{
		{`{"op" : "="}`, engine.Eq, nil},
		{`{"op" : "!="}`, engine.Neq, nil},
		{`{"op" : "<"}`, engine.Lt, nil},
		{`{"op" : "<="}`, engine.Lte, nil},
		{`{"op" : ">"}`, engine.Gt, nil},
		{`{"op" : ">="}`, engine.Gte, nil},
		{`{"op" : "BETWEEN"}`, engine.Btw, nil},
		{`{"op" : "invalid"}`, engine.Undefined, engine.ErrMalformedEntity},
		{`{"op" : null}`, engine.Undefined, engine.ErrMalformedEntity},
	}

	for i, tc := range cases {
		var cnd *testCond
		err := json.Unmarshal([]byte(tc.json), &cnd)

		assert.Equal(t, tc.err, err, fmt.Sprintf("failed at %d\n", i))
		assert.Equal(t, tc.op, cnd.Operator, fmt.Sprintf("failed at %d\n", i))
	}
}
