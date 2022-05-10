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
	expectedVersionMessage        = "Insights Results Aggregator Cleaner version 1.0"
	expectedAuthorsMessage        = "Pavel Tisnovsky"
	expectedCopyrightMessage      = "Red Hat Inc."
	expectedConfigurationMessage1 = "Driver"
	expectedConfigurationMessage2 = "Username"
	expectedConfigurationMessage3 = "Host"
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

	assert.Contains(t, output, expectedConfigurationMessage1)
	assert.Contains(t, output, expectedConfigurationMessage2)
	assert.Contains(t, output, expectedConfigurationMessage3)
}

func checkCapture(t *testing.T, err error) {
	if err != nil {
		t.Fatal("Unable to capture standard output", err)
	}
}

// TestDoSelectedOperationShowVersion checks the function showVersion called
// via doSelectedOperation function
func TestDoSelectedOperationShowVersion(t *testing.T) {
	// stub for structures needed to call the tested function
	configuration := main.ConfigStruct{}
	cliFlags := main.CliFlags{
		ShowVersion:       true,
		ShowAuthors:       false,
		ShowConfiguration: false,
	}

	// try to call the tested function and capture its output
	output, err := capture.StandardOutput(func() {
		code, err := main.DoSelectedOperation(&configuration, cliFlags)
		assert.Equal(t, code, main.ExitStatusOK)
		assert.Nil(t, err)
	})

	// check the captured text
	checkCapture(t, err)

	assert.Contains(t, output, expectedVersionMessage)
}

// TestDoSelectedOperationShowAuthors checks the function showAuthors called
// via doSelectedOperation function
func TestDoSelectedOperationShowAuthors(t *testing.T) {
	// stub for structures needed to call the tested function
	configuration := main.ConfigStruct{}
	cliFlags := main.CliFlags{
		ShowVersion:       false,
		ShowAuthors:       true,
		ShowConfiguration: false,
	}

	// try to call the tested function and capture its output
	output, err := capture.StandardOutput(func() {
		code, err := main.DoSelectedOperation(&configuration, cliFlags)
		assert.Equal(t, code, main.ExitStatusOK)
		assert.Nil(t, err)
	})

	// check the captured text
	checkCapture(t, err)

	assert.Contains(t, output, expectedAuthorsMessage)
	assert.Contains(t, output, expectedCopyrightMessage)
}

// TestDoSelectedOperationShowConfiguration checks the function
// showConfiguration called via doSelectedOperation function
func TestDoSelectedOperationShowConfiguration(t *testing.T) {
	// stub for structures needed to call the tested function
	configuration := main.ConfigStruct{}
	cliFlags := main.CliFlags{
		ShowVersion:       false,
		ShowAuthors:       false,
		ShowConfiguration: true,
	}

	// try to call the tested function and capture its output
	output, err := capture.ErrorOutput(func() {
		log.Logger = log.Output(zerolog.New(os.Stderr))
		code, err := main.DoSelectedOperation(&configuration, cliFlags)
		assert.Equal(t, code, main.ExitStatusOK)
		assert.Nil(t, err)
	})

	// check the captured text
	checkCapture(t, err)

	assert.Contains(t, output, expectedConfigurationMesage1)
	assert.Contains(t, output, expectedConfigurationMesage2)
	assert.Contains(t, output, expectedConfigurationMesage3)
}

// TestPrintTables checks the function printTables
func TestPrintTables(t *testing.T) {
	tables := []main.TableName{
		main.TableName("first"),
		main.TableName("second"),
		main.TableName("third"),
	}

	output, err := capture.ErrorOutput(func() {
		log.Logger = log.Output(zerolog.New(os.Stderr))
		main.PrintTables(tables)
	})

	// check the captured text
	checkCapture(t, err)

	assert.Contains(t, output, "\\\"table\\\":\\\"first\\\"")
	assert.Contains(t, output, "\\\"table\\\":\\\"second\\\"")
	assert.Contains(t, output, "\\\"table\\\":\\\"third\\\"")
}

// TestParseFlags is dummy test for parseFlags function
func TestParseFlags(t *testing.T) {
	flags := main.ParseFlags()
	assert.NotNil(t, flags)
}
