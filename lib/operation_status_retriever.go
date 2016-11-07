package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/deis/steward-framework"
)

type operationStatusRetriever struct {
	cl *restClient
}

func newOperationStatusRetriever(cl *restClient) *operationStatusRetriever {
	return &operationStatusRetriever{
		cl: cl,
	}
}

func (o *operationStatusRetriever) GetOperationStatus(
	ctx context.Context,
	req *framework.OperationStatusRequest,
) (*framework.OperationStatusResponse, error) {

	query := url.Values(map[string][]string{})
	query.Add(serviceIDQueryKey, req.ServiceID)
	query.Add(planIDQueryKey, req.PlanID)
	query.Add(operationQueryKey, req.Operation)
	apiReq, err := o.cl.Get(query, "v2", "service_instances", req.InstanceID, "last_operation")
	if err != nil {
		return nil, err
	}
	apiRes, err := o.cl.Do(ctx, apiReq)
	if err != nil {
		return nil, err
	}
	defer apiRes.Body.Close()
	// An HTTP response code of 410 (gone) is a distinct state that deprovision may wish to
	// interpret as success.
	if apiRes.StatusCode == http.StatusGone {
		return &framework.OperationStatusResponse{
			State: framework.OperationStateGone.String(),
		}, nil
	}
	res := &getLastOperationResponse{}
	if err := json.NewDecoder(apiRes.Body).Decode(res); err != nil {
		return nil, err
	}
	return res.getFrameworkOperationStatusResponse(), nil
}
