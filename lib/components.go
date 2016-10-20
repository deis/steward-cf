package lib

import (
	"github.com/deis/steward-framework"
)

func GetComponents() (framework.Cataloger, framework.Lifecycler, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, nil, err
	}
	cl := newRESTClient(cfg)
	return newCataloger(cl),
		newLifecycler(cl),
		nil
}
