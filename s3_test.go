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
