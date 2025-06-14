GO ?= go
GOBIN ?= $$($(GO) env GOPATH)/bin
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.56.2

APP_ID ?= activate.app
APP_NAME ?= activate
BUILD_NUMBER ?= 1
RELEASE ?= true
# TARGET_OS ?= android/arm64
TARGET_OS ?= darwin
APP_VERSION ?= "$(shell git describe --tags --abbrev=0 | cut -c2-)"
COMMIT_HASH ?= "$(shell git describe --long --dirty --always --match "" || true)"

.PHONY: lint
lint: linter
	$(GOLANGCI_LINT) run

.PHONY: linter
linter:
	test -f $(GOLANGCI_LINT) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$($(GO) env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

.PHONY: get-fyne
get-fyne:
	go get fyne.io/fyne/v2/cmd/fyne@latest

.PHONY: package
package:
	fyne package -os ${TARGET_OS} -appID ${APP_ID} -name ${APP_NAME}  -appVersion ${APP_VERSION} -appBuild=${BUILD_NUMBER} -release=${RELEASE} -metadata commithash=${COMMIT_HASH}
