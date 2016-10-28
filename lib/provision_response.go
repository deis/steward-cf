package lib

import (
	"github.com/deis/steward-framework"
)

// provisionResponse represents a response to a provisioning request. It is marked with JSON struct
// tags so that it can be encoded to, and decoded from the CloudFoundry provisioning response body
// format. See https://docs.cloudfoundry.org/services/api.html#provisioning for more details
type provisionResponse struct {
	Operation string     `json:"operation"`
	IsAsync   bool       `json:"-"`
	Extra     jsonObject `json:"extra,omitempty"`
}

func (p *provisionResponse) getFrameworkProvisionResponse() *framework.ProvisionResponse {
	return &framework.ProvisionResponse{
		Operation: p.Operation,
		IsAsync:   p.IsAsync,
		Extra:     map[string]interface{}(p.Extra),
	}
}
