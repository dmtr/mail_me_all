RESULT_DIR = ./k8s

.PHONY: build-yaml

all: build-yaml

build-yaml: build-yaml-web build-yaml-twproxy build-yaml-backend build-yaml-ingress

build-yaml-web:
	envsubst < web.tpl > $(RESULT_DIR)/web.yaml

build-yaml-twproxy:
	envsubst < twproxy.tpl > $(RESULT_DIR)/twproxy.yaml

build-yaml-backend:
	envsubst < backend.tpl > $(RESULT_DIR)/backend.yaml

build-yaml-ingress:
	envsubst < ingress.tpl > $(RESULT_DIR)/ingress.yaml

build-yaml-cert:
	envsubst < cert.tpl > $(RESULT_DIR)/cert.yaml
