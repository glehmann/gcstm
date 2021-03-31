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
	"github.com/glehmann/gcstm/lib"
	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore bucket time",
	Short: "Restore a bucket at a specific time",
	Run:   restoreRun,
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}

func restoreRun(cmd *cobra.Command, args []string) {
	restoreTime, err := lib.ParseTime(args[1])
	if err != nil {
		log.WithError(err).Fatal("Can't parse time")
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.WithError(err).Fatal("Connection to google cloud apis failed")
	}
	defer client.Close()
	bucket := client.Bucket(args[0])

	planElements, err := lib.PlanRestore(&ctx, bucket, restoreTime)
	if err != nil {
		log.WithError(err).Fatal("Planning for restore failed")
	}

	if lib.ApplyPlan(&ctx, bucket, planElements) != nil {
		log.WithError(err).Fatal("Planning for restore failed")
	}
}
