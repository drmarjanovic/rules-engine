package nats

import (
	"testing"
	"encoding/json"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/MainfluxLabs/rules-engine/engine"
	"github.com/gocql/gocql"
)

const (
	validRule = `
		{
		  "rules": [
			{
			  "conditions": [
				{
				  "operator": "=",
				  "property": "offline",
				  "deviceId": "a32db207-7236-4e75-abad-7c972f4cfd18",
				  "value": true
				}
			  ],
			  "userId": "a2dfc0dc-1f14-4935-a78b-92e77c0af7a1",
			  "name": "rule01",
			  "actions": [
				{
				  "name": "TURN OFF",
				  "deviceId": "937a7c3e-db39-4d75-b52b-b8442463761a"
				},
				{
				  "content": "<b>Alert!<\/b><br><p>High temperature in living room!<\/p>",
				  "recipient": "person01@home.com",
				  "name": "SEND EMAIL"
				},
				{
				  "content": "<b>Alert!<\/b><br><p>High temperature in living room!<\/p>",
				  "recipient": "person02@home.com",
				  "name": "SEND EMAIL"
				}
			  ]
			}
		  ]
		}`
	twoRules = `
		{
		  "rules": [
			{
			  "conditions": [
				{
				  "operator": ">=",
				  "property": "temperature",
				  "deviceId": "a32db207-7236-4e75-abad-7c972f4cfd18",
				  "value": 30
				}
			  ],
			  "userId": "a2dfc0dc-1f14-4935-a78b-92e77c0af7a1",
			  "name": "rule01",
			  "actions": [
				{
				  "name": "TURN OFF",
				  "deviceId": "937a7c3e-db39-4d75-b52b-b8442463761a"
				}
			  ]
			},
			{
			  "conditions": [
				{
				  "operator": "BETWEEN",
				  "property": "org",
				  "deviceId": "a32db207-7236-4e75-abad-7c972f4cfd18",
				  "value": {
					"from": 30,
					"to": 40
				  }
				}
			  ],
			  "userId": "a2dfc0dc-1f14-4935-a78b-92e77c0af7a1",
			  "name": "rule02",
			  "actions": [
				{
				  "name": "TURN OFF",
				  "deviceId": "937a7c3e-db39-4d75-b52b-b8442463761a"
				},
				{
				  "content": "<b>Alert!<\/b><br><p>High temperature in living room!<\/p>",
				  "recipient": "person02@home.com",
				  "name": "SEND EMAIL"
				}
			  ]
			}
		  ]
		}`
	ruleWithoutName = `
		{
		  "rules": [
			{
			  "conditions": [
				{
				  "operator": "=",
				  "property": "offline",
				  "deviceId": "a32db207-7236-4e75-abad-7c972f4cfd18",
				  "value": true
				}
			  ],
			  "userId": "a2dfc0dc-1f14-4935-a78b-92e77c0af7a1",
			  "actions": [
				{
				  "content": "<b>Alert!<\/b><br><p>High temperature in living room!<\/p>",
				  "recipient": "person02@home.com",
				  "name": "SEND EMAIL"
				}
			  ]
			}
		  ]
		}`
	invalidBtwVal = `
		{
			"rules": [
				{
					"conditions": [
						{
							"operator": ">=",
							"property": "temperature",
							"deviceId": "a32db207-7236-4e75-abad-7c972f4cfd18",
							"value": {
								"from": 30,
								"to": "wrong"
							}
						}
					],
					"userId": "a2dfc0dc-1f14-4935-a78b-92e77c0af7a1",
					"name": "rule01",
					"actions": [
						{
							"name": "TURN OFF",
							"deviceId": "937a7c3e-db39-4d75-b52b-b8442463761a"
						}
					]
				}
			]
		}`
	missingDeviceId = `
		{
		  "rules": [
			{
			  "conditions": [
				{
				  "operator": "=",
				  "property": "offline",
				  "value": true
				}
			  ],
			  "userId": "a2dfc0dc-1f14-4935-a78b-92e77c0af7a1",
			  "name": "rule01",
			  "actions": [
				{
				  "name": "TURN OFF",
				  "deviceId": "937a7c3e-db39-4d75-b52b-b8442463761a"
				}
			  ]
			}
		  ]
		}`
)

var (
	uuid             = gocql.TimeUUID().String()
	validAction      = action{name: sendEmail, content: "test", recipient: "test"}
	validCondition   = condition{uuid, "active", engine.Eq, true}
	invalidAction    = action{name: sendEmail, content: "", recipient: "test"}
	invalidCondition = condition{uuid, "active", engine.Gt, true}
)

