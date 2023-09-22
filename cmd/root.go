package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "note-cli",
	Short: "note-cli is a CLI command tool for creating, writing, and managing note articles",
	Long:  `note-cli is a CLI command tool for creating, writing, and managing note articles`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.note-cli.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
