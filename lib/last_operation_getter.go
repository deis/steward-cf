package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/deis/steward-framework"
)

type lastOperationGetter struct {
	cl *restClient
}

func newLastOperationGetter(cl *restClient) *lastOperationGetter {
	return &lastOperationGetter{
		cl: cl,
	}
}

func (l *lastOperationGetter) GetLastOperation(
	ctx context.Context,
	req *framework.GetLastOperationRequest,
) (*framework.GetLastOperationResponse, error) {

	query := url.Values(map[string][]string{})
	query.Add(serviceIDQueryKey, req.ServiceID)
	query.Add(planIDQueryKey, req.PlanID)
	query.Add(operationQueryKey, req.Operation)
	apiReq, err := l.cl.Get(query, "v2", "service_instances", req.InstanceID, "last_operation")
	if err != nil {
		return nil, err
	}
	apiRes, err := l.cl.Do(ctx, apiReq)
	if err != nil {
		return nil, err
	}
	defer apiRes.Body.Close()
	// An HTTP response code of 410 (gone) is a distinct state that deprovision may wish to
	// interpret as success.
	if apiRes.StatusCode == http.StatusGone {
		return &framework.GetLastOperationResponse{
			State: framework.LastOperationStateGone.String(),
		}, nil
	}
	res := new(framework.GetLastOperationResponse)
	if err := json.NewDecoder(apiRes.Body).Decode(res); err != nil {
		return nil, err
	}
	return res, nil
}
