package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

		// Create a worker for the CLI
		taskQueue := "cli-ticks-task-queue"
		w := worker.New(cl.GetTemporalClient(), taskQueue, worker.Options{})
		defer w.Stop()

		// Start the worker in a goroutine
		go func() {
			if err := w.Run(nil); err != nil {
				fmt.Printf("Worker error: %v\n", err)
			}
		}()

		// Create listener params
		requesterID := uuid.New()
		listenerParams := clients.ListenerParams{
			RequesterID: requesterID,
			Callback: func(_ workflow.Context, params api.ListenToTicksCallbackWorkflowParams) error {
				fmt.Println(params.Tick.String())
				return nil
			},
			Worker:    w,
			TaskQueue: taskQueue,
		}

		// Start listening to ticks
		if err := cl.ListenToTicks(cmd.Context(), listenerParams, ticksListenExchangeFlag, ticksListenPairFlag); err != nil {
			return fmt.Errorf("failed to start listening to ticks: %w", err)
		}

		fmt.Printf("Listening to ticks for %s/%s. Press Ctrl-C to stop.\n", ticksListenExchangeFlag, ticksListenPairFlag)

		// Wait for Ctrl-C signal
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		fmt.Println("\nStopping tick listener...")

		// Stop listening to ticks
		if err := cl.StopListeningToTicks(cmd.Context(), requesterID, ticksListenExchangeFlag, ticksListenPairFlag); err != nil {
			return fmt.Errorf("failed to stop listening to ticks: %w", err)
		}

		fmt.Println("Tick listenegit r stopped successfully.")
		return nil
	},
}

func setTicksCommands(cmd *cobra.Command) {
	ticksListenCmd.Flags().StringVarP(&ticksListenExchangeFlag, "exchange", "e", "binance", "Exchange")
	ticksListenCmd.Flags().StringVarP(&ticksListenPairFlag, "pair", "p", "BTC-USDT", "Pair")
	ticksCmd.AddCommand(ticksListenCmd)
	cmd.AddCommand(ticksCmd)
}
