# Change these variables as necessary.
MAIN_PACKAGE_PATH := .
BINARY_NAME := main

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	git diff --exit-code

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## run/live: run the application with reloading on file changes
.PHONY: run/live/iotest
run/live/iotest:
	go mod tidy;go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build/iotest" --build.bin "/tmp/bin/iotest" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"

## build/ssh: build the ssh command
.PHONY: build/iotest
build/iotest:
	go build -o=/tmp/bin/iotest cmd/iotest/main.go

## build/ssh: build the ssh command
.PHONY: build/ssh
build/ssh:
	go build -o=/tmp/bin/ssh cmd/ssh/main.go

## run/ssh: run the ssh command
.PHONY: run/ssh
run/ssh: build/ssh
	/tmp/bin/ssh

## run/live/ssh: run the ssh command with reloading on file changes
.PHONY: run/live/ssh
run/live/ssh:
	go mod tidy;go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build/ssh" --build.bin "/tmp/bin/ssh" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"

