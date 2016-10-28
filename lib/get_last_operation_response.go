package lib

import (
	"github.com/deis/steward-framework"
)

// getLastOperationResponse is the response body from a get last operation call
type getLastOperationResponse struct {
	State string `json:"state"`
}

func (g *getLastOperationResponse) getFrameworkGetLastOperationResponse() *framework.GetLastOperationResponse {
	return &framework.GetLastOperationResponse{
		State: g.State,
	}
}
