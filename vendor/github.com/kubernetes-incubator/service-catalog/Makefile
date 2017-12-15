# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

all: build test verify

# Some env vars that devs might find useful:
#  GOFLAGS      : extra "go build" flags to use - e.g. -v   (for verbose)
#  NO_DOCKER=1  : execute each step natively, not in a Docker container
#  TEST_DIRS=   : only run the unit tests from the specified dirs
#  UNIT_TESTS=  : only run the unit tests matching the specified regexp

# Define some constants
#######################
ROOT           = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BINDIR        ?= bin
BUILD_DIR     ?= build
COVERAGE      ?= $(CURDIR)/coverage.html
SC_PKG         = github.com/kubernetes-incubator/service-catalog
TOP_SRC_DIRS   = cmd contrib pkg plugin
SRC_DIRS       = $(shell sh -c "find $(TOP_SRC_DIRS) -name \\*.go \
                   -exec dirname {} \\; | sort | uniq")
TEST_DIRS     ?= $(shell sh -c "find $(TOP_SRC_DIRS) -name \\*_test.go \
                   -exec dirname {} \\; | sort | uniq")
VERSION       ?= $(shell git describe --always --abbrev=7 --dirty)

# Run stat against /dev/null and check if it has any stdout output.
# If stdout is blank, we are detecting bsd-stat because stat it has
# returned an error to stderr. If not bsd-stat, assume gnu-stat.
ifeq ($(shell stat -c"%U" /dev/null 2> /dev/null),)
STAT           = stat -f '%c %N'
else
STAT           = stat -c '%Y %n'
endif

NEWEST_GO_FILE = $(shell find $(SRC_DIRS) -name \*.go -exec $(STAT) {} \; \
                   | sort -r | head -n 1 | sed "s/.* //")
TYPES_FILES    = $(shell find pkg/apis -name types.go)
GO_VERSION     = 1.8

PLATFORM?=linux
ARCH?=amd64

GO_BUILD       = env GOOS=$(PLATFORM) GOARCH=$(ARCH) go build -i $(GOFLAGS) \
                   -ldflags "-X $(SC_PKG)/pkg.VERSION=$(VERSION)"
BASE_PATH      = $(ROOT:/src/github.com/kubernetes-incubator/service-catalog/=)
export GOPATH  = $(BASE_PATH):$(ROOT)/vendor

MUTABLE_TAG                      ?= canary
APISERVER_IMAGE                   = $(REGISTRY)apiserver:$(VERSION)
APISERVER_MUTABLE_IMAGE           = $(REGISTRY)apiserver:$(MUTABLE_TAG)
CONTROLLER_MANAGER_IMAGE          = $(REGISTRY)controller-manager:$(VERSION)
CONTROLLER_MANAGER_MUTABLE_IMAGE  = $(REGISTRY)controller-manager:$(MUTABLE_TAG)
USER_BROKER_IMAGE                 = $(REGISTRY)user-broker:$(VERSION)
USER_BROKER_MUTABLE_IMAGE         = $(REGISTRY)user-broker:$(MUTABLE_TAG)

# precheck to avoid kubernetes-incubator/service-catalog#361
$(if $(realpath vendor/k8s.io/kubernetes/vendor), \
	$(error the vendor directory exists in the kubernetes \
		vendored source and must be flattened. \
		run 'glide i -v'))

ifdef UNIT_TESTS
	UNIT_TEST_FLAGS=-run $(UNIT_TESTS) -v
endif

ifdef NO_DOCKER
	DOCKER_CMD =
	scBuildImageTarget =
else
	# Mount .pkg as pkg so that we save our cached "go build" output files
	DOCKER_CMD = docker run --rm -v $(PWD):/go/src/$(SC_PKG) \
	  -v $(PWD)/.pkg:/go/pkg scbuildimage
	scBuildImageTarget = .scBuildImage
endif

NON_VENDOR_DIRS = $(shell $(DOCKER_CMD) glide nv)

# This section builds the output binaries.
# Some will have dedicated targets to make it easier to type, for example
# "apiserver" instead of "bin/apiserver".
#########################################################################
build: .init .generate_files \
       $(BINDIR)/controller-manager $(BINDIR)/apiserver \
       $(BINDIR)/user-broker

user-broker: .init $(BINDIR)/user-broker
$(BINDIR)/user-broker: contrib/cmd/user-broker \
	  $(shell find contrib/cmd/user-broker -type f)
	$(DOCKER_CMD) $(GO_BUILD) -o $@ $(SC_PKG)/contrib/cmd/user-broker

# We'll rebuild apiserver if any go file has changed (ie. NEWEST_GO_FILE)
apiserver: .init $(BINDIR)/apiserver
$(BINDIR)/apiserver: .generate_files cmd/apiserver $(NEWEST_GO_FILE)
	$(DOCKER_CMD) $(GO_BUILD) -o $@ $(SC_PKG)/cmd/apiserver

