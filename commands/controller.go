// Package commands holds all the commands for the backendcontroller program.
// Flags and configurations are extracted here.
package commands

import (
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/scheedule/backendcontroller/server"
	"github.com/spf13/cobra"
)

// Main command
var controllerCmd = &cobra.Command{
	Use:   "backendcontroller",
	Short: "Service controller",
	Long:  "Provide proxy to backend services",
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()

		// Establish services
		var services = map[string]*url.URL{
			"course": {
				Scheme: "http",
				Host:   coursestoreHost + ":" + coursestorePort,
			},
			"schedule": {
				Scheme: "http",
				Host:   schedulestoreHost + ":" + schedulestorePort,
			},
		}

		server.New("sessionname", "sessionsecret", services, public)

		log.Fatal(http.ListenAndServe(":"+servePort, nil))
		log.Info("serving on port:", servePort)
	},
}

var verbose, public bool
var schedulestoreHost, schedulestorePort, coursestoreHost, coursestorePort, servePort string

// Initialize flags
func init() {
	controllerCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	controllerCmd.Flags().BoolVarP(&public, "public", "p", false, "Set user_id to \"\" and propogate")

	controllerCmd.Flags().StringVarP(
		&schedulestoreHost, "schedulestore_host", "", "schedulestore", "Hostname of schedule store")

	controllerCmd.Flags().StringVarP(
		&schedulestorePort, "schedulestore_port", "", "5000", "Port of schedule store")

	controllerCmd.Flags().StringVarP(
		&coursestoreHost, "coursestore_host", "", "coursestore", "Hostname of course store")

	controllerCmd.Flags().StringVarP(
		&coursestorePort, "coursestore_port", "", "7819", "Port of course store")

	controllerCmd.Flags().StringVarP(
		&servePort, "serve_port", "", "8080", "Port to serve endpoint on")
}

// Initialize configuration settings
func InitializeConfig() {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

// Execute controller command
func Execute() {
	if err := controllerCmd.Execute(); err != nil {
		panic(err)
	}
}
