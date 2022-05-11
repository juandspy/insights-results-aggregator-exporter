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
	"bytes"
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
	listOfTablesMsg        = "List of tables"
	tableNameMsg           = "Table name"
)

// Exit codes
const (
	// ExitStatusOK means that the tool finished with success
	ExitStatusOK = iota

	// ExitStatusStorageError is returned in case of any consumer-related
	// error
	ExitStatusStorageError

	// ExitStatusS3Error is returned in case of any error related with
	// S3/Minio connection
	ExitStatusS3Error

	// ExitStatusConfigurationError is returned in case user provide wrong
	// configuration on command line or in configuration file
	ExitStatusConfigurationError

	// ExitStatusIOError is returned in case of any I/O error (export data
	// into file failed etc.)
	ExitStatusIOError
)

const (
	configFileEnvVariableName = "INSIGHTS_RESULTS_AGGREGATOR_EXPORTER_CONFIG_FILE"
	defaultConfigFileName     = "config"
)

// output files or objects containing metadata
const (
	listOfTables  = "_tables.csv"
	metadataTable = "_metadata.csv"
	disabledRules = "_disabled_rules.csv"
	logFile       = "_logs.txt"
)

// messages
const (
	readDisabledRulesInfoFailed      = "Read disabled rules info failed"
	storeDisabledRulesIntoFileFailed = "Store disabled rules into file failed"
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
func showConfiguration(config *ConfigStruct) {
	storageConfig := GetStorageConfiguration(config)
	log.Info().
		Str("Driver", storageConfig.Driver).
		Str("DB Name", storageConfig.PGDBName).
		Str("Username", storageConfig.PGUsername). // password is omitted on purpose
		Str("Host", storageConfig.PGHost).
		Int("DB Port", storageConfig.PGPort).
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
		Uint("S3 Port", s3Configuration.EndpointPort).
		Str("AccessKeyID", s3Configuration.AccessKeyID).
		Str("SecretAccessKey", s3Configuration.SecretAccessKey).
		Bool("Use SSL", s3Configuration.UseSSL).
		Str("Bucket name", s3Configuration.Bucket).
		Msg("S3 configuration")
}

// performDataExport function exports all data into selected output
func performDataExport(configuration *ConfigStruct, cliFlags CliFlags) (int, error) {
	// prepare the storage
	storageConfiguration := GetStorageConfiguration(configuration)
	storage, err := NewStorage(&storageConfiguration)
	if err != nil {
		log.Err(err).Msg(operationFailedMessage)
		return ExitStatusStorageError, err
	}

	switch cliFlags.Output {
	case "S3":
		return performDataExportToS3(configuration, storage, cliFlags.ExportMetadata, cliFlags.ExportDisabledRules)
	case "file":
		return performDataExportToFiles(configuration, storage, cliFlags.ExportMetadata, cliFlags.ExportDisabledRules)
	default:
		return ExitStatusConfigurationError, fmt.Errorf("Unknown output type: %s", cliFlags.Output)
	}
}

// performDataExportToS3 exports all tables and metadata info configured S3
// bucket
func performDataExportToS3(configuration *ConfigStruct,
	storage *DBStorage, exportMetadata bool,
	ExportDisabledRules bool) (int, error) {
	minioClient, context, err := NewS3Connection(configuration)
	if err != nil {
		return ExitStatusS3Error, err
	}

	tableNames, err := storage.ReadListOfTables()
	if err != nil {
		log.Err(err).Msg(operationFailedMessage)
		return ExitStatusStorageError, err
	}

	log.Info().Int("tables count", len(tableNames)).Msg(listOfTablesMsg)
	printTables(tableNames)

	bucket := GetS3Configuration(configuration).Bucket
	log.Info().Str("bucket name", bucket).Msg("S3 bucket to write to")

	if exportMetadata {
		// export list of all tables into S3
		err = storeTableNames(context, minioClient,
			bucket, listOfTables, tableNames)
		if err != nil {
			log.Err(err).Msg("Store table list to S3 failed")
			return ExitStatusStorageError, err
		}

		// export tables metadata into S3
		err = storage.StoreTableMetadataIntoS3(context, minioClient,
			bucket, metadataTable, tableNames)
		if err != nil {
			log.Err(err).Msg("Store tables metadata to S3 failed")
			return ExitStatusStorageError, err
		}
	}

	if ExportDisabledRules {
		// export rules disabled by more users into CSV file
		disabledRulesInfo, err := storage.ReadDisabledRules()
		if err != nil {
			log.Err(err).Msg(readDisabledRulesInfoFailed)
			return ExitStatusStorageError, err
		}

		// export list of disabled rules
		err = storeDisabledRulesIntoS3(context, minioClient, bucket,
			disabledRules, disabledRulesInfo)
		if err != nil {
			log.Err(err).Msg(storeDisabledRulesIntoFileFailed)
			return ExitStatusIOError, err
		}
	}

	// read content of all tables and perform export
	for _, tableName := range tableNames {
		err = storage.StoreTable(context, minioClient, bucket, tableName)
		if err != nil {
			log.Err(err).Str(tableNameMsg, string(tableName)).Msg("Store table into S3 failed")
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

// performDataExportToFiles exports all tables and metadata info files
func performDataExportToFiles(configuration *ConfigStruct,
	storage *DBStorage, exportMetadata bool,
	exportDisabledRules bool) (int, error) {
	tableNames, err := storage.ReadListOfTables()
	if err != nil {
		log.Err(err).Msg(operationFailedMessage)
		return ExitStatusStorageError, err
	}

	log.Info().Int("count", len(tableNames)).Msg(listOfTablesMsg)

	// log into terminal
	printTables(tableNames)

	if exportMetadata {
		// export list of all tables into CSV file
		err = storeTableNamesIntoFile(listOfTables, tableNames)
		if err != nil {
			log.Err(err).Msg("Store table list to file failed")
			return ExitStatusStorageError, err
		}

		// export tables metadata into CSV file
		err = storage.StoreTableMetadataIntoFile(metadataTable, tableNames)
		if err != nil {
			log.Err(err).Msg("Store tables metadata to file failed")
			return ExitStatusStorageError, err
		}
	}

	if exportDisabledRules {
		// export rules disabled by more users into CSV file
		disabledRulesInfo, err := storage.ReadDisabledRules()
		if err != nil {
			log.Err(err).Msg(readDisabledRulesInfoFailed)
			return ExitStatusStorageError, err
		}

		// export list of disabled rules
		err = storeDisabledRulesIntoFile(disabledRules, disabledRulesInfo)
		if err != nil {
			log.Err(err).Msg(storeDisabledRulesIntoFileFailed)
			return ExitStatusIOError, err
		}
	}

	// read content of all tables and perform export
	for _, tableName := range tableNames {
		err = storage.StoreTableIntoFile(tableName)
		if err != nil {
			log.Err(err).Str(tableNameMsg, string(tableName)).Msg("Store table into file failed")
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

// checkS3Connection checks if connection to S3 is possible
func checkS3Connection(configuration *ConfigStruct) (int, error) {
	log.Info().Msg("Checking connection to S3")
	minioClient, context, err := NewS3Connection(configuration)
	if err != nil {
		return ExitStatusS3Error, err
	}

	exists, err := s3BucketExists(context, minioClient, GetS3Configuration(configuration).Bucket)
	if err != nil {
		return ExitStatusS3Error, err
	}

	if !exists {
		log.Error().Msg("Can not find expected bucket")
	} else {
		log.Info().Msg("Bucket has been found")
	}

	log.Info().Msg("Connection to S3 seems to be ok")
	return ExitStatusOK, nil
}

// doSelectedOperation function perform operation selected on command line.
// When no operation is specified, the Notification writer service is started
// instead.
func doSelectedOperation(configuration *ConfigStruct, cliFlags CliFlags, operationLogger zerolog.Logger) (int, error) {
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
	case cliFlags.CheckS3Connection:
		return checkS3Connection(configuration)
	default:
		// default operation - data export
		return performDataExport(configuration, cliFlags)
	}
	// this can not happen: return ExitStatusOK, nil
}

func parseFlags() (cliFlags CliFlags) {
	// define and parse all command line options
	flag.BoolVar(&cliFlags.ShowVersion, "version", false, "show version")
	flag.BoolVar(&cliFlags.ShowAuthors, "authors", false, "show authors")
	flag.BoolVar(&cliFlags.ShowConfiguration, "show-configuration", false, "show configuration")
	flag.BoolVar(&cliFlags.PrintSummaryTable, "summary", false, "print summary table after export")
	flag.StringVar(&cliFlags.Output, "output", "S3", "output to: file, S3")
	flag.BoolVar(&cliFlags.ExportMetadata, "metadata", false, "export metadata")
	flag.BoolVar(&cliFlags.ExportDisabledRules, "disabled-by-more-users", false, "export rules disabled by more users")
	flag.BoolVar(&cliFlags.CheckS3Connection, "check-s3-connection", false, "check S3 connection and exit")
	flag.BoolVar(&cliFlags.ExportLog, "export-log", false, "export log")

	// parse all command line flags
	flag.Parse()

	return
}

type DummyWriter struct{}

func (w DummyWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

// createOperationLog function constructs operation log instance
func createOperationLog(cliFlags CliFlags) (zerolog.Logger, bytes.Buffer, error) {
	dummyLogger := zerolog.New(DummyWriter{}).With().Logger()
	var buffer bytes.Buffer

	if cliFlags.ExportLog {
		switch cliFlags.Output {
		case "S3":
			memoryLogger := zerolog.New(&buffer).With().Logger()
			memoryLogger.Info().Msg("Memory logger initialized")
			return memoryLogger, buffer, nil
		case "file":
			logFile, err := os.Create(logFile)
			if err != nil {
				return dummyLogger, buffer, err
			}
			fileLogger := zerolog.New(logFile).With().Logger()
			fileLogger.Info().Msg("File logger initialized")
			return fileLogger, buffer, nil
		default:
			return dummyLogger, buffer, fmt.Errorf("Unknown output type: %s", cliFlags.Output)
		}
	}

	return dummyLogger, buffer, nil

}

func main() {
	log.Debug().Msg("Started")

	// parse all command line flags
	cliFlags := parseFlags()

	// config has exactly the same structure as *.toml file
	config, err := LoadConfiguration(configFileEnvVariableName, defaultConfigFileName)
	if err != nil {
		log.Err(err).Msg("Load configuration")
	}

	if config.Logging.Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	operationLogger, buffer, err := createOperationLog(cliFlags)
	if err != nil {
		log.Err(err).Msg("Create operation log")
		os.Exit(ExitStatusIOError)
		return
	}

	// perform selected operation
	exitStatus, err := doSelectedOperation(&config, cliFlags, operationLogger)
	if err != nil {
		log.Err(err).Msg("Do selected operation")
		os.Exit(exitStatus)
		return
	}

	log.Debug().Msg("Finished")
}
