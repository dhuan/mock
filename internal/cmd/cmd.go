package cmd

import (
	"fmt"
	"os"

	"github.com/dhuan/mock/internal/types"
	mocklib "github.com/dhuan/mock/pkg/mock"
	"github.com/spf13/cobra"
)

type MockConfig struct {
	Endpoints []types.EndpointConfig `json:"endpoints"`
}

type MockApiResponse struct {
	ValidationErrors *[]mocklib.ValidationError `json:"validation_errors"`
}

var (
	flagConfig             string
	flagPort               string
	flagCors               bool
	flagDelay              int64
	flagRoute              *[]string
	flagMethod             *[]string
	flagStatusCode         *[]int
	flagResponse           *[]string
	flagResponseFile       *[]string
	flagResponseFileServer *[]string
	flagFileServer         *[]string
	flagResponseSh         *[]string
	flagHeader             *[]string
	flagExec               *[]string
	flagResponseExec       *[]string
)

var rootCmd = &cobra.Command{
	Use: "mock",
}

func Execute() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)

	serveCmd.PersistentFlags().StringVarP(&flagConfig, "config", "c", "", "configuration file")
	serveCmd.PersistentFlags().StringVarP(&flagPort, "port", "p", "3000", "port to listen on")
	serveCmd.PersistentFlags().BoolVar(&flagCors, "cors", false, "enable CORS")
	serveCmd.PersistentFlags().Int64VarP(&flagDelay, "delay", "d", 0, "configuration file")
	flagRoute = serveCmd.PersistentFlags().StringArray("route", []string{}, "endpoint route")
	flagMethod = serveCmd.PersistentFlags().StringArray("method", []string{}, "endpoint method")
	flagStatusCode = serveCmd.PersistentFlags().IntSlice("status-code", []int{}, "endpoint response's status code")
	flagResponse = serveCmd.PersistentFlags().StringArray("response", []string{}, "endpoint response")
	flagResponseFile = serveCmd.PersistentFlags().StringArray("response-file", []string{}, "endpoint response file")
	flagResponseFileServer = serveCmd.PersistentFlags().StringArray("response-file-server", []string{}, "endpoint response file server")
	flagFileServer = serveCmd.PersistentFlags().StringArray("file-server", []string{}, "endpoint response file server")
	flagResponseSh = serveCmd.PersistentFlags().StringArray("response-sh", []string{}, "endpoint response script")
	flagHeader = serveCmd.PersistentFlags().StringArray("header", []string{}, "endpoint response header")
	flagExec = serveCmd.PersistentFlags().StringArray("exec", []string{}, "endpoint response exec")
	flagResponseExec = serveCmd.PersistentFlags().StringArray("response-exec", []string{}, "endpoint response exec")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
