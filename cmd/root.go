/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strconv"

	"github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/ion/cmd/file"
	"github.com/mas2020-golang/ion/cmd/security"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ion",
	Short: "Ion is your swiss knife for having with you a lot of useful commands",
	Long: `Ion is an all-in-one application to sum up a lot of useful tools in a single command.
The swiss knife for every SysAdmin/DevOps!. You can use the ion commands as you do with pipes,
standard input/output and a lot of other daily basis activities.

Some examples:
// to search some content into a file
$ ion search --no-colors "this" demo-file

// tail the last 10 rows
$ ion tail --rows 10 test.txt
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	setLogs()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ion.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Add the other commands
	rootCmd.AddCommand(file.NewTailCmd())
	rootCmd.AddCommand(file.NewTreeCmd())
	rootCmd.AddCommand(security.NewDecryptCmd())
	rootCmd.AddCommand(security.NewEncryptCmd())
	rootCmd.AddCommand(file.NewCountCmd())
	rootCmd.AddCommand(file.NewRmCmd())
	rootCmd.AddCommand(file.NewSearchCmd())
	rootCmd.AddCommand(file.NewSliceCmd())
}

// setLogs load the configuration for the logging system
func setLogs() {
	// PanicLevel: 0, FatalLevel: 1, ErrorLevel: 2, WarnLevel: 3, InfoLevel: 4, DebugLevel: 5, TraceLevel: 6
	logrus.SetLevel(0)
	// the log level is first taken from the env variable APP_LOGLEVEL. In case it doesn't exist it is loaded
	// the value in the utils.Config.Logging.Level variable.
	if len(os.Getenv("ION_LOGLEVEL")) > 0 {
		l, err := strconv.Atoi(os.Getenv("ION_LOGLEVEL"))
		output.CheckErrorAndExitLog("", "", err)
		logrus.SetLevel(logrus.Level(l))
	}
	// choose to colorize the log output
	if len(os.Getenv("ION_LOGLEVEL")) > 0 {
		if os.Getenv("ION_LOGCOLOR") == "true" {
			logrus.SetFormatter(&output.TextColorFormatter{})
		} else {
			logrus.SetFormatter(&output.TextFormatter{})
		}
	}
	logrus.SetOutput(os.Stdout)
}
