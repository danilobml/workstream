package httputils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func WriteJson[T any](w http.ResponseWriter, status int, data T) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        return fmt.Errorf("encode json: %w", err)
    }
    return nil
}

// Forwards auth - authorization should be: "Bearer <token>"
func CtxWithAuth(ctx context.Context, authorization string) context.Context {
    md := metadata.Pairs("authorization", authorization)
    return metadata.NewOutgoingContext(ctx, md)
}

