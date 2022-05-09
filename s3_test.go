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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/minio/minio-go/v7"

	main "github.com/RedHatInsights/insights-results-aggregator-exporter"
)

// mustConstructMinioClient helper function constructs an instance of Minio
// client or make the test fail
func mustConstructMinioClient(t *testing.T) *minio.Client {
	minioClient, err := minio.New("localhost:1234", &minio.Options{})
	assert.Nil(t, err)

	return minioClient
}

// Test case specification structure for function main.NewS3Connection
type newS3ConnectionTestSpecification struct {
	description   string
	configuration *main.ConfigStruct
	shouldFail    bool
	expectedError string
}

// TestNewS3Connection checks the function/constructor NewS3Connection
func TestNewS3Connection(t *testing.T) {
	// all test cases
	testCases := []newS3ConnectionTestSpecification{
		newS3ConnectionTestSpecification{
			description:   "nilConfiguration",
			configuration: nil,
			shouldFail:    true,
			expectedError: "Configuration is nil",
		},
		newS3ConnectionTestSpecification{
			description:   "emptyConfiguration",
			configuration: &main.ConfigStruct{},
			shouldFail:    true,
			expectedError: "Endpoint: :0 does not follow ip address or domain name standards.",
		},
		newS3ConnectionTestSpecification{
			description: "wrongConfiguration",
			configuration: &main.ConfigStruct{
				S3: main.S3Configuration{
					Type:            "",
					EndpointURL:     "",
					EndpointPort:    1234,
					AccessKeyID:     "",
					SecretAccessKey: "",
					UseSSL:          false,
					Bucket:          "",
				}},
			shouldFail:    true,
			expectedError: "Endpoint: :1234 does not follow ip address or domain name standards.",
		},
		newS3ConnectionTestSpecification{
			description: "correctConfiguration",
			configuration: &main.ConfigStruct{
				S3: main.S3Configuration{
					Type:            "minio",
					EndpointURL:     "localhost",
					EndpointPort:    1234,
					AccessKeyID:     "foobar",
					SecretAccessKey: "foobar",
					UseSSL:          false,
					Bucket:          "test",
				}},
			shouldFail: false,
		},
	}

	// run all specified test cases
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			// try to construct Minio client using nil
			// configuration
			client, _, err := main.NewS3Connection(testCase.configuration)

			// check for error
			if testCase.shouldFail {
				// client should not be constructed and error
				// should be returned
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.expectedError)
				assert.Nil(t, client)
			} else {
				// client should be constructed and error
				// should not be returned
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})

	}
}
