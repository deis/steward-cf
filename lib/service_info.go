package lib

import (
	"github.com/deis/steward-framework"
)

// serviceInfo represents all of the information about a service except for its plans
type serviceInfo struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Description string `json:"description"`
	// Tags          []string         `json:"tags"`
	// Requires      []string         `json:"requires"`
	// Bindable      bool             `json:"bindable"`
	// Metadata      ServicesMetadata `json:"metadata"`
	PlanUpdatable bool `json:"plan_updateable"`
}

func (s *serviceInfo) getFrameworkServiceInfo() framework.ServiceInfo {
	return framework.ServiceInfo{
		Name:          s.Name,
		ID:            s.ID,
		Description:   s.Description,
		PlanUpdatable: s.PlanUpdatable,
	}
}
