package nats

import (
	"encoding/json"
	"go.uber.org/zap"

	"github.com/nats-io/go-nats"
	"github.com/mainflux/mainflux/writer/cassandra"
	"github.com/mainflux/mainflux/writer"
)

var _ Subscriber = (*eventsSubscriber)(nil)

type eventsSubscriber struct {
	nc     *nats.Conn
	logger *zap.Logger
}

// NewEventSubscriber instantiates subscription handler for senML messages
func NewEventSubscriber(nc *nats.Conn, logger *zap.Logger) *eventsSubscriber {
	return &eventsSubscriber{nc, logger}
}

func (es *eventsSubscriber) Subscribe(subject string, queue string) (*nats.Subscription, error) {
	return es.nc.QueueSubscribe(subject, queue, func(m *nats.Msg) {
		var (
			events []writer.Message
			raw    writer.RawMessage
			err    error
		)

		if err = json.Unmarshal(m.Data, &raw); err != nil {
			es.logger.Error("Failed to unmarshal raw event message.", zap.Error(err))
			return
		}

		if events, err = cassandra.Normalize(raw); err != nil {
			es.logger.Error("Unable to parse SenML message.", zap.Error(err))
			return
		}

		sugar := es.logger.Sugar()
		sugar.Infof("Applying rules on %d events.", len(events))
	})
}
