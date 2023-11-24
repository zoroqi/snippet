package cmd

import "github.com/spf13/cobra"

var moveCmd = &cobra.Command{
	Use:     "move",
	Aliases: []string{"mv"},
	Short:   "Move snippet to another directory",
	Long:    `Move snippet to another directory`,
	RunE:    move,
}

var mvFlag struct {
	alias  string
	target string
}

func init() {
	moveCmd.Flags().StringVarP(&mvFlag.target, "target", "t", "", "move target")
	RootCmd.AddCommand(moveCmd)
}

func move(cmd *cobra.Command, args []string) error {
	return nil
}
