// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/itn/pkg/itn"
	"github.com/spf13/cobra"
)

// TODOs(bwagner5):
//   1. List view of valid instances to interrupt
//   2. Option to pass tags instead of instance IDs
//   3. Option to pass an OD instance and have this tool create a matching instance that is spot to test an interruption
//   4. Automated chaos - give this tool a tag or vpc and allow it to randomly interrupt spot instances at will

var version string

type Options struct {
	instanceIDs []string
	delay       time.Duration
	clean       bool
	version     bool
	region      string
	profile     string
}

func main() {
	options := Options{}
	rootCmd := &cobra.Command{
		Use:   "ec2-spot-interrupter",
		Short: "ec2-spot-interrupter is a simple CLI tool that triggers Amazon EC2 Spot Instance Interruption Notifications and Rebalance Recommendations.",
		Run: func(cmd *cobra.Command, _ []string) {
			if options.version {
				fmt.Println(version)
				os.Exit(0)
			}
			ctx := context.Background()
			cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(options.region), config.WithSharedConfigProfile(options.profile))
			if err != nil {
				fmt.Printf("❌ %s\n", err)
				os.Exit(1)
			}
			if err := itn.New(cfg).Interrupt(context.Background(), options.instanceIDs, options.delay, options.clean); err != nil {
				fmt.Printf("❌ %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Successfully sent spot rebalance recommendation and instance interruption to %v\n", options.instanceIDs)
		},
	}
	rootCmd.PersistentFlags().StringSliceVarP(&options.instanceIDs, "instance-ids", "i", []string{}, "instance IDs to interrupt")
	rootCmd.PersistentFlags().BoolVarP(&options.clean, "clean", "c", true, "clean up the underlying simulations")
	rootCmd.PersistentFlags().DurationVarP(&options.delay, "delay", "d", time.Second*15, "duration until the interruption notification is sent")
	rootCmd.PersistentFlags().BoolVarP(&options.version, "version", "v", false, "the version")
	rootCmd.PersistentFlags().StringVarP(&options.region, "region", "r", "", "the AWS Region")
	rootCmd.PersistentFlags().StringVarP(&options.profile, "profile", "p", "", "the AWS Profile")
	rootCmd.Execute()
}
