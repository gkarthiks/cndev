/*
Copyright © 2022 Karthikeyan Govindaraj
*/
package main

import (
	"github.com/gkarthiks/cndev/cmd"
	"github.com/gkarthiks/cndev/utils"
)

func main() {
	cmd.Execute()
}

var BuildVersion = "development"

func init() {
	utils.AppVersion = BuildVersion
}
