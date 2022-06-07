package cmd

import (
	"fmt"
	"os"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/types"
	"github.com/spf13/cobra"
)

type MockConfig struct {
	Endpoints []types.EndpointConfig `json:"endpoints"`
}

type MockApiResponse struct {
	ValidationErrors *[]mock.ValidationError `json:"validation_errors"`
}

var (
	flagConfig string
	flagPort   string
)

var rootCmd = &cobra.Command{
	Use: "mock",
}

func Execute() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().StringVarP(&flagConfig, "config", "c", "", "configuration file")
	serveCmd.PersistentFlags().StringVarP(&flagPort, "port", "p", "3000", "port to listen on")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
