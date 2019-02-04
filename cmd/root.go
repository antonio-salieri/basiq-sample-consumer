package cmd

import (
	"github.com/antonio-salieri/basiq-sample-consumer/client"
	"github.com/antonio-salieri/basiq-sample-consumer/service"
	"github.com/spf13/cobra"
)

var apiClient client.Client
var transactionService *service.TransactionService

var rootCmd = &cobra.Command{
	Use:   "basiq-sample-consumer [command] [subcommand]",
	Short: "Basiq provides a collection of APIs to help you build powerful financial solutions for a wide range of use cases.",
	Long: ` The most common use cases are:

	* Personal Financial Management. Enable your customers to aggregate all of their financial data in one place, identify expenses and gain valuable insight of their spending.
	* Wealth Management. Gain valuable insights and a clearer understanding of your customers’ financial positions to customize advice, recommendations, and product offerings.
	* Risk Insights. Gain real-time and comprehensive visibility of your customers' assets, income, non-credit payment patterns, and transactional details.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(userCmds)
	rootCmd.AddCommand(transactionCmds)
}

// Execute runs cmd application
func Execute(client client.Client, trService *service.TransactionService) error {

	apiClient = client
	transactionService = trService
	return rootCmd.Execute()
}
