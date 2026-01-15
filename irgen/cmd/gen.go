package cmd

import (
	generator "github.com/nagarajRPoojari/picasso/irgen/codegen"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen [source_file] [output_path]",
	Short: "Compiles Picasso source code into Intermediate Representation",
	Long: `gen generates IR files for given project directory
Example:
    picasso gen projectDir`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// args[0] is project directory
		c := generator.NewGenerator(args[0])
		c.BuildAll()
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
