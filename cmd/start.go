package cmd

import (
	"github.com/RGaius/octopus/pkg/bootstrap"
	"github.com/spf13/cobra"
)

var (
	configFilePath = ""

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "start running",
		Long:  "start running",
		Run: func(c *cobra.Command, args []string) {
			bootstrap.Start(configFilePath)
		},
	}
)

// init 解析命令参数
func init() {
	startCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "conf/octopus.yaml", "config file path")
}
