/*
Copyright © 2021, 2022 Red Hat, Inc.

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

// Unit test definitions for functions and methods defined in source file
// exporter.go

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/tisnik/go-capture"

	main "github.com/RedHatInsights/insights-results-aggregator-exporter"
)

const (
	expectedVersionMessage       = "Insights Results Aggregator Cleaner version 1.0"
	expectedAuthorsMessage       = "Pavel Tisnovsky"
	expectedCopyrightMessage     = "Red Hat Inc."
	expectedConfigurationMesage1 = "Driver"
	expectedConfigurationMesage2 = "Username"
	expectedConfigurationMesage3 = "Host"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// TestShowVersion checks the function showVersion
func TestShowVersion(t *testing.T) {
	// try to call the tested function and capture its output
	output, err := capture.StandardOutput(func() {
		main.ShowVersion()
	})

	// check the captured text
	checkCapture(t, err)

	assert.Contains(t, output, expectedVersionMessage)
}

// TestShowAuthors checks the function showAuthors
func TestShowAuthors(t *testing.T) {
	// try to call the tested function and capture its output
	output, err := capture.StandardOutput(func() {
		main.ShowAuthors()
	})

	// check the captured text
	checkCapture(t, err)

	assert.Contains(t, output, expectedAuthorsMessage)
	assert.Contains(t, output, expectedCopyrightMessage)
}

// TestShowConfiguration checks the function ShowConfiguration
func TestShowConfiguration(t *testing.T) {
	// fill in configuration structure
	configuration := main.ConfigStruct{}
	// try to call the tested function and capture its output
	output, err := capture.ErrorOutput(func() {
		log.Logger = log.Output(zerolog.New(os.Stderr))
		main.ShowConfiguration(&configuration)
	})

	// check the captured text
	checkCapture(t, err)

	assert.Contains(t, output, expectedConfigurationMesage1)
	assert.Contains(t, output, expectedConfigurationMesage2)
	assert.Contains(t, output, expectedConfigurationMesage3)
}

func checkCapture(t *testing.T, err error) {
	if err != nil {
		t.Fatal("Unable to capture standard output", err)
	}
}
