/*
Copyright © 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main_test

import (
	"testing"

	main "github.com/RedHatInsights/insights-results-aggregator-exporter"
	"github.com/stretchr/testify/assert"
)

func TestInitLoggingWrongSentryDSN(t *testing.T) {
	config, err := main.LoadConfiguration("", "tests/config2")
	assert.NoError(t, err, "unexpected error loading configuration")

	closer, err := main.InitLogging(&config)
	defer closer()
	assert.Error(t, err, "expecting an error due to an invalid Sentry DSN while initializing logging")
	assert.Contains(t, err.Error(), "DsnParseError")
}

func TestInitLoggingWithSentryDSN(t *testing.T) {
	envVar := "INSIGHTS_RESULTS_AGGREGATOR_EXPORTER__SENTRY__DSN"
	mustSetEnv(t, envVar, "https://public@sentry.example.com/1")

	config, err := main.LoadConfiguration("", "tests/config2")
	assert.NoError(t, err, "unexpected error loading configuration")

	closer, err := main.InitLogging(&config)
	defer closer()
	assert.NoError(t, err, "unexpected error initializing logging")
}
