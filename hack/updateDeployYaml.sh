#!/bin/bash

helm template cert-checker deploy/charts/cert-checker --no-hooks --set image.pullPolicy=Always  \
    | grep -vi helm \
    | grep -vi chart \
    | grep -v "# Source" \
    | grep -v "checksum/config" > deploy/yaml/deploy.yaml

helm template cert-checker deploy/charts/cert-checker --no-hooks -s templates/grafana-dashboard-cm.yaml --set grafanaDashboard.enabled=true  \
    | grep -vi helm \
    | grep -vi chart \
    | grep -v "# Source" \
    | grep -v "checksum/config" > deploy/yaml/grafana-dashboard-cm.yaml

helm template cert-checker deploy/charts/cert-checker --no-hooks -s templates/servicemonitor.yaml \
    --set serviceMonitor.enabled=true  \
    --set serviceMonitor.additionalLabels.release=prometheus  \
    | grep -vi helm \
    | grep -vi chart \
    | grep -v "# Source" \
    | grep -v "checksum/config" > deploy/yaml/servicemonitor.yaml
