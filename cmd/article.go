package cmd

import (
	"github.com/JY8752/note-cli/internal/run"
	"github.com/spf13/cobra"
)

var articleCmd = &cobra.Command{
	Use:   "article",
	Short: "Create a new article directory.",
	Long: `
Create a new article directory.
If you want to specify directory and file names, specify them as --name(-n) options.
You can also specify the -t(--time) option to make the current timestamp the file name.
If nothing is specified, the file is created with a unique file name by UUID.
	`,
	Args: cobra.NoArgs,
	RunE: run.CreateArticleFunc(&timeFlag, &noDirFlag, &name, &author),
	Example: `note-cli create article
note-cli create article --name article-a
note-cli create article -t`,
}

var (
	timeFlag  bool
	name      string
	author    string
	noDirFlag bool
)

func init() {
	articleCmd.Flags().BoolVarP(&timeFlag, "time", "t", false, "Create directory and file names with the current timestamp")
	articleCmd.Flags().StringVarP(&name, "name", "n", "", "Create a directory with the specified name")
	articleCmd.Flags().StringVarP(&author, "author", "a", "", "Author name")
	articleCmd.Flags().BoolVar(&noDirFlag, "no-dir", false, "Create an article file without directory.")

	articleCmd.MarkFlagsMutuallyExclusive("time", "name")

	createCmd.AddCommand(articleCmd)
}
