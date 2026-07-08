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
	os.Setenv("TEST_NUM_INVALID", "not-a-number")
	os.Setenv("TEST_DURATION", "100ms")
	os.Setenv("TEST_DURATION_INVALID", "not-a-duration")
	os.Setenv("TEST_EMPTY", "")
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
	var ss config.String

	assert.NoError(t, IntVar(&i, "TEST_NUM"))
	assert.Equal(t, 1048576, i)

	assert.NoError(t, Int64Var(&i64, "TEST_NUM"))
	assert.Equal(t, int64(1048576), i64)

	assert.NoError(t, UintVar(&u, "TEST_NUM"))
	assert.Equal(t, uint(1048576), u)

	assert.NoError(t, Uint64Var(&u64, "TEST_NUM"))
	assert.Equal(t, uint64(1048576), u64)

	assert.NoError(t, Float64Var(&f64, "TEST_NUM"))
	assert.Equal(t, float64(1048576), f64)

	assert.NoError(t, StringVar(&s, "TEST_NUM"))
	assert.Equal(t, "1048576", s)

	assert.NoError(t, DurationVar(&d, "TEST_DURATION"))
	assert.Equal(t, 100*time.Millisecond, d)

	assert.NoError(t, BoolVar(&b, "TEST_BOOL_TRUE"))
	assert.True(t, b)

	assert.NoError(t, BoolVar(&b, "TEST_BOOL_FALSE"))
	assert.False(t, b)

	assert.NoError(t, Var(&ss, "TEST_NUM"))
	assert.Equal(t, "1048576", ss.String())
}

func TestEnvMissing(t *testing.T) {
	i := 1
	i64 := int64(2)
	u := uint(3)
	u64 := uint64(4)
	f64 := float64(5)
	s := "original"
	b := true
	d := time.Second
	ss := config.String{}
	ss.Set("original")

	assert.NoError(t, IntVar(&i, "TEST_MISSING"))
	assert.Equal(t, 1, i)

	assert.NoError(t, Int64Var(&i64, "TEST_MISSING"))
	assert.Equal(t, int64(2), i64)

	assert.NoError(t, UintVar(&u, "TEST_MISSING"))
	assert.Equal(t, uint(3), u)

	assert.NoError(t, Uint64Var(&u64, "TEST_MISSING"))
	assert.Equal(t, uint64(4), u64)

	assert.NoError(t, Float64Var(&f64, "TEST_MISSING"))
	assert.Equal(t, float64(5), f64)

	assert.NoError(t, StringVar(&s, "TEST_MISSING"))
	assert.Equal(t, "original", s)

	assert.NoError(t, BoolVar(&b, "TEST_MISSING"))
	assert.True(t, b)

	assert.NoError(t, DurationVar(&d, "TEST_MISSING"))
	assert.Equal(t, time.Second, d)

	assert.NoError(t, Var(&ss, "TEST_MISSING"))
	assert.Equal(t, "original", ss.String())
}

func TestEnvEmpty(t *testing.T) {
	i := 1
	s := "original"
	b := true

	assert.NoError(t, IntVar(&i, "TEST_EMPTY"))
	assert.Equal(t, 1, i)

	assert.NoError(t, StringVar(&s, "TEST_EMPTY"))
	assert.Equal(t, "original", s)

	assert.NoError(t, BoolVar(&b, "TEST_EMPTY"))
	assert.True(t, b)
}

func TestEnvInvalid(t *testing.T) {
	var i int
	var i64 int64
	var u uint
	var u64 uint64
	var f64 float64
	var b bool
	var d time.Duration

	assert.Error(t, IntVar(&i, "TEST_NUM_INVALID"))
	assert.Error(t, Int64Var(&i64, "TEST_NUM_INVALID"))
	assert.Error(t, UintVar(&u, "TEST_NUM_INVALID"))
	assert.Error(t, Uint64Var(&u64, "TEST_NUM_INVALID"))
	assert.Error(t, Float64Var(&f64, "TEST_NUM_INVALID"))
	assert.Error(t, BoolVar(&b, "TEST_BOOL_INVALID"))
	assert.Error(t, DurationVar(&d, "TEST_DURATION_INVALID"))
}

// BoolVar assigns false before returning the error on a failed parse.
func TestBoolVarErrorSideEffect(t *testing.T) {
	b := true
	err := BoolVar(&b, "TEST_BOOL_INVALID")
	assert.Error(t, err)
	assert.False(t, b)
}

func TestEnvConfig(t *testing.T) {
	var i int
	var i64 int64
	var u uint
	var u64 uint64
	var f64 float64
	var s string
	var b bool
	var d time.Duration
	var ss config.String

	ec := NewConfig(ContinueOnError)

	assert.NoError(t, ec.IntVar(&i, "TEST_NUM"))
	assert.Equal(t, 1048576, i)

	assert.NoError(t, ec.Int64Var(&i64, "TEST_NUM"))
	assert.Equal(t, int64(1048576), i64)

	assert.NoError(t, ec.UintVar(&u, "TEST_NUM"))
	assert.Equal(t, uint(1048576), u)

	assert.NoError(t, ec.Uint64Var(&u64, "TEST_NUM"))
	assert.Equal(t, uint64(1048576), u64)

	assert.NoError(t, ec.Float64Var(&f64, "TEST_NUM"))
	assert.Equal(t, float64(1048576), f64)

	assert.NoError(t, ec.StringVar(&s, "TEST_NUM"))
	assert.Equal(t, "1048576", s)

	assert.NoError(t, ec.Var(&ss, "TEST_NUM"))
	assert.Equal(t, "1048576", ss.String())

	assert.NoError(t, ec.DurationVar(&d, "TEST_DURATION"))
	assert.Equal(t, 100*time.Millisecond, d)

	assert.NoError(t, ec.BoolVar(&b, "TEST_BOOL_TRUE"))
	assert.True(t, b)

	assert.NoError(t, ec.BoolVar(&b, "TEST_BOOL_FALSE"))
	assert.False(t, b)

	assert.Error(t, ec.BoolVar(&b, "TEST_BOOL_INVALID"))
}

func TestEnvConfigInvalid(t *testing.T) {
	var i int
	var i64 int64
	var u uint
	var u64 uint64
	var f64 float64
	var d time.Duration

	ec := NewConfig(ContinueOnError)

	assert.Error(t, ec.IntVar(&i, "TEST_NUM_INVALID"))
	assert.Error(t, ec.Int64Var(&i64, "TEST_NUM_INVALID"))
	assert.Error(t, ec.UintVar(&u, "TEST_NUM_INVALID"))
	assert.Error(t, ec.Uint64Var(&u64, "TEST_NUM_INVALID"))
	assert.Error(t, ec.Float64Var(&f64, "TEST_NUM_INVALID"))
	assert.Error(t, ec.DurationVar(&d, "TEST_DURATION_INVALID"))
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
