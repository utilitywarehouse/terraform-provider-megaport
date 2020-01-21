SWEEP?=staging
TEST?=./...
GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')
PKG_NAME=megaport
WEBSITE_REPO=github.com/hashicorp/terraform-website

default: build

build: fmtcheck
	go install

depscheck:
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)
	@echo "==> Checking source code with go mod vendor..."
	@go mod vendor
	@git diff --compact-summary --exit-code -- vendor || \
		(echo; echo "Unexpected difference in vendor/ directory. Run 'go mod vendor' command or revert any go.mod/go.sum/vendor changes and commit."; exit 1)

docscheck:
	@echo "==> Extracting provider json schema..."
	$(eval PROVIDER_SCHEMA=$(shell PROVIDER_NAME=$(PKG_NAME) $(CURDIR)/scripts/providerjsonschema.sh))
	@echo "==> Checking docs with tfproviderdocs..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--volume "$(shell pwd)/website:/terraform-provider-megaport/website" \
		--volume "$(PROVIDER_SCHEMA):/provider-schema" \
		--workdir /terraform-provider-megaport \
		bflad/tfproviderdocs \
		check \
		-providers-schema-json=/provider-schema/schema.json
	@rm -rf $(PROVIDER_SCHEMA)

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking code against linters..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--volume $(shell pwd):/src \
		--workdir /src \
		golangci/golangci-lint:latest-alpine \
		golangci-lint \
		run \
		-v \
		--timeout=2m \
		./...

providerlint:
	@echo "==> Checking provider with tfproviderlint..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--volume $(shell pwd):/src \
		bflad/tfproviderlint:latest \
		-c 0 \
		./...

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

test: fmtcheck
ifdef TEST_COVER
	$(eval TESTARGS=$(TESTARGS) -coverprofile=cover.out)
endif
	go test $(TEST) -v $(TESTARGS) -count=1 -timeout=1m -parallel=2
ifdef TEST_COVER
	go tool cover -html=cover.out
	rm cover.out
endif

testacc: fmtcheck
ifdef TEST_COVER
	$(eval TESTARGS=$(TESTARGS) -coverprofile=cover.out)
endif
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -count=1 -timeout=5m -parallel=1
ifdef TEST_COVER
	go tool cover -html=cover.out
	rm cover.out
endif

website: website-setup
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

websitefmtcheck:
	@echo "==> Checking docs formatting..."
	@docker run \
	    --interactive \
	    --rm \
	    --tty \
	    --volume $(shell pwd)/website:/src/website \
	    --volume $(shell pwd)/scripts:/src/scripts \
	    --workdir /src \
	    --entrypoint /bin/sh \
	    node:alpine \
	    -c 'apk add --quiet --no-cache bash terraform && scripts/websitefmtcheck.sh'

website-lint:
	@echo "==> Checking docs against linters..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--volume $(shell pwd)/website:/src/website \
		--workdir /src \
		golang:alpine \
		/bin/sh -c 'GO111MODULE=on go install github.com/client9/misspell/cmd/misspell && misspell -error -source=text website'

website-setup:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@ln -sf ../../../../ext/providers/$(PKG_NAME)/website/docs $(GOPATH)/src/$(WEBSITE_REPO)/content/source/docs/providers/$(PKG_NAME)
	@ln -sf ../../../ext/providers/$(PKG_NAME)/website/$(PKG_NAME).erb $(GOPATH)/src/$(WEBSITE_REPO)/content/source/layouts/

website-test: website-setup
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build depscheck docscheck fmt fmtcheck lint providerlint sweep test testacc website websitefmtcheck website-lint website-setup website-test
