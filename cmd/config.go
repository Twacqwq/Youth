package cmd

import (
	"github.com/Twacqwq/youth/pkg/utils"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "generate",
	Short: "生成配置文件",
	Run: func(cmd *cobra.Command, args []string) {
		utils.GenerateConfig()
	},
}
