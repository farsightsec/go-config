/*
 * Copyright 2017 Farsight Security, Inc.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

// Package config provides convenience wrappers for useful standard
// library types easing their use in JSON and YAML config files and as
// flag Values for command line tools.
//
// This allows you to Unmarshal JSON or YAML into:
//
//      type myConfig struct {
//              Server    config.URL
//              TLS       config.TLS
//      }
//
// and have a net/url url.URL value available as cfg.Server.URL, and
// a TLS configuration as cfg.TLS.Config, with the config package taking
// care of parsing and validation.
//
package config
