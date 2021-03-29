package delivery

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/opbls/scapo/petstore/domain"
	"github.com/opbls/scapo/petstore/openapi"
)

// Validate Fields.
func validatePathParam(p openapi.FindPetsParams) error {
	if p.Limit != nil && *p.Limit < 0 {
		return domain.Err400BadRequest
	}
	return nil
}

// Validate Fields.
func validatePet(p domain.Pet) error {
	err := validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
	)
	if err != nil {
		return domain.Err400BadRequest
	}
	return nil
}
