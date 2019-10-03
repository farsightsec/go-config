/*
 * Copyright 2018 Farsight Security, Inc.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// LoadYAML populates the configuration from the YAML-formatted contents
// of the file `filename`. If `required` is false, LoadYAML returns a nil
// error if the file does not exist.
func LoadYAML(i interface{}, filename string, required bool) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		if !required && os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return yaml.Unmarshal(b, i)
}

// LoadJSON populates the configuration from the JSON-formatted contents
// of the file `filename`. If `required` is false, LoadJSON returns a nil
// error if the file does not exist.
func LoadJSON(i interface{}, filename string, required bool) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		if !required && os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(b, i)
}
