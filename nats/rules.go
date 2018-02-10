package nats

import (
	"encoding/json"
	"go.uber.org/zap"

	"github.com/MainfluxLabs/rules-engine"
	"github.com/nats-io/go-nats"
)

var _ Subscriber = (*rulesSubscriber)(nil)

type rulesSubscriber struct {
	nc      *nats.Conn
	service rules.Service
	logger  *zap.Logger
}

// NewRulesSubscriber instantiates subscription handler for rule creation
func NewRulesSubscriber(nc *nats.Conn, service rules.Service, logger *zap.Logger) *rulesSubscriber {
	return &rulesSubscriber{nc, service, logger}
}

func (rs *rulesSubscriber) Subscribe(subject string, queue string) (*nats.Subscription, error) {
	return rs.nc.QueueSubscribe(subject, queue, func(m *nats.Msg) {
		var (
			rls []rules.Rule
			raw *RawRules
			err error
		)

		if err = json.Unmarshal(m.Data, &raw); err != nil {
			rs.logger.Error("Failed to unmarshal raw message.", zap.Error(err))
			return
		}

		if rls, err = raw.toRules(); err != nil {
			rs.logger.Error("Failed to parse rules.", zap.Error(err))
			return
		}

		sugar := rs.logger.Sugar()
		for _, r := range rls {
			if err = rs.service.SaveRule(r); err != nil {
				rs.logger.Error("Unable to save rule to cassandra.", zap.Error(err))
			} else {
				sugar.Infof("Rule %s successfully saved to cassandra.", r.Name)
			}
		}
	})
}
