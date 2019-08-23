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
	"fmt"
	"net"
	"strings"
)

// Addr is a generic network address with JSON and YAML Marshaler and
// Unmarshaler methods.
type Addr struct{ net.Addr }

type addr struct {
	Net  string
	Addr string
}

func (a *addr) Network() string {
	return a.Net
}

func (a *addr) String() string {
	return a.Addr
}

type errAddrFormatInvalid string

func (e errAddrFormatInvalid) Error() string {
	return fmt.Sprintf("Invalid address format '%s':"+
		" should be net:addr", string(e))
}

// Set satisfies flag.Value for use with command line flags.
func (a *Addr) Set(s string) error {
	l := strings.SplitN(s, ":", 2)
	if len(l) < 2 {
		return errAddrFormatInvalid(s)
	}
	a.Addr = &addr{l[0], l[1]}
	return nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (a *Addr) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return a.Set(s)
}

// UnmarshalYAML satisfies the yaml.Unmarshaler interface
func (a *Addr) UnmarshalYAML(u func(interface{}) error) error {
	var s string
	if err := u(&s); err != nil {
		return err
	}
	return a.Set(s)
}

// MarshalJSON satisfies the json.Marshaler interface
func (a Addr) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%s:%s", a.Network(), a.String()))
}

// MarshalYAML satisfies the yaml.Marshaler interface
func (a Addr) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%s:%s", a.Network(), a.String()), nil
}

type errInvalidUDPNetwork string

func (e errInvalidUDPNetwork) Error() string {
	return fmt.Sprintf("Invalid UDP network '%s'", string(e))
}

type errInvalidTCPNetwork string

func (e errInvalidTCPNetwork) Error() string {
	return fmt.Sprintf("Invalid TCP network '%s'", string(e))
}

type errInvalidUnixNetwork string

func (e errInvalidUnixNetwork) Error() string {
	return fmt.Sprintf("Invalid Unix network '%s'", string(e))
}

// UDPAddr is an address restricted to be in the network "udp",
// "udp4", or "udp6". Other networks are considered an error.
type UDPAddr struct{ *net.UDPAddr }

// Set satisfies flag.Value for use in command line arguments
func (u *UDPAddr) Set(s string) (err error) {
	var a Addr
	if err = a.Set(s); err != nil {
		return
	}
	n := a.Network()
	if n != "udp" && n != "udp4" && n != "udp6" {
		return errInvalidUDPNetwork(n)
	}
	u.UDPAddr, err = net.ResolveUDPAddr(n, a.String())
	return
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (u *UDPAddr) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return u.Set(s)
}

// UnmarshalYAML satisfies the yaml.Unmarshaler interface
func (u *UDPAddr) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	return u.Set(s)
}

// MarshalJSON satisfies the json.Marshaler interface
func (u UDPAddr) MarshalJSON() ([]byte, error) {
	a := Addr{u.UDPAddr}
	return a.MarshalJSON()
}

// MarshalYAML satisfies the yaml.Marshaler interface
func (u UDPAddr) MarshalYAML() (interface{}, error) {
	a := Addr{u.UDPAddr}
	return a.MarshalYAML()
}

// TCPAddr is an address restricted to be in the "tcp", "tcp4", or "tcp6"
// networks.
type TCPAddr struct{ *net.TCPAddr }

// Set satisfies flag.Value for use in command line arguments
func (t *TCPAddr) Set(s string) (err error) {
	var a Addr
	if err = a.Set(s); err != nil {
		return
	}
	n := a.Network()
	if n != "tcp" && n != "tcp4" && n != "tcp6" {
		return errInvalidTCPNetwork(n)
	}
	t.TCPAddr, err = net.ResolveTCPAddr(n, a.String())
	return
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *TCPAddr) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return t.Set(s)
}

// UnmarshalYAML satisfies the yaml.Unmarshaler interface
func (t *TCPAddr) UnmarshalYAML(u func(interface{}) error) error {
	var s string
	if err := u(&s); err != nil {
		return err
	}
	return t.Set(s)
}

// MarshalJSON satisfies the json.Marshaler interface
func (t TCPAddr) MarshalJSON() ([]byte, error) {
	a := Addr{t.TCPAddr}
	return a.MarshalJSON()
}

// MarshalYAML satisfies the yaml.Marshaler interface
func (t TCPAddr) MarshalYAML() (interface{}, error) {
	a := Addr{t.TCPAddr}
	return a.MarshalYAML()
}

// UnixAddr is a unix-domain socket address in the "unix", "unixpacket",
// or "unixgram" network
type UnixAddr struct{ *net.UnixAddr }

// Set satisfies flag.Value for use in command line arguments
func (u *UnixAddr) Set(s string) (err error) {
	var a Addr
	if err = a.Set(s); err != nil {
		return
	}
	n := a.Network()
	if n != "unix" && n != "unixpacket" && n != "unixgram" {
		return errInvalidUnixNetwork(n)
	}
	u.UnixAddr, err = net.ResolveUnixAddr(n, a.String())
	return
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (u *UnixAddr) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return u.Set(s)
}

// UnmarshalYAML satisfies the yaml.Unmarshaler interface
func (u *UnixAddr) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	return u.Set(s)
}

// MarshalJSON satisfies the json.Marshaler interface
func (u UnixAddr) MarshalJSON() ([]byte, error) {
	a := Addr{u.UnixAddr}
	return a.MarshalJSON()
}

// MarshalYAML satisfies the yaml.Marshaler interface
func (u UnixAddr) MarshalYAML() (interface{}, error) {
	a := Addr{u.UnixAddr}
	return a.MarshalYAML()
}
