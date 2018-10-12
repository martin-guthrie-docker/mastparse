package cmd

import (
	"github.com/martin-guthrie-docker/mastparse/pkg/mastparse"
	"github.com/spf13/cobra"

	"github.com/martin-guthrie-docker/mastparse/pkg/log"
)

func init() {
	rootCmd.AddCommand(cmdInspect)
}

var cmdInspect = &cobra.Command{
	Use:   "inspect <name>",
	Short: "parse deployment 'name' to the console",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inspectFunc(args)
	},
}

func inspectFunc(args []string) error {
	log.Term.Info("Start ", args)

	mp, err := mastparse.NewMparseClass(
		mastparse.MparseClassCfg{
			Log:  log.Term,
			Name: args[0],
			MastPath: GlobalConfigCfg.MastPath,
		})

	if err != nil {
		log.Term.Fatalf("mastparse.NewMparseClass failed")
		return err
	}

	err = mp.Open(nil)
	if err != nil {
		log.Term.Error("open failed")
		log.Term.Error("Please check your deployment name and/or path to mast datastore")
		return err
	}

	mp.ReadMastInventory()
	mp.PrintCLIs()

	mp.Close()
	log.Term.Debug("Done")
	return nil
}
