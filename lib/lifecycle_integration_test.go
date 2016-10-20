// +build integration

package lib

import (
	"context"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward-framework"
)

const (
	// cf-sample-broker isn't picky about inputs:
	fakeInstanceID = "fake-instance-id"
	fakeServiceID  = "fake-service-id"
	fakePlanID     = "fake-plan-id"
	fakeBindingID  = "fake-binding-id"
)

func TestProvision(t *testing.T) {
	resp, err := testLifecycler.Provision(context.Background(), &framework.ProvisionRequest{
		InstanceID:        fakeInstanceID,
		ServiceID:         fakeServiceID,
		PlanID:            fakePlanID,
		AcceptsIncomplete: true,
	})
	assert.NoErr(t, err)
	// Compare to known results from cf-sample-broker...
	assert.Equal(t, resp, &framework.ProvisionResponse{
		Operation: "create",
	}, "provision response")
}

func TestBind(t *testing.T) {
	resp, err := testLifecycler.Bind(context.Background(), &framework.BindRequest{
		InstanceID: fakeInstanceID,
		ServiceID:  fakeServiceID,
		PlanID:     fakePlanID,
		BindingID:  fakeBindingID,
	})
	assert.NoErr(t, err)
	// Compare to known results from cf-sample-broker...
	assert.Equal(t, len(resp.Creds), 10, "credentials count")
}

func TestUnbind(t *testing.T) {
	err := testLifecycler.Unbind(context.Background(), &framework.UnbindRequest{
		InstanceID: fakeInstanceID,
		ServiceID:  fakeServiceID,
		PlanID:     fakePlanID,
		BindingID:  fakeBindingID,
	})
	// Unbind returns no result except for any error that occurred...
	assert.NoErr(t, err)
}

func TestDeprovision(t *testing.T) {
	resp, err := testLifecycler.Deprovision(context.Background(), &framework.DeprovisionRequest{
		InstanceID:        fakeInstanceID,
		ServiceID:         fakeServiceID,
		PlanID:            fakePlanID,
		AcceptsIncomplete: true,
	})
	assert.NoErr(t, err)
	// Compare to known results from cf-sample-broker...
	assert.Equal(t, resp, &framework.DeprovisionResponse{
		Operation: "",
	}, "deprovision response")
}
