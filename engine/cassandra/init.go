package cassandra

import "github.com/gocql/gocql"

var tables []string = []string{
	`CREATE TABLE IF NOT EXISTS rules (
		id uuid,
		user_id uuid,
		name text,
		conditions blob,
		actions blob,
		PRIMARY KEY ((user_id), id)
	)`,
}

// Connect establishes connection to the Cassandra cluster.
func Connect(hosts []string, keyspace string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	return cluster.CreateSession()
}

// Initialize creates tables used by the service.
func Initialize(session *gocql.Session) error {
	for _, table := range tables {
		if err := session.Query(table).Exec(); err != nil {
			return err
		}
	}

	return nil
}
