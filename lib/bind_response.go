package lib

import (
	"github.com/deis/steward-framework"
)

// bindResponse represents a response to a binding request. It is marked with JSON struct tags so
// that it can be encoded to, and decoded from the CloudFoundry binding response body format.
// See https://docs.cloudfoundry.org/services/api.html#binding for more details
type bindResponse struct {
	Creds jsonObject `json:"credentials"`
}

func (b *bindResponse) getFrameworkBindResponse() *framework.BindResponse {
	return &framework.BindResponse{
		Creds: map[string]interface{}(b.Creds),
	}
}
