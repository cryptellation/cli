package main

import (
	"fmt"

	"github.com/cryptellation/go-clients/client"
	"github.com/cryptellation/ticks/api"
	"github.com/cryptellation/ticks/pkg/clients"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var (
	ticksListenExchangeFlag string
	ticksListenPairFlag     string
)

var ticksCmd = &cobra.Command{
	Use:     "ticks",
	Aliases: []string{"t"},
	Short:   "Listen to ticks",
}

var ticksListenCmd = &cobra.Command{
	Use:     "listen",
	Aliases: []string{"l"},
	Short:   "Listen to ticks",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Set client
		cl, err := client.New(client.WithTemporalAddress(temporalAddress))
		if err != nil {
			return err
		}
		defer cl.Close()

		// Create a worker
		id := uuid.New()
		tq := fmt.Sprintf("CryptellationCLI-%s-TaskQueue", id.String())
		w := worker.New(cl.GetTemporalClient(), tq, worker.Options{})
		defer w.Stop()

		if err := cl.ListenToTicks(cmd.Context(),
			clients.ListenerParams{
				RequesterID:        id,
				CallbackNamePrefix: "CryptellationCLI",
				Callback: func(_ workflow.Context, params api.ListenToTicksCallbackWorkflowParams) error {
					fmt.Println(params.Tick.String())
					return nil
				},
				Worker:    w,
				TaskQueue: tq,
			},
			ticksListenExchangeFlag,
			ticksListenPairFlag,
		); err != nil {
			return err
		}

		// Run the worker
		return w.Run(nil)
	},
}

func setTicksCommands(cmd *cobra.Command) {
	ticksListenCmd.Flags().StringVarP(&ticksListenExchangeFlag, "exchange", "e", "binance", "Exchange")
	ticksListenCmd.Flags().StringVarP(&ticksListenPairFlag, "pair", "p", "BTC-USDT", "Pair")
	ticksCmd.AddCommand(ticksListenCmd)
	cmd.AddCommand(ticksCmd)
}
