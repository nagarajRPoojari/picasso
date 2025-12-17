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
		c := generator.NewGenerator(args[0], args[1])
		c.BuildAll()
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
