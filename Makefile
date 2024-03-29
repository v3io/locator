# Copyright 2019 Iguazio
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
#

LABEL ?= unstable
REPOSITORY ?= gcr.io/iguazio
IMAGE = $(REPOSITORY)/locator:$(LABEL)

.PHONY: build
build:
	@docker build \
		--file cmd/locator/Dockerfile \
		--tag=$(IMAGE) \
		.

.PHONY: push
push:
	docker push $(IMAGE)

.PHONY: lint
lint:
	./hack/lint/install.sh
	./hack/lint/run.sh

.PHONY: fmt
fmt:
	@go fmt $(shell go list ./... | grep -v /vendor/)

.PHONY: test
test:
	go test -p1 -v ./pkg/...
