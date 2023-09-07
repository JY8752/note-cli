package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new article directory.",
	Long: `
Create a new article directory.
If you want to specify directory and file names, specify them as --name(-n) options.
You can also specify the -t(--time) option to make the current timestamp the file name.
If nothing is specified, the file is created with a unique file name by UUID.
	`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
