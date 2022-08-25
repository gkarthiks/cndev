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
	"github.com/gkarthiks/cndev/argocd"
	"github.com/gkarthiks/cndev/gogs"
	"github.com/gkarthiks/cndev/prompts"
	"github.com/gkarthiks/cndev/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:     "destroy",
	Aliases: []string{"delete", "clean"},
	Short:   "Cleans all the infrastructure installations",
	Long: `Cleans or deletes all the installation of argocd, 
gogs, grafana. Also deletes the repositories stored in the local git
server, argocd applications and configurations. Run destroy only if 
you must need a fresh start.`,
	Run: func(cmd *cobra.Command, args []string) {
		confirmDelete := prompts.PromptYesNo("Confirm deleting all the resources?")
		if confirmDelete == utils.StringYes {
			logrus.Info(utils.Red("Proceeding to delete the resources"))
			argocd.DeleteArgoCD()
			gogs.DeleteGogs()
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
