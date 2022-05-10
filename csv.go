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

package main

import (
	"encoding/csv"
	"io"
	"strconv"
)

// DisabledRulesToCSV function exports list of disabled rules + number of users
// who disabled rules to CSV file.
func DisabledRulesToCSV(buffer io.Writer, disabledRulesInfo []DisabledRuleInfo) error {
	writer := csv.NewWriter(buffer)

	var data = [][]string{{"Rule", "Count"}}

	err := writer.WriteAll(data)
	if err != nil {
		return err
	}

	for _, disabledRuleInfo := range disabledRulesInfo {
		err := writer.Write([]string{
			disabledRuleInfo.Rule,
			strconv.Itoa(disabledRuleInfo.Count)})
		if err != nil {
			return err
		}
	}

	writer.Flush()

	// check for any error during export to CSV
	err = writer.Error()
	if err != nil {
		return err
	}

	return nil
}
