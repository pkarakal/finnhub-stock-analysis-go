package cmd

import (
	"github.com/spf13/cobra"
)

type CLIOptions struct {
	Verbose bool
	Token   string
	Stocks  []string
}

var (
	rootCmd = &cobra.Command{
		Use:   "finnhub-ws",
		Short: "Finnhub stock analysis",
	}
	CLI CLIOptions
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVarP(&CLI.Verbose, "verbose", "v", false, "Run finnhub stock analysis verbosely")
	rootCmd.Flags().StringVarP(&CLI.Token, "token", "t", "", "Token to access finnhub ws api")
	rootCmd.Flags().StringArrayVarP(&CLI.Stocks, "stocks", "s", []string{}, "Stocks to run analysis on")
	rootCmd.MarkFlagRequired("token")
	rootCmd.MarkFlagRequired("stocks")
}
