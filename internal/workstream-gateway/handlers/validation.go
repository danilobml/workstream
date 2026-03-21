package handlers

import (
	"fmt"
	"net/http"

	"github.com/danilobml/workstream/internal/workstream-gateway/httputils"
	"github.com/go-playground/validator/v10"
)

func isInputValid(w http.ResponseWriter, structToValidate any) bool {
	validate := validator.New()
	err := validate.Struct(structToValidate)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		httputils.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", errors))
		return false
	}

	return true
}
