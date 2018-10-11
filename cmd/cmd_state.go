package cmd

import (
	"github.com/martin-guthrie-docker/mastparse/pkg/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdState)
}

var cmdState = &cobra.Command{
	Use:   "state",
	Short: "print state, configuration to the console/log",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		stateFunc(args)
	},
}

func stateFunc(args []string) error {
	log.Term.Debug("State:")
	// CmdConfig is from cmd_root_exec.go, global
	GlobalConfig.Dump(true)


	return nil
}