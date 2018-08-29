.PHONY: help deps generate qa test coverage format fmtcheck vet lint cyclo ineffassign misspell structcheck varcheck errcheck megacheck astscan check

SHELL=/bin/bash -o pipefail
DIR:=$(shell dirname $(lastword $(MAKEFILE_LIST)))

ifeq ($(GOPATH),)
	# extract the GOPATH
	GOPATH=$(shell go env GOPATH)
endif
export PATH := $(GOPATH)/bin:$(PATH)

DIR_RECURSE=$(DIR)/...

PROJECT=conform

GOENV=GOPATH=$(GOPATH) CGO_CFLAGS_ALLOW='-w'
 
# --- MAKE TARGETS ---

# Display general help about this command
help:
	@echo ""
	@echo "$(PROJECT) Makefile."
	@echo "GOPATH=$(GOPATH)"
	@echo "The following commands are available:"
	@echo ""
	@echo "    make deps        : Get the dependencies"
	@echo "    make generate    : Generate any required source code"
	@echo "    make qa          : Run all the tests and static analysis reports"
	@echo "    make test        : Run the unit tests"
	@echo "    make coverage    : Generate the coverage report"
	@echo ""
	@echo "    make format      : Format the source code"
	@echo "    make fmtcheck    : Check if the source code has been formatted"
	@echo "    make vet         : Check for suspicious constructs"
	@echo "    make lint        : Check for style errors"
	@echo "    make cyclo       : Generate the cyclomatic complexity report"
	@echo "    make ineffassign : Detect ineffectual assignments"
	@echo "    make misspell    : Detect commonly misspelled words in source files"
	@echo "    make structcheck : Find unused struct fields"
	@echo "    make varcheck    : Find unused global variables and constants"
	@echo "    make errcheck    : Check that error return values are used"
	@echo "    make megacheck   : Runs staticcheck, gosimple and unusued checks"
	@echo "    make astscan     : GO AST scanner"
	@echo "    make check       : Run all static code analysis checks"
	@echo ""

# Get the dependencies
deps:
	$(GOENV) go get ./... github.com/axw/gocov/gocov \
	github.com/client9/misspell/cmd/misspell \
	github.com/fzipp/gocyclo \
	github.com/securego/gosec/cmd/gosec \
	github.com/golang/lint/golint \
	github.com/gordonklaus/ineffassign \
	github.com/jstemmer/go-junit-report \
	github.com/kisielk/errcheck \
	github.com/opennota/check/cmd/structcheck \
	github.com/opennota/check/cmd/varcheck \
	github.com/stretchr/testify \
	honnef.co/go/tools/cmd/megacheck

# Generates any required source code
generate:
	go generate $(DIR_RECURSE)

# Alias to run targets to perform all static checks and tests
qa: check test coverage

# Run the unit tests
test: generate
	@mkdir -p target/test
	@mkdir -p target/report
	$(GOENV) go test \
	-covermode=atomic \
	-bench=. \
	-race \
	-cpuprofile=target/report/cpu.out \
	-memprofile=target/report/mem.out \
	-mutexprofile=target/report/mutex.out \
	-coverprofile=target/report/coverage.out \
	$(DIR) | \
	tee >(go-junit-report > target/test/report.xml)
	rm $(PROJECT).test

# Generate the coverage report
coverage:
	@mkdir -p target/report
	$(GOENV) go tool cover -html=target/report/coverage.out -o target/report/coverage.html

# Format the source code
format:
	@find . -type f -name "*.go" -exec gofmt -s -w {} \;

# Check if the source code has been formatted
fmtcheck:
	@mkdir -p target
	@find . -type f -name "*.go" -exec gofmt -s -d {} \; | tee target/format.diff
	@test ! -s target/format.diff || { echo "ERROR: the source code has not been formatted - please use 'make format' or 'gofmt'"; exit 1; }

# Check for syntax errors
vet:
	$(GOENV) go vet

# Check for style errors
lint:
	$(GOENV) golint -set_exit_status $(DIR_RECURSE)

# Report cyclomatic complexity
cyclo:
	@mkdir -p target/report
	$(GOENV) gocyclo -over 12 -avg $(DIR) | tee target/report/cyclo.txt

# Detect ineffectual assignments
ineffassign:
	@mkdir -p target/report
	$(GOENV) ineffassign $(DIR) | tee target/report/ineffassign.txt

# Detect commonly misspelled words in source files
misspell:
	@mkdir -p target/report
	$(GOENV) misspell -error $(DIR_RECURSE)  | tee target/report/misspell.txt

# Find unused struct fields.
structcheck:
	@mkdir -p target/report
	$(GOENV) structcheck $(DIR_RECURSE)  | tee target/report/structcheck.txt

# Find unused global variables and constants.
varcheck:
	@mkdir -p target/report
	$(GOENV) varcheck -e $(DIR_RECURSE)  | tee target/report/varcheck.txt

# Check that error return values are used.
errcheck:
	@mkdir -p target/report
	$(GOENV) errcheck $(DIR_RECURSE)  | tee target/report/errcheck.txt

# staticcheck, gosimple, unusued
megacheck:
	@mkdir -p target/report
	$(GOENV) megacheck -tests -simple.exit-non-zero=true -staticcheck.exit-non-zero=true -unused.exit-non-zero=true  $(DIR_RECURSE) | tee target/report/megacheck.txt
	$(GOENV) megacheck -tests=false -simple.enabled=false -staticcheck.enabled=false -unused.exit-non-zero=true $(DIR_RECURSE) | tee target/report/megacheck.txt

# AST scanner
astscan:
	@mkdir -p target/report
	$(GOENV) gosec $(DIR_RECURSE) | tee target/report/astscan.txt

# Alias to run static code analysis checks
check: fmtcheck vet lint cyclo ineffassign misspell structcheck varcheck errcheck megacheck astscan
