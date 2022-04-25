package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	version = "0.0.3"
)

// rootCmd represents the base command when called without any subcommands
var (
	lang        string
	APIBaseURL  string
	exsentences string
	exintro     bool

	rootCmd = &cobra.Command{
		Use:   "wpdia-go",
		Short: "Simple cli used to get the description of a given text in Wikipedia",
		Long: `wpdia-go is a simple cli used to get the description of a given text in Wikipedia.
It takes in argument a given text and will retrieve the extract of page content using the 
TextExtracts API (https://www.mediawiki.org/wiki/Extension:TextExtracts#API).

For multi-word search, enclose them using double quotes: "<multi word search>".


The source code is available at https://github.com/lescactus/wpedia-go.`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Fprintf(os.Stderr, "Error: expected 1 argument, got %d\n", len(args))
				os.Exit(1)
			}

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

			displayExtract(extract.Query.Pages[fmt.Sprint(id)])

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

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Set the API base URL corresponding to the desired language
	APIBaseURL = fmt.Sprintf("https://%s.wikipedia.org/w/api.php", lang)
}
