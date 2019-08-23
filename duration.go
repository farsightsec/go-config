/*
 * Copyright 2017 Farsight Security, Inc.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package config

import (
	"encoding/json"
	"time"
)

// Duration provides JSON Marshaling and Unmarshaling for time.Duration
// values. The JSON string format is that supported by the Parse() and
// String() methods of time.Duration, e.g., "1m30s", "100ms", etc.
type Duration struct{ time.Duration }

// Set satisfies the flag.Value interface for use as a command line
// flag.
func (d *Duration) Set(s string) (err error) {
	d.Duration, err = time.ParseDuration(s)
	return
}

// MarshalJSON satisfies json.Marshaler
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// MarshalYAML satisfies yaml.Marshaler
func (d Duration) MarshalYAML() (interface{}, error) {
	return d.String(), nil
}

// UnmarshalJSON satisfies json.Unmarshaler
func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return d.Set(s)
}

// UnmarshalYAML satisfies yaml.Unmarshaler
func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	return d.Set(s)
}
