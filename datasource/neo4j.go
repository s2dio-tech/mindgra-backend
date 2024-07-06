package datasource

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/s2dio-tech/mindgra-backend/common"
)

type Neo4J struct {
	Driver neo4j.Driver
}

func InitNeo4J() Neo4J {
	driver, err := neo4j.NewDriver(
		"neo4j://"+common.AppConfig.DBHost+":"+common.AppConfig.DBPort,
		neo4j.BasicAuth(common.AppConfig.DBUsername, common.AppConfig.DBPassword, ""),
	)
	if err != nil {
		panic(err)
	}
	return Neo4J{
		Driver: driver,
	}
}

func (m *Neo4J) ExecRead(query string, params map[string]any) ([]*neo4j.Record, error) {
	session := m.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close()

	result, err := session.Run(query, params)

	if err != nil {
		return nil, err
	}

	records, err := result.Collect()
	if err != nil {
		return nil, err
	}

	return records, err
}

func (m *Neo4J) ExecWrite(query string, params map[string]any) ([]*neo4j.Record, error) {
	session := m.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		res, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}

		records, err := res.Collect()
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

func (m *Neo4J) Disconnect() {
	m.Driver.Close()
}
