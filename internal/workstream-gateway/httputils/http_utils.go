package httputils

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/danilobml/workstream/internal/platform/errs"
	"google.golang.org/grpc/metadata"
)

// Forwards auth - authorization should be: "Bearer <token>"
func CtxWithAuth(ctx context.Context, authorization string) context.Context {
    md := metadata.Pairs("authorization", authorization)
    return metadata.NewOutgoingContext(ctx, md)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}

func WriteJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteErrorsResponse(w http.ResponseWriter, err error) {
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			WriteJSONError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, errs.ErrAlreadyExists) {
			WriteJSONError(w, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, errs.ErrInvalidCredentials) {
			WriteJSONError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, errs.ErrUnauthorized) {
			WriteJSONError(w, http.StatusUnauthorized, err.Error())
			return
		}

		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}


