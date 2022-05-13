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
	"bytes"
	"testing"

	main "github.com/RedHatInsights/insights-results-aggregator-exporter"

	"github.com/stretchr/testify/assert"
)

// TestDisabledRulesToCSVNilBuffer check how nil buffer is handled by
// DisabledRulesToCSV function
func TestDisabledRulesToCSVNilBuffer(t *testing.T) {
	// empty list
	disabledRules := []main.DisabledRuleInfo{}

	err := main.DisabledRulesToCSV(nil, disabledRules)
	assert.Error(t, err, "Buffer is nil")
}

// TestDisabledRulesToCSVEmptyListOfRules check exporting empty list of
// disabled rules into CSV
func TestDisabledRulesToCSVEmptyListOfRules(t *testing.T) {
	// buffer
	buffer := new(bytes.Buffer)

	// empty list
	disabledRules := []main.DisabledRuleInfo{}

	err := main.DisabledRulesToCSV(buffer, disabledRules)
	assert.Nil(t, err, "Error is not expected")

	content := buffer.String()
	expected := "Rule,Count\n"
	assert.Equal(t, expected, content)
}

// TestDisabledRulesToCSVE check exporting non-empty list of disabled rules
// into CSV
func TestDisabledRulesToCSV(t *testing.T) {
	// buffer
	buffer := new(bytes.Buffer)

	// empty list
	disabledRules := []main.DisabledRuleInfo{
		main.DisabledRuleInfo{"first", 1},
		main.DisabledRuleInfo{"second", 2},
		main.DisabledRuleInfo{"third", 3},
	}

	err := main.DisabledRulesToCSV(buffer, disabledRules)
	assert.Nil(t, err, "Error is not expected")

	content := buffer.String()
	expected := "Rule,Count\nfirst,1\nsecond,2\nthird,3\n"
	assert.Equal(t, expected, content)
}
