package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init db.json",
	Long:  `Init db.json`,
	RunE:  initFunc,
}

var initFlag struct {
	target string
}

func init() {
	initCmd.Flags().StringVarP(&initFlag.target, "target", "t", "",
		"init db.json, and save it to directory")
	RootCmd.AddCommand(initCmd)
}

func initFunc(cmd *cobra.Command, args []string) error {
	if initFlag.target == "" {
		return errors.New("no target")
	}
	return db.CreateDbFile(initFlag.target)
}
