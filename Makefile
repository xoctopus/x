PACKAGES=$(shell go list ./... | grep -E -v 'pb$|testdata|mock|proto|example|testx/internal')
MOD=$(shell cat go.mod | grep ^module -m 1 | awk '{ print $$2; }' || '')
MOD_NAME=$(shell basename $(MOD))

GOTEST=go
GOBUILD=go

# dependencies
DEP_FMT=$(shell type goimports-reviser > /dev/null 2>&1 && echo $$?)
DEP_XGO=$(shell type xgo > /dev/null 2>&1 && echo $$?)
DEP_INEFFASSIGN=$(shell type ineffassign > /dev/null 2>&1 && echo $$?)
DEP_GOCYCLO=$(shell type gocyclo > /dev/null 2>&1 && echo $$?)
DEP_LINTER=$(shell type golangci-lint > /dev/null 2>&1 && echo $$?)

# git info
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --tags --abbrev=0)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_AT=$(shell date "+%Y%m%d%H%M%S")

show:
	@echo "packages:"
	@for item in $(PACKAGES); do echo "    $$item"; done
	@echo "module:"
	@echo "    path=$(MOD)"
	@echo "    module=$(MOD_NAME)"
	@echo "tools:"
	@echo "    build=$(GOBUILD)"
	@echo "    test=$(GOTEST)"
	@echo "    goimports-reviser=$(shell which goimports-reviser)"
	@echo "    xgo=$(shell which xgo)"
	@echo "    ineffassign=$(shell which ineffassign)"
	@echo "    gocyclo=$(shell which gocyclo)"
	@echo "git:"
	@echo "    commit_id=$(GIT_COMMIT)"
	@echo "    tag=$(GIT_TAG)"
	@echo "    branch=$(GIT_BRANCH)"
	@echo "    build_time=$(BUILD_AT)"
	@echo "    name=$(MOD_NAME)"

# install dependencies
dep:
	@echo "==> installing dependencies"
	@if [ "${DEP_FMT}" != "0" ]; then \
		echo "    goimports-reviser for format sources"; \
		go install github.com/incu6us/goimports-reviser/v3@latest; \
	fi
	@if [ "${GOTEST}" = "xgo" ] && [ "${DEP_XGO}" != "0" ]; then \
		echo "    xgo for unit test"; \
		go install github.com/xhd2015/xgo/cmd/xgo@latest; \
	fi
	@if [ "${DEP_INEFFASSIGN}" != "0" ]; then \
		echo "    ineffassign for detecting ineffectual assignments"; \
		go install github.com/gordonklaus/ineffassign@latest; \
	fi
	@if [ "${DEP_GOCYCLO}" != "0" ]; then \
		echo "\tgocyclo for calculating cyclomatic complexities of functions"; \
		go install github.com/fzipp/gocyclo/cmd/gocyclo@latest; \
	fi
	@if [ "${DEP_LINTER}" != "0" ]; then \
		echo "\tgolangci-lint for code static checking"; \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest ;\
	fi

upgrade-dep:
	@echo "==> upgrading dependencies"
	@echo "    goimports-reviser for format sources"
	@go install github.com/incu6us/goimports-reviser/v3@latest
	@echo "    ineffassign for detecting ineffectual assignments"
	@go install github.com/gordonklaus/ineffassign@latest
	@echo "    gocyclo for calculating cyclomatic complexities of functions"
	@go install github.com/gordonklaus/ineffassign@latest
	@echo "    golangci-lint for code static checking"
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@if [ "${GOTEST}" = "xgo" ]; then \
		echo "    xgo for unit test"; \
		go install github.com/xhd2015/xgo/cmd/xgo@latest; \
	fi

update:
	@go get -u all

tidy:
	@echo "==> tidy"
	@go mod tidy

test: dep tidy
	@echo "==> run unit test"
	@if [ "${GOTEST}" = "xgo" ]; then \
		$(GOTEST) test -failfast -parallel 1 -gcflags="all=-N -l" ${PACKAGES}; # xgo mock functions may cause data race \
	else \
		$(GOTEST) test -race -failfast -parallel 1 -gcflags="all=-N -l" ${PACKAGES}; \
	fi

cover: dep tidy
	@echo "==> run unit test with coverage"
	@$(GOTEST) test -failfast -parallel 1 -gcflags="all=-N -l" ${PACKAGES} -covermode=count -coverprofile=cover.out

ci-cover:
	@if [ "${GOTEST}" = "xgo" ]; then \
		go install github.com/xhd2015/xgo/cmd/xgo@latest; \
	fi
	@$(GOTEST) test -failfast -parallel 1 -gcflags="all=-N -l" ${PACKAGES} -covermode=count -coverprofile=cover.out

view-cover: cover
	@echo "==> run unit test with coverage and view"
	@go tool cover -html cover.out

fmt: dep clean
	@echo "==> format code"
	@goimports-reviser -rm-unused -set-alias -output write \
		-imports-order 'std,general,company,project' \
		-excludes '.git/,.xgo/,*.pb.go,*_generated.go' ./...

lint: dep
	@echo "==> static check"
	@echo "    >>>static checking"
	@go vet ./...
	@echo "    done"
	@echo "    >>>detecting ineffectual assignments"
	@ineffassign ./...
	@echo "    done"
	@echo "    >>>detecting cyclomatic complexities over 10 and average"
	@gocyclo -over 10 -avg -ignore '_test|_test.go|vendor|pb' . || true
	@echo "    done"

pre-commit: dep lint fmt cover clean

clean:
	@find . -name cover.out | xargs rm -rf
	@find . -name .xgo | xargs rm -rf
	@rm -rf build/*
