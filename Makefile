# Copyright © 2021 Kris Nóva <kris@nivenly.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#   ███╗   ██╗ █████╗ ███╗   ███╗██╗
#   ████╗  ██║██╔══██╗████╗ ████║██║
#   ██╔██╗ ██║███████║██╔████╔██║██║
#   ██║╚██╗██║██╔══██║██║╚██╔╝██║██║
#   ██║ ╚████║██║  ██║██║ ╚═╝ ██║███████╗
#   ╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝
#

all: compile

# Global release version.
# Change this to bump the build version!
version="1.0.2"

compile: ## Compile for the local architecture ⚙
	@echo "Compiling..."
	./_embed.sh
	go build -ldflags "-X 'github.com/kris-nova/naml.Version=$(version)'" -o naml cmd/*.go

install: ## Install your naml 🎉
	@echo "Installing..."
	sudo cp naml /usr/local/bin/naml

test: clean compile install ## 🤓 Test is used to test your naml
	@echo "Testing..."
	go test -v ./...

clean: ## Clean your artifacts 🧼
	@echo "Cleaning..."
	rm -rf release
	rm -rf embed_*
	rm -rf naml
	rm -rf app
	rm -rf out/*
	rm -rf tmp/*

.PHONY: release
release: ## Make the binaries for a GitHub release 📦
	mkdir -p release
	GOOS="linux" GOARCH="amd64" go build -ldflags "-X 'github.com/kris-nova/naml.Version=$(version)'" -o release/naml-linux-amd64 cmd/*.go
	GOOS="linux" GOARCH="arm" go build -ldflags "-X 'github.com/kris-nova/naml.Version=$(version)'" -o release/naml-linux-arm cmd/*.go
	GOOS="linux" GOARCH="arm64" go build -ldflags "-X 'github.com/kris-nova/naml.Version=$(version)'" -o release/naml-linux-arm64 cmd/*.go
	GOOS="linux" GOARCH="386" go build -ldflags "-X 'github.com/kris-nova/naml.Version=$(version)'" -o release/naml-linux-386 cmd/*.go
	GOOS="darwin" GOARCH="amd64" go build -ldflags "-X 'github.com/kris-nova/naml.Version=$(version)'" -o release/naml-darwin-amd64 cmd/*.go

.PHONY: aur
aur: ## Publish to AUR using my SSH key
	@echo "Publishing to AUR using Kris Nóva's key (if exists)..."
	cd aur && ./publish

.PHONY: help
help:  ## 🤔 Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'
