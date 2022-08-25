// Package cmd
/* Copyright © 2022 Karthikeyan Govindaraj

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
	"github.com/gkarthiks/cndev/utils"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the ArgoCD, git-server in private env",
	Long:  `Initializes or deploys the ArgoCD, gogs git-server in the private environment`,
	Run: func(cmd *cobra.Command, args []string) {
		argocdNamespace, err := cmd.Flags().GetString("argocd-namespace")
		if err != nil || len(argocdNamespace) < 1 {
			argocdNamespace = utils.ArgoNamespace
		}
		argocd.DeployArgoCD(argocdNamespace)
		gogs.DeployGogs()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("argocd-namespace", "", "argocd", "ArgoCD Deployment Namespace; defaults to argocd")
}
