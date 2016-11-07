package lib

import (
	"github.com/deis/steward-framework"
)

// getLastOperationResponse is the response body from a get last operation call
type getLastOperationResponse struct {
	State string `json:"state"`
}

func (g *getLastOperationResponse) getFrameworkOperationStatusResponse() *framework.OperationStatusResponse {
	return &framework.OperationStatusResponse{
		State: g.State,
	}
}
