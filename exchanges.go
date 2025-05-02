package main

import (
	"fmt"
	"strings"

	"github.com/cryptellation/exchanges/api"
	"github.com/cryptellation/go-client/client"
	"github.com/spf13/cobra"
)

var exchangesCmd = &cobra.Command{
	Use:     "exchanges",
	Aliases: []string{"e"},
	Short:   "Manage info on exchanges",
}

var exchangesShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"s"},
	Short:   "Show exchange",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set client
		cl, err := client.New(client.WithTemporalAddress(temporalAddress))
		if err != nil {
			return err
		}
		defer cl.Close()

		// Get exchanges
		res, err := cl.GetExchange(cmd.Context(), api.GetExchangeWorkflowParams{
			Name: args[0],
		})
		if err != nil {
			return err
		}

		// Display result
		switch {
		case jsonOutput:
			return displayJSON(res.Exchange)
		default:
			var pairs string
			if len(res.Exchange.Pairs) > 7 {
				res.Exchange.Pairs = res.Exchange.Pairs[:5]
				pairs = strings.Join(res.Exchange.Pairs, ",") + ",..."
			} else {
				pairs = strings.Join(res.Exchange.Pairs, ",")
			}

			fmt.Println("Name:", res.Exchange.Name)
			fmt.Println("Periods:", strings.Join(res.Exchange.Periods, ","))
			fmt.Println("Pairs:", pairs)
			fmt.Println("Fees:", res.Exchange.Fees)
		}

		return nil
	},
}

var exchangesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List exchanges",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Set client
		cl, err := client.New(client.WithTemporalAddress(temporalAddress))
		if err != nil {
			return err
		}
		defer cl.Close()

		// List exchanges
		res, err := cl.ListExchanges(cmd.Context(), api.ListExchangesWorkflowParams{})
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(res.List)
		default:
			fmt.Println("NAME")
			for i := range res.List {
				fmt.Println(res.List[i])
			}
		}

		return nil
	},
}

func setExchangesCommands(cmd *cobra.Command) {
	exchangesCmd.AddCommand(exchangesShowCmd)
	exchangesCmd.AddCommand(exchangesListCmd)
	cmd.AddCommand(exchangesCmd)
}
