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

// DBDriver type for db driver enum
type DBDriver int

// TableName type represents table name
type TableName string

// DisabledRuleInfo contains information about rules disabled by user
type DisabledRuleInfo struct {
	Rule  string
	Count int
}

// CliFlags represents structure holding all command line arguments and flags.
type CliFlags struct {
	ShowVersion         bool
	ShowAuthors         bool
	ShowConfiguration   bool
	PrintSummaryTable   bool
	Output              string
	CheckS3Connection   bool
	ExportMetadata      bool
	ExportDisabledRules bool
	ExportLog           bool
}

// M represents a map with string keys and any value
type M map[string]interface{}
