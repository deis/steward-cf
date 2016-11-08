# Steward CF

[![Build Status](https://travis-ci.com/deis/steward-cf.svg?token=UQsxfwHAz3NPyVqxkrrp&branch=master)]
(https://travis-ci.com/deis/steward-cf)

Steward CF is a Kubernetes service catalog controller implementation that
translates standard Kubernetes service catalog resources into [CloudFoundry Broker
API][cfbroker] calls. This project is built on top of
[steward-framework](https://gihub.com/deis/steward-framework), the common
service catalog controller framework for Kubernetes.

Specifically, Steward CF's high-level goals are to:

1. Decouple the provider of the service from its consumers
2. Allow operators to independently scale and manage applications and the services they depend on
3. Provide a standard way for operators to:
  - Publish a catalog of services to operators or other interested parties
  - Provision a service
  - Bind an application to a service
  - Configure the application to consume the service through standard Kubernetes resources

## Glossary

* **Consumer**: An application and/or developer that requires access to some service provided by a
third party.

* **Requestable Service**: A service offered by a third party that may be provisioned for or
exposed to a **consumer**. **Requestable services** are not related to Kubernetes services.
Examples include:
    * account and access credentials for an off-cluster SaaS service like Sendgrid
    * access credentials for a relational data store like MySQL or Postgres

* **Service Plan**: is a specific "configuration" of a **requestable services**, which may be
expressed in semantically meaningful terms like a version number or "small", "medium", or "large".

* **Service Catalog**: is a registry of **requestable services** and **service plans** available to
**consumers** within a Kubernetes cluster.

* **Service Plan Claim**: is a Kubernetes `ConfigMap` which represents the desire of a **consumer**
to gain access to a **requestable service** (a specific **service plan** thereof). The **service
plan claim** references both the **requestable service** and **service plan** by ID and also
informs Steward CF where the consuming application expects to read **service
credentials/configuration** that are created after processing the claim.

* **Service Credentials/Configuration**: Is the configuration (hostnames, usernames, passwords,
etc.) meant for the **consumer** to use for connection and authentication to a **service
instance**.

* **Service Instance**: is the entity or entities provisioned or exposed on behalf of a
**consumer** that made a **service plan claim**. Examples of **service instances** include:
    * a provisioned AWS RDS service, a logical database and credentials
    * a logical database, username and password created on a shared RDBMS

* **Service Provider**: A system that lives either on or off-cluster and is directly capable of
providing a **requested service** or provisioning **service instances**. For example, in the case
of Sendgrid, their SaaS platform that provides email delivery services qualifies as a service
provider. Amazon's RDS (Relational Database Service) which is capable of provisioning managed
relation databases also qualifies as a service provider.

* **Cloud Foundry Service Broker API**: An API definition created by Cloud Foundry, broadly
describing a uniform interface for provisioning and deprovisioning services from **service
providers** and for binding **consumers** to and unbinding them from those services.

* **Cloud Foundry Service Broker**: A concrete implementation of the **Cloud Foundry Service Broker
API**, e.g. <https://github.com/cloudfoundry/cf-mysql-release>.

* **Backing Broker**: With respect to Steward CF, a **Cloud Foundry service broker** which Steward
CF will utilize to provision and deprovision services (and credentials) from **service providers**.

Putting it all together, Steward CF watches the Kubernetes event stream for **service plan claims**
and submits those claims to a **backing broker**, which in turn delegates discreet actions like
provisioning to a **service provider** to create **service instances** and **service
credentials/configuration** for use by **consumers**.

# Deploying Steward CF

Please see [INSTALLATION.md](./doc/INSTALLATION.md) for full instructions, including sample
Kubernetes manifests, on how to deploy Steward CF to your cluster.

Once deployed, you can view logs for each Steward CF instance via the standard `kubectl logs`
command:

```console
kubectl logs -f ${STEWARD_CF_POD_NAME} --namespace=${STEWARD_NAMESPACE}
```

# Concepts

Steward CF (via the Steward framework) runs a control loop to watch the Kubernetes event stream for
a set of [`ThirdPartyResource`][3pr]s (called 3PRs hereafter) in one, some, or all available
namespaces. It uses these 3PRs to communicate with an operator that requests a service.

A single Steward CF process is responsible for talking to a single **backing broker**. If a cluster
operator wishes to expose multiple **backing broker**, he or she would deploy additional instances
of Steward CF.

## Available Services

On startup, a Steward CF process publishes its service data as a set of `ServiceCatalogEntry` 3PRs
that indicate the availability of each of a **Requestable Service**. Each **Service Catalog Entry**
contains the name of the Steward CF instance (specified in configuration), the **Requestable
Service**, and at least one **Service Plan**. Here are some example 3PRs:

- `firststeward-mysql-small`
- `secondsteward-mysql-large`
- `thirdsteward-memcache-xlarge`

Once published, an operator (or other interested party) will be able to see these
`ServiceCatalogEntry` 3PRs to determine what services are available to applications in the cluster.
Currently, we recommend simply using the `kubectl` command to list the catalog:

```console
kubectl get servicecatalogentries --namespace=steward
```

Please see [DATA_STRUCTURES.md](./doc/DATA_STRUCTURES.md) for a complete example of a
`ServiceCatalogEntry`.

## Requesting a Service from the Catalog

Once an operator has found a service and plan they would like to use, they should submit a
[ConfigMap][configMap] containing a [`ServicePlanClaim`](./doc/DATA_STRUCTURES.md) data structure
(just called `ServicePlanClaim`s hereafter).

Steward CF constantly watches for `ServicePlanClaim`s in its control loop. Upon finding a new
`ServicePlanClaim`, it executes the following algorithm:

1. Looks for the `ServiceCatalogEntry` 3PR in the catalog
  - If not found, sets the `status` field to `Failed`, adds an appropriate explanation to the
  `statusDescription` field field to a human-readable description of the error, and stops
  processing
1. Looks in the `action` field of the claim and takes the appropriate action
  - Valid values are `provision`, `bind`, `unbind`, `deprovision`, `create` and `delete`. See
  [`ServicePlanClaim` documentation](./doc/DATA_STRUCTURES.md#serviceplanclaim) for details on each
  value
  - If the action failed, Steward CF sets the `status` field to `Failed` and adds an appropriate
  explanation to the `statusDescription` field
1. On success, writes values appropriate to the `action` that was submitted. See below for details
on each `action`

### `provision`
- `status: provisioned`
- `instance-id: $UUID` (where `$UUID` is the instance ID returned by the provision operation)

### `bind`
- `status: bound`
- `bind-id: $UUID` (where `$UUID` is the bind ID returned by the bind operation)
- Also creates a [Secret][secrets] with the credentials data for the service. The Secret's name and
namespace will be created according to the `target-name` and `target-namespace` fields passed in
the `ServicePlanClaim`. See [`ServicePlanClaim` documentation]
(./doc/DATA_STRUCTURES.md#serviceplanclaim) for more information

### `unbind`
- `status: unbound`
- Removes the Secret created as a result of the `bind` action

### `deprovision`
- `status: deprovisioned`

### `create`

This action produces results equivalent to claims with `action: provision`, then `action: bind`

### `delete`

This action produces results equivalent to claims with `action: unbind`, then `action: deprovision`


# Development & Testing

Steward CF is written in Go and tested with [Go unit tests](https://godoc.org/testing).

If you'd like to contribute to this project, simply fork the repository, make your changes, and
submit a pull request. Please make sure to follow [these guidelines](CONTRIBUTING.md) when
contributing.

[cfbroker]: https://docs.cloudfoundry.org/services/overview.html
[3pr]: https://github.com/kubernetes/kubernetes/blob/master/docs/design/extending-api.md
[rds]: https://aws.amazon.com/rds
[configMap]: http://kubernetes.io/docs/user-guide/configmap/
[secrets]: http://kubernetes.io/docs/user-guide/secrets/
[servicePlanCreation]: ./DATA_STRUCTURES.md#serviceplancreation
