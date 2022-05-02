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

package main

import (
	"encoding/csv"
	"os"

	"github.com/rs/zerolog/log"
)

// storeTableNamesIntoFile function stores names of all tables into the
// specified file
func storeTableNamesIntoFile(fileName string, tableNames []TableName) error {
	// open new CSV file to be filled in
	fout, err := os.Create(fileName)
	if err != nil {
		return err
	}

	// initialize CSV writer
	writer := csv.NewWriter(fout)
	var data = [][]string{{"Table name"}}

	// header
	err = writer.WriteAll(data)
	if err != nil {
		return err
	}

	// table names
	for _, tableName := range tableNames {
		err := writer.Write([]string{string(tableName)})
		if err != nil {
			log.Error().Err(err).Msg("Write to CSV")
		}
	}

	writer.Flush()

	// check for any error during export to CSV
	err = writer.Error()
	if err != nil {
		return err
	}

	// close the file and check if close operation was ok
	err = fout.Close()
	if err != nil {
		return err
	}

	return nil
}
