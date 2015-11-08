// Package commands holds all the commands for the backend_controller program.
// Flags and configurations are extracted here.
package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/scheedule/backend_controller/server"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
)

// Main command
var ControllerCmd = &cobra.Command{
	Use:   "backend_controller",
	Short: "Service controller",
	Long:  "Provide proxy to backend services",
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()

		// Establish services
		var services = map[string]*url.URL{
			"course": {
				Scheme: "http",
				Host:   coursestore_host + ":" + coursestore_port,
			},
			"schedule": {
				Scheme: "http",
				Host:   schedulestore_host + ":" + schedulestore_port,
			},
		}

		server.New("sessionname", "sessionsecret", services)

		log.Fatal(http.ListenAndServe(":"+serve_port, nil))
		log.Info("Listening")
	},
}

var Verbose bool
var schedulestore_host, schedulestore_port, coursestore_host, coursestore_port, serve_port string

// Initialize flags
func init() {
	ControllerCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	ControllerCmd.Flags().StringVarP(
		&schedulestore_host, "schedulestore_host", "", "schedulestore", "Hostname of schedule store")

	ControllerCmd.Flags().StringVarP(
		&schedulestore_port, "schedulestore_port", "", "", "Port of schedule store")

	ControllerCmd.Flags().StringVarP(
		&coursestore_host, "coursestore_host", "", "coursestore", "Hostname of course store")

	ControllerCmd.Flags().StringVarP(
		&coursestore_port, "coursestore_port", "", "", "Port of course store")

	ControllerCmd.Flags().StringVarP(
		&serve_port, "serve_port", "", "8080", "Port to serve endpoint on")
}

// Initialize configuration settings
func InitializeConfig() {
	if Verbose {
		log.SetLevel(log.DebugLevel)
	}
}

// Execute controller command
func Execute() {
	if err := ControllerCmd.Execute(); err != nil {
		panic(err)
	}
}
