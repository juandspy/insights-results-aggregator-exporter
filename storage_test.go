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

import (
	"testing"

	"github.com/stretchr/testify/assert"

	main "github.com/RedHatInsights/insights-results-aggregator-exporter"
)

// TestNewStorage checks whether constructor for new storage returns error for improper storage configuration
func TestNewStorageError(t *testing.T) {
	_, err := main.NewStorage(&main.StorageConfiguration{
		Driver: "non existing driver",
	})
	assert.EqualError(t, err, "driver non existing driver is not supported")
}

// TestNewStoragePostgreSQL function tests creating new storage with logs
func TestNewStoragePostgreSQL(t *testing.T) {
	_, err := main.NewStorage(&main.StorageConfiguration{
		Driver:        "postgres",
		PGUsername:    "user",
		PGPassword:    "password",
		PGHost:        "nowhere",
		PGPort:        1234,
		PGDBName:      "test",
		PGParams:      "",
		LogSQLQueries: true,
	})

	// we just happen to make connection without trying to actually connect
	assert.Nil(t, err)
}

// TestNewStorageSQLite3 function tests creating new storage with logs
func TestNewStorageSQLite3(t *testing.T) {
	_, err := main.NewStorage(&main.StorageConfiguration{
		Driver:        "sqlite3",
		LogSQLQueries: true,
	})

	// we just happen to make connection without trying to actually connect
	assert.Nil(t, err)
}

// TestClose function tests database close operation.
func TestClose(t *testing.T) {
	storage, err := main.NewStorage(&main.StorageConfiguration{
		Driver:        "sqlite3",
		LogSQLQueries: true,
	})

	// we just happen to make connection without trying to actually connect
	assert.Nil(t, err)

	// try to close the storage
	err = storage.Close()

	// it should not fail
	assert.Nil(t, err)
}