controller-manager: .init $(BINDIR)/controller-manager
$(BINDIR)/controller-manager: .generate_files cmd/controller-manager $(NEWEST_GO_FILE)
	$(DOCKER_CMD) $(GO_BUILD) -o $@ $(SC_PKG)/cmd/controller-manager

# This section contains the code generation stuff
#################################################
.generate_exes: $(BINDIR)/defaulter-gen \
                $(BINDIR)/deepcopy-gen \
                $(BINDIR)/conversion-gen \
                $(BINDIR)/client-gen \
                $(BINDIR)/lister-gen \
                $(BINDIR)/informer-gen \
                $(BINDIR)/openapi-gen
	touch $@

$(BINDIR)/defaulter-gen: .init
	$(DOCKER_CMD) go build -o $@ $(SC_PKG)/vendor/k8s.io/kubernetes/cmd/libs/go2idl/defaulter-gen

$(BINDIR)/deepcopy-gen: .init
	$(DOCKER_CMD) go build -o $@ $(SC_PKG)/vendor/k8s.io/kubernetes/cmd/libs/go2idl/deepcopy-gen

$(BINDIR)/conversion-gen: .init
	$(DOCKER_CMD) go build -o $@ $(SC_PKG)/vendor/k8s.io/kubernetes/cmd/libs/go2idl/conversion-gen

$(BINDIR)/client-gen: .init
	$(DOCKER_CMD) go build -o $@ $(SC_PKG)/vendor/k8s.io/kubernetes/cmd/libs/go2idl/client-gen

$(BINDIR)/lister-gen: .init
	$(DOCKER_CMD) go build -o $@ $(SC_PKG)/vendor/k8s.io/kubernetes/cmd/libs/go2idl/lister-gen

$(BINDIR)/informer-gen: .init
	$(DOCKER_CMD) go build -o $@ $(SC_PKG)/vendor/k8s.io/kubernetes/cmd/libs/go2idl/informer-gen

$(BINDIR)/openapi-gen: vendor/k8s.io/kubernetes/cmd/libs/go2idl/openapi-gen
	$(DOCKER_CMD) go build -o $@ $(SC_PKG)/$^

$(BINDIR)/e2e.test: .init
	$(DOCKER_CMD) go test -c -o $@ $(SC_PKG)/test/e2e

# Regenerate all files if the gen exes changed or any "types.go" files changed
.generate_files: .init .generate_exes $(TYPES_FILES)
	# Generate defaults
	$(DOCKER_CMD) $(BINDIR)/defaulter-gen \
		--v 1 --logtostderr \
		--go-header-file "vendor/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt" \
		--input-dirs "$(SC_PKG)/pkg/apis/servicecatalog" \
		--input-dirs "$(SC_PKG)/pkg/apis/servicecatalog/v1alpha1" \
	  	--extra-peer-dirs "$(SC_PKG)/pkg/apis/servicecatalog" \
		--extra-peer-dirs "$(SC_PKG)/pkg/apis/servicecatalog/v1alpha1" \
		--output-file-base "zz_generated.defaults"
	# Generate deep copies
	$(DOCKER_CMD) $(BINDIR)/deepcopy-gen \
		--v 1 --logtostderr \
		--go-header-file "vendor/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt" \
		--input-dirs "$(SC_PKG)/pkg/apis/servicecatalog" \
		--input-dirs "$(SC_PKG)/pkg/apis/servicecatalog/v1alpha1" \
		--bounding-dirs "github.com/kubernetes-incubator/service-catalog" \
		--output-file-base zz_generated.deepcopy
	# Generate conversions
	$(DOCKER_CMD) $(BINDIR)/conversion-gen \
		--v 1 --logtostderr \
		--go-header-file "vendor/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt" \
		--input-dirs "$(SC_PKG)/pkg/apis/servicecatalog" \
		--input-dirs "$(SC_PKG)/pkg/apis/servicecatalog/v1alpha1" \
		--output-file-base zz_generated.conversion
	# generate all pkg/client contents
	$(DOCKER_CMD) $(BUILD_DIR)/update-client-gen.sh
	# generate openapi
	$(DOCKER_CMD) $(BINDIR)/openapi-gen \
		--v 1 --logtostderr \
		--go-header-file "vendor/github.com/kubernetes/repo-infra/verify/boilerplate/boilerplate.go.txt" \
		--input-dirs "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1alpha1,k8s.io/client-go/pkg/api/v1,k8s.io/apimachinery/pkg/apis/meta/v1" \
		--output-package "github.com/kubernetes-incubator/service-catalog/pkg/openapi"
	# generate codec
	$(DOCKER_CMD) $(BUILD_DIR)/update-codecgen.sh
	touch $@

# Some prereq stuff
###################

.init: $(scBuildImageTarget)
	touch $@

.scBuildImage: build/build-image/Dockerfile
	sed "s/GO_VERSION/$(GO_VERSION)/g" < build/build-image/Dockerfile | \
	  docker build -t scbuildimage -
	touch $@

