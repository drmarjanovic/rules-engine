package nats

import (
	"encoding/json"
	"go.uber.org/zap"

	"github.com/nats-io/go-nats"
	"github.com/MainfluxLabs/rules-engine/engine"
)

var _ Subscriber = (*rulesSubscriber)(nil)

type rulesSubscriber struct {
	nc      *nats.Conn
	service engine.Service
	logger  *zap.Logger
}

// NewRulesSubscriber instantiates subscription handler for rule creation
func NewRulesSubscriber(nc *nats.Conn, service engine.Service, logger *zap.Logger) *rulesSubscriber {
	return &rulesSubscriber{nc, service, logger}
}

func (rs *rulesSubscriber) Subscribe(subject string, queue string) (*nats.Subscription, error) {
	return rs.nc.QueueSubscribe(subject, queue, func(m *nats.Msg) {
		var (
			rls []engine.Rule
			raw *rulesMsg
			err error
		)

		if err = json.Unmarshal(m.Data, &raw); err != nil {
			rs.logger.Error("Failed to unmarshal raw message.", zap.Error(err))
			return
		}

		if rls, err = raw.toDomain(); err != nil {
			rs.logger.Error("Failed to toDomain rules.", zap.Error(err))
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
