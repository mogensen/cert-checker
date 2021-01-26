KIND_CLUSTER_NAME="cert-checker"
BINDIR ?= $(CURDIR)/bin
TMPDIR ?= $(CURDIR)/tmp
ARCH   ?= amd64

help:  ## display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: help build docker all clean

test: ## test cert-checker
	go test ./...

build: ## build cert-checker
	mkdir -p $(BINDIR)
	CGO_ENABLED=0 go build -o ./bin/cert-checker ./cmd/.

verify: test build ## tests and builds cert-checker

image: ## build docker image
	docker build -t mogensen/cert-checker:v0.0.1 .

clean: ## clean up created files
	rm -rf \
		$(BINDIR) \
		$(TMPDIR)

all: test build docker ## runs test, build and docker

test-coverage: ## Generate test coverage report
	mkdir -p $(TMPDIR)
	go test ./... --coverprofile $(TMPDIR)/outfile
	go tool cover -html=$(TMPDIR)/outfile

report-card: ## Generate static analysis report
	goreportcard-cli -v

dev.kind.delete: ## Delete local kubernetes cluster
	kind delete clusters $(KIND_CLUSTER_NAME)

dev.kind.create: ## Create local cluster
	kind create cluster --name $(KIND_CLUSTER_NAME) --config deploy/kind/kind-cluster-config.yaml || true
	kubectl apply --wait -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update
	sleep 120
	helm upgrade --wait --install prometheus prometheus-community/kube-prometheus-stack \
	 --set grafana.ingress.enabled=enabled										 		\
	 --set grafana.ingress.hosts='{grafana.localtest.me}'  						 		\
	 --set prometheus.ingress.enabled=enabled										 	\
	 --set prometheus.ingress.hosts='{prometheus.localtest.me}'  						\
	 --set grafana.sidecar.dashboards.searchNamespace=ALL

dev.kind.install: image ## Install cert-checker on kind cluster
	kind --name $(KIND_CLUSTER_NAME) load docker-image   mogensen/cert-checker:v0.0.1
	kubectl apply -n cert-checker -f deploy/yaml/
	kubectl delete pod -l app=cert-checker -n cert-checker
