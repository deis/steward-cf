# Demo-ing Steward CF

This document outlines how to show a demo for Steward CF.

## Setup

To run a demo, you need the following:

1. A running service broker that implements the [Cloud Foundry Service Broker API][cfbrokerapi]
2. A Steward CF instance running and configured to talk to the broker

### Broker Setup

First, broker setup and run instructions from one broker to the next, but we've provided a
[sample broker][sample-broker]. This broker is easily installable and requires no configuration,
but also does nothing useful. It's best used for testing the Steward-to-broker communication, but
for real demos we recommend using a broker that does something more useful.

For more advanced users, the Cloud Foundry community's [AWS RDS service broker](https://github.com/cloudfoundry-community/pe-rds-broker) is an excellent example of a
non-trivial broker. Configure that broker according to its README and deploy it anywhere. To run it
in your Kubernetes cluster, it can be installed using Deis Workflow's `deis pull` workflow.

### Configuring and Running Steward CF

Once a broker is deployed, configuring and starting Steward CF is fairly simple. See the [pre-built
manifest](https://github.com/deis/steward-cf/blob/master/manifests/steward-cf-template.yaml)
that contains all the configuration necessary to run Steward CF. Simply change each config field in
the `spec.template.spec.containers[0].env` field and run `kubectl create -f
manifests/steward-cf-template.yaml` to install.

After Steward CF starts, it will query the backing broker's catalog, convert each catalog entry
into a list of Kubernetes [Third Party Resource][3pr]s called `ServiceCatalogEntries` and write
each into the `steward` namespace.

## Running the demo

### Inspect the Catalog
After a successful startup, view the catalog by running the following command:

```console
kubectl get servicecatalogentries --namespace=steward
```

You'll see a list of entries. Choose one and run the following command to see further details:

```console
kubectl get servicecatalogentry $ENTRY_NAME --namespace=steward -o yaml
```

You should see some YAML output. Make note of the `service_info.id` and `plan_info.id` fields, as
you'll use them in the service plan claim that you'll create next.

### Make a Claim

Interact with Steward CF using service plan claims (claims hereafter). Claims are config maps with
the following properties:

- A label with key `type` and value `service-plan-claims`
- `data` with the following keys:
  - `claim-id` - an consumer-specified identifier for the claim. Currently, these need not follow
  any format nor are they required to be unique
  - `service-id` - the ID of a service that's present in the catalog
  - `plan-id` - the ID of a plan for the above service
  - `action` - one of the actions in the below bulleted list
  - `target-name` - the name of the Kubernetes secret that Steward CF will create and to which it
  will write bound credentials. The secret will always be created in the same namespace as the
  claim itself. Note that Steward CF deletes the secret named here during an unbind operation, so
  this field should never be changed

Steward CF watches the `default` namespace for claims that have been added or modified. The
following actions all trigger transitions within Steward CF's internal state machine:

- `provision` - the creation action of a requestable service
- `bind` - the process of getting credentials to access a provisioned requestable service
- `unbind` - the process of invalidating (or "giving up") credentials that were acquired by a bind
operation
- `deprovision` - the deletion action of a requestable service that was previously provisioned
- `create` - the compound action of `provision`, then `bind`
- `delete` - the compound action of `unbind`, then `deprovision`

To submit a claim, complete the following steps:

- Open `manifests/sample-service-plan-claim-cm.yaml`
  - Change the `metadata.name` value to `"steward-demo"`
  - Change the `data.service-id` value to `"$SERVICE_ID"`
  - Change the `data.plan-id` field to `"$PLAN_ID"`
  - Change the `data.target-name` field to `"steward-demo"`
- Submit the claim: `kubectl create -f manifests/sample-service-plan-claim-cm.yaml`
- Inspect the created secret: `kubectl get secret steward-demo -o yaml`
- When done, delete the provisioned, bound resource:
  - `kubectl edit configmap claim-1`
  - Change the `data.action` field to `"delete"`
- Ensure that the previously-created secret is missing: `kubect get secret steward-demo`

A few extra notes:

- Recall that `create` and `delete` are compound actions (`provision`/`bind` and `unbind`/
`deprovision`, respectively). Some audiences may want to see all four discreet actions
(`provision`, `bind`, `unbind`, `deprovision`). If that's the case, simply modify the instructions
above to:
  - Change the `data.action` field to `"provision"` before submitting the claim
  - Call `kubectl edit configmap claim-1` three times to change the `action` field to `bind`, then
  `unbind`, then finally `deprovision`
  - Ensure that the appropriate resources are in the appropriate state after each action. For
  example, the appropriate command should be run to successful completion in command mode.

[sample-broker]: https://github.com/deis/cf-sample-broker
[cfbrokerapi]: https://docs.Cloud Foundry.org/services/api.html
