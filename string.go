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
	"io/ioutil"
	"os"
	"strings"
)

// String is a string value which can optionally be read from an environment
// variable or file.
//
// A string value beginning with "$" is replaced by the value of the environment
// variable named by the rest of the string. If the value starts with "/", "./",
// or "../", it is replaced by the contents of the file named by the path.
// Otherwise, the string value is used as is.
//
// Marshaling a String marshals the original form (environment variable or file,
// if applicable) in all cases.
type String struct {
	source, value string
}

// String() returns the string value of the string.
func (s *String) String() string {
	return s.value
}

// Set sets the String to the value v, expanding v if it is
// an environment variable or file.
func (s *String) Set(v string) (err error) {
	s.source = v
	if strings.HasPrefix(v, "$") {
		s.value = os.Getenv(v[1:])
	} else if strings.HasPrefix(v, "/") || strings.HasPrefix(v, "./") || strings.HasPrefix(v, "../") {
		buf, err := ioutil.ReadFile(v)
		if err != nil {
			return err
		}
		s.value = strings.TrimSpace(string(buf))
	} else {
		s.value = v
	}
	return nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (s *String) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	return s.Set(v)
}

// UnmarshalYAML satisfies the yaml.Unmarshaler interface
func (s *String) UnmarshalYAML(u func(interface{}) error) error {
	var v string
	if err := u(&v); err != nil {
		return err
	}
	return s.Set(v)
}

// MarshalJSON satisfies the json.Marshaler interface
func (s *String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.source)
}

// MarshalYAML satisfies the yaml.Marshaler interface
func (s *String) MarshalYAML() (interface{}, error) {
	return s.source, nil
}
