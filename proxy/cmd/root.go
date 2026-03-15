package main

import (
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/cmd/cmds"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "proxy",
	Short: "proxy is a proxy server",
	Long:  `proxy is a proxy server`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(cmds.ListenCMD)
}

func main() {
	Execute()
}
