# Copyright (c) 2023 Saglabs SA. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

MFS3VER=3.0.117
MFS4VER=4.56.6
DRIVER_VERSION ?= 0.9.4
#MFS3TAGCE=$(DRIVER_VERSION)-$(MFS3VER)
#MFS3TAGPRO=$(DRIVER_VERSION)-$(MFS3VER)-pro
MFS4TAGCE=$(DRIVER_VERSION)-$(MFS4VER)
MFS4TAGPRO=$(DRIVER_VERSION)-$(MFS4VER)-pro
DEVTAG=$(DRIVER_VERSION)-dev

NAME=moosefs-csi-plugin
USERNAME=cbpowell
REPO=cbpowell
DOCKER_REGISTRY=ghcr.io

ready: clean compile
publish-dev: clean compile build-dev push-dev
publish-prod: clean compile build-prod push-prod

compile:
	@echo "==> Building the project"
	@env CGO_ENABLED=0 GOCACHE=/tmp/go-cache GOOS=linux GOARCH=amd64 go build -a -o cmd/moosefs-csi-plugin/${NAME} cmd/moosefs-csi-plugin/main.go

build-dev:
	@echo "==> Building DEV docker images"
	@docker build -t $(DOCKER_REGISTRY)/$(REPO)/moosefs-csi:$(DEVTAG) cmd/moosefs-csi-plugin

push-dev:
	@echo "==> Logging into repo"
	# @docker login --username $(USERNAME) --password ${GH_TOKEN} $(DOCKER_REGISTRY)
	@echo "==> Publishing DEV $(DOCKER_REGISTRY)/moosefs-csi-plugin:$(DEVTAG)"
	@docker push $(DOCKER_REGISTRY)/$(REPO)/moosefs-csi:$(DEVTAG)
	@echo "==> Your DEV image is now available at $(DOCKER_REGISTRY)/moosefs-csi-plugin:$(DEVTAG)"

build-prod:
	#@docker build -t $(DOCKER_REGISTRY)/$(REPO)/moosefs-csi:$(MFS3TAGCE) cmd/moosefs-csi-plugin -f cmd/moosefs-csi-plugin/Dockerfile-mfs3-ce
	#@docker build -t $(DOCKER_REGISTRY)/$(REPO)/moosefs-csi:$(MFS3TAGPRO) cmd/moosefs-csi-plugin -f cmd/moosefs-csi-plugin/Dockerfile-mfs3-pro
	@docker build -t $(DOCKER_REGISTRY)/$(REPO)/moosefs-csi:$(MFS4TAGCE) cmd/moosefs-csi-plugin -f cmd/moosefs-csi-plugin/Dockerfile-mfs4-ce
	#@docker build -t $(DOCKER_REGISTRY)/$(REPO)/moosefs-csi:$(MFS4TAGPRO) cmd/moosefs-csi-plugin -f cmd/moosefs-csi-plugin/Dockerfile-mfs4-pro

push-prod:
	@echo "==> Logging into repo"
	#@docker login --username $(USERNAME) --password $(PASSWORD) $(DOCKER_REGISTRY)
	@echo "==> Publishing $(DOCKER_REGISTRY)/moosefs-csi"
	#@docker push $(DOCKER_REGISTRY)/$(REPO)/moosefs-csi:$(MFS3TAGCE)
	#@docker push $(DOCKER_REGISTRY)/$(REPO)/moosefs-csi:$(MFS3TAGPRO)
	@docker push $(DOCKER_REGISTRY)/$(REPO)/moosefs-csi:$(MFS4TAGCE)
	#@docker push $(DOCKER_REGISTRY)/$(REPO)/moosefs-csi:$(MFS4TAGPRO)

clean:
	@echo "==> Cleaning releases"
	@GOOS=linux go clean -i -x ./...

.PHONY: clean
