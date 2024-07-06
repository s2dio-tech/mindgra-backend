package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/s2dio-tech/mindgra-backend/datasource"
	"github.com/s2dio-tech/mindgra-backend/domain"
)

type tokenRepository struct {
	Datasource *datasource.Neo4J
}

func InitTokenRepository(db *datasource.Neo4J) domain.TokenRepository {
	return &tokenRepository{
		Datasource: db,
	}
}

func (repo *tokenRepository) Store(r *domain.Token) (*string, error) {
	session := repo.Datasource.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `
		MATCH (u:User {id: $userId})
		CREATE (t:Token {
			id: apoc.create.uuid(),
			userId: $userId,
			type: $type,
			token: $token,
			createdAt: $createdAt
		})
		CREATE (t)-[:BELONGS_TO]->(u)
		RETURN t.id as id;`
		parameters := map[string]interface{}{
			"userId":    r.UserId,
			"type":      string(r.Type),
			"token":     r.Token,
			"createdAt": neo4j.LocalDateTimeOf(r.CreatedAt),
		}
		res, err := tx.Run(query, parameters)
		if err != nil || !res.Next() {
			return nil, err
		}

		record := res.Record().Values[0]
		return record, nil
	})

	if err != nil {
		return nil, err
	}

	id := result.(string)
	return &id, nil
}

func (repo *tokenRepository) FindToken(tokenType domain.TokenType, token string, userId string) (*domain.Token, error) {
	session := repo.Datasource.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})

	defer session.Close()

	result, err := session.Run(
		`MATCH (t:Token {type: $type, token: $token, userId: $userId})
		 RETURN t.id as id, t.createdAt as createdAt`,
		map[string]interface{}{
			"type":   string(tokenType),
			"token":  token,
			"userId": userId,
		},
	)
	if err != nil {
		return nil, err
	}
	if !result.Next() {
		return nil, nil
	}

	record := result.Record()
	if err != nil {
		return nil, err
	}

	id, _ := record.Get("id")
	createdAt, _ := record.Get("createdAt")
	return &domain.Token{
		Id:        id.(string),
		CreatedAt: createdAt.(neo4j.LocalDateTime).Time(),
	}, nil
}

func (repo *tokenRepository) FindOne(tokenType domain.TokenType, userId string) (*domain.Token, error) {
	session := repo.Datasource.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})

	defer session.Close()

	result, err := session.Run(
		`MATCH (t:Token {type: $type, userId: $userId})
		 RETURN t.id as id, t.token as token, t.createdAt as createdAt
		 LIMIT 1;`,
		map[string]interface{}{
			"type":   string(tokenType),
			"userId": userId,
		},
	)
	if err != nil {
		return nil, err
	}
	if !result.Next() {
		return nil, nil
	}

	record := result.Record()
	if err != nil {
		return nil, err
	}

	id, _ := record.Get("id")
	token, _ := record.Get("token")
	createdAt, _ := record.Get("createdAt")
	return &domain.Token{
		Id:        id.(string),
		Token:     token.(string),
		CreatedAt: createdAt.(neo4j.LocalDateTime).Time(),
	}, nil
}

func (r *tokenRepository) DeleteByTypeAndUserId(tType domain.TokenType, userId string) error {
	_, err := r.Datasource.ExecWrite(
		`MATCH (t:Token {type: $type, userId: $userId}) DETACH DELETE t;`,
		map[string]interface{}{
			"type":   tType,
			"userId": userId,
		},
	)
	return err
}
