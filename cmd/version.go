package cmd

import (
	"fmt"
	"github.com/RGaius/octopus/pkg/common/version"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print version",
		Long:  "print version",
		Run: func(c *cobra.Command, args []string) {
			fmt.Printf("version: %v\n", version.Get())
		},
	}

	revisionCmd = &cobra.Command{
		Use:   "revision",
		Short: "print revision with building date",
		Long:  "print revision with building date",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("revision: %v\n", version.GetRevision())
		},
	}
)
