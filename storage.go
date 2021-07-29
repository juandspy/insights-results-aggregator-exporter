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

package main

import (
	"fmt"

	"database/sql"

	_ "github.com/lib/pq"           // PostgreSQL database driver
	_ "github.com/mattn/go-sqlite3" // SQLite database driver

	"github.com/rs/zerolog/log"
)

// Driver types
const (
	// DBDriverSQLite3 shows that db driver is sqlite
	DBDriverSQLite3 DBDriver = iota
	// DBDriverPostgres shows that db driver is postgres
	DBDriverPostgres
	// DBDriverGeneral general sql(used for mock now)
	DBDriverGeneral
)

// Error messages
const (
	canNotConnectToDataStorageMessage = "Can not connect to data storage"
	unableToCloseDBRowsHandle         = "Unable to close the DB rows handle"
	sqlStatementExecutionError        = "SQL statement execution error"
)

// SQL statements
const (
	// Select all public tables from open database
	selectListOfTables = `
           SELECT tablename
             FROM pg_catalog.pg_tables
            WHERE schemaname != 'information_schema'
              AND schemaname != 'pg_catalog';
   `
)

// Storage represents an interface to almost any database or storage system
type Storage interface {
	Close() error

	ReadListOfTables() ([]TableName, error)
	ReadTable(tableName string) error
}

// DBStorage is an implementation of Storage interface that use selected SQL like database
// like SQLite, PostgreSQL, MariaDB, RDS etc. That implementation is based on the standard
// sql package. It is possible to configure connection via Configuration structure.
// SQLQueriesLog is log for sql queries, default is nil which means nothing is logged
type DBStorage struct {
	connection   *sql.DB
	dbDriverType DBDriver
}

// NewStorage function creates and initializes a new instance of Storage interface
func NewStorage(configuration StorageConfiguration) (*DBStorage, error) {
	log.Info().Msg("Initializing connection to storage")

	driverType, driverName, dataSource, err := initAndGetDriver(configuration)
	if err != nil {
		log.Error().Err(err).Msg("Unsupported driver")
		return nil, err
	}

	log.Info().
		Str("driver", driverName).
		Str("datasource", dataSource).
		Msg("Making connection to data storage")

	// prepare connection
	connection, err := sql.Open(driverName, dataSource)
	if err != nil {
		log.Error().Err(err).Msg("Can not connect to data storage")
		return nil, err
	}

	log.Info().Msg("Connection to storage established")
	return NewFromConnection(connection, driverType), nil
}

// NewFromConnection function creates and initializes a new instance of Storage interface from prepared connection
func NewFromConnection(connection *sql.DB, dbDriverType DBDriver) *DBStorage {
	return &DBStorage{
		connection:   connection,
		dbDriverType: dbDriverType,
	}
}

// initAndGetDriver initializes driver(with logs if logSQLQueries is true),
// checks if it's supported and returns driver type, driver name, dataSource and error
func initAndGetDriver(configuration StorageConfiguration) (driverType DBDriver, driverName string, dataSource string, err error) {
	//var driver sql_driver.Driver
	driverName = configuration.Driver

	switch driverName {
	case "sqlite3":
		driverType = DBDriverSQLite3
		//driver = &sqlite3.SQLiteDriver{}
		// dataSource = configuration.SQLiteDataSource
	case "postgres":
		driverType = DBDriverPostgres
		//driver = &pq.Driver{}
		dataSource = fmt.Sprintf(
			"postgresql://%v:%v@%v:%v/%v?%v",
			configuration.PGUsername,
			configuration.PGPassword,
			configuration.PGHost,
			configuration.PGPort,
			configuration.PGDBName,
			configuration.PGParams,
		)
	default:
		err = fmt.Errorf("driver %v is not supported", driverName)
		return
	}

	return
}

// Close method closes the connection to database. Needs to be called at the
// end of application lifecycle.
func (storage DBStorage) Close() error {
	log.Info().Msg("Closing connection to data storage")
	if storage.connection != nil {
		err := storage.connection.Close()
		if err != nil {
			log.Error().Err(err).Msg("Can not close connection to data storage")
			return err
		}
	}
	return nil
}

// Read list of tables reads names of all public tables stored in opened
// database.
func (storage DBStorage) ReadListOfTables() ([]TableName, error) {
	// slice to make list of tables
	var tableList = make([]TableName, 0)

	rows, err := storage.connection.Query(selectListOfTables)
	if err != nil {
		return tableList, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg(unableToCloseDBRowsHandle)
		}
	}()

	// read all table names
	for rows.Next() {
		var tableName TableName

		err := rows.Scan(&tableName)
		if err != nil {
			if closeErr := rows.Close(); closeErr != nil {
				log.Error().Err(closeErr).Msg(unableToCloseDBRowsHandle)
			}
			return tableList, err
		}
		tableList = append(tableList, tableName)
	}

	return tableList, nil
}

func (storage DBStorage) ReadTable(tableName TableName) error {
	// it is not possible to use parameter for table name or a key
	// disable "G201 (CWE-89): SQL string concatenation (Confidence: HIGH, Severity: MEDIUM)"
	// #nosec G201
	sqlStatement := fmt.Sprintf("SELECT * FROM %s", tableName)
	log.Info().Str("SQL statement", sqlStatement).Msg("Performing")

	rows, err := storage.connection.Query(sqlStatement)
	if err != nil {
		log.Error().Err(err).Msg(sqlStatementExecutionError)
		return err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg(unableToCloseDBRowsHandle)
		}
	}()

	columns, err := rows.Columns()

	log.Info().Str("table", string(tableName)).Int("columns", len(columns)).Msg("table metadata")

	// prepare data structure to hold raw values
	values := make([]interface{}, len(columns))
	for i, _ := range columns {
		values[i] = new(sql.RawBytes)
	}

	// iterate over all rows
	for rows.Next() {
		// read raw values
		err = rows.Scan(values...)
		if err != nil {
			log.Error().Err(err).Msg("Unable to scan row")
		}
		// it is now needed to check each element of values for nil
		// then to use type introspection and type assertion to be
		// able to fetch the column into a typed variable if needed
	}
	return nil
}
