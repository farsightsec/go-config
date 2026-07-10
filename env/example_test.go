/*
 * Copyright (c) 2018, 2026 DomainTools LLC
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package env

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/farsightsec/go-config"
	"github.com/stretchr/testify/require"
)

type ExampleConfig struct {
	Title   string
	Version int
	URL     config.URL
}

// Test_Options_Precedence demonstrates the four-level precedence order for user options:
//  1. built-in defaults (lowest)
//  2. configuration file parameters
//  3. environment parameters
//  4. command line parameters (highest)

func Test_Options_Precedence(t *testing.T) {
	var conf ExampleConfig

	// Set defaults and bind command line flags
	fs := flag.NewFlagSet("example", flag.ContinueOnError)

	fs.StringVar(&conf.Title, "title", "default-title", "Application title")
	fs.IntVar(&conf.Version, "version", 1, "App version")
	require.NoError(t, conf.URL.Set("http://default.example/"))
	fs.Var(&conf.URL, "url", "App URL")

	require.Equal(t, "default-title", conf.Title)
	require.Equal(t, 1, conf.Version)
	require.Equal(t, "http://default.example/", conf.URL.String())

	// Allow the config file path to be overridden by the environment
	// This must happen before LoadYAML because the command line is not parsed until later
	confFile := filepath.Join(t.TempDir(), "example.conf")
	require.NoError(t, os.WriteFile(confFile, []byte(
		"title: file-title\nversion: 2\nurl: http://file.example/\n",
	), 0600))

	t.Setenv("EXAMPLE_CONF", confFile)
	require.NoError(t, StringVar(&confFile, "EXAMPLE_CONF"))

	// Load values from the configuration file, overriding defaults
	require.NoError(t, config.LoadYAML(&conf, confFile, false))
	require.Equal(t, "file-title", conf.Title)
	require.Equal(t, 2, conf.Version)
	require.Equal(t, "http://file.example/", conf.URL.String())

	// Import values from environment, overriding file values
	t.Setenv("EXAMPLE_TITLE", "env-title")
	t.Setenv("EXAMPLE_VERSION", "3")
	t.Setenv("EXAMPLE_URL", "http://env.example/")

	require.NoError(t, StringVar(&conf.Title, "EXAMPLE_TITLE"))
	require.NoError(t, IntVar(&conf.Version, "EXAMPLE_VERSION"))
	require.NoError(t, Var(&conf.URL, "EXAMPLE_URL"))

	require.Equal(t, "env-title", conf.Title)
	require.Equal(t, 3, conf.Version)
	require.Equal(t, "http://env.example/", conf.URL.String())

	// Parse command line flags, overriding environment variables
	require.NoError(t, fs.Parse([]string{
		"-title", "cli-title",
		"-version", "4",
		"-url", "http://cli.example/",
	}))

	require.Equal(t, "cli-title", conf.Title)
	require.Equal(t, 4, conf.Version)
	require.Equal(t, "http://cli.example/", conf.URL.String())
}
