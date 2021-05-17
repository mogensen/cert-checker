#!/bin/bash


if ! helm version -c --short | grep -E "v3." >/dev/null; then
    echo "Helm v3 is needed!"
    exit 1
fi

helm template cert-checker deploy/charts/cert-checker --no-hooks --set image.pullPolicy=Always  \
    --set ingress.enabled=true  \
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

cp deploy/charts/cert-checker/dashboards/cert-checker.json deploy/docker-compose/grafana/provisioning/dashboards/cert-checker.json
