package lib

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/deis/steward-framework"
)

type cataloger struct {
	cl *restClient
}

func newCataloger(cl *restClient) *cataloger {
	return &cataloger{
		cl: cl,
	}
}

func (c *cataloger) List(ctx context.Context) ([]*framework.Service, error) {
	req, err := c.cl.Get(emptyQuery, "v2", "catalog")
	if err != nil {
		return nil, err
	}
	res, err := c.cl.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	serviceList := &serviceList{}
	// TODO: drain the response body to avoid a connection leak
	if err := json.NewDecoder(res.Body).Decode(serviceList); err != nil {
		logger.Debugf("error decoding JSON response body from backend CF broker (%s)", err)
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errUnexpectedResponseCode{URL: req.URL.String(), Expected: http.StatusOK, Actual: res.StatusCode}
	}
	return serviceList.getFrameworkServices(), nil
}
