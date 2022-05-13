/*
Copyright © 2022 Red Hat, Inc.

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
	"io/ioutil"
	"os"
	"testing"

	main "github.com/RedHatInsights/insights-results-aggregator-exporter"

	"github.com/stretchr/testify/assert"
)

// mustCreateTemporaryDirectory helper function creates temporary directory
// that will be cleaned up after tests
func mustCreateTemporaryDirectory(t *testing.T) string {
	directory, err := ioutil.TempDir(os.TempDir(), "exporter")
	if err != nil {
		t.Fatal(err)
	}

	return directory
}

// mustReadFile helper function tries to read specified file and return its
// content as a string
func mustReadFile(t *testing.T, filename string) string {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	return string(fileContent)
}
