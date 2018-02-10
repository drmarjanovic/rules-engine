package nats

import (
	"testing"
	"encoding/json"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/MainfluxLabs/rules-engine"
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
		{invalidBtwVal, []string{}, []int{}, []int{}, rules.ErrMalformedEntity},
		{missingDeviceId, []string{}, []int{}, []int{}, rules.ErrMalformedEntity},
	}

	for i, tc := range cases {
		var raw *RawRules
		json.Unmarshal([]byte(tc.msg), &raw)
		rls, err := raw.toRules()

		for i, r := range rls {
			assert.Equal(t, tc.ruleNames[i], r.Name, fmt.Sprintf("failed at %d\n", i))
			assert.Equal(t, tc.actionNum[i], len(r.Actions), fmt.Sprintf("failed at %d\n", i))
			assert.Equal(t, tc.condNum[i], len(r.Conditions), fmt.Sprintf("failed at %d\n", i))
		}

		assert.Equal(t, tc.err, err, fmt.Sprintf("failed at %d\n", i))
	}
}
