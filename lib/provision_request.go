package lib

import (
	"github.com/deis/steward-framework"
	"github.com/pborman/uuid"
)

// provisionRequest represents a request to do a service provision operation. This struct is
// JSON-compatible with the request body detailed at
// https://docs.cloudfoundry.org/services/api.html#provisioning
type provisionRequest struct {
	OrganizationGUID  string     `json:"organization_guid"`
	PlanID            string     `json:"plan_id"`
	ServiceID         string     `json:"service_id"`
	SpaceGUID         string     `json:"space_guid"`
	AcceptsIncomplete bool       `json:"accepts_incomplete"`
	Parameters        jsonObject `json:"parameters"`
}

func getAPIProvisionRequest(req *framework.ProvisionRequest) *provisionRequest {
	return &provisionRequest{
		OrganizationGUID:  uuid.New(),
		PlanID:            req.PlanID,
		ServiceID:         req.ServiceID,
		SpaceGUID:         uuid.New(),
		AcceptsIncomplete: true,
		Parameters:        jsonObject(req.Parameters),
	}
}
