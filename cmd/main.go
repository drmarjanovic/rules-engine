package main

import (
	"fmt"
	"os"
	"strings"
	"net/http"

	"github.com/MainfluxLabs/rules-engine"
	"github.com/MainfluxLabs/rules-engine/api"
	"github.com/MainfluxLabs/rules-engine/cassandra"
)

const (
	sep         string = ","
	defPort     string = "9000"
	envPort     string = "PORT"
	defCluster  string = "127.0.0.1"
	envCluster  string = "RULES_ENGINE_DB_CLUSTER"
	defKeyspace string = "rules_engine"
	envKeyspace string = "RULES_ENGINE_DB_KEYSPACE"
)

type config struct {
	Port     string
	Cluster  string
	Keyspace string
}

func main() {
	cfg := config{
		Port:     getenv(envPort, defPort),
		Cluster:  getenv(envCluster, defCluster),
		Keyspace: getenv(envKeyspace, defKeyspace),
	}

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
