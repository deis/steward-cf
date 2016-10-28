package lib

import (
	"github.com/deis/steward-framework"
)

// service is the represntation of a steward service. It also is compatible with the CloudFoundry
// catalog API. See https://docs.cloudfoundry.org/services/api.html#catalog-mgmt for more detail
type service struct {
	serviceInfo
	Plans []servicePlan `json:"plans"`
}

func (s *service) getFrameworkService() *framework.Service {
	plans := make([]framework.ServicePlan, len(s.Plans))
	for i, plan := range s.Plans {
		plans[i] = plan.getFrameworkServicePlan()
	}
	service := &framework.Service{
		ServiceInfo: s.serviceInfo.getFrameworkServiceInfo(),
		Plans:       plans,
	}
	return service
}
