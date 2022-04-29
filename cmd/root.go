package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	version = "0.0.6"
)

var (
	// Base URL of the Wikipedia API
	APIBaseURL string

	// Flags
	timeout     time.Duration // http client timeout
	lang        string        // language of the Wikipedia page
	output      string        // output formatter of the program
	exsentences string        // number of sentences to return from a page
	exintro     bool          // whether or not to only the intro of a page
	fullOutput  bool          // whether or not to output also the page namespace and page id

	// validOutputs represents the authorized values for the 'output' flag
	validOutputs = []string{"plain", "pretty", "json", "yaml"}

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "wpdia-go",
		Short: "Simple cli used to get the description of a given text in Wikipedia",
		Long: `wpdia-go is a simple cli used to get the description of a given text in Wikipedia.
It takes in argument a given text and will retrieve the extract of page content using the 
TextExtracts API (https://www.mediawiki.org/wiki/Extension:TextExtracts#API).

For multi-word search, enclose them using double quotes: "<multi word search>".


The source code is available at https://github.com/lescactus/wpedia-go.`,

		// Ensure the 'output' flag value is valid
		PreRun: validateOutputFlag,

		// Only one argument is allowed
		Args: cobra.ExactArgs(1),

		// Main work function
		Run: func(cmd *cobra.Command, args []string) {
			w, err := NewWikiClient(APIBaseURL, "")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			// Get the id of the page requested
			id, err := w.SearchTitle(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			// If the search was unsuccessful
			if id == 0 {
				fmt.Fprintln(os.Stderr, "Error: no page found on Wikipedia for the given query: "+args[0])
				os.Exit(1)
			}

			// User has set 'exsentences' which is mutually exclusive with 'exintro'
			// Disable 'exintro'
			if cmd.Flag("exsentences").Changed {
				exintro = false
			}

			extract, err := w.GetExtract(id)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			page := extract.Query.Pages[fmt.Sprint(id)]

			// Output formatter options
			var d Displayer
			switch output {
			case "plain":
				d = NewPlainFormat()
			case "pretty":
				d = NewPrettyFormat(100)
			case "json":
				d = NewJsonFormat("", "    ")
			case "yaml":
				d = NewYamlFormat()
			}

			// Write extract to the terminal
			d.Write(os.Stdout, &page, fullOutput)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&lang, "lang", "l", "en", "Language. This will set the API endpoint used to retrieve data.")
	rootCmd.PersistentFlags().StringVarP(&exsentences, "exsentences", "s", "10", "How many sentences to return from Wikipedia. Must be between 1 and 10. If > 10, then default to 10. Mutually exclusive with 'exintro'.")
	rootCmd.PersistentFlags().BoolVarP(&exintro, "exintro", "i", true, "Return only content before the first section. Mutually exclusive with 'exsentences'.")
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 15*time.Second, "Timeout value of the http client to the Wikipedia API. Examples values: '10s', '500ms'")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "plain", fmt.Sprintf("Output type. Valid choices are %v.", validOutputs))
	rootCmd.PersistentFlags().BoolVarP(&fullOutput, "full", "f", false, "Also print the page Namespace and page ID.")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Set the API base URL corresponding to the desired language
	APIBaseURL = fmt.Sprintf("https://%s.wikipedia.org/w/api.php", lang)
}

// validateOutputFlag will determine whether the given value of the 'output' flag is valid.
// It exit the program with an error if not.
func validateOutputFlag(cmd *cobra.Command, args []string) {
	if !isPresent(validOutputs, output) {
		fmt.Fprintln(os.Stderr, errors.New(fmt.Sprintf("error: invalid value for flag 'output'. Valid values are %v", validOutputs)))
		os.Exit(1)
	}
}
