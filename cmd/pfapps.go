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
	"github.com/gkarthiks/cndev/k8s"
	"github.com/gkarthiks/cndev/utils"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// pfappsCmd represents the pfapps command
var pfappsCmd = &cobra.Command{
	Use:     "pfapps",
	Aliases: []string{"port-forward-apps", "pf", "pfa"},
	Short:   "Using this command port-forward the apps to your local machine",
	Long: `This command is used to port-forward the applications that are 
running in the cluster which are deployed via cndev. For example:

cndev pfapps --list
cndev pfapps --forward <app-name> --port <local-port>
`,
	Run: func(cmd *cobra.Command, args []string) {

		if cmd.Flags().NFlag() < 1 {
			cmd.Help()
			os.Exit(0)
		}

		validateMutuallyExclusiveFlags(cmd)

		list, _ := cmd.Flags().GetBool("list")
		if list == true {
			var terminalTableData [][]string
			header := []string{"SNO", "NAME", "DEFAULT_PORT"}
			terminalTableData = append(terminalTableData, header)
			for idx, app := range utils.PFApps {
				sno := strconv.Itoa(idx + 1)
				appRecord := []string{sno, app.Name, app.DefaultPort}
				terminalTableData = append(terminalTableData, appRecord)
			}
			pterm.DefaultTable.WithHasHeader(true).WithData(terminalTableData).Render()
		}

		forwardApp, _ := cmd.Flags().GetString("forward")
		forwardPort, _ := cmd.Flags().GetString("port")
		isAvailable, appSpec := utils.GetAppSpec(forwardApp)
		if len(forwardApp) > 0 && isAvailable {
			k8s.Forward(appSpec, forwardPort)
		}
	},
}

func init() {
	rootCmd.AddCommand(pfappsCmd)

	pfappsCmd.Flags().BoolP("list", "l", false, "lists all pods capable of port-forwarding")
	pfappsCmd.Flags().StringP("forward", "f", "", "port-forwards the given application")
	pfappsCmd.Flags().StringP("port", "p", "", "local port number for given application")
}

func validateMutuallyExclusiveFlags(cmd *cobra.Command) {
	list, _ := cmd.Flags().GetBool("list")
	forwardApp, _ := cmd.Flags().GetString("forward")
	forwardPort, _ := cmd.Flags().GetString("port")

	if list && (len(forwardApp) > 0 || len(forwardPort) > 0) {
		logrus.Errorf("flag --list cannot be used with the --forward and/or --port")
		os.Exit(0)
	}
}
