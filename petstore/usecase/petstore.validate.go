package usecase

import (
	"github.com/opbls/scapo/petstore/domain"
)

// ValidatePathParam validate Request Parameter PetID.
func validatePathParamPetID(id int) error {
	// open api
	if id < 0 {
		return domain.Err400BadRequest
	}
	return nil
}
