package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/deis/steward-framework"
)

type binder struct {
	cl *restClient
}

func newBinder(cl *restClient) *binder {
	return &binder{
		cl: cl,
	}
}

func (b *binder) Bind(
	ctx context.Context,
	req *framework.BindRequest,
) (*framework.BindResponse, error) {

	bodyBytes := new(bytes.Buffer)
	if err := json.NewEncoder(bodyBytes).Encode(req); err != nil {
		return nil, err
	}

	apiReq, err := b.cl.Put(
		emptyQuery,
		bodyBytes,
		"v2",
		"service_instances",
		req.InstanceID,
		"service_bindings",
		req.BindingID,
	)
	if err != nil {
		return nil, err
	}

	apiRes, err := b.cl.Do(ctx, apiReq)
	if err != nil {
		return nil, err
	}
	defer apiRes.Body.Close()
	if apiRes.StatusCode != http.StatusOK && apiRes.StatusCode != http.StatusCreated {
		return nil, errUnexpectedResponseCode{
			URL:      apiReq.URL.String(),
			Expected: http.StatusOK,
			Actual:   apiRes.StatusCode,
		}
	}

	res := new(framework.BindResponse)
	if err := json.NewDecoder(apiRes.Body).Decode(res); err != nil {
		return nil, err
	}
	logger.Debugf("got response %+v from backing broker", *res)
	return res, nil
}
