package main

import (
	"fmt"
	"time"

	"github.com/cryptellation/candlesticks/api"
	"github.com/cryptellation/candlesticks/pkg/period"
	"github.com/cryptellation/go-clients/client"
	"github.com/spf13/cobra"
)

var (
	candlesticksListExchangeFlag string
	candlesticksListPairFlag     string
	candlesticksListPeriodFlag   string
	candlesticksListStartFlag    string
	candlesticksListEndFlag      string
)

var candlesticksCmd = &cobra.Command{
	Use:     "candlesticks",
	Aliases: []string{"c"},
	Short:   "List candlesticks",
}

var candlesticksListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List candlesticks",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Set client
		cl, err := client.New(client.WithTemporalAddress(temporalAddress))
		if err != nil {
			return err
		}
		defer cl.Close()

		// Parse start date
		start, err := time.Parse(time.RFC3339, candlesticksListStartFlag)
		if err != nil {
			return err
		}

		// Parse end date
		end, err := time.Parse(time.RFC3339, candlesticksListEndFlag)
		if err != nil {
			return err
		}

		// Parse period
		per, err := period.FromString(candlesticksListPeriodFlag)
		if err != nil {
			return err
		}

		// Execute call
		res, err := cl.ListCandlesticks(cmd.Context(), api.ListCandlesticksWorkflowParams{
			Exchange: candlesticksListExchangeFlag,
			Pair:     candlesticksListPairFlag,
			Period:   per,
			Start:    &start,
			End:      &end,
		})
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(res.List)
		default:
			for _, cs := range res.List {
				fmt.Println(cs.String())
			}
		}

		return nil
	},
}

func setCandlesticksCommands(cmd *cobra.Command) {
	candlesticksListCmd.Flags().StringVarP(&candlesticksListExchangeFlag, "exchange", "e", "binance", "Exchange")
	candlesticksListCmd.Flags().StringVarP(&candlesticksListPairFlag, "pair", "p", "ETH-USDT", "Pair")
	candlesticksListCmd.Flags().StringVarP(&candlesticksListPeriodFlag, "period", "P", "H1", "Period")
	candlesticksListCmd.Flags().StringVarP(
		&candlesticksListStartFlag, "start", "s", time.Now().AddDate(0, 0, -8).Format(time.RFC3339), "Start")
	candlesticksListCmd.Flags().StringVarP(
		&candlesticksListEndFlag, "end", "E", time.Now().AddDate(0, 0, -1).Format(time.RFC3339), "End")
	candlesticksCmd.AddCommand(candlesticksListCmd)
	cmd.AddCommand(candlesticksCmd)
}
