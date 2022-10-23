package cmd

import (
	"github.com/spf13/cobra"
	"github.com/syrull/dnsc2-server/pkg"
)

var rootCmd = &cobra.Command{
	Use:   "dns-server",
	Short: "A C2 DNS Server",
	Long:  pkg.PrintHeader(),
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
