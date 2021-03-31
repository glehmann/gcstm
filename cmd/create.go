/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create projejct bucket",
	Short: "create a bucket with versioning enabled",
	Args:  cobra.ExactArgs(2),
	Run:   createRun,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().Int32P("retention", "r", 30, "Retention time in days")
}

func createRun(cmd *cobra.Command, args []string) {
	retention, err := cmd.Flags().GetInt32("retention")
	if err != nil {
		log.WithError(err).Fatal("Getting flag 'retention' failed")
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.WithError(err).Fatal("Connection to google cloud apis failed")
	}
	defer client.Close()
	bucket := client.Bucket(args[1])

	if err := bucket.Create(ctx, args[0], &storage.BucketAttrs{
		VersioningEnabled: true,
		Lifecycle: storage.Lifecycle{
			Rules: []storage.LifecycleRule{
				{
					Action: storage.LifecycleAction{
						Type: "Delete",
					},
					Condition: storage.LifecycleCondition{
						AgeInDays: int64(retention),
						Liveness:  storage.Archived,
					},
				},
			},
		},
	}); err != nil {
		log.WithError(err).Fatal("Bucket creation failed")
	}
}
