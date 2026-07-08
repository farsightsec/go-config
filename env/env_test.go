/*
 * Copyright (c) 2026 DomainTools LLC
 * Copyright 2023 Farsight Security, Inc.
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
	"time"

	"github.com/farsightsec/go-config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	os.Setenv("TEST_BOOL_FALSE", "false")
	os.Setenv("TEST_BOOL_TRUE", "true")
	os.Setenv("TEST_BOOL_INVALID", "maybe?")
	os.Setenv("TEST_NUM", "1048576")
	os.Setenv("TEST_DURATION", "100ms")
}

func checkOK(t *testing.T, err error, ok bool) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("check failed")
	}
}

func TestEnvTypes(t *testing.T) {
	var i int
	var i64 int64
	var u uint
	var u64 uint64
	var f64 float64
	var s string
	var b bool
	var d time.Duration

	checkOK(t, IntVar(&i, "TEST_NUM"), i == 1048576)
	checkOK(t, Int64Var(&i64, "TEST_NUM"), i64 == 1048576)
	checkOK(t, UintVar(&u, "TEST_NUM"), u == 1048576)
	checkOK(t, Uint64Var(&u64, "TEST_NUM"), u64 == 1048576)
	checkOK(t, Float64Var(&f64, "TEST_NUM"), f64 == 1048576)
	checkOK(t, StringVar(&s, "TEST_NUM"), s == "1048576")
	checkOK(t, DurationVar(&d, "TEST_DURATION"), d == 100*time.Millisecond)
	checkOK(t, BoolVar(&b, "TEST_BOOL_TRUE"), b)
	checkOK(t, BoolVar(&b, "TEST_BOOL_FALSE"), !b)
}

func TestEnvMissing(t *testing.T) {
	i := 10
	b := true
	checkOK(t, IntVar(&i, "TEST_MISSING"), i == 10)
	checkOK(t, BoolVar(&b, "TEST_MISSING"), b)
}

func TestEnvConfig(t *testing.T) {
	var i int = 0
	var i64 int64 = 0
	var u uint = 0
	var u64 uint64 = 0
	var f64 float64 = 0
	var s string = ""
	var b bool = false
	var d time.Duration = time.Second
	var ss = config.String{}

	ec := NewConfig(ContinueOnError)

	checkOK(t, ec.IntVar(&i, "TEST_NUM"), i == 1048576)
	checkOK(t, ec.Int64Var(&i64, "TEST_NUM"), i64 == 1048576)
	checkOK(t, ec.UintVar(&u, "TEST_NUM"), u == 1048576)
	checkOK(t, ec.Uint64Var(&u64, "TEST_NUM"), u64 == 1048576)
	checkOK(t, ec.Float64Var(&f64, "TEST_NUM"), f64 == 1048576)
	checkOK(t, ec.StringVar(&s, "TEST_NUM"), s == "1048576")
	checkOK(t, ec.Var(&ss, "TEST_NUM"), ss.String() == "1048576")
	checkOK(t, ec.DurationVar(&d, "TEST_DURATION"), d == 100*time.Millisecond)
	checkOK(t, ec.BoolVar(&b, "TEST_BOOL_TRUE"), b)
	checkOK(t, ec.BoolVar(&b, "TEST_BOOL_FALSE"), !b)

	if ec.BoolVar(&b, "TEST_BOOL_INVALID") == nil {
		t.Error("Expected parse error")
	}
}

// TestPrecedence verifies the documented loading order:
//
//  1. built-in default (lowest)
//  2. config file
//  3. environment variable
//  4. command line flag (highest)
//
// Each subtest adds one layer and asserts it wins over all lower layers.
// A flag.FlagSet is used instead of the global flag package to avoid
// conflicting with the test runner's own flags.
func TestPrecedence(t *testing.T) {
	// Write a temporary YAML config file used by all subtests.
	cfgFile := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(cfgFile, []byte("value: from-file\n"), 0600))

	type cfg struct {
		Value string `yaml:"value"`
	}

	applyAll := func(t *testing.T, envVal string, flagVal string) string {
		t.Helper()
		c := cfg{Value: "default"}

		// Layer 2: config file overwrites default.
		require.NoError(t, config.LoadYAML(&c, cfgFile, true))

		// Layer 3: env overwrites file (only when set).
		if envVal != "" {
			t.Setenv("PREC_TEST_VALUE", envVal)
		} else {
			os.Unsetenv("PREC_TEST_VALUE")
		}
		require.NoError(t, StringVar(&c.Value, "PREC_TEST_VALUE"))

		// Layer 4: flag overwrites env (only when provided).
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.StringVar(&c.Value, "value", c.Value, "")
		if flagVal != "" {
			require.NoError(t, fs.Parse([]string{"-value", flagVal}))
		} else {
			require.NoError(t, fs.Parse(nil))
		}

		return c.Value
	}

	t.Run("default", func(t *testing.T) {
		c := cfg{Value: "default"}
		assert.Equal(t, "default", c.Value)
	})

	t.Run("file_overrides_default", func(t *testing.T) {
		res := applyAll(t, "", "")
		assert.Equal(t, "from-file", res)
	})

	t.Run("env_overrides_file", func(t *testing.T) {
		res := applyAll(t, "from-env", "")
		assert.Equal(t, "from-env", res)
	})

	t.Run("flag_overrides_env", func(t *testing.T) {
		res := applyAll(t, "from-env", "from-flag")
		assert.Equal(t, "from-flag", res)
	})
}
