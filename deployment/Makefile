RESULT_DIR = ./k8s

.PHONY: build-yaml

all: build-yaml

build-yaml: build-yaml-web build-yaml-twproxy

build-yaml-web:
	envsubst < web.tpl > $(RESULT_DIR)/web.yaml

build-yaml-twproxy:
	envsubst < twproxy.tpl > $(RESULT_DIR)/twproxy.yaml
