package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/danilobml/workstream/internal/platform/dtos"
	"github.com/danilobml/workstream/internal/platform/httputils"
	"github.com/danilobml/workstream/internal/workstream-gateway/services/ports"
	"github.com/danilobml/workstream/internal/workstream-identity/helpers"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IdentityHandler struct {
	identityService ports.IdentityServicePort
	apiKey          string
}

func NewIdentityHandler(identityService ports.IdentityServicePort, apiKey string) *IdentityHandler {
	return &IdentityHandler{
		identityService: identityService,
		apiKey:          apiKey,
	}
}

func (ih *IdentityHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	registerReq := dtos.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !ih.isInputValid(w, registerReq) {
		return
	}

	registerReq.Password = strings.TrimSpace(registerReq.Password)
	registerReq.Email = strings.TrimSpace(registerReq.Email)

	resp, err := ih.identityService.Register(ctx, registerReq)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusCreated, resp)
}

func (ih *IdentityHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	loginReq := dtos.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !ih.isInputValid(w, loginReq) {
		return
	}

	loginReq.Password = strings.TrimSpace(loginReq.Password)
	loginReq.Email = strings.TrimSpace(loginReq.Email)

	resp, err := ih.identityService.Login(ctx, loginReq)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, resp)
}

func (ih *IdentityHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	ctx := httputils.CtxWithAuth(r.Context(), auth)

	users, err := ih.identityService.ListAllUsers(ctx)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, users)
}

func (ih *IdentityHandler) UnregisterUser(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	ctx := httputils.CtxWithAuth(r.Context(), auth)

	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "no valid user id supplied")
		return
	}

	err = ih.identityService.Unregister(ctx, dtos.UnregisterRequest{Id: userId})
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusNoContent, "unregistered")
}

func (ih *IdentityHandler) RemoveUser(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	ctx := httputils.CtxWithAuth(r.Context(), auth)

	idString := r.PathValue("id")
	userId, err := uuid.Parse(idString)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "no valid user id supplied")
		return
	}

	err = ih.identityService.RemoveUser(ctx, dtos.RemoveUserRequest{Id: userId})
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusNoContent, "removed")
}

/* 
func (ih *IdentityHandler) GetUserData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := ih.identityService.GetUserData(ctx)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, resp)
}

func (ih *IdentityHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "no valid user id supplied")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	updateReq := dtos.UpdateUserRequest{}
	err = json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !ih.isInputValid(w, updateReq) {
		return
	}

	updateReq.ID = userId
	updateReq.Email = strings.TrimSpace(updateReq.Email)

	err = ih.identityService.UpdateUserData(ctx, updateReq)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, "updated successfully")
}

func (ih *IdentityHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	requestPassResetReq := dtos.RequestPasswordResetRequest{}
	err := json.NewDecoder(r.Body).Decode(&requestPassResetReq)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !ih.isInputValid(w, requestPassResetReq) {
		return
	}

	requestPassResetReq.Email = strings.TrimSpace(requestPassResetReq.Email)

	err = ih.identityService.RequestPasswordReset(ctx, requestPassResetReq)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusNoContent, "")
}

func (ih *IdentityHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	resetPassReq := dtos.ResetPasswordRequest{}
	err := json.NewDecoder(r.Body).Decode(&resetPassReq)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !ih.isInputValid(w, resetPassReq) {
		return
	}

	resetPassReq.Password = strings.TrimSpace(resetPassReq.Password)
	resetPassReq.Email = strings.TrimSpace(resetPassReq.Email)
	resetPassReq.ResetToken = strings.TrimSpace(resetPassReq.ResetToken)

	err = ih.identityService.ResetPassword(ctx, resetPassReq)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusNoContent, "")
}

func (ih *IdentityHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	apiKey := strings.TrimSpace(r.Header.Get("User-Api-Key"))

	if apiKey != ih.apiKey {
		helpers.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	checkUserReq := dtos.CheckUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&checkUserReq)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !ih.isInputValid(w, checkUserReq) {
		return
	}

	resp, err := ih.identityService.CheckUser(ctx, checkUserReq)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, resp)
}
 */

// Validation Helper:
func (ih *IdentityHandler) isInputValid(w http.ResponseWriter, structToValidate any) bool {
	validate := validator.New()
	err := validate.Struct(structToValidate)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", errors))
		return false
	}

	return true
}
