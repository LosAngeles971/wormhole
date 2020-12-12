package cmd

import (
	"it/losangeles971/wormhole/internal"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open a wormhole between source and target",
	Long: `Open a wormhole between source and target
Usage:
	wormhole open --source <bind address>:<port> --target <address>:<port>
`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.Open()
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
	rootCmd.PersistentFlags().StringVarP(&internal.Source, "source", "s", "", "Source in form of <bind address>:<port>")
	rootCmd.PersistentFlags().StringVarP(&internal.Target, "target", "t", "", "Target in form of <address>:<port>ss")
	rootCmd.PersistentFlags().IntVarP(&internal.MaxConnections, "max", "m", 20, "Number of max connections (default is 20)")
}
