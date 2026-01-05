package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJson[T any](w http.ResponseWriter, status int, data T) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        return fmt.Errorf("encode json: %w", err)
    }
    return nil
}
