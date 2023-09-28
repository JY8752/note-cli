package cmd

import (
	"github.com/JY8752/note-cli/internal/run"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Create title image.",
	Long:  `Create title image`,
	Args:  cobra.NoArgs,
	RunE:  run.CreateImageFunc(&templateNo, &iconPath, &outputPath),
}

var (
	templateNo int16
	iconPath   string
	outputPath string
)

func init() {
	createCmd.AddCommand(imageCmd)

	imageCmd.Flags().Int16Var(&templateNo, "template", 1, "Template files can be specified by number")
	imageCmd.Flags().StringVarP(&iconPath, "icon", "i", "", "Icons can be included in the generated image by specifying the path where the icon is located")
	imageCmd.Flags().StringVarP(&outputPath, "output", "o", "", "You can specify the path to output the generated images")
}
