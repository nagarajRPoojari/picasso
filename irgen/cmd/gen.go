/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	generator "github.com/nagarajRPoojari/niyama/irgen/codegen"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate IR",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c := generator.NewGenerator()
		c.Build(args[0])
		c.Compile()
		c.Dump(args[1])
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
