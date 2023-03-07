PROJECT=icepay::svc
PREFIX=$(shell pwd)
VERSION=$(shell git describe --match 'v[0-9]*'  --always)
DEFAULT_BRANCH=$(shell git symbolic-ref --short -q HEAD)

ifndef OS
	OS=linux
	unameOut=$(shell uname -s)
	ifeq ($(unameOut),Darwin)
		OS=darwin
	endif

	ifeq ($(OSTYPE),win32)
		OS=windows
	endif
endif

ifndef ARCH
	ARCH=amd64
	unameOut=$(shell uname -m)
	ifeq ($(unameOut),i386)
		ARCH=386
	endif

	ifeq ($(unameOut),i686)
		ARCH=386
	endif
endif

ifndef BRANCH
	BRANCH=$(DEFAULT_BRANCH)
endif

ifdef CI_COMMIT_REF_SLUG
	BRANCH=$(CI_COMMIT_REF_SLUG)
endif

ifndef DEPLOY_REPLICA
	DEPLOY_REPLICA=1
endif

ifndef GO
	GO=go
endif

ifndef GOFMT
	GOFMT=gofmt
endif

ifndef PROTOC
	PROTOC=protoc
endif

ifndef GIT
	GIT=git
endif

ifndef SWAG
	SWAG=swag
endif

ifndef DOCKER
	DOCKER=docker
endif

SOURCE_DIR=$(PREFIX)
BINARY_DIR=$(PREFIX)/bin
BINARY_LINKER_SVC=svc
DOCKER_TAG=useccr.ccs.tencentyun.com/icepay/svc:latest

.PHONY: all summary svc test upgrade swag doc docker push
.DEFAULT: all

# Targets
all: summary fmt doc svc

summary:
	@printf "\033[1;37m  == \033[1;32m$(PROJECT) \033[1;33m$(VERSION) \033[1;37m==\033[0m\n"
	@printf "    Platform    : \033[1;37m$(shell uname -sr)\033[0m\n"
	@printf "    Target OS   : \033[1;37m$(OS)\033[0m\n"
	@printf "    Target Arch : \033[1;37m$(ARCH)\033[0m\n"
	@printf "    Go          : \033[1;37m`$(GO) version`\033[0m\n"
	@printf "    Docker      : \033[1;37m`$(DOCKER) -v`\033[0m\n"
	@printf "    Git         : \033[1;37m$(shell git version)\033[0m\n"
	@printf "    Branch      : \033[1;37m$(BRANCH)\033[0m\n"
	@echo

fmt:
	@printf "\033[1;36m  Gofmt - Code syntax & format check\033[0m\n"
	@test -z "$$($(GOFMT) -s -l ${SOURCE_DIR} 2>&1 | tee /dev/stderr)" || (echo >&2 " - Format check failed" && false)
	@printf "\033[1;32m    ... Passed\033[0m\n"
	@echo

svc:
	@printf "\033[1;36m  Compiling $(BINARY_LINKER_SVC) ...\033[0m\n"
	@mkdir -p $(BINARY_DIR)
	@printf "    \033[1;34mTarget : \033[1;35m$(BINARY_DIR)/$(BINARY_LINKER_SVC)\033[0m\n"
	@GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 $(GO) build -a -ldflags '-extldflags "-static"' -o $(BINARY_DIR)/$(BINARY_LINKER_SVC) $(SOURCE_DIR)
	@echo

test:
	@$(GO) clean --testcache
	@GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 $(GO) test ${SOURCE_DIR}/internal/generator --cover -v
	@echo

upgrade:
	@printf "\033[1;36m  Upgrading dependences ...\033[0m\n"
	@GOOS=$(OS) GOARCH=$(ARCH) $(GO) get -u ./...
	@GOOS=$(OS) GOARCH=$(ARCH) $(GO) mod tidy
	@echo

swag:
	@$(GO) install github.com/swaggo/swag/cmd/swag@latest

doc:
	@printf "\033[1;36m  Generating swagger docs ...\033[0m\n"
	@$(SWAG) init -d ./,./handler
	@echo

docker:
	@printf "\033[1;36m  Docker build ...\033[0m\n"
	@$(DOCKER) build --rm -t $(DOCKER_TAG) .
	@echo

push:
	@printf "\033[1;36m  Docker push ...\033[0m\n"
	@$(DOCKER) push $(DOCKER_TAG)
	@echo
