package cmd

import (
	"fmt"
	"os"

	"github.com/dhuan/mock/internal/types"
	mocklib "github.com/dhuan/mock/pkg/mock"
	"github.com/spf13/cobra"
)

type MockConfig struct {
	Base        string                   `json:"base"`
	Endpoints   []types.EndpointConfig   `json:"endpoints"`
	Middlewares []types.MiddlewareConfig `json:"middlewares"`
}

type MockApiResponse struct {
	ValidationErrors *[]mocklib.ValidationError `json:"validation_errors"`
}

var (
	flagConfig    string
	flagPort      string
	flagCors      bool
	flagRegex     bool
	flagAppend    bool
	flagJson      bool
	flagValueOnly bool
	flagDelay     int64
	flagBaseApi   string
)

var rootCmd = &cobra.Command{
	Use: "mock",
}

func Execute() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(forwardCmd)
	rootCmd.AddCommand(replaceCmd)
	rootCmd.AddCommand(getRouteParamCmd)
	rootCmd.AddCommand(wipeHeadersCmd)
	rootCmd.AddCommand(setHeaderCmd)
	rootCmd.AddCommand(writeCmd)
	rootCmd.AddCommand(getQueryCmd)
	rootCmd.AddCommand(getPayloadCmd)
	rootCmd.AddCommand(getHeaderCmd)
	rootCmd.AddCommand(setStatusCmd)
	rootCmd.AddCommand(versionCmd)

	serveCmd.PersistentFlags().StringVarP(&flagConfig, "config", "c", "", "configuration file")
	serveCmd.PersistentFlags().StringVarP(&flagPort, "port", "p", "UNSET", "port to listen on")
	serveCmd.PersistentFlags().BoolVar(&flagCors, "cors", false, "enable CORS")
	serveCmd.PersistentFlags().Int64VarP(&flagDelay, "delay", "d", 0, "configuration file")
	serveCmd.PersistentFlags().StringArray("route", []string{}, "endpoint route")
	serveCmd.PersistentFlags().StringArray("method", []string{}, "endpoint method")
	serveCmd.PersistentFlags().IntSlice("status-code", []int{}, "endpoint response's status code")
	serveCmd.PersistentFlags().StringArray("response", []string{}, "endpoint response")
	serveCmd.PersistentFlags().StringArray("response-file", []string{}, "endpoint response file")
	serveCmd.PersistentFlags().StringArray("response-file-server", []string{}, "endpoint response file server")
	serveCmd.PersistentFlags().StringArray("file-server", []string{}, "endpoint response file server")
	serveCmd.PersistentFlags().StringArray("response-sh", []string{}, "endpoint response script")
	serveCmd.PersistentFlags().StringArray("shell-script", []string{}, "endpoint response script")
	serveCmd.PersistentFlags().StringArray("header", []string{}, "endpoint response header")
	serveCmd.PersistentFlags().StringArray("exec", []string{}, "endpoint response exec")
	serveCmd.PersistentFlags().StringArray("response-exec", []string{}, "endpoint response exec")
	serveCmd.PersistentFlags().StringArray("middleware", []string{}, "middleware")
	serveCmd.PersistentFlags().StringArray("route-match", []string{}, "filter middleware by route")
	serveCmd.PersistentFlags().StringVarP(&flagBaseApi, "base", "b", "", "base API")

	writeCmd.PersistentFlags().BoolVarP(&flagAppend, "append", "a", false, "append instead of overwriting")
	writeCmd.PersistentFlags().BoolVar(&flagJson, "json", false, "treats received data as JSON and adds necessary JSON header.")
	wipeHeadersCmd.PersistentFlags().BoolVar(&flagRegex, "regex", false, "enable regular expression for seaching")
	replaceCmd.PersistentFlags().BoolVar(&flagRegex, "regex", false, "enable regular expression for seaching")
	getHeaderCmd.PersistentFlags().BoolVarP(&flagValueOnly, "value", "v", false, "get value only")
	getHeaderCmd.PersistentFlags().BoolVar(&flagRegex, "regex", false, "enable regular expression for seaching")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
