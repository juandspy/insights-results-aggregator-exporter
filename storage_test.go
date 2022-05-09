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
	"errors"
	"testing"

	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
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

// mustCreateMockConnection function tries to create a new mock connection and
// checks if the operation was finished without problems.
func mustCreateMockConnection(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	// try to initialize new mock connection
	connection, mock, err := sqlmock.New()

	// check the status
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return connection, mock
}

// checkConnectionClose function perform mocked DB closing operation and checks
// if the connection is properly closed from unit tests.
func checkConnectionClose(t *testing.T, connection *sql.DB) {
	// connection to mocked DB needs to be closed properly
	err := connection.Close()

	// check the error status
	if err != nil {
		t.Fatalf("error during closing connection: %v", err)
	}
}

// checkAllExpectations function checks if all database-related operations have
// been really met.
func checkAllExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	// check if all expectations were met
	err := mock.ExpectationsWereMet()

	// check the error status
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Expected queries
const (
	readRecordCountQuery = "SELECT count\\(\\*\\) FROM TESTED_TABLE"
)

// check the function ReadRecordCount
func TestReadRecordCount(t *testing.T) {
	// prepare new mocked connection to database
	connection, mock := mustCreateMockConnection(t)

	// prepare mocked result for SQL query
	rowsCount := sqlmock.NewRows([]string{"count"})
	expected := 100
	rowsCount.AddRow(expected)

	// expected query performed by tested function
	mock.ExpectQuery(readRecordCountQuery).WillReturnRows(rowsCount)
	mock.ExpectClose()

	// prepare connection to mocked database
	storage := main.NewFromConnection(connection, 1)

	// call the tested method
	count, err := storage.ReadRecordsCount("TESTED_TABLE")
	if err != nil {
		t.Errorf("error was not expected %s", err)
	}

	if count != expected {
		t.Errorf("wrong number records returned: %d", count)
	}

	// connection to mocked DB needs to be closed properly
	checkConnectionClose(t, connection)

	// check if all expectations were met
	checkAllExpectations(t, mock)
}

func TestReadRecordCountScanError(t *testing.T) {
	// prepare new mocked connection to database
	connection, mock := mustCreateMockConnection(t)

	// prepare mocked result for SQL query
	rowsCount := sqlmock.NewRows([]string{"count"})
	rowsCount.AddRow("this is not integer")

	// expected query performed by tested function
	mock.ExpectQuery(readRecordCountQuery).WillReturnRows(rowsCount)
	mock.ExpectClose()

	// prepare connection to mocked database
	storage := main.NewFromConnection(connection, 1)

	// call the tested method
	_, err := storage.ReadRecordsCount("TESTED_TABLE")
	if err == nil {
		t.Errorf("error is expected")
	}

	// connection to mocked DB needs to be closed properly
	checkConnectionClose(t, connection)

	// check if all expectations were met
	checkAllExpectations(t, mock)
}

func TestReadRecordCountOnError(t *testing.T) {
	// error to be thrown
	mockedError := errors.New("mocked error")

	// prepare new mocked connection to database
	connection, mock := mustCreateMockConnection(t)

	// expected query performed by tested function
	mock.ExpectQuery(readRecordCountQuery).WillReturnError(mockedError)
	mock.ExpectClose()

	// prepare connection to mocked database
	storage := main.NewFromConnection(connection, 1)

	// call the tested method
	count, err := storage.ReadRecordsCount("TESTED_TABLE")
	if err != mockedError {
		t.Errorf("different error was returned: %v", err)
	}

	if count != -1 {
		t.Errorf("wrong number records returned: %d", count)
	}

	// connection to mocked DB needs to be closed properly
	checkConnectionClose(t, connection)

	// check if all expectations were met
	checkAllExpectations(t, mock)
}
