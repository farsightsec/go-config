/*
 * Copyright 2023 Farsight Security, Inc.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package env

import (
	"os"
	"testing"
	"time"
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

	ec := NewConfig(ContinueOnError)

	checkOK(t, ec.Var(&i, "TEST_NUM"), i == 1048576)
	checkOK(t, ec.Var(&i64, "TEST_NUM"), i64 == 1048576)
	checkOK(t, ec.Var(&u, "TEST_NUM"), u == 1048576)
	checkOK(t, ec.Var(&u64, "TEST_NUM"), u64 == 1048576)
	checkOK(t, ec.Var(&f64, "TEST_NUM"), f64 == 1048576)
	checkOK(t, ec.Var(&s, "TEST_NUM"), s == "1048576")
	checkOK(t, ec.Var(&d, "TEST_DURATION"), d == 100*time.Millisecond)
	checkOK(t, ec.Var(&b, "TEST_BOOL_TRUE"), b)
	checkOK(t, ec.Var(&b, "TEST_BOOL_FALSE"), !b)

	if ec.Var(&b, "TEST_BOOL_INVALID") == nil {
		t.Error("Expected parse error")
	}
}
