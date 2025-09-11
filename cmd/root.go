/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	stdJSON "encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dcaravel/jsonfinder/pkg/config"
	"github.com/dcaravel/jsonfinder/pkg/json"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jsonfinder <search_term>",
	Short: "Search JSON documents",
	Long: `Utility for searching large JSON documents.

Will output the matching values in the original JSON structure (with unmatched items removed). 

Optionally - specific extra fields (a.k.a. context) can be included in the output as well.

Will attempt config 'jsonfinder-config.json' from current directory if not explicit config specified.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := loadConfig()
		if err != nil {
			return err
		}

		if rootCmdOpts.context != nil {
			cleanContext := make([]string, 0, len(rootCmdOpts.context))
			for _, c := range rootCmdOpts.context {
				cleanContext = append(cleanContext, strings.TrimSpace(c))
			}

			c.Context = cleanContext
		}

		if rootCmdOpts.file != "" {
			c.FilePath = rootCmdOpts.file
		}

		c.SearchTerm = args[0]
		c.AddIndexes = rootCmdOpts.addIndexes

		items, err := json.Search(c)
		if err != nil {
			return fmt.Errorf("searching: %w", err)
		}

		if len(items) == 0 {
			fmt.Printf("Nothing found")
			return nil
		}

		switch rootCmdOpts.output {
		case "json":
			json.PrintAsJson(c, items)
		case "list":
			json.PrintAsTable(c, items)
		default:
			return fmt.Errorf("unknown output format: %v", rootCmdOpts.output)
		}
		return nil
	},
}

type rootCmdOptions struct {
	file       string
	output     string
	addIndexes bool
	context    []string
	config     string
}

var rootCmdOpts *rootCmdOptions

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func loadConfig() (*config.Config, error) {
	if rootCmdOpts.config != "" {
		// a config file specified, load it.
		return loadConfigFromFile(rootCmdOpts.file)
	}

	// try to find a default config in the current directory
	path := "jsonfinder-config.json"
	c, err := loadConfigFromFile(path)
	if err == nil {
		return c, nil
	} else if !os.IsNotExist(err) {
		// There was an error, and the error is NOT not found
		return nil, fmt.Errorf("loading local config %q: %w", path, err)
	}

	return &config.Config{}, nil
}

func loadConfigFromFile(path string) (*config.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %q: %w", path, err)
	}

	// var c *config.Config
	c := new(config.Config)
	err = stdJSON.Unmarshal(data, c)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling to json: %w", err)
	}

	return c, nil
}

func init() {
	rootCmdOpts = new(rootCmdOptions)

	cmd := rootCmd
	cmd.Flags().StringVarP(&rootCmdOpts.file, "file", "f", "", "File path to parse")
	cmd.Flags().StringVarP(&rootCmdOpts.output, "output", "o", "json", "Output format [json, list]")
	cmd.Flags().StringVarP(&rootCmdOpts.config, "config", "c", "", "Config file to determine parsing and output behavior, CLI flags have precedence to any values in the config")
	cmd.Flags().BoolVarP(&rootCmdOpts.addIndexes, "indexes", "", false, "Add _oindex element to objects in an array that contains the objects array index in the original JSON document")
	cmd.Flags().StringSliceVarP(&rootCmdOpts.context, "context", "", nil, "Comma separated list of json keys to include in the output if a finding traverses the path")

	cmd.MarkFlagRequired("file")
}
