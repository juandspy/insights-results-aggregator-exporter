# insights-results-aggregator-exporter

[![GoDoc](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator-exporter?status.svg)](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator-exporter)
[![GitHub Pages](https://img.shields.io/badge/%20-GitHub%20Pages-informational)](https://redhatinsights.github.io/insights-results-aggregator-exporter/)
[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/insights-results-aggregator-exporter)](https://goreportcard.com/report/github.com/RedHatInsights/insights-results-aggregator-exporter)
[![Build Status](https://travis-ci.com/RedHatInsights/insights-results-aggregator-exporter.svg?branch=master)](https://travis-ci.com/RedHatInsights/insights-results-aggregator-exporter)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/RedHatInsights/insights-results-aggregator-exporter)
[![License](https://img.shields.io/badge/license-Apache-blue)](https://github.com/RedHatInsights/insights-results-aggregator-exporter/blob/master/LICENSE)

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
