package main

import (
	"fmt"
	"github.com/nimling/openapi-converter/internal"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	internal.PrintBanner()

	rootCmd := &cobra.Command{
		Use:   "openapi-converter",
		Short: "OpenAPI Converter - Transform API specifications into actionable documentation",
		Long: `OpenAPI Converter (OAC) is a powerful CLI tool designed to streamline API documentation workflows.

It provides two main capabilities:
- Convert OpenAPI/Swagger specifications into Nginx configurations and VitePress documentation
- Synchronize documentation files across projects using pattern-based mapping

Perfect for maintaining consistent API documentation across microservices and documentation pages`,
		Version: "1.0.0",
	}

	rootCmd.AddCommand(internal.NewConvertCommand())
	rootCmd.AddCommand(internal.NewSyncCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
