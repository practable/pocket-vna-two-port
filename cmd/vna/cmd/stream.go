/*
Copyright © 2021 Tim Drysdale <timothy.d.drysdale@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ory/viper"
	"github.com/practable/pocket-vna-two-port/pkg/middle"
	"github.com/practable/pocket-vna-two-port/pkg/pocket"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream connects a pocketVNA to a websocket server",
	Long: `Stream connects the first available pocketVNA to a websocket server. The websocket server is specified via an environment variable

export VNA_DESTINATION=ws://localhost:8888/ws/vna
export VNA_CALIBRATION=ws://localhost:8888/ws/calibration
export VNA_RFSWITCH=ws://localhost:8888/ws/rfswitch

vna stream 

Note that development can be enabled by setting environment variable VNA_DEVELOPMENT
export VNA_DEVELOPMENT=true

`,
	Run: func(cmd *cobra.Command, args []string) {

		viper.SetEnvPrefix("VNA")
		viper.AutomaticEnv()

		viper.SetDefault("addr", "localhost:9001")
		viper.SetDefault("baud", 57600)
		viper.SetDefault("log_file", "/var/log/vna/vna.log")
		viper.SetDefault("log_format", "json")
		viper.SetDefault("log_level", "warn")
		viper.SetDefault("port", "/dev/ttyUSB0")
		viper.SetDefault("timeout_usb", "30s")
		viper.SetDefault("timeout_request", "3m")
		viper.SetDefault("topic", "ws://localhost:8888/ws/data")

		addr := viper.GetString("addr")
		baud := viper.GetInt("baud")
		logFile := viper.GetString("log_file")
		logFormat := viper.GetString("log_format")
		logLevel := viper.GetString("log_level")
		port := viper.GetString("port")
		timeoutUSBStr := viper.GetString("timeout_usb")
		timeoutRequestStr := viper.GetString("timeout_request")
		topic := viper.GetString("topic")

		// parse durations

		timeoutRequest, err := time.ParseDuration(timeoutRequestStr)

		if err != nil {
			fmt.Print("cannot parse duration in VNA_TIMEOUT_REQUEST=" + timeoutRequestStr)
			os.Exit(1)
		}

		timeoutUSB, err := time.ParseDuration(timeoutUSBStr)

		if err != nil {
			fmt.Print("cannot parse duration in VNA_TIMEOUT_USB=" + timeoutUSBStr)
			os.Exit(1)
		}

		// set up logging
		switch strings.ToLower(logLevel) {
		case "trace":
			log.SetLevel(log.TraceLevel)
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		case "fatal":
			log.SetLevel(log.FatalLevel)
		case "panic":
			log.SetLevel(log.PanicLevel)
		default:
			fmt.Println("BOOK_LOG_LEVEL can be trace, debug, info, warn, error, fatal or panic but not " + logLevel)
			os.Exit(1)
		}

		switch strings.ToLower(logFormat) {
		case "json":
			log.SetFormatter(&log.JSONFormatter{})
		case "text":
			log.SetFormatter(&log.TextFormatter{})
		default:
			fmt.Println("BOOK_LOG_FORMAT can be json or text but not " + logLevel)
			os.Exit(1)
		}

		if strings.ToLower(logFile) == "stdout" {

			log.SetOutput(os.Stdout) //

		} else {

			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				log.SetOutput(file)
			} else {
				log.Infof("Failed to log to %s, logging to default stderr", logFile)
			}
		}

		// Report useful info
		log.Infof("vna version: %s", versionString())
		log.Infof("addr: [%s]", addr)
		log.Infof("baud: [%d]", baud)
		log.Infof("log file: [%s]", logFile)
		log.Infof("log format: [%s]", logFormat)
		log.Infof("log level: [%s]", logLevel)
		log.Infof("port: [%s]", port)
		log.Infof("topic: [%s]", topic)
		log.Infof("timeoutRequest: [%s]", timeoutUSB)
		log.Infof("timeoutUSB: [%s]", timeoutUSB)

		ctx, cancel := context.WithCancel(context.Background())

		c := make(chan os.Signal, 1)

		signal.Notify(c, os.Interrupt)

		go func() {
			for range c {
				cancel()
				os.Exit(0)
			}
		}()

		// connect to VNA
		v, disconnect, err := pocket.NewHardware()
		defer disconnect()

		m := middle.New(ctx, addr, port, baud, timeoutUSB, timeoutRequest, topic, &v)
		go m.Run()

		<-ctx.Done()

	},
}

func init() {
	rootCmd.AddCommand(streamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// streamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// streamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
