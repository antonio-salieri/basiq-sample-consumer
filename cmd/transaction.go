package cmd

import (
	"log"

	"github.com/antonio-salieri/basiq-sample-consumer/entity"
	"github.com/spf13/cobra"
)

var transactionCmds = &cobra.Command{
	Use:   "transaction [command]",
	Short: "Manipulates with user entity",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	averageTransaction := &cobra.Command{
		Use:   "average-transaction",
		Short: "Shows average transaction amount per transaction type",
		Run: func(cmd *cobra.Command, args []string) {

			aggregatedTransactions, err := transactionService.AggregateTransactionPerDebitCategory(userID, entity.ConnectionData{
				InstitutionID: institutionID,
				LoginID:       institutionLoginID,
				LoginPassword: institutionLoginPassword,
			})
			if err != nil {
				log.Fatal(err)
			}
			if transactionType == "" {
				aggregatedTransactions.GetAverageAmounts(nil).Print()
			} else {
				aggregatedTransactions.GetAverageAmounts(&transactionType).Print()
			}

		},
	}

	averageTransaction.Flags().StringVarP(&userID, "userID", "u", defaultUserID, "User ID")
	averageTransaction.Flags().StringVarP(&institutionID, "institutionID", "i", defaultinstitutionID, "Institution id")
	averageTransaction.Flags().StringVarP(&institutionLoginID, "loginID", "l", defaultUserBankLoginID, "Institution login id")
	averageTransaction.Flags().StringVarP(&institutionLoginPassword, "institutionPassword", "p", defaultUserBankPassword, "Institution login password")
	averageTransaction.Flags().StringVarP(&transactionType, "transactionType", "t", "", "Transaction type to select (if not specified all debit transactions will be used)")
	averageTransaction.MarkFlagRequired("userID")

	transactionCmds.AddCommand(averageTransaction)
}
