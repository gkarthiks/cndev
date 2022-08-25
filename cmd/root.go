// Package cmd
/*
Copyright © 2022 Karthikeyan Govindaraj
*/
package cmd

import (
	"fmt"
	"github.com/gkarthiks/cndev/utils"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cndev",
	Short: "A cli to mimic live developer experience",
	Long: fmt.Sprint(utils.Magenta(`
                  __         
  _________  ____/ /__ _   __
 / ___/ __ \/ __  / _ \ | / /
/ /__/ / / / /_/ /  __/ |/ / 
\___/_/ /_/\__,_/\___/|___/  `)) + `A cli tool to create the simulated
environment of the DEV environment in the private space with the 
private (local) Kubernetes Cluster, Git Server and ArgoCD. A valid
Kubernetes Cluster is needed and must be configured correctly in
the local kubeconfig file. Head over to project's README document
to see how to use.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

//func init() {
//	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
//}
