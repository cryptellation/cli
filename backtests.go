package main

import (
	"fmt"

	"github.com/cryptellation/backtests/api"
	"github.com/cryptellation/go-clients/client"
	"github.com/spf13/cobra"
)

var backtestsCmd = &cobra.Command{
	Use:     "backtests",
	Aliases: []string{"b"},
	Short:   "Manage backtests",
}

var backtestsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List backtests",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Set client
		cl, err := client.New(client.WithTemporalAddress(temporalAddress))
		if err != nil {
			return err
		}
		defer cl.Close()

		// Execute call
		res, err := cl.ListBacktests(cmd.Context(), api.ListBacktestsWorkflowParams{})
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(res)
		default:
			fmt.Println("ID")
			for _, b := range res {
				fmt.Println(b.ID)
			}
		}

		return nil
	},
}

func setBacktestsCommands(cmd *cobra.Command) {
	backtestsCmd.AddCommand(backtestsListCmd)
	cmd.AddCommand(backtestsCmd)
}
