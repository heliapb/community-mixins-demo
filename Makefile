GOCMD=go
GOMAIN=main.go
GOBUILD=$(GOCMD) build
GOOS?=$(shell go env GOOS)
ENVVARS=GOOS=$(GOOS) CGO_ENABLED=0

.PHONY: build
build:
	@echo "Removing old rules and dashboards"
	@rm -rf ./rules/
	@rm -rf ./dashboards/
	@echo "Building rules and dashboard"
	@$(ENVVARS) $(GOCMD) run $(GOMAIN) --output-rules-dir="./rules/prometheus" --output-rules="yaml" --output-dir="./dashboards/perses" --output="yaml" --project="perses-dev" --datasource="prometheus-datasource"

.PHONY: apply-dashboards
apply-dashboards:
	@echo "Applying dashboards"
	@percli login http://localhost:8080
	@percli apply -d ./dashboards/perses/perses
	@percli apply -d ./dashboards/perses/prometheus
	@percli apply -d ./dashboards/perses/blackbox-exporter
	@percli apply -d ./dashboards/perses/demo-app
