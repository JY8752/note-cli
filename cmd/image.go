package cmd

import (
	"github.com/JY8752/note-cli/internal/run"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image [markdown file path]",
	Short: "Create title image.",
	Long: heredoc.Doc(`
		Generate a title image using the title and author name provided in the article file.
		Execute the command in the directory where the article file exists.
		If the article file does not exist, image generation will fail.
	`),
	Example: heredoc.Doc(`
		note-cli create image ./article.md --template 1 -i ./icon.png -o ./output.png
	`),
	Args: cobra.ExactArgs(1),
	RunE: run.CreateImageFunc(&templateNo, &iconPath, &outputPath),
}

var (
	templateNo int16
	iconPath   string
	outputPath string
)

func init() {
	createCmd.AddCommand(imageCmd)

	imageCmd.Flags().Int16Var(&templateNo, "template", 1, "Template files can be specified by number.catalog -> https://github.com/JY8752/note-cli/blob/main/docs/templates/templates.md")
	imageCmd.Flags().StringVarP(&iconPath, "icon", "i", "", "Icons can be included in the generated image by specifying the path where the icon is located")
	imageCmd.Flags().StringVarP(&outputPath, "output", "o", "", "You can specify the path to output the generated images")
}
