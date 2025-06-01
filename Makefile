PACKAGES=$(shell go list ./... | grep -E -v 'pb$|example|proto|testdata|mock')
MOD=$(shell cat go.mod | grep ^module -m 1 | awk '{ print $$2; }' || '')

GOTEST=go
GOBUILD=go

# dependencies
GOIMPORTS=$(shell type goimports-reviser > /dev/null 2>&1 && echo $$?)
XGO=$(shell type xgo > /dev/null 2>&1 && echo $$?)

# git info
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --tags --abbrev=0)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_AT=$(shell date "+%Y%m%d%H%M%S")

show:
	@echo "packages:"
	@for item in $(PACKAGES); do echo "\t$$item"; done
	@echo "module name:"
	@echo "\t$(MOD)"
	@echo "tools:"
	@echo "\tbuild=$(BUILD_TOOL) test=$(TEST_TOOL) fmt=$(GOIMPORTS) xgo=$(XGO)"
	@echo "git:"
	@echo "\tcommit_id=$(GIT_COMMIT) tag=$(GIT_TAG) branch=$(GIT_BRANCH) build_time=$(BUILD_AT) name=$(MOD_NAME)"

# install dependencies
dependencies:
	@if [ "${GOIMPORTS}" != "0" ]; then \
		echo "installing goimports-reviser for format sources"; \
		go install -v github.com/incu6us/goimports-reviser/v3@latest; \
	fi
	@if [ "${TEST_TOOL}" == "xgo" ] && [ "${XGO}" != "0" ]; then \
		echo "installing xgo for unit test"; \
		go install github.com/xhd2015/xgo/cmd/xgo@latest; \
	fi

update:
	@echo "installing goimports-reviser for format sources"
	@go install -v github.com/incu6us/goimports-reviser/v3@latest
	@if [ "${TEST_TOOL}" == "xgo" ]; then \
		echo "installing xgo for unit test"; \
		go install github.com/xhd2015/xgo/cmd/xgo@latest; \
	fi


tidy:
	@echo "==> go mod tidy"
	@go mod tidy

cover: tidy vet
	@echo "==> run unit test with coverage"
	@go test -race -failfast -parallel 1 -gcflags="all=-N -l" ${PACKAGES} -covermode=atomic -coverprofile=cover.out

test: tidy vet
	@echo "==> run unit test"
	@go test -race -failfast -parallel 1 -gcflags="all=-N -l" ${PACKAGES}

vet:
	@go vet ${PACKAGES}

fmt:
	@echo "==> format code"
	@goimports-reviser -rm-unused -list-diff -set-alias -imports-order 'std,general,company,project' -excludes '.git/,*.pb.go' ./...

clean:
	@find . \( -name .xgo -o -name cover.out \) | xargs rm -rf
	@rm -rf build/*

report:
	@echo ">>>static checking"
	@go vet ./...
	@echo "done\n"
	@echo ">>>detecting ineffectual assignments"
	@ineffassign ./...
	@echo "done\n"
	@echo ">>>detecting icyclomatic complexities over 10 and average"
	@gocyclo -over 10 -avg -ignore '_test|vendor' . || true
	@echo "done\n"

check: tidy vet test
