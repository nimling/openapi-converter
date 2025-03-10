package main

import (
	"fmt"
	"github.com/nimling/openapi-converter/converter"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var (
	cmdOutputDir            string
	cmdDocsDir              string
	cmdIndexPath            string
	cmdFilePrefix           string
	cmdCommonPrefix         string
	cmdWriteIntroduction    bool
	cmdMergeResponsesInline bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "api-converter",
		Short: "Convert OpenAPIDoc specifications into Nginx Configurations or Vitrepress docs",
		Long:  `A CLI tool that converts OpenAPIDoc 3.1.X YAML files to Nginx location definitions and outputs docs pages for vitepress`,
		RunE:  run,
	}

	rootCmd.Flags().StringVarP(&cmdOutputDir, "output", "o", "", "Output directory for nginx configuration files")
	rootCmd.Flags().StringVarP(&cmdDocsDir, "docs", "d", "", "Output directory for API docs")
	rootCmd.Flags().StringVarP(&cmdIndexPath, "index", "i", "", "Output index.md path to write Vitepress features")
	rootCmd.Flags().StringVarP(&cmdFilePrefix, "file-prefix", "", "", "The output filename prefix to use")
	rootCmd.Flags().StringVarP(&cmdCommonPrefix, "common-prefix", "", "", "The output path prefix to use for the Vitepress links")
	rootCmd.Flags().BoolVarP(&cmdWriteIntroduction, "write-introduction", "", false, "Writes API docs introduction file")
	rootCmd.Flags().BoolVarP(&cmdMergeResponsesInline, "merge-responses-inline", "", false, "Merges inline allOf definitions into a single inline object")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("at least one input file or directory is required")
	}

	if cmdOutputDir != "" {
		if err := os.MkdirAll(cmdOutputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	for _, path := range args {
		if err := processPath(path); err != nil {
			return err
		}
	}

	return nil
}

func processPath(pattern string) error {

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		return fmt.Errorf("no matches found for pattern: %s", pattern)
	}

	// Process each matched path
	for _, path := range matches {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml")) {
					return processFile(path)
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml") {
			err = processFile(path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func processFile(filePath string) error {
	conv, err := converter.NewOpenApiConverter(filePath)
	if err != nil {
		return fmt.Errorf("failed to convert OpenAPIDoc to conv config: %w", err)
	}

	conv.FilePrefix = cmdFilePrefix
	conv.WriteIntroduction = cmdWriteIntroduction
	conv.CommonPrefix = cmdCommonPrefix

	if err = conv.ValidateDocument(); err != nil {
		return fmt.Errorf("error: %s", err)
	}

	if cmdMergeResponsesInline {
		err = conv.MergeResponsesInline()
		if err != nil {
			return fmt.Errorf("error: %s", err)
		}
		fmt.Printf("Successfully merged responses inline for spec %s\n", filePath)
	}

	if cmdOutputDir != "" {
		config, err := conv.WriteNginxConfiguration()
		if err != nil {
			return fmt.Errorf("failed to write nginx config: %w", err)
		}

		outputPath := filepath.Join(cmdOutputDir, filepath.Base(filePath[:len(filePath)-len(filepath.Ext(filePath))])+".conf.template")
		if err := os.WriteFile(outputPath, []byte(config), 0644); err != nil {
			return fmt.Errorf("failed to write nginx config: %w", err)
		}

		fmt.Printf("Successfully converted %s to %s\n", filePath, outputPath)
	}

	if len(cmdDocsDir) > 0 {
		err = conv.WriteVitePressDocs(cmdDocsDir)
		if err != nil {
			return fmt.Errorf("failed to write Vitepress docs: %w", err)
		}
	}

	if cmdIndexPath != "" {
		err = conv.WriteVitePressFeatures(cmdIndexPath)
		if err != nil {
			return fmt.Errorf("error: failed to write markdown: %w", err)
		}

		fmt.Printf("Successfully updated index.md features in %s\n", cmdIndexPath)
	}

	return nil
}
