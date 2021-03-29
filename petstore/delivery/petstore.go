package delivery

import (
	"encoding/json"
	"log"
	"net/http"

	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/opbls/scapo/petstore/domain"
	"github.com/opbls/scapo/petstore/openapi"
	"github.com/opbls/scapo/petstore/usecase"
)

type (
	// PetStoreDelivery interface.
	PetStoreDelivery openapi.ServerInterface

	// PetStoreDeliveryImpl struct.
	PetStoreDeliveryImpl struct {
		Usecase usecase.PetStoreUsecase
	}
)

// NewPetStoreDelivery returns Petstore ServerInterface.
func NewPetStoreDelivery(usecase usecase.PetStoreUsecase) PetStoreDelivery {
	return &PetStoreDeliveryImpl{
		Usecase: usecase,
	}
}

// FindPets Impl.
// OpenAPI 3 defines default serialization method
//  - style: form
//  - explode: true
// and default behavior of reserved chars are percent-encoded
//  - allowReserved: false
func (impl *PetStoreDeliveryImpl) FindPets(w http.ResponseWriter, r *http.Request, params openapi.FindPetsParams) {

	// validate
	if err := validatePathParam(params); err != nil {
		writeError(w, err)
		return
	}

	condition := domain.QueryCondition{}
	if params.Tags != nil && len(*params.Tags) > 0 {
		condition["tags"] = *params.Tags
	}
	if params.Limit == nil {
		condition["limit"] = 100
	} else {
		condition["limit"] = params.Limit
	}

	pets, err := impl.Usecase.FindPets(&condition)
	if err != nil {
		writeError(w, err)
		return
	}

	write200OK(w, pets)
}

// AddPet Impl.
func (impl *PetStoreDeliveryImpl) AddPet(w http.ResponseWriter, r *http.Request) {

	np := domain.Pet{}
	if err := json.NewDecoder(r.Body).Decode(&np); err != nil {
		writeError(w, err)
		return
	}

	// validate
	if err := validatePet(np); err != nil {
		writeError(w, err)
		return
	}

	p, err := impl.Usecase.AddPet(&np)
	if err != nil {
		writeError(w, err)
		return
	}

	write200OK(w, p)
}

// DeletePet Impl
func (impl *PetStoreDeliveryImpl) DeletePet(w http.ResponseWriter, r *http.Request, id int64) {

	did := int(id)

	i, err := impl.Usecase.DeletePet(did)
	if err != nil {
		writeError(w, err)
		return
	}
	log.Println("delete success affected row: ", i)

	//act as not found
	if i == 0 {
		writeError(w, domain.Err404NotFound)
		return
	}

	write204NoContent(w)

}

// FindPetById Impl.
func (impl *PetStoreDeliveryImpl) FindPetById(w http.ResponseWriter, r *http.Request, id int64) {

	fid := int(id)

	rslt, err := impl.Usecase.FindPetById(fid)
	if err != nil {
		writeError(w, err)
		return
	}

	if rslt == nil {
		// write204NoContent(w)
		writeError(w, domain.Err404NotFound)
		return
	}
	// response
	write200OK(w, &rslt)
}

func write200OK(w http.ResponseWriter, objects interface{}) {
	writeSuccess(w, http.StatusOK, objects)
}

func write201(w http.ResponseWriter, objects interface{}) {
	writeSuccess(w, http.StatusCreated, objects)
}

func write204NoContent(w http.ResponseWriter) {
	writeSuccess(w, http.StatusNoContent, nil)
}

func writeSuccess(w http.ResponseWriter, code int, objects interface{}) {
	w.WriteHeader(code)
	if objects != nil {
		writer := json.NewEncoder(w)
		writer.Encode(objects)
	}
}

func writeError(w http.ResponseWriter, err error) {
	code := getStatusCode(err)
	commonError := openapi.Error{
		Code:    int32(code),
		Message: err.Error(),
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(commonError)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	log.Println(err)
	switch err {
	case domain.Err500InternalServerError:
		return http.StatusInternalServerError
	case domain.Err400BadRequest:
		return http.StatusBadRequest
	case domain.Err404NotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
