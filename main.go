/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/mas2020-golang/ion/cmd"
	"github.com/mas2020-golang/ion/packages/utils"
)

var (
	GitCommit, BuildDate, GoVersion string = "N/A", "", "go1.21.1"
)

func main() {
	utils.GitCommit = GitCommit
	utils.BuildDate = BuildDate
	utils.GoVersion = GoVersion
	cmd.Execute()
}