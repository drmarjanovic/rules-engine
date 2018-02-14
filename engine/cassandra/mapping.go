package cassandra

import (
	"github.com/fatih/structs"
	"github.com/asaskevich/govalidator"
	"github.com/MainfluxLabs/rules-engine/engine"
)

const (
	sendEmail string = "SEND EMAIL"
	turnOff   string = "TURN OFF"

	name      string = "Name"
	content   string = "Content"
	recipient string = "Recipient"
	deviceId  string = "DeviceId"
)

type dbAction map[string]interface{}

func fromDomain(actions []engine.Action) ([]dbAction) {
	var dbActions []dbAction
	for _, a := range actions {
		dbActions = append(dbActions, structs.Map(a))
	}

	return dbActions
}

func toDomain(dbActs []dbAction) ([]engine.Action, error) {
	var actions []engine.Action

	for _, dba := range dbActs {
		a, err := dba.toDomain()
		if err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}

	return actions, nil
}

func (action dbAction) toDomain() (engine.Action, error) {
	if err := action.validate(); err != nil {
		return nil, err
	}

	name, ok := action[name]
	if !ok {
		return nil, engine.ErrMalformedEntity
	}

	switch name {
	case sendEmail:
		if se, err := action.toSendEmail(); err == nil {
			return se, nil
		}
	case turnOff:
		if a, err := action.toTurnOff(); err == nil {
			return a, nil
		}
	}

	return nil, engine.ErrMalformedEntity
}

func (action dbAction) validate() error {
	name, ok := action[name]

	if !ok {
		return engine.ErrMalformedEntity
	}

	switch name {
	case sendEmail:
		if _, err := requireStrProp(action, content); err != nil {
			return err
		}
		if _, err := requireStrProp(action, recipient); err != nil {
			return err
		}
	case turnOff:
		if id, err := requireStrProp(action, deviceId); err != nil || !govalidator.IsUUID(*id) {
			return engine.ErrMalformedEntity
		}
	default:
		return engine.ErrMalformedEntity
	}

	return nil
}

func (action dbAction) toSendEmail() (*engine.SendEmailAction, error) {
	se := &engine.SendEmailAction{
		Name: name,
	}

	c, err := requireStrProp(action, content)
	if err != nil {
		return nil, err
	}
	se.Content = *c

	r, err := requireStrProp(action, recipient)
	if err != nil {
		return nil, err
	}
	se.Recipient = *r

	return se, nil
}

func (action dbAction) toTurnOff() (*engine.TurnOffAction, error) {
	se := &engine.TurnOffAction{
		Name: name,
	}

	id, err := requireStrProp(action, deviceId)
	if err != nil {
		return nil, err
	}
	se.DeviceId = *id

	return se, nil
}

func requireStrProp(object map[string]interface{}, prop string) (*string, error) {
	if p, ok := object[prop]; ok {
		if sp, ok := p.(string); ok && sp != "" {
			return &sp, nil
		}
	}

	return nil, engine.ErrMalformedEntity
}