# Util targets
##############
.PHONY: verify verify-client-gen
verify: .init .generate_files verify-client-gen
	@echo Running gofmt:
	@$(DOCKER_CMD) gofmt -l -s $(TOP_SRC_DIRS) > .out 2>&1 || true
	@bash -c '[ "`cat .out`" == "" ] || \
	  (echo -e "\n*** Please 'gofmt' the following:" ; cat .out ; echo ; false)'
	@rm .out
	@#
	@echo Running golint and go vet:
	@# Exclude the generated (zz) files for now, as well as defaults.go (it
	@# observes conventions from upstream that will not pass lint checks).
	@$(DOCKER_CMD) sh -c \
	  'for i in $$(find $(TOP_SRC_DIRS) -name *.go \
	    | grep -v generated \
	    | grep -v ^pkg/client/ \
	    | grep -v v1alpha1/defaults.go); \
	  do \
	   golint --set_exit_status $$i || exit 1; \
	  done'
	@#
	$(DOCKER_CMD) go vet $(NON_VENDOR_DIRS)
	@echo Running repo-infra verify scripts
	@$(DOCKER_CMD) vendor/github.com/kubernetes/repo-infra/verify/verify-boilerplate.sh --rootdir=. | grep -v generated > .out 2>&1 || true
	@bash -c '[ "`cat .out`" == "" ] || (cat .out ; false)'
	@rm .out
	@#
	@echo Running href checker:
	@$(DOCKER_CMD) verify-links.sh .
	@echo Running errexit checker:
	@$(DOCKER_CMD) build/verify-errexit.sh

verify-client-gen: .init .generate_files
	$(DOCKER_CMD) $(BUILD_DIR)/verify-client-gen.sh

format: .init
	$(DOCKER_CMD) gofmt -w -s $(TOP_SRC_DIRS)

coverage: .init
	$(DOCKER_CMD) contrib/hack/coverage.sh --html "$(COVERAGE)" \
	  $(addprefix ./,$(TEST_DIRS))

test: .init build test-unit test-integration

test-unit: .init build
	@echo Running tests:
	$(DOCKER_CMD) go test -race $(UNIT_TEST_FLAGS) \
	  $(addprefix $(SC_PKG)/,$(TEST_DIRS))

test-integration: .init $(scBuildImageTarget) build
	# test kubectl
	contrib/hack/setup-kubectl.sh
	contrib/hack/test-apiserver.sh
	# golang integration tests
	$(DOCKER_CMD) test/integration.sh

test-e2e: .generate_files $(BINDIR)/e2e.test
	$(BINDIR)/e2e.test

clean: clean-bin clean-build-image clean-generated clean-coverage

clean-bin:
	rm -rf $(BINDIR)
	rm -f .generate_exes

clean-build-image:
	rm -f .scBuildImage
	docker rmi -f scbuildimage > /dev/null 2>&1 || true

clean-generated:
	rm -f .generate_files
	find $(TOP_SRC_DIRS) -name zz_generated* -exec rm {} \;

clean-coverage:
	rm -f $(COVERAGE)

# Building Docker Images for our executables
############################################
images: user-broker-image \
    controller-manager-image apiserver-image

user-broker-image: contrib/build/user-broker/Dockerfile $(BINDIR)/user-broker
	mkdir -p contrib/build/user-broker/tmp
	cp $(BINDIR)/user-broker contrib/build/user-broker/tmp
	docker build -t $(USER_BROKER_IMAGE) contrib/build/user-broker
	docker tag $(USER_BROKER_IMAGE) $(USER_BROKER_MUTABLE_IMAGE)
	rm -rf contrib/build/user-broker/tmp

apiserver-image: build/apiserver/Dockerfile $(BINDIR)/apiserver
	mkdir -p build/apiserver/tmp
	cp $(BINDIR)/apiserver build/apiserver/tmp
	docker build -t $(APISERVER_IMAGE) build/apiserver
	docker tag $(APISERVER_IMAGE) $(APISERVER_MUTABLE_IMAGE)
	rm -rf build/apiserver/tmp

controller-manager-image: build/controller-manager/Dockerfile $(BINDIR)/controller-manager
	mkdir -p build/controller-manager/tmp
	cp $(BINDIR)/controller-manager build/controller-manager/tmp
	docker build -t $(CONTROLLER_MANAGER_IMAGE) build/controller-manager
	docker tag $(CONTROLLER_MANAGER_IMAGE) $(CONTROLLER_MANAGER_MUTABLE_IMAGE)
	rm -rf build/controller-manager/tmp

# Push our Docker Images to a registry
######################################
push: user-broker-push controller-manager-push apiserver-push

user-broker-push: user-broker-image
	docker push $(USER_BROKER_IMAGE)
	docker push $(USER_BROKER_MUTABLE_IMAGE)

controller-manager-push: controller-manager-image
	docker push $(CONTROLLER_MANAGER_IMAGE)
	docker push $(CONTROLLER_MANAGER_MUTABLE_IMAGE)

apiserver-push: apiserver-image
	docker push $(APISERVER_IMAGE)
	docker push $(APISERVER_MUTABLE_IMAGE)
