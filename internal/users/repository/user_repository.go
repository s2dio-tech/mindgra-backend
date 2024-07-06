package repository

import (
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/s2dio-tech/mindgra-backend/common"
	"github.com/s2dio-tech/mindgra-backend/datasource"
	"github.com/s2dio-tech/mindgra-backend/domain"
)

type userRepository struct {
	Datasource *datasource.Neo4J
}

func InitUserRepository(db *datasource.Neo4J) domain.UserRepository {
	return &userRepository{
		Datasource: db,
	}
}

func (repo *userRepository) Create(u *domain.User) (*string, error) {
	user := *u
	user.CreatedAt = time.Now()

	session := repo.Datasource.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `CREATE (u:User {
			id: apoc.create.uuid(),
			name: $name,
			email: $email,
			password: $password,
			role: $role,
			createdAt: $createdAt
		})
		RETURN u.id as id;`
		parameters := map[string]interface{}{
			"name":      user.Name,
			"email":     user.Email,
			"password":  user.Password,
			"role":      user.Role,
			"createdAt": neo4j.LocalDateTimeOf(user.CreatedAt),
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

	return common.Nullable{Value: result}.ToStringPtr(), nil
}

func (repo *userRepository) FindById(id string) (user *domain.User, err error) {

	session := repo.Datasource.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})

	defer session.Close()

	result, err := session.Run(
		`MATCH (u:User {id: $id})
		RETURN u.name AS name,
			u.email AS email,
			u.role AS role;`,
		map[string]interface{}{
			"id": id,
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
	email, _ := record.Get("email")
	name, _ := record.Get("name")
	role, _ := record.Get("role")
	return &domain.User{
		Id:    id,
		Email: email.(string),
		Name:  name.(string),
		Role:  domain.RoleMap[role.(string)],
	}, nil
}

func (repo *userRepository) FindByEmail(email string) (*domain.User, error) {
	session := repo.Datasource.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})

	defer session.Close()

	result, err := session.Run(
		`MATCH (u:User {email: $email})
			RETURN u.id AS id,
				u.name AS name,
				u.email AS email,
				u.password AS password,
				u.role AS role;`,
		map[string]interface{}{
			"email": email,
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
	name, _ := record.Get("name")
	mail, _ := record.Get("email")
	pass, _ := record.Get("password")
	role, _ := record.Get("role")
	return &domain.User{
		Id:       id.(string),
		Name:     name.(string),
		Email:    mail.(string),
		Password: pass.(string),
		Role:     domain.RoleMap[role.(string)],
	}, nil
}

func (repo *userRepository) Update(id string, data map[string]interface{}) error {
	session := repo.Datasource.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `MATCH (u:User {id: $id})
			SET u.password = $password;`
		parameters := map[string]interface{}{
			"id":        id,
			"password":  data["password"],
			"updatedAt": neo4j.LocalDateTimeOf(time.Now()),
		}
		res, err := tx.Run(query, parameters)
		if err != nil || !res.Next() {
			return nil, err
		}
		return nil, nil
	})

	return err
}
