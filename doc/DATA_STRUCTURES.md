# `ServiceCatalogEntry`

This object is written to the `steward` namespace and represents a single service + plan pair that
at least one Steward CF instance (through its backing broker) can provision and bind. It has the
following fields:

- `service_info` - an object containing information about the service. See below for a description
of the fields in this object

- `service_plan` - an object containing information about the service's plan. See below for a
description of the fields in this Object

## `service_info`

This object contains information about a service offered by Steward CF. It has the fields listed
below. Each is a string unless otherwise indicated.

- `id` - a unique identifier for the service
- `name` - the name of the service
- `description` - a description of the service
- `plan_updateable` - a boolean indicating whether the service's plans are updateable

## `service_plan`

This object contains information about an individual plan for a service offered by Steward CF. It
has the fields listed below. Each is a string unless otherwise indicated.

- `id` - a unique identifier for the plan
- `name` - the name of the plan
- `descripton` - a description of the plan
- `free` - a boolean indicating whether this plan is free

## Example

Below is an example of a `ServiceCatalogEntry`, in yaml format:

```yaml
apiVersion: steward.deis.io/v1
description: RDS PostgreSQL service (RDS PostgreSQL 9.4 (db.m3.medium, 10Gb))
kind: ServiceCatalogEntry
metadata:
  creationTimestamp: 2016-10-31T20:45:35Z
  labels:
    broker: cf-rds
    plan-id: d03b544e-3be5-4aca-bb3b-11544247f313
    plan-name: 9.4-medium
    service-id: a2c9adda-6511-462c-9934-b3fd8236e9f0
    service-name: rdspostgres
  name: cf-rds-rdspostgres-9-4-medium
  namespace: steward
  resourceVersion: "403"
  selfLink: /apis/steward.deis.io/v1/namespaces/steward/servicecatalogentries/cf-rds-rdspostgres-9-4-medium
  uid: f923f469-9faa-11e6-80ea-feb8bdc266e3
service_info:
  Description: RDS PostgreSQL service
  ID: a2c9adda-6511-462c-9934-b3fd8236e9f0
  Name: rdspostgres
  PlanUpdatable: true
service_plan:
  Description: RDS PostgreSQL 9.4 (db.m3.medium, 10Gb)
  Free: false
  ID: d03b544e-3be5-4aca-bb3b-11544247f313
  Name: 9.4-medium
```

# `ServicePlanClaim`

This object is submitted by a service consumer in the form of a [`ConfigMap`][configMap] when that
consumer wants Steward CF to create a new service for its use. Steward CF then mutates the object
to communicate the status of the service creation operation. Consumers may watch the event stream
for this object to assert progress of service creation. `ServicePlanClaim`s contain the following
fields:

- `claim-id` - a unique consumer-generated [UUID][uuid]

- `service-id` - the identifier for the desired service

- `plan-id` - the identifier for the desired plan

- `target-name` - the name of a Kubernetes [`Secret`][secret] into which Steward CF should write
  the resulting credentials after binding to a service

- `action` - the consumer-specified action to take. Valid values are `provision`,`bind`, `unbind`,
  `deprovision`, `create` and `delete`. A few more notes:
  - Steward CF will never modify this value
  - `create` will execute both the `provision` and `bind` actions, in that order
  - `delete` will execute both the `unbind` and `deprovision` actions, in that order
  - All new `ServicePlanClaim`s submitted by consumer must have `action` set to `provision` or
    `create`
  - If Steward CF encounters an error, any actions it has already completed will not be rolled
    back. See the following examples:
    - If you submit a claim with `action: create` and the bind step fails, the provision step will
      not be rolled back
    - If you submit a claim with `action: bind` with a `target-name` that points to a Secret that
      already exists, Steward CF will execute the bind action on the backing broker, but will fail
      to write the new credentials Secret. The backing broker bind action will not be rolled back
  - It is an error for this field to be empty

- `status` - the current status of the claim. Steward CF will modify this value, but will ignore
  any modifications by the consumer. Valid values with short descriptions are listed below:
  - `provisioning` - immediately after `action` is set to `provision` if the backing broker carries
    out provisioning _synchronously_
  - `provisioning-async` - immediately after `action` is set to `provision` if the backing broker
    carries out provisioning _asynchronously_
  - `provisioned` - after `action` is set to `provision` and the provisioning process has succeeded
  - `binding` - immediately after `action` is set to `bind`
  - `bound` - after `action` is set to `bind` and the binding process succeeded
  - `unbinding` - immediately after `action` is set to `unbind`
  - `unbound` - after `action` is set to `unbind` and the unbinding process succeeded
  - `deprovisioning` - immediately after `action` is set to `deprovision` if the backing broker
    carries out deprovisioning _synchronously_
  - `deprovisioning-async` - immediately after `action` is set to `deprovision` if the backing
    broker carries out deprovisioning _asynchronously_
  - `deprovisioned` - after `action` is set to `deprovision` and the deprovisioning process
    succeeded
  - `failed` - after any `action` failed
- `status-description` - a human-readable explanation of the current `status`. Steward CF will
  modify this value, but will ignore any modifications by consumer
- `instance-id` - for internal use only. The consumer should not modify this field
- `bind-id` - for internal use only. The consumer should not modify this field
- `extra` - for internal use only. The consumer should not modify this field

[3pr]: https://github.com/kubernetes/kubernetes/blob/master/docs/design/extending-api.md
[uuid]: https://tools.ietf.org/html/rfc4122
[configMap]: http://kubernetes.io/docs/user-guide/configmap/
[secret]: http://kubernetes.io/docs/user-guide/secrets/
