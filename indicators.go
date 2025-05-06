package main

import (
	"fmt"
	"time"

	"github.com/cryptellation/candlesticks/pkg/candlestick"
	"github.com/cryptellation/candlesticks/pkg/period"
	"github.com/cryptellation/go-clients/client"
	"github.com/cryptellation/sma/api"
	"github.com/spf13/cobra"
)

var indicatorsCmd = &cobra.Command{
	Use:     "indicators",
	Aliases: []string{"i"},
	Short:   "Manage indicators",
}

var (
	smaListExchangeFlag     string
	smaListPairFlag         string
	smaListPeriodFlag       string
	smaListStartFlag        string
	smaListEndFlag          string
	smaListPeriodNumberFlag int
	smaListPriceTypeFlag    string
)

var smaCmd = &cobra.Command{
	Use:     "sma",
	Aliases: []string{"s"},
	Short:   "Manage Simple Moving Average (SMA)",
}

var smaListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List simple moving averages points",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Set client
		cl, err := client.New(client.WithTemporalAddress(temporalAddress))
		if err != nil {
			return err
		}
		defer cl.Close()

		// Parse start date
		start, err := time.Parse(time.RFC3339, smaListStartFlag)
		if err != nil {
			return err
		}

		// Parse end date
		end, err := time.Parse(time.RFC3339, smaListEndFlag)
		if err != nil {
			return err
		}

		// Parse period
		per, err := period.FromString(smaListPeriodFlag)
		if err != nil {
			return err
		}

		// Parse price type
		priceType := candlestick.PriceType(smaListPriceTypeFlag)
		if err := priceType.Validate(); err != nil {
			return fmt.Errorf("%w: %s", err, smaListPriceTypeFlag)
		}

		// Execute call
		res, err := cl.ListSMA(cmd.Context(), api.ListWorkflowParams{
			Exchange:     smaListExchangeFlag,
			Pair:         smaListPairFlag,
			Period:       per,
			Start:        start,
			End:          end,
			PeriodNumber: smaListPeriodNumberFlag,
			PriceType:    priceType,
		})
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(res.Data.ToArray())
		default:
			fmt.Println(res.Data.String())
		}

		return nil
	},
}

func setIndicatorsCommands(cmd *cobra.Command) {
	smaListCmd.Flags().StringVarP(&smaListExchangeFlag, "exchange", "e", "binance", "Exchange")
	smaListCmd.Flags().StringVarP(&smaListPairFlag, "pair", "p", "ETH-USDT", "Pair")
	smaListCmd.Flags().StringVarP(&smaListPeriodFlag, "period", "P", "H1", "Period")
	smaListCmd.Flags().IntVarP(&smaListPeriodNumberFlag, "period-number", "N", 10, "Period number")
	smaListCmd.Flags().StringVarP(&smaListPriceTypeFlag, "price-type", "T", "close", "Price type")
	smaListCmd.Flags().StringVarP(
		&smaListStartFlag, "start", "s", time.Now().AddDate(0, 0, -8).Format(time.RFC3339), "Start")
	smaListCmd.Flags().StringVarP(
		&smaListEndFlag, "end", "E", time.Now().AddDate(0, 0, -1).Format(time.RFC3339), "End")
	smaCmd.AddCommand(smaListCmd)
	indicatorsCmd.AddCommand(smaCmd)
	cmd.AddCommand(indicatorsCmd)
}
