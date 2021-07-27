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
)

// initDatabaseConnection initializes driver, checks if it's supported and
// initializes connection to the storage.
func initDatabaseConnection(configuration StorageConfiguration) (*sql.DB, error) {
	driverName := configuration.Driver
	dataSource := ""
	log.Info().Str("driverName", configuration.Driver).Msg("DB connection configuration")

	// initialize connection into selected database using the right driver
	switch driverName {
	case "sqlite3":
		//driverType := DBDriverSQLite3
		//driver = &sqlite3.SQLiteDriver{}
		dataSource = configuration.SQLiteDataSource
	case "postgres":
		//driverType := DBDriverPostgres
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
		err := fmt.Errorf("driver %v is not supported", driverName)
		log.Err(err).Msg(canNotConnectToDataStorageMessage)
		return nil, err
	}

	// try to initialize connection to the storage
	connection, err := sql.Open(driverName, dataSource)

	// check if establishing a connection was successful
	if err != nil {
		log.Err(err).Msg(canNotConnectToDataStorageMessage)
		return nil, err
	}

	return connection, nil
}
