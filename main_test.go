package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deepmap/oapi-codegen/examples/petstore-expanded/chi/api"
	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/deepmap/oapi-codegen/pkg/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/opbls/scapo/petstore/delivery"
	"github.com/opbls/scapo/petstore/openapi"
	"github.com/opbls/scapo/petstore/repository"
	"github.com/opbls/scapo/petstore/usecase"
	"github.com/stretchr/testify/assert"
)

// add some comment.
func TestHandler(t *testing.T) {
	var err error

	r := chi.NewRouter()
	swagger, _ := openapi.GetSwagger()
	swagger.Servers = nil
	r.Use(middleware.OapiRequestValidator(swagger))

	// database
	ddl := `CREATE TABLE IF NOT EXISTS petstore(
		id integer PRIMARY KEY autoincrement
		, name text NOT NULL
		, tag text
	);`

	db, _ := sqlx.Connect("sqlite3", ":memory:")
	//db, _ := sqlx.Connect("sqlite3", "test.db")
	defer db.Close()

	db.MustExec(ddl)

	// handlers
	repo := repository.NewPetStoreRepository(db)
	usecase := usecase.NewPetStoreUsecase(repo)
	handler := delivery.NewPetStoreDelivery(usecase)
	openapi.HandlerFromMux(handler, r)

	//////////////////
	// TEST DATA
	//////////////////
	testname := "testname"
	testtag := "testtag"

	petData := []openapi.Pet{}
	newPetData := []openapi.NewPet{}
	imax := 10
	for i := 1; i <= imax; i++ {
		name := fmt.Sprintf(testname+"%d", i)
		tag := fmt.Sprintf(testtag+"%d", i%5)
		newPetData = append(newPetData, popNewPet(name, tag))
		petData = append(petData, popPet(i, name, tag))

		dml := `insert into petstore(name, tag) values("` + name + `", "` + tag + `");`
		db.MustExec(dml)
	}
	////////////////////
	// TEST
	////////////////////
	//	AddPet
	//	DeletePet
	//	FindPetById

	//////////
	//	FindPets
	//////////
	t.Run("SUCCESS_FindPets", func(t *testing.T) {
		var rp []openapi.Pet

		url := fmt.Sprintf("/pets")
		rr := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusOK, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&rp)
		assert.NoError(t, err, "error unmarshal response")
		assert.Equal(t, len(petData), len(rp))
	})

	t.Run("SUCCESS_FindPets_Tags", func(t *testing.T) {
		var rp []openapi.Pet

		url := fmt.Sprintf("/pets?tags=%s1", testtag)
		rr := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusOK, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&rp)
		assert.NoError(t, err, "error unmarshal response")
		assert.Equal(t, 2, len(rp))
	})

	t.Run("SUCCESS_FindPets_Tags_Tags", func(t *testing.T) {
		var rp []openapi.Pet

		url := fmt.Sprintf("/pets?tags=%s&tags=%s", testtag+"1", testtag+"3")
		rr := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusOK, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&rp)
		assert.NoError(t, err, "error unmarshal response")
		assert.Equal(t, 4, len(rp))
	})

	t.Run("SUCCESS_FindPets_Limit", func(t *testing.T) {
		var rp []openapi.Pet

		url := fmt.Sprintf("/pets?limit=%d", 5)
		rr := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusOK, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&rp)
		assert.NoError(t, err, "error unmarshal response")
		assert.Equal(t, 5, len(rp))
	})

	t.Run("SUCCESS_FindPets_Limit_Tags", func(t *testing.T) {
		var rp []openapi.Pet

		url := fmt.Sprintf("/pets?limit=%d&Tags=%s", 1, testtag+"1")
		rr := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusOK, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&rp)
		assert.NoError(t, err, "error unmarshal response")
		assert.Equal(t, 1, len(rp))
	})

	t.Run("SUCCESS_FindPets_Tags_Limit_Tags", func(t *testing.T) {
		var rp []openapi.Pet

		url := fmt.Sprintf("/pets?Tags=%s&limit=%d&Tags=%s", testtag+"1", 3, testtag+"3")
		rr := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusOK, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&rp)
		assert.NoError(t, err, "error unmarshal response")
		assert.Equal(t, 3, len(rp))
	})

	// abnormal 200
	t.Run("ABNORMAL_FindPets_Limit0", func(t *testing.T) {
		var rp []openapi.Pet

		url := fmt.Sprintf("/pets?limit=%d", 0)
		rr := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusOK, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&rp)
		assert.NoError(t, err, "error unmarshal response")
		assert.Equal(t, 0, len(rp))
	})

	// abnormal 200
	t.Run("ABNORMAL_FindPets_Tags0", func(t *testing.T) {
		var rp []openapi.Pet

		url := fmt.Sprintf("/pets?tags=%s", "foo")
		rr := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusOK, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&rp)
		assert.NoError(t, err, "error unmarshal response")
		assert.Equal(t, 0, len(rp))
	})

	//////////
	//	FindPetById
	//////////
	t.Run("SUCCESS_FindPetById", func(t *testing.T) {
		p := petData[0]
		var rp openapi.Pet

		url := fmt.Sprintf("/pets/%d", p.Id)
		rr := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder

		err = json.NewDecoder(rr.Body).Decode(&rp)
		assert.NoError(t, err, "error unmarshal response")
		assert.Equal(t, p, rp)
	})
	// abnormal 400
	t.Run("ABNORMAL_FindPetById_Nagative", func(t *testing.T) {
		var rp openapi.Pet

		url := fmt.Sprintf("/pets/%d", -1)
		rr := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&rp)
		assert.NoError(t, err, "error unmarshal response")
	})
	// abnormal 404
	t.Run("ABNORMAL_FindPetById_NotExist", func(t *testing.T) {
		var rp openapi.Pet

		url := fmt.Sprintf("/pets/%d", 1000000)
		rr := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusNotFound, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&rp)
		assert.NoError(t, err, "error unmarshal response")
	})

	//////////
	//  AddPet
	//////////
	t.Run("SUCEESS_AddPet", func(t *testing.T) {
		np := popNewPet("foo", "bar")
		p := popPet(1+len(petData), np.Name, *np.Tag)
		petData = append(petData, p)
		var rp openapi.Pet

		url := fmt.Sprintf("/pets")
		rr := testutil.NewRequest().Post(url).WithJsonBody(np).GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusOK, rr.Code)

		json.NewDecoder(rr.Body).Decode(&rp)
		assert.NoError(t, err, "error unmarshal response")
		assert.Equal(t, np.Name, rp.Name)
		assert.Equal(t, *rp.Tag, *rp.Tag)
	})
	//////////
	//	DeletePet
	//////////
	t.Run("SUCCESS_DeletePets", func(t *testing.T) {

		url := fmt.Sprintf("/pets/%d", -1+len(petData))
		rr := testutil.NewRequest().Delete(url).GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusNoContent, rr.Code)

		// slice update
		petData = petData[:-1+len(petData)]

		var rp []api.Pet
		url = fmt.Sprintf("/pets")
		rr = testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusOK, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&rp)
		assert.NoError(t, err, "error unmarshal response")
		assert.Equal(t, len(petData), len(rp))

	})
	// abnormal 400
	t.Run("ABNORMAL_DeletePets_Negative", func(t *testing.T) {

		url := fmt.Sprintf("/pets/%d", -1)
		rr := testutil.NewRequest().Delete(url).GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusBadRequest, rr.Code)

	})
	// abnormal 404
	t.Run("ABNORMAL_DeletePets_NotFound", func(t *testing.T) {

		url := fmt.Sprintf("/pets/%d", 10000)
		rr := testutil.NewRequest().Delete(url).GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusNotFound, rr.Code)

	})
}

func doGet(t *testing.T, mux *chi.Mux, url string) *httptest.ResponseRecorder {
	response := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, mux)
	return response.Recorder
}

func popNewPet(name string, tag string) openapi.NewPet {
	return openapi.NewPet{Name: name, Tag: &tag}
}

func popPet(id int, name string, tag string) openapi.Pet {
	ret := openapi.Pet{}
	ret.Id = int64(id)
	ret.Name = name
	ret.Tag = &tag
	return ret
}
