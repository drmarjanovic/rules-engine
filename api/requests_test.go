package api

import (
	"testing"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
	"github.com/MainfluxLabs/rules-engine"
)

func TestViewRuleReqValidation(t *testing.T) {
	cases := []struct {
		userId string
		ruleId string
		err    error
	}{
		{gocql.TimeUUID().String(), gocql.TimeUUID().String(), nil},
		{"malformed user id", gocql.TimeUUID().String(), rules.ErrMalformedUrl},
		{gocql.TimeUUID().String(), "malformed rule id", rules.ErrMalformedUrl},
		{"malformed user id", "malformed rule id", rules.ErrMalformedUrl},
	}

	for i, tc := range cases {
		req := viewRuleReq{tc.userId, tc.ruleId}
		err := req.validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("failed at %d\n", i))
	}
}

func TestListRulesReqValidation(t *testing.T) {
	cases := []struct {
		userId string
		err    error
	}{
		{gocql.TimeUUID().String(), nil},
		{"malformed user-id", rules.ErrMalformedUrl},
	}

	for i, tc := range cases {
		req := listRulesReq{tc.userId}
		err := req.validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("failed at %d\n", i))
	}
}
