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
	"net/url"
)

// URL provides JSON Marshaling and Unmarshaling of URLs, internally
// representing them as *url.URL from net/url.
type URL struct{ *url.URL }

// Set satisfies the flag.Value interface along with the net.URL String()
// method.
func (u *URL) Set(s string) (err error) {
	u.URL, err = url.Parse(s)
	return
}

// UnmarshalJSON satisfies json.Unmarshaler
func (u *URL) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return u.Set(s)
}

// MarshalJSON satisfies json.Marshaler
func (u URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalYAML satisfies yaml.Unmarshaler
func (u *URL) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	return u.Set(s)
}

// MarshalYAML satisfies yaml.Marshaler
func (u URL) MarshalYAML() (interface{}, error) {
	return u.String(), nil
}
