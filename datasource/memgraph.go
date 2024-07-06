package datasource

import (
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type MemGraph struct {
	Driver neo4j.Driver
}

func InitMemGraph() MemGraph {
	uri, found := os.LookupEnv("MEMGRAPH_URI")
	if !found {
		panic("MEMGRAPH_URI not set")
	}
	username, found := os.LookupEnv("MEMGRAPH_USERNAME")
	if !found {
		panic("MEMGRAPH_USERNAME not set")
	}
	password, found := os.LookupEnv("MEMGRAPH_PASSWORD")
	if !found {
		panic("MEMGRAPH_PASSWORD not set")
	}
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		panic(err)
	}
	return MemGraph{
		Driver: driver,
	}
}

func (m *MemGraph) ExecRead(query string, params map[string]any) ([]*neo4j.Record, error) {
	session := m.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close()

	result, err := session.Run(query, params)

	if err != nil {
		return nil, err
	}

	records, err := result.(neo4j.Result).Collect()
	if err != nil {
		return nil, err
	}

	return records, err
}

func (m *MemGraph) ExecWrite(query string, params map[string]any) ([]*neo4j.Record, error) {
	session := m.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		res, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}

		records, err := res.(neo4j.Result).Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*neo4j.Record), err
}

func (m *MemGraph) Disconnect() {
	m.Driver.Close()
}
