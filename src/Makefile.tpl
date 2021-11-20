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
version=$(shell git rev-parse HEAD)

compile: ## Compile for the local architecture ⚙
	@echo "Compiling..."
	go build -o app .

install: ## Install your naml 🎉
	@echo "Installing..."
	sudo cp app /usr/local/bin/app

test: clean compile install ## 🤓 Test is used to test your naml
	@echo "Testing..."
	go test -v ./...

clean: ## Clean your artifacts 🧼
	@echo "Cleaning..."
	rm -rvf app
	
release: ## Make the binaries for a GitHub release 📦
	mkdir -p release
	GOOS="linux" GOARCH="amd64" go build -o release/app-linux-amd64 main.go
	GOOS="linux" GOARCH="arm" go build -o release/app-linux-arm main.go
	GOOS="linux" GOARCH="arm64" go build -o release/app-linux-arm64 main.go
	GOOS="linux" GOARCH="386" go build -o release/app-linux-386 main.go
	GOOS="darwin" GOARCH="amd64" go build -o release/app-darwin-amd64 main.go


.PHONY: help
help:  ## 🤔 Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'
