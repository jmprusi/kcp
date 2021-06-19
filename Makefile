all: build
.PHONY: all

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

build:
	go build -ldflags "-X k8s.io/client-go/pkg/version.gitVersion=$$(git describe --abbrev=8 --dirty --always)" -o bin/kcp ./cmd/kcp
	go build -ldflags "-X k8s.io/client-go/pkg/version.gitVersion=$$(git describe --abbrev=8 --dirty --always)" -o bin/syncer ./cmd/syncer
	go build -ldflags "-X k8s.io/client-go/pkg/version.gitVersion=$$(git describe --abbrev=8 --dirty --always)" -o bin/cluster-controller ./cmd/cluster-controller
	go build -ldflags "-X k8s.io/client-go/pkg/version.gitVersion=$$(git describe --abbrev=8 --dirty --always)" -o bin/deployment-splitter ./cmd/deployment-splitter
.PHONY: build

vendor:
	go mod tidy
	go mod vendor
.PHONY: vendor

codegen:
	./hack/update-codegen.sh
.PHONY: codegen

KIND = $(shell pwd)/bin/kind
kind:
	$(call go-get-tool,$(KIND),sigs.k8s.io/kind@v0.10.0)

.PHONY: local-setup
local-setup: build kind
	./utils/local-setup.sh

.PHONY: clean
clean:
	-rm -f ./bin/*
	-rm -rf ./tmp
	-rm -rf ./.kcp
