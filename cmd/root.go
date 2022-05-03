package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/lescactus/wpdia-go/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	version = "0.1.0"
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

	log       *logrus.Logger
	logLevel  string
	logFormat string

	randomPage bool // whether or not to look for a random page

	// validOutputs represents the authorized values for the 'output' flag
	validOutputs = []string{"plain", "pretty", "json", "yaml"}

	// validLogLevel represents the authorized values for the 'loglevel' flag
	validLogLevels = []string{"debug", "info", "warn", "error"}

	// validLogLevel represents the authorized values for the 'logformat' flag
	validLogFormats = []string{"text", "json"}

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "wpdia-go",
		Short: "Simple cli used to get the description of a given text in Wikipedia",
		Long: `wpdia-go is a simple cli used to get the description of a given text in Wikipedia.
It takes in argument a given text and will retrieve the extract of page content using the 
TextExtracts API (https://www.mediawiki.org/wiki/Extension:TextExtracts#API).

For multi-word search, enclose them using double quotes: "<multi word search>".


The source code is available at https://github.com/lescactus/wpedia-go.`,

		// Ensure the 'output' and 'loglevel' flags values are valid
		PreRunE: validateFlags,

		// Main work function
		Run: func(cmd *cobra.Command, args []string) {
			var title string

			// When the '--random' flag is set, we don't need anything in argument
			// Oherwise we do
			if randomPage {
				if len(args) > 0 {
					log.WithFields(logrus.Fields{
						"level": logLevel,
					}).Warn(fmt.Sprintf("The --random is set, the given arguments will be ignored: %v", args))
				}
			} else {
				title = args[0]
			}

			log.WithFields(logrus.Fields{
				"level": logLevel,
				"url":   APIBaseURL,
			}).Info("Creating new Wiki client...")

			w, err := NewWikiClient(APIBaseURL, "")
			if err != nil {
				log.SetOutput(os.Stderr)
				log.WithFields(logrus.Fields{
					"level": logLevel,
					"url":   APIBaseURL,
				}).Error(err)
				os.Exit(1)
			}

			log.WithFields(logrus.Fields{
				"level": logLevel,
				"url":   APIBaseURL,
			}).Debug("New Wiki client created")

			log.WithFields(logrus.Fields{
				"level": logLevel,
			}).Debug("Disabling 'exintro'...")

			// User has set 'exsentences' which is mutually exclusive with 'exintro'
			// Disable 'exintro'
			if cmd.Flag("exsentences").Changed {
				exintro = false
			}

			log.WithFields(logrus.Fields{
				"level":  logLevel,
				"title":  title,
				"random": randomPage,
			}).Info("Getting text extract...")

			var extract *WikiTextExtractResponse
			if randomPage {
				// Call the Random API
				extract, err = w.GetExtractRandom()
			} else {
				log.WithFields(logrus.Fields{
					"level": logLevel,
					"title": title,
				}).Info("Searching title...")

				// Get the id of the page requested
				var id uint64
				id, err = w.SearchTitle(title)
				if err != nil {
					log.SetOutput(os.Stderr)
					log.WithFields(logrus.Fields{
						"level": logLevel,
						"url":   APIBaseURL,
						"title": title,
					}).Error(err)
					os.Exit(1)
				}

				log.WithFields(logrus.Fields{
					"level": logLevel,
					"title": title,
				}).Debug("Title found")

				// If the search was unsuccessful
				if id == 0 {
					log.SetOutput(os.Stderr)
					log.WithFields(logrus.Fields{
						"level": logLevel,
						"title": title,
					}).Error("Error: no page found on Wikipedia for the given query")
					os.Exit(1)
				}

				// Call the TextExtracts API for the requested page id
				extract, err = w.GetExtract(id)
			}

			if err != nil {
				log.SetOutput(os.Stderr)
				log.WithFields(logrus.Fields{
					"level":  logLevel,
					"title":  title,
					"random": randomPage,
				}).Errorf("Error: %s", err.Error())
				os.Exit(1)
			}

			log.WithFields(logrus.Fields{
				"level":  logLevel,
				"title":  title,
				"random": randomPage,
			}).Debug("Text extract found")

			// Because we request only 1 page from Wikipedia's API,
			// extract.Query.Pages **should be** a map of only one element
			// If it is unexpectedly not the case, exit the program immediately with an error.
			if len(extract.Query.Pages) != 1 {
				log.SetOutput(os.Stderr)
				log.WithFields(logrus.Fields{
					"level": logLevel,
				}).Fatal(fmt.Sprintf("Expected an anwser of 1 page, got %d", len(extract.Query.Pages)))
				os.Exit(1)
			}

			var page Page
			for _, v := range extract.Query.Pages {
				page = v
			}

			// Ensure the page isn't a disambiguation
			// In the case it is, simply print a message saying to refine the search
			if page.IsDisambiguation() {
				log.WithFields(logrus.Fields{
					"level": logLevel,
					"title": title,
					"id":    *page.Pageid,
				}).Warn("The requested page is a disambiguation page")

				page.Extract = `/!\ The requested page is a disambiguation page /!\

A disambiguation page is Wikipedia's way of resolving conflicts that arise when a potential article title is ambiguous - most often because it refers to more than one subject covered by Wikipedia.
For example, Mercury can refer to a chemical element, a planet, a Roman god, and many other things.

Try to refine the search in a more precise manner. Example:
	'Nancy France' instead of 'Nancy' - or 'Go verb' instead of 'Go'`
			}

			log.WithFields(logrus.Fields{
				"level": logLevel,
			}).Debug("Setting formatter...")

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

			log.WithFields(logrus.Fields{
				"level": logLevel,
			}).Debugf("Formatter set to %s", output)

			// Write extract to the terminal
			d.Write(os.Stdout, &page, fullOutput)
			if err != nil {
				log.SetOutput(os.Stderr)
				log.WithFields(logrus.Fields{
					"level": logLevel,
				}).Error(err)
				os.Exit(1)
			}
		},

		Version: version,
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
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "e", "error", fmt.Sprintf("Log level verbosity. Accepted values are %v.", validLogLevels))
	rootCmd.PersistentFlags().StringVarP(&logFormat, "logformat", "a", "text", fmt.Sprintf("Log format. Accepted values are %v.", validLogFormats))
	rootCmd.PersistentFlags().BoolVarP(&randomPage, "random", "r", false, "Return a random article.")

	cobra.OnInitialize(initConfig, setLogger)
}

func initConfig() {
	// Set the API base URL corresponding to the desired language
	APIBaseURL = fmt.Sprintf("https://%s.wikipedia.org/w/api.php", lang)
}

func setLogger() {
	var err error
	log, err = logger.NewLogger(logLevel, logFormat)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// validateFlags will determine whether the given value of the 'output', 'loglevel' and 'logformat' flags are valid.
// It exit the program with an error if not.
func validateFlags(cmd *cobra.Command, args []string) error {
	if !isPresent(validOutputs, output) {
		return errors.New(fmt.Sprintf("error: invalid value for flag 'output'. Valid values are %v", validOutputs))
	}

	if !isPresent(validLogLevels, logLevel) {
		return errors.New(fmt.Sprintf("error: invalid value for flag 'loglevel'. Valid values are %v", validLogLevels))
	}

	if !isPresent(validLogFormats, logFormat) {
		return errors.New(fmt.Sprintf("error: invalid value for flag 'logformat'. Valid values are %v", validLogFormats))
	}

	return nil
}

// hasOneArg will verify whether the slice passed in agument contains only one element.
// It will return true if the slice contains only one element, false otherwise.
func hasOneArg(args []string) bool {
	if len(args) != 1 {
		return false
	}
	return true
}
