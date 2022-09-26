package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "youth",
	Short: "提交最新一期青年大学习",
	Run: func(cmd *cobra.Command, args []string) {
		Youth()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
