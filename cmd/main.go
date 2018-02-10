package main

import (
	"fmt"
	"os"
	"strings"
	"net/http"

	subscribers "github.com/MainfluxLabs/rules-engine/nats"
	"github.com/MainfluxLabs/rules-engine"
	"github.com/MainfluxLabs/rules-engine/api"
	"github.com/MainfluxLabs/rules-engine/cassandra"
	"github.com/nats-io/go-nats"
	"go.uber.org/zap"
)

const (
	sep           string = ","
	defPort       string = "9000"
	envPort       string = "PORT"
	defCluster    string = "127.0.0.1"
	envCluster    string = "RULES_ENGINE_DB_CLUSTER"
	defKeyspace   string = "rules_engine"
	envKeyspace   string = "RULES_ENGINE_DB_KEYSPACE"
	defNatsURL    string = nats.DefaultURL
	envNatsURL    string = "NATS_URL"
	eventsSubject string = "msg.*"
	eventsQueue   string = "event.consumer"
	rulesSubject  string = "rules"
	rulesQueue    string = "rule.consumer"
)

type config struct {
	Port     string
	Cluster  string
	Keyspace string
	NatsURL  string
}

func main() {
	cfg := config{
		Port:     getenv(envPort, defPort),
		Cluster:  getenv(envCluster, defCluster),
		Keyspace: getenv(envKeyspace, defKeyspace),
		NatsURL:  getenv(envNatsURL, defNatsURL),
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	session, err := cassandra.Connect(strings.Split(cfg.Cluster, sep), cfg.Keyspace)
	if err != nil {
		os.Exit(1)
	}
	defer session.Close()

	if err := cassandra.Initialize(session); err != nil {
		os.Exit(1)
	}

	rulesRepo := cassandra.NewRuleRepository(session)
	svc := rules.NewService(rulesRepo)

	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		os.Exit(1)
	}
	defer nc.Close()

	eventsSubscriber := subscribers.NewEventSubscriber(nc, logger)
	if _, err = eventsSubscriber.Subscribe(eventsSubject, eventsQueue); err != nil {
		logger.Error("Unable to subscribe on Mainflux msg.* topics.", zap.Error(err))
		os.Exit(1)
	}

	rulesSubscriber := subscribers.NewRulesSubscriber(nc, svc, logger)
	if _, err = rulesSubscriber.Subscribe(rulesSubject, rulesQueue); err != nil {
		logger.Error("Unable to subscribe on rules topic.", zap.Error(err))
		os.Exit(1)
	}

	p := fmt.Sprintf(":%s", cfg.Port)
	http.ListenAndServe(p, api.MakeHandler(svc))
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
