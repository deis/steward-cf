package lib

import (
	"github.com/deis/steward-framework"
)

// deprovisionResponse represents a response to a provisioning request. It is marked with JSON
// struct tags so that it can be encoded to, and decoded from the CloudFoundry deprovisioning
// response body format. See https://docs.cloudfoundry.org/services/api.html#deprovisioning for
// more details
type deprovisionResponse struct {
	Operation string `json:"operation"`
	IsAsync   bool   `json:"-"`
}

func (d *deprovisionResponse) getFrameworkDeprovisionResponse() *framework.DeprovisionResponse {
	return &framework.DeprovisionResponse{
		Operation: d.Operation,
		IsAsync:   d.IsAsync,
	}
}
