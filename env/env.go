/*
 * Copyright 2018 Farsight Security, Inc.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

// Package env provides primitives for loading configuration from environment
// variables. It is loosely modeled after the flag package, in that it can
// associate environment variables with the same types of values as flag can
// associate command line flags with.
//
// Configuration can come from hardcoded defaults, environment, configuration
// files, and command line flags. To implement the usual precedence of:
//
//      1) built-in defaults (lowest)
//      2) environment parameters
//      3) configuration file parameters
//      4) command line parameters (highest)
//
// define defaults and command line bindings (with the flag package) first,
// followed by environment bindings with this package. Then load values from
// configuration files and finally parse the command line flags. See the
// example included in this package for illustration.
//
// Note that in this scheme, an alternate config file location must be provided
// in the environment, and not on the command line, as the command line is
// not parsed until after the configuration file is read.
package env

import (
	"os"
	"strconv"
	"time"
)

// A Value can be converted from a string.
//
// Note that this Value interface is a subset of flag.Value, so types
// which can be used as command line flags values can be used identically
// with environment values.
type Value interface {
	Set(string) error
}

// Var loads Value v's value from the environment variable key, if key is
// set and nonempty.
func Var(v Value, key string) error {
	val := os.Getenv(key)
	if val == "" {
		return nil
	}
	return v.Set(val)
}

// StringVar loads s's string value from environment variable key, if key
// is set and nonempty.
func StringVar(s *string, key string) error {
	v := os.Getenv(key)
	if v != "" {
		*s = v
	}
	return nil
}

// StringVar loads a boolean value from environment variable key, if key
// is set and nonempty. The value associated with key is parsed with
// strconv.ParseBool.
func BoolVar(b *bool, key string) error {
	v := os.Getenv(key)
	if v == "" {
		return nil
	}
	val, err := strconv.ParseBool(v)
	*b = val
	return err
}

type intValue int

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 10, strconv.IntSize)
	*i = intValue(v)
	return err
}

type int64Value int64

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 10, 64)
	*i = int64Value(v)
	return err
}

type uintValue uint

func (u *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 10, strconv.IntSize)
	*u = uintValue(v)
	return err
}

type uint64Value uint64

func (u *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 10, 64)
	*u = uint64Value(v)
	return err
}

type float64Value float64

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = float64Value(v)
	return err
}

// IntVar loads an integer from the environment variable key into *i
func IntVar(i *int, key string) error {
	return Var((*intValue)(i), key)
}

// Int64Var loads a 64-bit integer from the environment variable key into *i
func Int64Var(i *int64, key string) error {
	return Var((*int64Value)(i), key)
}

// UintVar loads an unsigned integer from the environment variable key into *u
func UintVar(u *uint, key string) error {
	return Var((*uintValue)(u), key)
}

// UintVar loads an unsigned 64-bit integer from the environment variable key into *u
func Uint64Var(u *uint64, key string) error {
	return Var((*uint64Value)(u), key)
}

// Float64var loads a floating-point quantity from the environment variable key into *f
func Float64Var(f *float64, key string) error {
	return Var((*float64Value)(f), key)
}

type durationValue time.Duration

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

// DurationVar loads a duration value from the environment variable key into *d.
// The value associated with key may be in any format recognized by time.ParseDuration
func DurationVar(d *time.Duration, key string) error {
	return Var((*durationValue)(d), key)
}
