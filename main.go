/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/mas2020-golang/ion/cmd"
	"github.com/mas2020-golang/ion/packages/utils"
)

var (
	GitCommit, BuildDate string = "N/A", ""
)

func main() {
	utils.GitCommit = GitCommit
	utils.BuildDate = BuildDate
	cmd.Execute()
}
