package api

import "context"

// Server implements the handlers defined by the OpenAPI contract.
type Server struct{}

var _ StrictServerInterface = Server{}

// GetHealth reports that the API is running.
func (s Server) GetHealth(ctx context.Context, request GetHealthRequestObject) (GetHealthResponseObject, error) {
	return GetHealth200JSONResponse{Status: Ok}, nil
}
