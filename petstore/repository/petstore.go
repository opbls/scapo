package repository

import (
	"bytes"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/opbls/scapo/petstore/domain"
)

type (
	// PetStoreRepository interface.
	PetStoreRepository interface {
		QueryPets(condition *domain.QueryCondition) (*domain.Pets, error)
		QueryPet(id int) (*domain.Pet, error)
		CreatePet(pet *domain.Pet) (*domain.Pet, error)
		DeletePet(id int) (int, error)
	}

	// PetStoreRepositoryImpl struct.
	PetStoreRepositoryImpl struct {
		DB *sqlx.DB
	}
)

// NewPetStoreRepository instantiate PetStoreRepository.
func NewPetStoreRepository(db *sqlx.DB) PetStoreRepository {
	return &PetStoreRepositoryImpl{
		DB: db,
	}
}

// QueryPets return Pets from db.
func (impl PetStoreRepositoryImpl) QueryPets(condition *domain.QueryCondition) (*domain.Pets, error) {
	/*
		SELECT id, name, tag FROM petstore WHERE tag IN ('foo', 'bar') LIMIT 10;
	*/

	// build sql
	SQL := `SELECT id, name, tag FROM petstore `
	if _, ok := (*condition)["tags"]; ok {
		SQL += `WHERE tag IN (:tags) `
	}
	SQL += `LIMIT :limit`

	// build bind parameter
	query, binds, err := sqlx.Named(SQL, asMap(condition))
	if err != nil {
		return nil, domain.Err500InternalServerError
	}
	query, binds, err = sqlx.In(query, binds...)
	if err != nil {
		return nil, domain.Err500InternalServerError
	}
	query = impl.DB.Rebind(query)

	// access db
	rslts := domain.Pets{}
	err = impl.DB.Select(&rslts, query, binds...)
	if err != nil {
		return nil, domain.Err500InternalServerError
	}

	return &rslts, nil
}

// QueryPet return Pet from db.
func (impl PetStoreRepositoryImpl) QueryPet(id int) (*domain.Pet, error) {
	/*
		SELECT id, name, tag FROM petstore WHERE id = 1 LIMIT 1;
	*/

	// build sql
	SQL := `SELECT id, name, tag FROM petstore WHERE id = :id LIMIT 1`

	// access db
	rows, err := impl.DB.Queryx(SQL, id)
	if err != nil {
		return nil, domain.Err500InternalServerError
	}
	defer rows.Close()

	if rows.Next() {
		rslt := domain.Pet{}
		rows.StructScan(&rslt)
		return &rslt, nil
	}

	return nil, nil
}

// CreatePet provide Pet to db.
func (impl PetStoreRepositoryImpl) CreatePet(p *domain.Pet) (*domain.Pet, error) {
	/*
		INSERT INTO petstore(name, tag) VALUES('foo', 'bar');
	*/

	SQL := `INSERT INTO petstore(name, tag) VALUES(:name, :tag)`

	// access db
	stmt, err := impl.DB.Preparex(SQL)
	if err != nil {
		return nil, domain.Err500InternalServerError
	}
	defer stmt.Close()

	rslt, err := stmt.Exec(p.Name, p.Tag)
	if err != nil {
		return nil, domain.Err500InternalServerError
	}
	i, err := rslt.LastInsertId()
	if err != nil {
		return nil, domain.Err500InternalServerError
	}

	p.Id = i

	return p, nil
}

// DeletePet delete Pet from db.
func (impl PetStoreRepositoryImpl) DeletePet(id int) (int, error) {
	/*
		DELETE FROM petstore WHERE id = 0
	*/

	notaffected := -1
	SQL := `DELETE FROM petstore WHERE id = :id`

	// access db
	stmt, err := impl.DB.Preparex(SQL)
	if err != nil {
		return notaffected, domain.Err500InternalServerError
	}
	defer stmt.Close()

	rslt, err := stmt.Exec(id)
	if err != nil {
		return notaffected, domain.Err500InternalServerError
	}

	i, err := rslt.RowsAffected()
	if err != nil {
		return notaffected, domain.Err500InternalServerError
	}

	return int(i), nil
}

// asMap cast QueryCondition to map[string]interface{}.
func asMap(object *domain.QueryCondition) map[string]interface{} {
	var i interface{}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(object)
	json.Unmarshal(b.Bytes(), &i)
	ret, _ := i.(map[string]interface{})
	return ret
}
