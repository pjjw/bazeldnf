package main

import (
	"github.com/rmohr/bazeldnf/pkg/logger"
	"github.com/spf13/cobra"
)

func InitLogger(cmd *cobra.Command) {
	if json, _ := cmd.Flags().GetBool("json"); json {
		logger.SetJSON()
	}
	if quiet, _ := cmd.Flags().GetBool("quiet"); quiet {
		logger.SetQuiet()
	}
	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		logger.SetDebug()
	}
}

func main() {
	Execute()
}
