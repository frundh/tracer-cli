package cmd

import (
	"github.com/spf13/cobra"
)

var traceCmd = &cobra.Command{
	Use:   "trace",
	Short: "Send sample traces using different protocols",
}

func init() {
	rootCmd.AddCommand(traceCmd)
}
