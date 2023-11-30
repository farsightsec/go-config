/*
 * Copyright 2017 Farsight Security, Inc.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// TLSClientAuth provides a convenience wrapper for tls.ClientAuthType and
// conversion to and from string format.
//
// Supported string values are:
//
//	"none":           tls.NoClientCert  (default)
//	"request":        tls.RequestClientCert
//	"require":        tls.RequireAnyClientCert
//	"verify":         tls.VerifyClientCertIfGiven
//	"require+verify": tls.RequireAndVerifyClientCert
type TLSClientAuth struct{ tls.ClientAuthType }

var clientAuthTypes = map[string]tls.ClientAuthType{
	"none":           tls.NoClientCert,
	"request":        tls.RequestClientCert,
	"require":        tls.RequireAnyClientCert,
	"verify":         tls.VerifyClientCertIfGiven,
	"require+verify": tls.RequireAndVerifyClientCert,
}

// String satisfies the flag.Value interface
func (auth *TLSClientAuth) String() string {
	for name, typ := range clientAuthTypes {
		if typ == auth.ClientAuthType {
			return name
		}
	}
	return ""
}

type invalidClientAuthType string

func (i invalidClientAuthType) Error() string {
	return fmt.Sprintf(`Invalid ClientAuthType "%s".`, string(i))
}

type invalidClientAuthTypeValue tls.ClientAuthType

func (i invalidClientAuthTypeValue) Error() string {
	return fmt.Sprintf("Invalid ClientAuthType value %v", int(i))
}

// Set satisfies the flag.Value interface.
func (auth *TLSClientAuth) Set(s string) error {
	if a, ok := clientAuthTypes[s]; ok {
		auth.ClientAuthType = a
		return nil
	}
	return invalidClientAuthType(s)
}

// MarshalJSON satisfies the json.Marshaler interface
func (auth TLSClientAuth) MarshalJSON() ([]byte, error) {
	s := auth.String()
	if s == "" {
		return nil, invalidClientAuthTypeValue(auth.ClientAuthType)
	}
	return json.Marshal(s)
}

// MarshalYAML satisfies the yaml.Marshaler interface
func (auth TLSClientAuth) MarshalYAML() (interface{}, error) {
	s := auth.String()
	if s == "" {
		return nil, invalidClientAuthTypeValue(auth.ClientAuthType)
	}
	return s, nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (auth *TLSClientAuth) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return auth.Set(strings.ToLower(s))
}

// UnmarshalYAML satisfies the yaml.Unmarshaler interface
func (auth *TLSClientAuth) UnmarshalYAML(u func(interface{}) error) error {
	var s string
	if err := u(&s); err != nil {
		return err
	}
	return auth.Set(strings.ToLower(s))
}

// TLSConfig contains the configuration for TLS as it appears on the JSON
// or YAML config. Values parsed from the config are translated and loaded
// into corresponding fields in tls.Config.
type TLSConfig struct {
	RootCAFiles   []string      `json:"rootCAFiles,omitempty" yaml:"rootCAFiles,omitempty"`
	ClientCAFiles []string      `json:"clientCAFiles,omitempty" yaml:"clientCAFiles,omitempty"`
	ClientAuth    TLSClientAuth `json:"clientAuth,omitempty" yaml:"clientAuth,omitempty"`
	Certificates  []struct {
		CertFile string `json:"certFile" yaml:"certFile"`
		KeyFile  string `json:"keyFile" yaml:"keyFile"`
	} `json:"certificates,omitempty" yaml:"certificates,omitempty"`
}

// TLS provides JSON and YAML Marshalers and Unmarshalers for loading
// values into tls.Config.
//
// The JSON and YAML configuration format is provided by the embedded
// type TLSConfig.
type TLS struct {
	TLSConfig
	*tls.Config
}

// MarshalJSON satisfies the json.Marshaler interface
func (t TLS) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.TLSConfig)
}

// MarshalYAML satisfies the yaml.Marshaler interface
func (t TLS) MarshalYAML() (interface{}, error) {
	return t.TLSConfig, nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *TLS) UnmarshalJSON(b []byte) (err error) {
	if err = json.Unmarshal(b, &t.TLSConfig); err != nil {
		return
	}

	t.Config, err = loadTLSConfig(t.TLSConfig)
	return
}

// UnmarshalYAML satisfies the yaml.Unmarshaler interface
func (t *TLS) UnmarshalYAML(u func(interface{}) error) (err error) {
	if err = u(&t.TLSConfig); err != nil {
		return
	}
	t.Config, err = loadTLSConfig(t.TLSConfig)
	return
}

func loadTLSConfig(jc TLSConfig) (*tls.Config, error) {
	var err error
	tc := new(tls.Config)

	tc.ClientAuth = jc.ClientAuth.ClientAuthType

	if len(jc.RootCAFiles) > 0 {
		tc.RootCAs, err = loadCertPool(jc.RootCAFiles)
		if err != nil {
			return nil, err
		}
	}

	if len(jc.ClientCAFiles) > 0 {
		tc.ClientCAs, err = loadCertPool(jc.ClientCAFiles)
		if err != nil {
			return nil, err
		}
	}

	for _, kp := range jc.Certificates {
		cert, err := tls.LoadX509KeyPair(kp.CertFile, kp.KeyFile)
		if err != nil {
			return nil, err
		}
		tc.Certificates = append(tc.Certificates, cert)
	}

	tc.BuildNameToCertificate()
	return tc, nil
}

func loadCertPool(files []string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	for _, f := range files {
		if pem, err := ioutil.ReadFile(f); err == nil {
			pool.AppendCertsFromPEM(pem)
		} else {
			return nil, err
		}
	}
	return pool, nil
}
