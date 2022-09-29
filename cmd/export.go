package cmd

import "github.com/spf13/cobra"

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出相关数据",
}

func init() {
	exportCmd.AddCommand(ScreenshotCmd)
}
