# insights-results-aggregator-exporter
Exporter for Insights Results data stored by Insights Results Aggregator

## Usage

```
Usage of ./irae:
  -authors
        show authors
  -output string
        output to: CSV, S3
  -show-configuration
        show configuration
  -summary
        print summary table after export
  -version
        show version

```

## Makefile usage

```
Usage: make <OPTIONS> ... <TARGETS>

Available targets are:

clean                Run go clean
build                Keep this rule for compatibility
fmt                  Run go fmt -w for all sources
lint                 Run golint
vet                  Run go vet. Report likely mistakes in source code
cyclo                Run gocyclo
ineffassign          Run ineffassign checker
shellcheck           Run shellcheck
errcheck             Run errcheck
goconst              Run goconst checker
gosec                Run gosec checker
abcgo                Run ABC metrics checker
style                Run all the formatting related commands (fmt, vet, lint, cyclo) + check shell scripts
run                  Build the project and executes the binary
test                 Run the unit tests
bdd_tests            Run BDD tests
before_commit        Checks done before commit
help                 Show this help screen
```
