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
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/apex/log"
	"github.com/glehmann/gcstm/lib"
	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan bucket time",
	Short: "Show what a restore would do at a specific time",
	Run:   planRun,
}

func init() {
	rootCmd.AddCommand(planCmd)
}

func planRun(cmd *cobra.Command, args []string) {
	planTime, err := lib.ParseTime(args[1])
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

	planElements, err := lib.PlanRestore(&ctx, bucket, planTime)
	if err != nil {
		log.WithError(err).Fatal("Planning for restore failed")
	}

	for name, planElement := range planElements {
		fmt.Println(name, planElement.Action)
	}
}
