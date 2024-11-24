package cmd

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:          "octopus",
		Short:        "octopus",
		Long:         "octopus",
		SilenceUsage: true,
	}
)

// init 初始化命令行工具
func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(revisionCmd)
}

// Execute 执行命令行解析
func Execute() {
	_ = rootCmd.Execute()
}
