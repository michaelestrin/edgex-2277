/*******************************************************************************
 * Copyright 2017 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * @microservice: core-data-go library
 * @author: Ryan Comer, Dell
 * @version: 0.5.0
 *******************************************************************************/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/edgexfoundry/edgex-go"
	"github.com/edgexfoundry/edgex-go/core/data"
	"github.com/edgexfoundry/edgex-go/support/logging-client"
)

var loggingClient logger.LoggingClient

func main() {
	start := time.Now()
	var (
		useConsul = flag.String("consul", "", "Should the service use consul?")
		useProfile = flag.String("profile", "default", "Specify a profile other than default.")
	)
	flag.Parse()

	// Load configuration data from file.
	// Right now, we always do this first because it contains the Consul endpoint host/port
	configFile := determineConfigFile(*useProfile)
	configuration, err := readConfigurationFile(configFile)
	if err != nil {
		logBeforeTermination(fmt.Errorf("could not load configuration file (%s): %v", configFile, err.Error()))
		return
	}

	//Determine if configuration should be overridden from Consul
	var consulMsg string
	if *useConsul == "y" {
		consulMsg = "Loading configuration from Consul..."
		err := data.ConnectToConsul(*configuration)
		if err != nil {
			logBeforeTermination(err)
			return //end program since user explicitly told us to use Consul.
		}
	} else {
		consulMsg = "Bypassing Consul configuration..."
	}

	logTarget := setLoggingTarget(*configuration)
	// Create Logger (Default Parameters)
	loggingClient = logger.NewClient(configuration.Applicationname, configuration.EnableRemoteLogging, logTarget)

	loggingClient.Info(consulMsg)
	loggingClient.Info(fmt.Sprintf("Starting %s %s ", data.COREDATASERVICENAME, edgex.Version))

	err = data.Init(*configuration, loggingClient)
	if err != nil {
		loggingClient.Error(fmt.Sprintf("call to init() failed: %v", err.Error()))
		return
	}

	r := data.LoadRestRoutes()
	http.TimeoutHandler(nil, time.Millisecond*time.Duration(5000), "Request timed out")
	loggingClient.Info(configuration.Appopenmsg, "")

	startHeartbeat(configuration.Heartbeatmsg, configuration.Heartbeattime)
	// Time it took to start service
	loggingClient.Info("Service started in: "+time.Since(start).String(), "")
	loggingClient.Info("Listening on port: " + strconv.Itoa(configuration.Serverport))
	loggingClient.Error(http.ListenAndServe(":"+strconv.Itoa(configuration.Serverport), r).Error())
}

func startHeartbeat(msg string, interval int) {
	chBeats := make(chan string)
	go data.Heartbeat(msg, interval, chBeats)
	go func() {
		for {
			msg, ok := <-chBeats
			if !ok {
				break
			}
			loggingClient.Info(msg)
		}
		close(chBeats)
	}()
}

func logBeforeTermination(err error) {
	loggingClient = logger.NewClient(data.COREDATASERVICENAME, false, "")
	loggingClient.Error(err.Error())
}

func determineConfigFile(profile string) string {
	switch profile {
		case "docker":
			return "./res/configuration-docker.json"
	    default:
		    return "./res/configuration.json"
	}
}

// Read the configuration file and update configuration struct
func readConfigurationFile(path string) (*data.ConfigurationStruct, error) {
	var configuration data.ConfigurationStruct
	// Read the configuration file
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Decode the configuration as JSON
	err = json.Unmarshal(contents, &configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}

func setLoggingTarget(conf data.ConfigurationStruct) string {
	logTarget := conf.Loggingremoteurl
	if !conf.EnableRemoteLogging {
		return conf.Loggingfile
	}
	return logTarget
}