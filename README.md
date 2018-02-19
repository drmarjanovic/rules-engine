# Mainflux rules engine

[![Build Status][travis-img]][travis-url] [![codecov][codecov-img]][codecov-url]

The service exposes DSL for specifying alarming rules over an HTTP.

### Documentation
Documentation of the DSL syntax for the rules can be found [here](doc/DSLSYNTAX.md).

### Running

Make sure to start **Cassandra** and **Nats** first. From the project's root execute following command:
```
docker-compose -f docker-compose.infrastructure.yml up
```

In order to service successfully start, create **keyspace** in **Cassandra** named as exported environment variable `RULES_ENGINE_DB_KEYSPACE` (default **"rules_engine"**):
```
docker exec -it mainflux-rules-engine-cassandra cqlsh -e "CREATE KEYSPACE IF NOT EXISTS rules_engine WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };"
```

Execute command to run service:
```
go run cmd/main.go
```

It runs service on `127.0.0.1:9000` by default, or on port exported in `PORT` environment variable.
To verify setup, go to the browser and check `127.0.0.1:9000/health` URL.

If you want to run both services as Docker containers then build Docker images via:
```
./build-docker-images.sh
```

Run `docker-compose.mainflux.yml` via:
```
docker-compose -f docker-compose.mainflux.yml up
```

[codecov-img]: https://codecov.io/gh/MainfluxLabs/rules-engine/branch/dev/graph/badge.svg
[codecov-url]: https://codecov.io/gh/MainfluxLabs/rules-engine
[travis-img]: https://travis-ci.org/MainfluxLabs/rules-engine.svg?branch=dev
[travis-url]: https://travis-ci.org/MainfluxLabs/rules-engine
