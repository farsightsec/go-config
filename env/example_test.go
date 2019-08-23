/*
 * Copyright 2018 Farsight Security, Inc.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package env

import (
	"flag"
	"log"

	"github.com/farsightsec/go-config"
)

type ExampleConfig struct {
	Title   string
	Version int
	URL     config.URL
}

func Example() {
	var conf ExampleConfig
	confFile := "/etc/example/example.conf"

	// First, set up command line arg bindings and default values with
	// the flag package.
	//
	flag.StringVar(&conf.Title, "title", "Default Title", "Application title")
	flag.IntVar(&conf.Version, "version", 2, "App version")

	// Default values not associated with a flag (or for flag.Var) can
	// be set explicitly here
	conf.URL.Set("http://www.farsightsecurity.com/")
	flag.Var(&conf.URL, "url", "App URL")

	// Next, import new defaults from the environment with this package.
	StringVar(&conf.Title, "EXAMPLE_TITLE")
	if err := IntVar(&conf.Version, "EXAMPLE_VERSION"); err != nil {
		log.Fatal("Invalid EXAMPLE_VERSION value: ", err)
	}
	if err := Var(&conf.URL, "EXAMPLE_URL"); err != nil {
		log.Fatal("Invalid EXAMPLE_URL value: ", err)
	}

	// Allow config file to be overridden by environment
	StringVar(&confFile, "EXAMPLE_CONF")

	// Load values from configuration file
	err := config.LoadYAML(&conf, confFile, false)
	if err != nil {
		log.Fatalf("Failed to load config from %s: %v", confFile, err)
	}

	// Finally, read values from command line arguments
	flag.Parse()
}
