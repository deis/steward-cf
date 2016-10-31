# Installing Steward CF

Installation of Steward CF into a running Kubernetes cluster is (currently) facilitated through a
Make target. (Eventually, we will provide a Helm chart for this purpose.) Users wishing to
familiarize themselves with the particulars of the deployment will want to examine this
repository's `Makefile` and artifacts in the `manifests/` directory.

All subsequent sections of this document assume that you have a running Kubernetes cluster and that
your `kubectl` client is properly configured to interact with that cluster.

## Prerequisites

Steward CF requires a backing broker that implements the [CloudFoundry Service Broker
API][cfbroker]. The backing broker is the middle-man between Steward CF and a service provider. The
backing broker may be deployed within your Kubernetes cluster or elsewhere.

If you're trying Steward CF for the first time or are hacking on Steward CF, the Steward team has
provided a trivial [sample broker][cf-sample-broker]. See that project's
[README.md](https://github.com/deis/cf-sample-broker/blob/master/README.md) for installation
instructions.

## Deploy Steward CF

The installation process requires the following environment variables to be set in advance:

- `BROKER_NAME` - a name by which you can identify the backing broker that the Steward CF instance
will communicate with

- `BROKER_ACCESS_SCHEME` - the scheme (`http` or `https`) by which to access the backing broker

- `BROKER_HOST` - the IP or DNS name of the backing broker

- `BROKER_PORT` - the TCP port of the backing broker

- `BROKER_USERNAME` - the username Steward CF should use to authenticate with the backing broker

- `BROKER_PASSWORD` - the password Steward CF should use to authenticate with the backing broker

Note that Steward CF _also_ responds to the following environment variables, although the Make
target currently used to facilitate installation does _not_ pass these through to Steward CF pods.
Advanced users wishing to experiment with any of the following settings may do so by editing
Steward CF's deployment template directly.

- `BROKER_REQUEST_TIMEOUT_SEC` - the timeout (in seconds) after which Steward should fail a request
to the broker for any individual operation. If not set, this defaults to `5`.

- `WATCH_NAMESPACES` - a comma-delimited list of namespaces in which the Steward CF instance should
watch for service plan claims. If not set, this defaults to `"default"`.

- `MAX_ASYNC_MINUTES` - the maximum number of minutes to wait while asynchronous provisioning or
deprovisioning operations execute. After this time has elapsed, Steward CF will consider such
operations to have timed out. If not set, this defaults to `60`.

- `API_PORT` - the TCP port for Steward CF's HTTP-based health check endpoint to bind to. If not
set, this defaults to `8080`.

- `LOG_LEVEL` - the log level to be applied to filtering Steward CF log messages by severity. If
not set, this defaults to `"info"`.

With all configuration now set, Steward CF can be deployed as follows:

```
$ make deploy
```

Or build and deploy Steward CF from source using:

```
$ make dev-deploy
```

For details on Steward CF's pure Kubernetes-based workflow, please refer to
[README.md](./README.md).

[cf-sample-broker]: https://github.com/deis/cf-sample-broker
[cfbroker]: https://docs.cloudfoundry.org/services/overview.html
