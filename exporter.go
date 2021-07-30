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
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Messages
const (
	versionMessage         = "Insights Results Aggregator Cleaner version 1.0"
	authorsMessage         = "Pavel Tisnovsky, Red Hat Inc."
	operationFailedMessage = "Operation failed"
)

// Exit codes
const (
	// ExitStatusOK means that the tool finished with success
	ExitStatusOK = iota

	// ExitStatusStorageError is returned in case of any consumer-related error
	ExitStatusStorageError
)

const (
	configFileEnvVariableName = "INSIGHTS_RESULTS_AGGREGATOR_EXPORTER_CONFIG_FILE"
	defaultConfigFileName     = "config"
)

// showVersion function displays version information.
func showVersion() {
	fmt.Println(versionMessage)
}

// showAuthors function displays information about authors.
func showAuthors() {
	fmt.Println(authorsMessage)
}

// showConfiguration function displays actual configuration.
func showConfiguration(config ConfigStruct) {
	storageConfig := GetStorageConfiguration(config)
	log.Info().
		Str("Driver", storageConfig.Driver).
		Str("DB Name", storageConfig.PGDBName).
		Str("Username", storageConfig.PGUsername). // password is omitted on purpose
		Str("Host", storageConfig.PGHost).
		Int("Port", storageConfig.PGPort).
		Bool("LogSQLQueries", storageConfig.LogSQLQueries).
		Msg("Storage configuration")

	loggingConfig := GetLoggingConfiguration(config)
	log.Info().
		Str("Level", loggingConfig.LogLevel).
		Bool("Pretty colored debug logging", loggingConfig.Debug).
		Msg("Logging configuration")

	s3Configuration := GetS3Configuration(config)
	log.Info().
		Str("Type", s3Configuration.Type).
		Str("URL", s3Configuration.EndpointURL).
		Uint("Port", s3Configuration.EndpointPort).
		Str("AccessKeyID", s3Configuration.AccessKeyID).
		Str("SecretAccessKey", s3Configuration.SecretAccessKey).
		Bool("Use SSL", s3Configuration.UseSSL).
		Str("Bucket", s3Configuration.Bucket).
		Msg("S3 configuration")
}

// performDataExport function exports all data into selected output
func performDataExport(configuration ConfigStruct, cliFlags CliFlags) (int, error) {
	// prepare the storage
	storageConfiguration := GetStorageConfiguration(configuration)
	storage, err := NewStorage(storageConfiguration)
	if err != nil {
		log.Err(err).Msg(operationFailedMessage)
		return ExitStatusStorageError, err
	}

	tableNames, err := storage.ReadListOfTables()
	if err != nil {
		log.Err(err).Msg(operationFailedMessage)
		return ExitStatusStorageError, err
	}

	log.Info().Int("count", len(tableNames)).Msg("List of tables")
	printTables(tableNames)
	for _, tableName := range tableNames {
		err = storage.ReadTable(tableName)
		if err != nil {
			log.Err(err).Msg("Read table content failed")
			return ExitStatusStorageError, err
		}
	}

	// we have finished, let's close the connection to database
	err = storage.Close()
	if err != nil {
		log.Err(err).Msg(operationFailedMessage)
		return ExitStatusStorageError, err
	}

	// default exit value + no error
	return ExitStatusOK, nil
}

func printTables(tableNames []TableName) {
	for i, tableName := range tableNames {
		log.Info().Int("#", i+1).Str("table", string(tableName)).Msg("Table in database")
	}
}

// doSelectedOperation function perform operation selected on command line.
// When no operation is specified, the Notification writer service is started
// instead.
func doSelectedOperation(configuration ConfigStruct, cliFlags CliFlags) (int, error) {
	switch {
	case cliFlags.ShowVersion:
		showVersion()
		return ExitStatusOK, nil
	case cliFlags.ShowAuthors:
		showAuthors()
		return ExitStatusOK, nil
	case cliFlags.ShowConfiguration:
		showConfiguration(configuration)
		return ExitStatusOK, nil
	default:
		// default operation - data export
		return performDataExport(configuration, cliFlags)
	}
	// this can not happen: return ExitStatusOK, nil
}

func main() {
	// command line flags
	var cliFlags CliFlags

	// define and parse all command line options
	flag.BoolVar(&cliFlags.ShowVersion, "version", false, "show version")
	flag.BoolVar(&cliFlags.ShowAuthors, "authors", false, "show authors")
	flag.BoolVar(&cliFlags.ShowConfiguration, "show-configuration", false, "show configuration")
	flag.BoolVar(&cliFlags.PrintSummaryTable, "summary", false, "print summary table after export")
	flag.StringVar(&cliFlags.Output, "output", "", "output to: CSV, S3")

	// parse all command line flags
	flag.Parse()

	// config has exactly the same structure as *.toml file
	config, err := LoadConfiguration(configFileEnvVariableName, defaultConfigFileName)
	if err != nil {
		log.Err(err).Msg("Load configuration")
	}

	if config.Logging.Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// perform selected operation
	exitStatus, err := doSelectedOperation(config, cliFlags)
	if err != nil {
		log.Err(err).Msg("Do selected operation")
		os.Exit(exitStatus)
		return
	}

	log.Debug().Msg("Started")

	log.Debug().Msg("Finished")
}
