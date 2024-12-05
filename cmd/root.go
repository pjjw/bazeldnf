package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	debug bool
	json  bool
	quiet bool
)

var rootCmd = &cobra.Command{
	Use:          "bazeldnf",
	Short:        "bazeldnf is a tool which can query RPM repos and determine package dependencies",
	Long:         `The tool allows resolving package dependencies mainly for the purpose to create custom-built SCRATCH containers consisting of RPMs, trimmed down to the absolute necessary`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		InitLogger(cmd)
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress output")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "log debug output")
	rootCmd.PersistentFlags().BoolVarP(&json, "json", "", false, "format log output as json")
	rootCmd.AddCommand(NewXATTRCmd())
	rootCmd.AddCommand(NewSandboxCmd())
	rootCmd.AddCommand(NewFetchCmd())
	rootCmd.AddCommand(NewInitCmd())
	rootCmd.AddCommand(NewRpmTreeCmd())
	rootCmd.AddCommand(NewResolveCmd())
	rootCmd.AddCommand(NewReduceCmd())
	rootCmd.AddCommand(NewRpm2TarCmd())
	rootCmd.AddCommand(NewPruneCmd())
	rootCmd.AddCommand(NewTar2FilesCmd())
	rootCmd.AddCommand(NewLddCmd())
	rootCmd.AddCommand(NewVerifyCmd())
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
