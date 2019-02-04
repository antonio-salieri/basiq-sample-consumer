package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var userCmds = &cobra.Command{
	Use:   "user [command]",
	Short: "Manipulates with user entity",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {

	userGetCommand := &cobra.Command{
		Use:   "get",
		Short: "Fetch existing user",
		Run: func(cmd *cobra.Command, args []string) {
			user, err := apiClient.GetUser(userID)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Found user: %+v", user)
		},
	}
	userGetCommand.Flags().StringVarP(&userID, "userID", "u", "", "User ID")
	userGetCommand.MarkFlagRequired("userID")

	userCreateCommand := &cobra.Command{
		Use:   "create",
		Short: "Create new user",
		Run: func(cmd *cobra.Command, args []string) {
			user, err := apiClient.CreateUser(userEmail, userMobile)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Created user: %+v", user)
		},
	}
	userCreateCommand.Flags().StringVarP(&userEmail, "email", "e", "", "User email address")
	userCreateCommand.Flags().StringVarP(&userMobile, "mobile", "m", "", "User mobile phone")

	userDeleteCommand := &cobra.Command{
		Use:   "delete",
		Short: "Delete existing user",
		Run: func(cmd *cobra.Command, args []string) {
			err := apiClient.DeleteUser(userID)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("User deleted")
		},
	}
	userDeleteCommand.Flags().StringVarP(&userID, "userID", "u", "", "User ID")
	userDeleteCommand.MarkFlagRequired("userID")

	userCmds.AddCommand(userCreateCommand)
	userCmds.AddCommand(userDeleteCommand)
	userCmds.AddCommand(userGetCommand)
}
