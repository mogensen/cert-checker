# cert-checker

[![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fmogensen%2Fcert-checker%2Fbadge%3Fref%3Dmain&style=flat)](https://actions-badge.atrox.dev/mogensen/cert-checker/goto?ref=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mogensen/cert-checker)](https://goreportcard.com/report/github.com/mogensen/cert-checker)
[![codecov](https://codecov.io/gh/mogensen/cert-checker/branch/main/graph/badge.svg)](https://codecov.io/gh/mogensen/cert-checker)

cert-checker is a certificate monitoring utility for watching tls certificates. These
checks get exposed as Prometheus metrics to be viewed on a dashboard, or _soft_
alert cluster operators.

This tool is heavily inspired by the awesome [version-checker by jetstack](https://github.com/jetstack/version-checker/).

## Table of contents

- [cert-checker](#cert-checker)
  * [Table of contents](#table-of-contents)
  * [Features](#features)
    + [Testing for Certificate Errors](#testing-for-certificate-errors)
    + [Testing for minimal TLS Version](#testing-for-minimal-tls-version)
    + [Permissions](#permissions)
  * [Installation](#installation)
    + [Run in Docker](#run-in-docker)
    + [Using docker-compose](#using-docker-compose)
    + [In Kubernetes as static manifests](#in-kubernetes-as-static-manifests)
    + [Helm](#helm)
    + [Kustomize](#kustomize)
  * [Web dashboard](#web-dashboard)
  * [Metrics](#metrics)
    + [Grafana Dashboard](#grafana-dashboard)
  * [Options](#options)
  * [Development](#development)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>

## Features

### Testing for Certificate Errors

cert-checker supports the following types of certificate errors (and possible more):

- Expired certificates
- Wrong host
- Bad root certificates
- Revoked certificate
- Cipher suites not allowed
    * `dh480`
    * `dh512`
    * `null`
    * `rc4`

If cert-checker finds any certificate errors, these are displayed on the Grafana dashboard.

### Testing for minimal TLS Version

cert-checker checks the minimum supported SSL/TLS version for the endpoints.

The following SSL/TLS versions are tested:
 - SSL 3.0 - Deprecated in 2015
 - TLS 1.0 - Deprecated in 2020
 - TLS 1.1 - Deprecated in 2020
 - TLS 1.2
 - TLS 1.3

See [Transport Layer Security](https://en.wikipedia.org/wiki/Transport_Layer_Security) for more info.

The minimum supported versions are displayed on the Grafana dashboard.

### Permissions

A great bonus of how the cert-checker is implemented is that it can run without `root`, and without `CAP_NET_RAW` capability.
And without Administrator privileges in Windows.

---

## Installation

cert-checker can be installed as a standalone static binary from the release page

[latest release](https://github.com/mogensen/cert-checker/releases/latest/)

Create a config file like the below example:

`config.yaml`:

```yaml
loglevel: debug
port: 8080  # Optional
intervalminutes: 10
certificates:
    - dns: google.com
    - dns: expired.badssl.com
```

```bash
./cert-checker -c config.yaml
DEBU[2021-05-17T17:27:44+02:00] Probing all
INFO[2021-05-17T17:27:44+02:00] serving ui on 0.0.0.0:8081
INFO[2021-05-17T17:27:44+02:00] serving metrics on 0.0.0.0:8080/metrics
DEBU[2021-05-17T17:27:44+02:00] Probing: google.com
...
# Now open browser at:
#   -  http://localhost:8081/
#   -  http://localhost:8080/metrics
```

### Run in Docker

You can use the published docker image like this:

First create a config file as above, or download the demo file:

```bash
curl https://raw.githubusercontent.com/mogensen/cert-checker/main/config.yaml -O
```


```bash
# Start docker container (mounting the config file may be different on OSX and Windows)
docker run -p 8081:8081 -p 8080:8080 -v ${PWD}/config.yaml:/app/config.yaml mogensen/cert-checker:latest
# Now open browser at:
#   -  http://localhost:8081/
#   -  http://localhost:8080/metrics
```

See released docker images on [DockerHub](https://hub.docker.com/r/mogensen/cert-checker)

### Using docker-compose

This repository contains an example of deploying the entire Prometheus, Grafana and cert-checker stack, using docker-compose.

```bash
cd deploy/docker-compose/
docker-compose up -d
```

| Service           | URL                                                                                   |
|-------------------|---------------------------------------------------------------------------------------|
| cert-checker      | ui endpoint http://localhost:8081/                                                    |
| cert-checker      | metrics endpoint http://localhost:8080/metrics                                        |
| Prometheus        | example query http://localhost:9090/graph?g0.expr=cert_checker_expire_time{}&g0.tab=0 |
| Grafana           | Dashboard  http://localhost:3000/d/cert-checker/certificate-checker                   |

Remember to edit the `deploy/docker-compose/cert-checker/config.yaml` with the actual domains you want to monitor..

See [stefanprodan/dockprom](https://github.com/stefanprodan/dockprom) for more Prometheus, Grafana, AlertManager examples using Docker-compose


### In Kubernetes as static manifests

cert-checker can be installed as static manifests:

```sh
$ kubectl create namespace cert-checker

# Deploy cert-checker, with kubernetes services and demo configuration
$ kubectl apply -n cert-checker -f deploy/yaml/deploy.yaml

# If you are using the Grafana sidecar for loading dashboards
$ kubectl apply -n cert-checker -f deploy/yaml/grafana-dashboard-cm.yaml

# If you are using the Prometheus CRDs for setting up scrape targets
$ kubectl apply -n cert-checker -f deploy/yaml/servicemonitor.yaml
```

Remember to edit the configmap with the actual domains you want to monitor..

### Helm

cert-checker can be installed as as helm release:

```bash
$ kubectl create namespace cert-checker
$ helm install cert-checker deploy/charts/cert-checker --namespace cert-checker
```

Depending on your setup, you may need to modify the `ServiceMonitor` to get Prometheus to scrape it in a particular namespace.
See [this](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack#prometheusioscrape).

You may also need to add additional labels to the `ServiceMonitor`.
If you have installed the `prometheus-community/kube-prometheus-stack` with the name of `prometheus` the following should work:

```bash
$ helm upgrade cert-checker deploy/charts/cert-checker \
    --namespace cert-checker            \
    --set=grafanaDashboard.enabled=true \
    --set=serviceMonitor.enabled=true   \
    --set=serviceMonitor.additionalLabels.release=prometheus
```

### Kustomize

cert-checker can be installed using [kustomize](https://kustomize.io/):

Create a `kustomization.yaml` file:
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: cert-checker
resources:
- github.com/mogensen/cert-checker/deploy/yaml
# optionally pin to a specific git tag
# - github.com/mogensen/cert-checker/deploy/yaml?ref=cert-checker-0.0.5

# override confimap with your required settings
patchesStrategicMerge:
- |-
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: cert-checker
      namespace: cert-checker
    data:
      config.yaml: |
        loglevel: info
        intervalminutes: 60
        certificates:
            - dns: my-very-own-domain.com
```
Use the `kustomization.yaml` file to preview and deploy cert-checker:
```bash
$ kustomize build kustomization.yaml | less # preview yaml manifests
$ kustomize build kustomization.yaml | kubectl apply --dry-run=client -f - # dry-run apply manifests
$ kustomize build kustomization.yaml | kubectl apply -f - # deploy manifests
```

## Web dashboard

By default, cert-checker will expose a web ui on `http://0.0.0.0:8081/`.

![](img/web-ui.jpg)
<center></center>
<p align="center">
  <b>Web dashboard</b><br>
</p>


## Metrics

By default, cert-checker will expose the version information as Prometheus
metrics on `http://0.0.0.0:8080/metrics`.

### Grafana Dashboard

A Grafana dashboard is also included in this repository.
It is located in the deployment folder: `deploy/yaml/grafana-dashboard-cm.yaml`

![](img/grafana.jpg)
<center></center>
<p align="center">
  <b>Grafana Dashboard</b><br>
</p>

The dashboard shows the following

 - Number of Broken Certificates
 - Number of Certificates about to expire
 - Number of Good Certificates
 - A list with Certificates with errors
 - A list of Certificates Expirations for valid certificates
 - Minimum TLS versions supported

The conventions used on the dashboard are:

 - Red (text or background): Something is broken, and should be fixed!
 - Orange (text or background): Something smells, and should properly be fixed!
 - Green (text or background): All is good! Go drink coffee!

---

## Options

By default, without the flag `-c, --config`, cert-checker will
use a config file located next to the binary named `config.yaml`.

This is currently the only flag / option available.

```bash
$ cert-checker -h
Certificate monitoring utility for watching tls certificates and reporting the result as metrics.

Usage:
  version-checker [flags]

Flags:
  -c, --config string   config file (default is config.yaml) (default "config.yaml")
  -h, --help            help for version-checker
```

---

## Development

Test the full setup in Kubernetes with Prometheus and Grafana dashboards:

```bash
# First create a new kind cluster locally, and install prometheus
make dev-kind-create
# Build a docker image, load it into kind and deploy cert-checker and promeheus/grafana stuff
make image dev-kind-install
```

Access the local infrastructure here:

| System             | URL                                                                                                        |
| ------------------ |------------------------------------------------------------------------------------------------------------|
| Prometheus         | http://prometheus.localtest.me/graph?g0.expr=cert_checker_is_valid&g0.tab=1&g0.stacked=0&g0.range_input=1h |
| Grafana            | http://grafana.localtest.me/d/cert-checker/certificate-checker                                             |
| build-in dashboard | http://cert-checker.localtest.me/                                                                          |