func TestParsingRules(t *testing.T) {
	cases := []struct {
		msg       string
		ruleNames []string
		condNum   []int
		actionNum []int
		err       error
	}{
		{validRule, []string{"rule01"}, []int{1}, []int{3}, nil},
		{twoRules, []string{"rule01", "rule02"}, []int{1, 1}, []int{1, 2}, nil},
		{ruleWithoutName, []string{""}, []int{1}, []int{1}, nil},
		{invalidBtwVal, []string{}, []int{}, []int{}, engine.ErrMalformedEntity},
		{missingDeviceId, []string{}, []int{}, []int{}, engine.ErrMalformedEntity},
	}

	for i, tc := range cases {
		var raw *rulesMsg
		json.Unmarshal([]byte(tc.msg), &raw)
		rls, err := raw.toDomain()

		for i, r := range rls {
			assert.Equal(t, tc.ruleNames[i], r.Name, fmt.Sprintf("failed at %d\n", i))
			assert.Equal(t, tc.actionNum[i], len(r.Actions), fmt.Sprintf("failed at %d\n", i))
			assert.Equal(t, tc.condNum[i], len(r.Conditions), fmt.Sprintf("failed at %d\n", i))
		}

		assert.Equal(t, tc.err, err, fmt.Sprintf("failed at %d\n", i))
	}
}

func TestValidateRule(t *testing.T) {
	cases := []struct {
		r   rule
		err error
	}{
		{rule{uuid, "", []condition{validCondition}, []action{validAction}}, nil},
		{rule{uuid, "test", []condition{validCondition}, []action{validAction}}, nil},
		{rule{"test", "", []condition{validCondition}, []action{validAction}}, engine.ErrMalformedEntity},
		{rule{uuid, "", []condition{}, []action{validAction}}, engine.ErrMalformedEntity},
		{rule{uuid, "", []condition{validCondition}, []action{}}, engine.ErrMalformedEntity},
		{rule{uuid, "", []condition{validCondition, invalidCondition}, []action{validAction}}, engine.ErrMalformedEntity},
		{rule{uuid, "", []condition{validCondition}, []action{validAction, invalidAction}}, engine.ErrMalformedEntity},
	}

	for i, tc := range cases {
		err := tc.r.validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("failed at %d\n", i))
	}

}

func TestValidateCondition(t *testing.T) {
	cases := []struct {
		cnd condition
		err error
	}{
		{condition{uuid, "active", engine.Eq, true}, nil},
		{condition{uuid, "active", engine.Gt, true}, engine.ErrMalformedEntity},
		{condition{uuid, "name", engine.Eq, "test"}, nil},
		{condition{uuid, "active", engine.Btw, "test"}, engine.ErrMalformedEntity},
		{condition{uuid, "temp", engine.Neq, float64(5)}, nil},
		{condition{uuid, "active", engine.Btw, float64(5)}, engine.ErrMalformedEntity},
		{condition{uuid, "active", engine.Btw, map[string]interface{}{from: float64(5), to: float64(10)}}, nil},
		{condition{uuid, "active", engine.Btw, map[string]interface{}{from: "5", to: "10"}}, engine.ErrMalformedEntity},
		{condition{uuid, "active", engine.Btw, true}, engine.ErrMalformedEntity},
		{condition{uuid, "active", engine.Btw, map[string]interface{}{from: float64(10), to: float64(5)}}, engine.ErrMalformedEntity},
		{condition{uuid, "active", engine.Btw, map[string]interface{}{from: float64(10), to: float64(10)}}, engine.ErrMalformedEntity},
		{condition{"invalid", "active", engine.Eq, true}, engine.ErrMalformedEntity},
		{condition{uuid, "", engine.Eq, true}, engine.ErrMalformedEntity},
		{condition{"", "test", engine.Eq, true}, engine.ErrMalformedEntity},
	}

	for i, tc := range cases {
		err := tc.cnd.validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("failed at %d\n", i))
	}
}

func TestValidateAction(t *testing.T) {
	cases := []struct {
		action action
		err    error
	}{
		{action{name: sendEmail, content: "test", recipient: "test"}, nil},
		{action{name: sendEmail, content: "", recipient: "test"}, engine.ErrMalformedEntity},
		{action{name: sendEmail, content: 5, recipient: "test"}, engine.ErrMalformedEntity},
		{action{name: sendEmail, content: "test", recipient: ""}, engine.ErrMalformedEntity},
		{action{name: "invalidName", content: "test", recipient: "test"}, engine.ErrMalformedEntity},
		{action{"invalidProperty": "test", content: "test", recipient: "test"}, engine.ErrMalformedEntity},
		{action{name: turnOff, deviceId: uuid}, nil},
		{action{name: turnOff, deviceId: "test"}, engine.ErrMalformedEntity},
		{action{name: turnOff, content: "test", recipient: "test"}, engine.ErrMalformedEntity},
	}

	for i, tc := range cases {
		err := tc.action.validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("failed at %d\n", i))
	}
}
