package main

import (
	"fmt"

	"github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/go-clients/client"
	"github.com/spf13/cobra"
)

var forwardtestsCmd = &cobra.Command{
	Use:     "forwardtests",
	Aliases: []string{"f"},
	Short:   "Manage forwardtests",
}

var forwardtestsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List forwardtests",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Set client
		cl, err := client.New(client.WithTemporalAddress(temporalAddress))
		if err != nil {
			return err
		}
		defer cl.Close()

		// Execute call
		res, err := cl.ListForwardtests(cmd.Context(), api.ListForwardtestsWorkflowParams{})
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(res)
		default:
			fmt.Println("ID")
			for _, f := range res {
				fmt.Println(f.ID)
			}
		}

		return nil
	},
}

func setForwardtestsCommands(cmd *cobra.Command) {
	forwardtestsCmd.AddCommand(forwardtestsListCmd)
	cmd.AddCommand(forwardtestsCmd)
}
