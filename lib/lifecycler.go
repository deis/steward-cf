package lib

type lifecycler struct {
	*provisioner
	*deprovisioner
	*binder
	*unbinder
	*operationStatusRetriever
}

func newLifecycler(cl *restClient) *lifecycler {
	return &lifecycler{
		provisioner:              newProvisioner(cl),
		deprovisioner:            newDeprovisioner(cl),
		binder:                   newBinder(cl),
		unbinder:                 newUnbinder(cl),
		operationStatusRetriever: newOperationStatusRetriever(cl),
	}
}
