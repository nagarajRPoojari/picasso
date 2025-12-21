package cmd

import (
	generator "github.com/nagarajRPoojari/niyama/irgen/codegen"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen [source_file] [output_path]",
	Short: "Compiles Niyama source code into Intermediate Representation",
	Long: `The gen command initiates the full compilation pipeline:
1. It initializes a NewGenerator with the provided source file.
2. It executes BuildAll, which runs the Lexer, Parser, and Semantic Analyzer.
3. It emits the final IR or target code to the specified output path.

Example:
    niyama gen ./main.niy ./build/main.ll`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// args[0] is the input source file path.
		c := generator.NewGenerator(args[0])
		c.BuildAll()
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
