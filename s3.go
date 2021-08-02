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
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// NewS3Connection function initializes connection to S3/Minio storage.
func NewS3Connection(configuration ConfigStruct) (*minio.Client, context.Context, error) {
	s3Configuration := GetS3Configuration(configuration)

	endpoint := fmt.Sprintf("%s:%d",
		s3Configuration.EndpointURL, s3Configuration.EndpointPort)

	log.Info().Str("S3 endpoint", endpoint).Msg("Preparing connection")

	ctx := context.Background()

	// Initialize minio client object
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			s3Configuration.AccessKeyID,
			s3Configuration.SecretAccessKey, ""),
		Secure: s3Configuration.UseSSL,
	})
	if err != nil {
		log.Error().Err(err).Msg("Unable to initialize connection to S3")
		return nil, nil, err
	}

	log.Info().Msg("Connection established")
	return minioClient, ctx, nil
}

// s3BucketExists checks if bucket with given name exists and can be retrieved
func s3BucketExists(minioClient *minio.Client, ctx context.Context, bucketName string) (bool, error) {
	found, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Error().Err(err).Str("bucket", bucketName).Msg("Bucket can not be found")
		return false, err
	}

	return found, nil
}
