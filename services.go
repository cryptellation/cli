package main

import (
	"fmt"
	"sort"

	"github.com/cryptellation/go-clients/client"
	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:     "services",
	Aliases: []string{"s"},
	Short:   "Manage services",
}

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i"},
	Short:   "Read info from services",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Set client
		cl, err := client.New(client.WithTemporalAddress(temporalAddress))
		if err != nil {
			return err
		}
		defer cl.Close()

		// Get info
		res, err := cl.ServicesInfo(cmd.Context())
		if err != nil {
			return err
		}

		switch {
		case jsonOutput:
			return displayJSON(res)
		default:
			keys := make([]string, 0, len(res))
			for k := range res {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			format := "%-20s %+v\n"
			fmt.Printf(format, "NAME", "DATA")
			for _, k := range keys {
				fmt.Printf(format, k, res[k])
			}
		}

		return nil
	},
}

func setServicesCommands(cmd *cobra.Command) {
	servicesCmd.AddCommand(infoCmd)
	cmd.AddCommand(servicesCmd)
}
