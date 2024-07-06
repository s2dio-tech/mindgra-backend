package repository

import (
	"strconv"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/s2dio-tech/mindgra-backend/common"
	"github.com/s2dio-tech/mindgra-backend/datasource"
	"github.com/s2dio-tech/mindgra-backend/domain"
	"golang.org/x/exp/slog"
)

type wordRepository struct {
	Datasource *datasource.Neo4J
}

func InitWordRepository(db *datasource.Neo4J) domain.WordRepository {
	return &wordRepository{
		Datasource: db,
	}
}

func recordToWord(record map[string]any) *domain.Word {
	w := domain.Word{
		Id:          record["id"].(string),
		UserId:      record["userId"].(string),
		Content:     record["content"].(string),
		Description: common.Nullable{Value: record["description"]}.ToStringPtr(),
		Refs:        common.Nullable{Value: record["refs"]}.ToStringArrayPtr(),
	}
	if record["createdAt"] != nil {
		w.CreatedAt = common.ToPointer(record["createdAt"].(neo4j.LocalDateTime).Time())
	}
	return &w
}

func (repo *wordRepository) Store(w domain.Word, graphId string, linkWordId *string) (*string, error) {
	query := `MATCH (u:User {id: $userId})
		MATCH (s:Graph {id: $graphId})
		CREATE (w:Word {
			id: apoc.create.uuid(),
			userId: $userId,
			graphId: $graphId,
			content: $content,
			description: $description,
			refs: $refs,
			createdAt: $createdAt
		})
		CREATE (u)-[:OWN]->(w)
		CREATE (s)-[:WORD]->(w)`
	params := map[string]interface{}{
		"userId":      w.UserId,
		"graphId":     graphId,
		"content":     w.Content,
		"description": w.Description,
		"refs":        w.Refs,
		"createdAt":   neo4j.LocalDateTimeOf(time.Now()),
	}

	if linkWordId != nil {
		query = "MATCH (r:Word {id: $linkWordId})\n" + query + "\nCREATE (w)-[:CONCERN]->(r)"
		params["linkWordId"] = linkWordId
	}
	query += "\nRETURN w.id AS id;"

	result, err := repo.Datasource.ExecWrite(query, params)

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

func (r *wordRepository) Update(w domain.Word) error {
	_, err := r.Datasource.ExecWrite(
		`MATCH (w:Word {id: $id})
		SET w.content= $content,
			w.description= $description,
			w.refs= $refs,
			w.updatedAt= $updatedAt;`,
		map[string]interface{}{
			"id":          w.Id,
			"content":     w.Content,
			"description": w.Description,
			"refs":        w.Refs,
			"updatedAt":   neo4j.LocalDateTimeOf(time.Now()),
		},
	)
	return err
}

func (r *wordRepository) Delete(id string) error {
	// remove word and links
	_, err := r.Datasource.ExecWrite(
		`MATCH (w:Word {id: $id})
			MATCH (r:Link) WHERE r.word1Id = $id OR r.word2Id = $id
			DETACH DELETE w,r;`,
		map[string]interface{}{
			"id": id,
		},
	)
	return err
}

func (r *wordRepository) FindById(id string) (*domain.Word, error) {
	result, err := r.Datasource.ExecRead(
		`MATCH (w:Word {id: $id})
			RETURN w.id as id,
				w.userId AS userId,
				w.content AS content,
				w.description AS description,
				w.refs AS refs,
				w.createdAt AS createdAt;`,
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

	return recordToWord(record.AsMap()), nil
}

func (r *wordRepository) FindByRandomId() (*domain.Word, error) {
	result, err := r.Datasource.ExecRead(
		`MATCH (w:Word) 
			RETURN rand() as r,
				w.id as id,
				w.userId AS userId,
				w.content AS content,
				w.description AS description,
				w.refs AS refs,
				w.createdAt AS createdAt
			ORDER BY r
			LIMIT 1;`,
		nil,
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

	return recordToWord(record.AsMap()), nil
}

func (r *wordRepository) FindByIds(ids []string) ([]domain.Word, error) {
	result, err := r.Datasource.ExecRead(
		`MATCH (w:Word) WHERE w.id IN $ids
			RETURN w.id as id,
				w.userId AS userId,
				w.content AS content;`,
		map[string]interface{}{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}

	words := []domain.Word{}
	for _, record := range result {
		words = append(words, *recordToWord(record.AsMap()))
	}
	return words, nil
}

func (r *wordRepository) FindByGraphId(graphId string) ([]domain.Word, []domain.WordsLink, error) {
	// get words
	ws, err := r.Datasource.ExecRead(`
		MATCH(s:Graph {id: $graphId})-[:WORD]-(w:Word)
		RETURN distinct w as ws`,
		map[string]interface{}{
			"graphId": graphId,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	// get links
	ls, err := r.Datasource.ExecRead(`
		MATCH(s:Graph {id: $graphId})-[:WORD]-(w:Word)-[r:CONCERN]-()
		WITH collect(r) as rls
		UNWIND rls as list
		UNWIND list as ritem
		RETURN distinct ritem as ls`,
		map[string]interface{}{
			"graphId": graphId,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	wordIds := map[int64]string{}
	words := []domain.Word{}
	links := []domain.WordsLink{}

	for _, record := range ws {
		if n, ok := record.AsMap()["ws"].(dbtype.Node); ok {
			w := *recordToWord(n.GetProperties())
			wordIds[n.GetId()] = w.Id
			words = append(words, w)
		}
	}

	if len(ls) == 0 {
		return words, links, nil
	}

	for _, record := range ls {
		if r, ok := record.AsMap()["ls"].(dbtype.Relationship); ok {
			links = append(links, domain.WordsLink{
				SourceId: wordIds[r.StartId],
				TargetId: wordIds[r.EndId],
			})
		}
	}
	return words, links, nil
}

func (r *wordRepository) FindNeighborIds(id string, depth int) ([]domain.WordsLink, error) {
	result, err := r.Datasource.ExecRead(
		`MATCH path = (w1:Word {id: $id})-[:CONCERN*0..`+strconv.Itoa(depth)+`]-(w2:Word)
			UNWIND relationShips(path) AS r
			RETURN startNode(r).id AS id1, endNode(r).id as id2;
		`,
		map[string]interface{}{
			"id": id,
			// "depth": depth,
		},
	)
	if err != nil {
		slog.Error("Error in FindNeighborIds", err)
		return nil, err
	}

	res := []domain.WordsLink{}
	for _, record := range result {
		id1, _ := record.Get("id1")
		id2, _ := record.Get("id2")
		res = append(res, domain.WordsLink{
			SourceId: id1.(string),
			TargetId: id2.(string),
		})
	}
	return res, nil
}

func (r *wordRepository) FindByContentOrDescription(search string, limit int) ([]domain.Word, error) {
	result, err := r.Datasource.ExecRead(
		`CALL db.index.fulltext.queryNodes("contentAndDescriptions", $text) YIELD node, score
			RETURN node.id as id,
				node.content as content,
				node.description as description,
				score
			ORDER BY score DESC
			LIMIT $limit;
		`,
		map[string]interface{}{
			"text":  search,
			"limit": limit,
		},
	)
	if err != nil {
		return nil, err
	}
	words := []domain.Word{}
	for _, record := range result {
		words = append(words, *recordToWord(record.AsMap()))
	}
	return words, nil
}

func (r *wordRepository) FindPath(fromId string, toId string) ([]domain.Word, []domain.WordsLink, error) {
	result, err := r.Datasource.ExecRead(
		`MATCH
			(w1:Word {id: $fromId}),
			(w2:Word {id: $toId}),
			p = shortestPath((w1)-[:CONCERN*]-(w2))
		RETURN nodes(p) as nodes, relationships(p) as relationships`,
		map[string]interface{}{
			"fromId": fromId,
			"toId":   toId,
			// "userId": userId,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	if len(result) == 0 {
		return []domain.Word{}, []domain.WordsLink{}, nil
	}

	wordIds := map[int64]string{}
	words := []domain.Word{}
	links := []domain.WordsLink{}

	res := result[0].AsMap()
	for _, record := range res["nodes"].([]any) {
		n := record.(dbtype.Node)
		w := *recordToWord(n.GetProperties())
		wordIds[n.GetId()] = w.Id
		words = append(words, w)
	}
	for _, record := range res["relationships"].([]any) {
		r := record.(dbtype.Relationship)
		links = append(links, domain.WordsLink{
			SourceId: wordIds[r.StartId],
			TargetId: wordIds[r.EndId],
		})
	}
	return words, links, nil
}

func (r *wordRepository) StoreRelation(sourceId string, targetId string) error {
	// remove word and links
	_, err := r.Datasource.ExecWrite(
		`MATCH (w1:Word {id: $id1})
		MATCH (w2:Word {id: $id2})
		CREATE (w1)-[:CONCERN]->(w2);`,
		map[string]interface{}{
			"id1": sourceId,
			"id2": targetId,
		},
	)
	return err
}
