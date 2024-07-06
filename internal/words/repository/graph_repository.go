package repository

import (
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/s2dio-tech/mindgra-backend/common"
	"github.com/s2dio-tech/mindgra-backend/datasource"
	"github.com/s2dio-tech/mindgra-backend/domain"
)

type graphRepository struct {
	Datasource *datasource.Neo4J
}

func InitGraphRepository(db *datasource.Neo4J) domain.GraphRepository {
	return &graphRepository{
		Datasource: db,
	}
}

func recordToGraph(record map[string]any) *domain.Graph {
	w := domain.Graph{
		Id:     record["id"].(string),
		Name:   record["name"].(string),
		UserId: record["userId"].(string),
		Type:   record["type"].(string),
	}
	if record["createdAt"] != nil {
		w.CreatedAt = common.ToPointer(record["createdAt"].(neo4j.LocalDateTime).Time())
	}
	return &w
}

func (repo *graphRepository) Select(userId string) ([]domain.Graph, error) {
	res, err := repo.Datasource.ExecRead(`
		MATCH (u:User {id: $userId})-[:OWN]->(s:Graph {deleteFlag: false})
		RETURN s.id AS id, s.name AS name, s.userId as userId, s.type as type, s.createdAt as createdAt;`,
		map[string]interface{}{
			"userId": userId,
		},
	)

	if err != nil {
		return nil, err
	}

	graphs := []domain.Graph{}
	for _, record := range res {
		graphs = append(graphs, *recordToGraph(record.AsMap()))
	}
	return graphs, nil
}

func (repo *graphRepository) Store(w domain.Graph) (*string, error) {
	query := `MATCH (u:User {id: $userId})
		CREATE (w:Graph {
			id: apoc.create.uuid(),
			userId: $userId,
			name: $name,
			type: "3d",
			createdAt: $createdAt,
			deleteFlag: false
		})
		CREATE (u)-[:OWN]->(w)
		RETURN w.id AS id;`
	params := map[string]interface{}{
		"userId":    w.UserId,
		"name":      w.Name,
		"createdAt": neo4j.LocalDateTimeOf(time.Now()),
	}

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

func (repo *graphRepository) Update(id string, s domain.Graph) error {
	_set := "w.updatedAt= $updatedAt"
	params := map[string]interface{}{
		"id":        id,
		"updatedAt": neo4j.LocalDateTimeOf(time.Now()),
	}
	if s.Name != "" {
		_set += ", w.name= $name"
		params["name"] = s.Name
	}
	if s.Type != "" {
		_set += ", w.type= $type"
		params["type"] = s.Type
	}
	query := `MATCH (w:Graph {id: $id}) SET ` + _set

	_, err := repo.Datasource.ExecWrite(query, params)
	return err
}

func (r *graphRepository) Delete(id string) error {
	// remove graph and links
	_, err := r.Datasource.ExecWrite(
		`MATCH (s:Graph {id: $id})
			SET s.deleteFlag = true;`,
		map[string]interface{}{
			"id": id,
		},
	)
	return err
}

func (r *graphRepository) SelectOne(id string) (*domain.Graph, error) {
	result, err := r.Datasource.ExecRead(
		`MATCH (g:Graph {id: $id, deleteFlag: false})
			RETURN g.id as id,
				g.userId AS userId,
				g.name AS name,
				g.type AS type,
				g.createdAt AS createdAt;`,
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

	return recordToGraph(record.AsMap()), nil
}
