package lib

import (
	"github.com/deis/steward-framework"
)

// servicePlan is the steward representation of a service plan. It's also compatible with the
// CloudFoundtry service plan object. See https://docs.cloudfoundry.org/services/api.html#PObject
// for more detail
type servicePlan struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// Metadata ServiceMetadata `json:"metadata"`
	Free bool `json:"free"`
}

func (s *servicePlan) getFrameworkServicePlan() framework.ServicePlan {
	return framework.ServicePlan{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Free:        s.Free,
	}
}
