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

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls bucket time",
	Short: "list the content of a bucket at a point in time",
	Args:  cobra.ExactArgs(2),
	RunE:  lsRun,
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolP("long", "l", false, "Use a long listing format")
}

func lsRun(cmd *cobra.Command, args []string) error {
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

	objectsAtRestoreTime, err := lib.ListObjectAtRestoreTime(&ctx, bucket, restoreTime)
	if err != nil {
		log.WithError(err).Fatal("Listing bucket's objects failed")
	}

	if longStatus, err := cmd.Flags().GetBool("long"); err != nil {
		log.WithError(err).Fatal("Getting flag 'long' failed")
	} else if longStatus {
		for name, attrs := range objectsAtRestoreTime {
			fmt.Println(name, attrs.Deleted.IsZero(), attrs.Updated)
		}
	} else {
		for name := range objectsAtRestoreTime {
			fmt.Println(name)
		}
	}

	return nil
}
