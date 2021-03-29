package usecase

import (
	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/opbls/scapo/petstore/domain"
	"github.com/opbls/scapo/petstore/repository"
)

type (
	// PetStoreUsecase interface.
	PetStoreUsecase interface {
		FindPets(condition *domain.QueryCondition) (*domain.Pets, error)
		AddPet(np *domain.Pet) (*domain.Pet, error)
		DeletePet(id int) (int, error)
		FindPetById(id int) (*domain.Pet, error)
	}

	// PetStoreUsecaseImpl impl.
	PetStoreUsecaseImpl struct {
		Repository repository.PetStoreRepository
	}
)

// NewPetStoreUsecase returns Petstore Usecase.
func NewPetStoreUsecase(repo repository.PetStoreRepository) PetStoreUsecase {
	return &PetStoreUsecaseImpl{
		Repository: repo,
	}
}

// FindPets Impl.
func (impl *PetStoreUsecaseImpl) FindPets(condition *domain.QueryCondition) (*domain.Pets, error) {
	return impl.Repository.QueryPets(condition)
}

// AddPet Impl.
func (impl *PetStoreUsecaseImpl) AddPet(np *domain.Pet) (*domain.Pet, error) {
	return impl.Repository.CreatePet(np)
}

// DeletePet Impl
func (impl *PetStoreUsecaseImpl) DeletePet(id int) (int, error) {
	// validate
	if err := validatePathParamPetID(id); err != nil {
		return -1, domain.Err400BadRequest
	}

	return impl.Repository.DeletePet(id)
}

// FindPetById Impl.
func (impl *PetStoreUsecaseImpl) FindPetById(id int) (*domain.Pet, error) {
	// validate
	if err := validatePathParamPetID(id); err != nil {
		return nil, domain.Err400BadRequest
	}

	return impl.Repository.QueryPet(id)
}
