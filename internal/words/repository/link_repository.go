package repository

import (
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/s2dio-tech/mindgra-backend/common"
	"github.com/s2dio-tech/mindgra-backend/datasource"
	"github.com/s2dio-tech/mindgra-backend/domain"
)

type linkRepository struct {
	Datasource *datasource.Neo4J
}

func InitLinkRepository(db *datasource.Neo4J) domain.LinkRepository {
	return &linkRepository{
		Datasource: db,
	}
}

func (repo *linkRepository) Store(r domain.Link) (*string, error) {

	result, err := repo.Datasource.ExecWrite(
		`MATCH (u:User {id: $userId})
			MATCH (w1:Word {id: $word1Id})
			MATCH (w2:Word {id: $word2Id})
			MATCH (w1)-[r]-(w2)
			CREATE (w:Link {
				id: apoc.create.uuid(),
				userId: $userId,
				word1Id: $word1Id,
				word2Id: $word2Id,
				content: $content,
				description: $description,
				refs: $refs,
				createdAt: $createdAt
			})
			CREATE (u)-[:OWN]->(w)
			SET r.id = w.id
			RETURN w.id as id;`,
		map[string]interface{}{
			"userId":      r.UserId,
			"word1Id":     r.Word1Id,
			"word2Id":     r.Word2Id,
			"content":     r.Content,
			"description": r.Description,
			"refs":        r.Refs,
			"createdAt":   neo4j.LocalDateTimeOf(r.CreatedAt),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, common.ErrInternalServerError
	}
	record := result[0]
	if err != nil {
		return nil, err
	}

	_id, _ := record.Get("id")
	return common.Nullable{Value: _id}.ToStringPtr(), nil
}

func (r *linkRepository) FindById(id string) (*domain.Link, error) {
	result, err := r.Datasource.ExecRead(
		`MATCH (r:Link {id: $id})
			RETURN r.id AS id,
				r.userId AS userId,
				r.content AS content,
				r.description AS description,
				r.refs AS refs,
				r.createdAt AS createdAt;`,
		map[string]interface{}{
			"id": id,
		},
	)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	record := result[0]
	if err != nil {
		return nil, err
	}

	_id, _ := record.Get("id")
	_userId, _ := record.Get("userId")
	_content, _ := record.Get("content")
	_description, _ := record.Get("description")
	_refs, _ := record.Get("refs")
	_createdAt, _ := record.Get("createdAt")
	return &domain.Link{
		Id:          _id.(string),
		UserId:      _userId.(string),
		Content:     _content.(string),
		Description: common.Nullable{Value: _description}.ToStringPtr(),
		Refs:        common.Nullable{Value: _refs}.ToStringArrayPtr(),
		CreatedAt:   _createdAt.(neo4j.LocalDateTime).Time(),
	}, nil
}

func (r *linkRepository) FindByWordIds(w1Id string, w2Id string) (*domain.Link, error) {
	result, err := r.Datasource.ExecRead(
		`MATCH (r:Link)
			WHERE (r.word1Id = $w1Id AND r.word2Id = $w2Id)
				 OR (r.word1Id = $w2Id AND r.word2Id = $w1Id)
			RETURN r.id AS id,
				r.userId AS userId,
				r.content AS content,
				r.description AS description,
				r.refs AS refs,
				r.createdAt AS createdAt;`,
		map[string]interface{}{
			"w1Id": w1Id,
			"w2Id": w2Id,
		},
	)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	record := result[0]
	if err != nil {
		return nil, err
	}

	_id, _ := record.Get("id")
	_userId, _ := record.Get("userId")
	_content, _ := record.Get("content")
	_description, _ := record.Get("description")
	_refs, _ := record.Get("refs")
	_createdAt, _ := record.Get("createdAt")
	return &domain.Link{
		Id:          _id.(string),
		UserId:      _userId.(string),
		Content:     _content.(string),
		Description: common.Nullable{Value: _description}.ToStringPtr(),
		Refs:        common.Nullable{Value: _refs}.ToStringArrayPtr(),
		CreatedAt:   _createdAt.(neo4j.LocalDateTime).Time(),
	}, nil
}

func (repo *linkRepository) Update(id string, link domain.Link) error {
	query := `MATCH (r:Link {id: $id})
	SET r.content = $content,
			r.description = $description,
			r.refs = $refs,
			r.updatedAt = $updatedAt;`
	params := map[string]interface{}{
		"id":          id,
		"content":     link.Content,
		"description": link.Description,
		"refs":        link.Refs,
		"updatedAt":   neo4j.LocalDateTimeOf(time.Now()),
	}

	_, err := repo.Datasource.ExecWrite(query, params)

	if err != nil {
		return err
	}
	return nil
}

func (r *linkRepository) Delete(w1Id string, w2Id string) error {
	// _, err := r.Datasource.ExecWrite(
	// 	`MATCH (r:Link {id: $id})
	// 		MATCH (w1:Word) WHERE Id(w1) = r.word1Id
	// 		MATCH (w2:Word) WHERE Id(w2) = r.word2Id
	// 		MATCH (w1)-[l]-(w2)
	// 		DETACH DELETE r
	// 		DELETE l;`,
	// 	map[string]interface{}{
	// 		"id": id,
	// 	},
	// )
	_, err := r.Datasource.ExecWrite(
		`MATCH (w1:Word {id: $w1Id})-[r:CONCERN]-(w2:Word {id: $w2Id})
		DELETE r;
		`,
		map[string]interface{}{
			"w1Id": w1Id,
			"w2Id": w2Id,
		},
	)
	return err
}
