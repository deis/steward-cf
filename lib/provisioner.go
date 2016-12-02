package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/deis/steward-framework"
)

type provisioner struct {
	cl *restClient
}

func newProvisioner(cl *restClient) *provisioner {
	return &provisioner{
		cl: cl,
	}
}

func (p *provisioner) Provision(
	ctx context.Context,
	brokerSpec framework.ServiceBrokerSpec,
	req *framework.ProvisionRequest,
) (*framework.ProvisionResponse, error) {

	query := url.Values(map[string][]string{})
	query.Add(asyncQueryKey, "true")
	bodyBytes := new(bytes.Buffer)
	if err := json.NewEncoder(bodyBytes).Encode(getAPIProvisionRequest(req)); err != nil {
		return nil, err
	}
	apiReq, err := p.cl.Put(
		brokerSpec,
		query,
		bodyBytes,
		"v2",
		"service_instances",
		req.InstanceID,
	)
	if err != nil {
		return nil, err
	}
	apiRes, err := p.cl.Do(ctx, apiReq)
	if err != nil {
		return nil, err
	}

	res := &provisionResponse{}
	switch apiRes.StatusCode {
	case http.StatusOK, http.StatusCreated:
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
	return res.getFrameworkProvisionResponse(), nil
}
