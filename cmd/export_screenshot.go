package cmd

import (
	pkg "github.com/Twacqwq/youth/pkg/youth"
	"github.com/spf13/cobra"
)

var (
	fileName string
	filePath string
)

var ScreenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "导出最新一期完成截图",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.CompleteJPG(fileName, filePath)
	},
}

func init() {
	ScreenshotCmd.Flags().StringVarP(&fileName, "name", "n", "screenshot", "指定截图名称")
	ScreenshotCmd.Flags().StringVarP(&filePath, "path", "p", ".", "指定截图导出路径")
}
