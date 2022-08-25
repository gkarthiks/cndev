// Package cmd
/*
Copyright © 2022 Karthikeyan Govindaraj

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
	"fmt"
	"github.com/gkarthiks/cndev/k8s"
	"github.com/gkarthiks/cndev/utils"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of cndev",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cndev version:    %s\n", utils.AppVersion)
		fmt.Printf("Kubernetes version: %s \n", k8s.K8sVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
