package lib

import (
	"context"
	"net/http"
	"net/url"

	"github.com/deis/steward-framework"
)

type unbinder struct {
	cl *restClient
}

func newUnbinder(cl *restClient) *unbinder {
	return &unbinder{
		cl: cl,
	}
}

func (u *unbinder) Unbind(
	ctx context.Context,
	brokerSpec framework.BrokerSpec,
	req *framework.UnbindRequest,
) error {

	query := url.Values(map[string][]string{})
	query.Add(serviceIDQueryKey, req.ServiceID)
	query.Add(planIDQueryKey, req.PlanID)
	apiReq, err := u.cl.Delete(
		brokerSpec,
		query,
		"v2",
		"service_instances",
		req.InstanceID,
		"service_bindings",
		req.BindingID,
	)
	if err != nil {
		return err
	}
	apiRes, err := u.cl.Do(ctx, apiReq)
	if err != nil {
		return err
	}
	if apiRes.StatusCode != http.StatusOK {
		return errUnexpectedResponseCode{
			URL:      apiReq.URL.String(),
			Expected: http.StatusOK,
			Actual:   apiRes.StatusCode,
		}
	}
	return nil
}
