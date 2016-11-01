package lib

import (
	"github.com/deis/steward-framework"
)

// GetComponents returns a cataloger and lifecycler based on given configuration. Returns
// nil, nil and a non-nil error if any step in config or component creation failed.
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
