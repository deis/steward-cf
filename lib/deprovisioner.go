package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/deis/steward-framework"
)

type deprovisioner struct {
	cl *restClient
}

func newDeprovisioner(cl *restClient) *deprovisioner {
	return &deprovisioner{
		cl: cl,
	}
}

func (d *deprovisioner) Deprovision(
	ctx context.Context,
	brokerSpec framework.BrokerSpec,
	req *framework.DeprovisionRequest,
) (*framework.DeprovisionResponse, error) {

	query := url.Values(map[string][]string{})
	query.Add(serviceIDQueryKey, req.ServiceID)
	query.Add(planIDQueryKey, req.PlanID)
	query.Add(asyncQueryKey, "true")
	apiReq, err := d.cl.Delete(brokerSpec, query, "v2", "service_instances", req.InstanceID)
	if err != nil {
		return nil, err
	}
	apiRes, err := d.cl.Do(ctx, apiReq)
	if err != nil {
		return nil, err
	}
	defer apiRes.Body.Close()

	res := &deprovisionResponse{}
	switch apiRes.StatusCode {
	case http.StatusOK:
		res.IsAsync = false
	case http.StatusAccepted:
		res.IsAsync = true
	default:
		return nil, errUnexpectedResponseCode{
			URL:      apiReq.URL.String(),
			Expected: http.StatusOK,
			Actual:   apiRes.StatusCode,
		}
	}

	if err := json.NewDecoder(apiRes.Body).Decode(res); err != nil {
		return nil, err
	}
	return res.getFrameworkDeprovisionResponse(), nil
}
