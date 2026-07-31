/*
 * Copyright 2026 DomainTools LLC
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package config

import (
	"crypto/tls"
	"fmt"
	"os"
	"testing"

	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

func getTempFile(t *testing.T, ext string, content string) string {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "tmpdata."+ext)
	if err != nil {
		t.Errorf("Error creating temp file of ext %s", ext)
		return ""
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		t.Errorf("Error writing to tmp file of ext %s", ext)
		return ""
	}

	tmpfile.Close()
	return tmpfile.Name()
}

func testMarshalling(t *testing.T, value interface{}, teststr string, expSuccess bool) {
	t.Helper()

	var mb []byte
	var mb2 interface{}
	var err error = nil
	var err2 error = nil

	switch value.(type) {
	case *TCPAddr:
		a := value.(*TCPAddr)
		mb, err = a.MarshalJSON()
		mb2, err2 = a.MarshalYAML()
	case *UDPAddr:
		a := value.(*UDPAddr)
		mb, err = a.MarshalJSON()
		mb2, err2 = a.MarshalYAML()
	case *Addr:
		a := value.(*Addr)
		mb, err = a.MarshalJSON()
		mb2, err2 = a.MarshalYAML()
	case *Duration:
		d := value.(*Duration)
		mb, err = d.MarshalJSON()
		mb2, err2 = d.MarshalYAML()
	case *String:
		s := value.(*String)
		mb, err = s.MarshalJSON()
		mb2, err2 = s.MarshalYAML()
	case *URL:
		u := value.(*URL)
		mb, err = u.MarshalJSON()
		mb2, err2 = u.MarshalYAML()
	default:
		t.Errorf("Test function received unexpected type: %T", value)
		return
	}

	jstr := fmt.Sprintf("\"%s\"", teststr)

	if !expSuccess {
		if err == nil && string(mb) == jstr {
			t.Errorf("Unexpected successful marshaling of bad %T in JSON form", value)
		}
	} else if err != nil {
		t.Errorf("Error marshaling %T in JSON form: %v", value, err)
	} else if string(mb) != jstr {
		t.Errorf("Marshaled JSON %T was not expected value: %v", value, string(mb))
	}

	if !expSuccess {
		if err2 == nil {
			if ystr, ok := mb2.(string); ok != true || teststr == ystr {
				t.Errorf("Unexpected successful marshaling of bad %T in YAML form", value)
			}
		}
	} else if err2 != nil {
		t.Errorf("Error marshaling %T in YAML form: %v", value, err2)
	} else if ystr, ok := mb2.(string); ok != true || teststr != ystr {
		t.Errorf("Marshaled YAML %T was not expected value: %v", value, ystr)
	}

	if expSuccess {
		var err error = nil
		var err2 error = nil

		badjson := []byte("BAD, JSON, STRING")

		switch value.(type) {
		case *TCPAddr:
			a := &TCPAddr{}
			err = a.UnmarshalJSON(mb)
			err2 = a.UnmarshalJSON(badjson)
		case *UDPAddr:
			a := &UDPAddr{}
			err = a.UnmarshalJSON(mb)
			err2 = a.UnmarshalJSON(badjson)
		case *Addr:
			a := &Addr{}
			err = a.UnmarshalJSON(mb)
			err2 = a.UnmarshalJSON(badjson)
		case *Duration:
			d := &Duration{}
			err = d.UnmarshalJSON(mb)
			err2 = d.UnmarshalJSON(badjson)
		case *String:
			s := &String{}
			err = s.UnmarshalJSON(mb)
			err2 = s.UnmarshalJSON(badjson)
		case *URL:
			u := &URL{}
			err = u.UnmarshalJSON(mb)
			err2 = u.UnmarshalJSON(badjson)
		}

		if err != nil {
			t.Errorf("Unmarshaling of JSON data into type %T produced error: %v", value, err)
		}
		if err2 == nil {
			t.Errorf("Unmarshaling of bad JSON data into type %T failed to produce error", value)
		}
	}
}

func TestBadAddr(t *testing.T) {
	t.Run("Addr/valid", func(t *testing.T) {
		a := &Addr{}
		if err := a.Set("udp:1.2.3.4"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testMarshalling(t, a, "udp:1.2.3.4", true)
	})

	cases := []struct {
		name string
		set  func(string) error
		in   string
	}{
		{
			"Addr/malformed",
			func(s string) error { return (&Addr{}).Set(s) },
			"BADADDR",
		},
		{
			"UDPAddr/wrong_network",
			func(s string) error { return (&UDPAddr{}).Set(s) },
			"tcp:1.2.3.4",
		},
		{
			"UDPAddr/malformed",
			func(s string) error { return (&UDPAddr{}).Set(s) },
			"BADADDR",
		},
		{
			"TCPAddr/wrong_network",
			func(s string) error { return (&TCPAddr{}).Set(s) },
			"udp:1.2.3.4",
		},
		{
			"TCPAddr/malformed",
			func(s string) error { return (&TCPAddr{}).Set(s) },
			"BADADDR",
		},
		{
			"UnixAddr/wrong_network",
			func(s string) error { return (&UnixAddr{}).Set(s) },
			"udp:1.2.3.4",
		},
		{
			"UnixAddr/malformed",
			func(s string) error { return (&UnixAddr{}).Set(s) },
			"BADADDR",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.set(tc.in); err == nil {
				t.Errorf("expected error for input %q", tc.in)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	ustr := "100ms"
	d := &Duration{}
	d.Set(ustr)
	testMarshalling(t, d, ustr, true)

	ustr = "1m30duif"
	d.Set(ustr)
	testMarshalling(t, d, ustr, false)
}

func TestString(t *testing.T) {
	ustr := "bla bla bla"
	s := &String{}
	s.Set(ustr)

	if s.String() != ustr {
		t.Errorf("Unexpected string set error")
	}

	testMarshalling(t, s, ustr, true)

	envname := "ENV_DOES_NOT_EXIST"
	envval := "abc123"
	// os.Unsetenv(envname)
	os.Setenv(envname, envval)
	if err := s.Set("$" + envname); err != nil {
		t.Errorf("Unexpected error setting string from environment variable")
	} else if s.String() != envval {
		t.Errorf("Expected environment variable string to contain different data: %v", s.String())
	}

	if err := s.Set("/file/does/not/exist"); err == nil {
		t.Errorf("Expected string loaded from non-existent file to generate erro rbut it did not")
	}

	tmpdata := "this is some test data"

	if tmpf := getTempFile(t, "txt", tmpdata); tmpf != "" {
		defer os.Remove(tmpf)

		if err := s.Set(tmpf); err != nil {
			t.Errorf("Error loading string from file: %v", err)
		} else if s.String() != tmpdata {
			t.Errorf("Unexpected data in string loaded from file: %s", s.String())
		}

	} else {
		t.Errorf("Failed to get tmpfile for String test")
	}

}

func TestURL(t *testing.T) {
	ustr := "http://www.google.com:443/hello"
	u := &URL{}
	u.Set(ustr)

	testMarshalling(t, u, ustr, true)
}

func TestFile(t *testing.T) {
	var dummystr interface{}

	if dne_path, err := os.Getwd(); err == nil {
		dne_path += "/doooooooooes_nooooooot_exist"

		if err := LoadJSON(&dummystr, dne_path, true); err == nil {
			t.Errorf("Expected non-existent JSON config load to produce error")
		}

		if err := LoadJSON(&dummystr, dne_path, false); err != nil {
			t.Errorf("Unexpected error produced by non-existent non-required JSON config: %v", err)
		}

		if err := LoadYAML(&dummystr, dne_path, true); err == nil {
			t.Errorf("Expected non-existent YAML config load to produce error")
		}

		if err := LoadYAML(&dummystr, dne_path, false); err != nil {
			t.Errorf("Unexpected error produced by non-existent non-required YAML config: %v", err)
		}

	} else {
		t.Errorf("Error getting current working directory")
	}

	jdata := "{ \"str1\":\"an example string\", \"num1\":31337 }"
	ydata := "str1: an example string\nnum1: 31337"

	if tmpj := getTempFile(t, "json", jdata); tmpj != "" {
		defer os.Remove(tmpj)

		if err := LoadJSON(&dummystr, tmpj, true); err != nil {
			t.Errorf("Error loading dummy JSON config: %v", err)
		}

	} else {
		t.Errorf("Failed to get JSON tmpfile")
	}

	if tmpy := getTempFile(t, "yml", ydata); tmpy != "" {
		defer os.Remove(tmpy)

		if err := LoadYAML(&dummystr, tmpy, true); err != nil {
			t.Errorf("Error loading dummy YAML config: %v", err)
		}

	} else {
		t.Errorf("Failed to get YAML tmpfile")
	}
}

type ExConf struct {
	Str1      String
	EnvString String
	Num1      uint64
	Url       URL
	Dur       Duration
	Addr      Addr
	UAddr     UnixAddr
	TAddr     TCPAddr
	UDAddr    UDPAddr
}

func TestFile2(t *testing.T) {
	var dummyconf1, dummyconf2 ExConf

	jdata := `{ "str1": "an example string", "envstring": "$HOME", "num1": 31337, "url": "http://www.google.com", "dur": "90s", "taddr": "tcp:1.2.3.4:80", "udaddr": "udp:4.2.2.4:53", "uaddr": "unix:/tmp/sock", "addr": "tcp:1.2.3.4:443" }`
	ydata := "str1: an example string\nenvstring: $HOME\nnum1: 31337\nurl: http://www.google.com\ndur: 90s\ntaddr: tcp:1.2.3.4:80\nudaddr: udp:4.2.2.4:53\nuaddr: unix:/tmp/sock\naddr: tcp:1.2.3.4:443"

	if tmpj := getTempFile(t, "json", jdata); tmpj != "" {
		defer os.Remove(tmpj)

		if err := LoadJSON(&dummyconf1, tmpj, true); err != nil {
			t.Errorf("Error loading dummy JSON config: %v", err)
		}

	} else {
		t.Errorf("Failed to get JSON tmpfile")
	}

	if _, err := json.Marshal(dummyconf1); err != nil {
		t.Errorf("Failed to marshal JSON config: %v", err)
	}

	if tmpy := getTempFile(t, "yml", ydata); tmpy != "" {
		defer os.Remove(tmpy)

		if err := LoadYAML(&dummyconf2, tmpy, true); err != nil {
			t.Errorf("Error loading dummy YAML config: %v", err)
		}
	} else {
		t.Errorf("Failed to get YAML tmpfile")
	}

	if _, err := yaml.Marshal(dummyconf2); err != nil {
		t.Errorf("Failed to marshal YAML config: %v", err)
	}
}

type TLSConf struct {
	Conf  TLS
	Cauth TLSClientAuth
}

func TestTLS(t *testing.T) {
	ustr := "BADAUTH"
	ta := &TLSClientAuth{}

	if err := ta.Set(ustr); err == nil {
		t.Errorf("Expected Bad TLS client auth type to produce error")
	}

	ustr = "require"
	if err := ta.Set(ustr); err != nil {
		t.Errorf("Error in setting TLSClientAuth: %v", err)
	}

	jdata := `{ "cauth": "require" }`
	ydata := "cauth: require"

	if tmpj := getTempFile(t, "json", jdata); tmpj != "" {
		var testconf TLSConf
		defer os.Remove(tmpj)

		if err := LoadJSON(&testconf, tmpj, true); err != nil {
			t.Errorf("Error loading dummy JSON config: %v", err)
		}

		if _, err := json.Marshal(testconf); err != nil {
			t.Errorf("TLSClientAuth JSON marshaling error: %v", err)
		}
	} else {
		t.Errorf("Failed to get JSON tmpfile")
	}

	if tmpy := getTempFile(t, "yml", ydata); tmpy != "" {
		var testconf TLSConf
		defer os.Remove(tmpy)

		if err := LoadYAML(&testconf, tmpy, true); err != nil {
			t.Errorf("Error loading dummy YAML config: %v", err)
		}

		if _, err := yaml.Marshal(testconf); err != nil {
			t.Errorf("TLSClientAuth YAML marshaling error: %v", err)
		}
	} else {
		t.Errorf("Failed to get YAML tmpfile")
	}

	jdata = `{ "conf": { "rootCAFiles": [ "test/root_ca_file.pem" ], "clientCAfiles": [ "test/client.crt" ], "certificates": [ { "certFile": "test/client.crt", "keyFile": "test/client.key" } ] } }`
	ydata = "---\nconf:\n  rootCAFiles: []\n  clientCAfiles: []\n\n"

	if tmpj := getTempFile(t, "json", jdata); tmpj != "" {
		var testconf TLSConf
		defer os.Remove(tmpj)

		if err := LoadJSON(&testconf, tmpj, true); err != nil {
			t.Errorf("Error loading dummy JSON config: %v", err)
		}

		if _, err := json.Marshal(testconf); err != nil {
			t.Errorf("TLS JSON marshaling error: %v", err)
		}
	} else {
		t.Errorf("Failed to get JSON tmpfile")
	}

	if tmpy := getTempFile(t, "yml", ydata); tmpy != "" {
		var testconf TLSConf
		defer os.Remove(tmpy)

		if err := LoadYAML(&testconf, tmpy, true); err != nil {
			t.Errorf("Error loading dummy YAML config: %v", err)
		}

		if _, err := yaml.Marshal(testconf); err != nil {
			t.Errorf("TLS YAML marshaling error: %v", err)
		}
	} else {
		t.Errorf("Failed to get YAML tmpfile")
	}

	badconfs := []string{
		`{ "conf": { "rootCAFiles": [ "test/DOES_NOT_EXIST.pem" ], "clientCAfiles": [ "test/client.crt" ], "certificates": [ { "certFile": "test/client.crt", "keyFile": "test/client.key" } ] } }`,
		`{ "conf": { "rootCAFiles": [ "test/root_ca_file.pem" ], "clientCAfiles": [ "test/DOES_NOT_EXIST.crt" ], "certificates": [ { "certFile": "test/client.crt", "keyFile": "test/client.key" } ] } }`,
		`{ "conf": { "rootCAFiles": [ "test/root_ca_file.pem" ], "clientCAfiles": [ "test/client.crt" ], "certificates": [ { "certFile": "test/DOES_NOT_EXIST.crt", "keyFile": "test/client.key" } ] } }`,
	}

	for _, i := range badconfs {
		if tmpj := getTempFile(t, "json", i); tmpj != "" {
			var testconf TLSConf
			defer os.Remove(tmpj)

			if err := LoadJSON(&testconf, tmpj, true); err == nil {
				t.Errorf("Expected TLS conf with missing file to produce error but it did not")
			}

		} else {
			t.Errorf("Failed to get JSON tmpfile")
		}
	}
}

func TestUnixAddr(t *testing.T) {
	u := &UnixAddr{}
	if err := u.Set("unix:/tmp/sock"); err != nil {
		t.Errorf("Error setting UnixAddr: %v", err)
	}

	mb, err := u.MarshalJSON()
	if err != nil {
		t.Errorf("UnixAddr MarshalJSON error: %v", err)
	}
	if string(mb) != `"unix:/tmp/sock"` {
		t.Errorf("UnixAddr MarshalJSON unexpected value: %s", string(mb))
	}

	my, err := u.MarshalYAML()
	if err != nil {
		t.Errorf("UnixAddr MarshalYAML error: %v", err)
	}
	if my != "unix:/tmp/sock" {
		t.Errorf("UnixAddr MarshalYAML unexpected value: %v", my)
	}

	u2 := &UnixAddr{}
	if err := u2.UnmarshalJSON(mb); err != nil {
		t.Errorf("UnixAddr UnmarshalJSON error: %v", err)
	}
	if u2.String() != u.String() {
		t.Errorf("UnixAddr UnmarshalJSON value mismatch: %s", u2.String())
	}

	if err := u2.UnmarshalJSON([]byte("BAD, JSON")); err == nil {
		t.Errorf("Expected UnixAddr UnmarshalJSON bad input to produce error")
	}
}

func TestStringSourcePreserved(t *testing.T) {
	s := &String{}
	envname := "TEST_STRING_SOURCE"
	if err := os.Setenv(envname, "envvalue"); err != nil {
		t.Errorf("os.Setenv error: %v", err)
		return
	}
	if err := s.Set("$" + envname); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if s.String() != "envvalue" {
		t.Errorf("Expected expanded value, got: %s", s.String())
	}
	mb, err := s.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON error: %v", err)
	}
	if string(mb) != fmt.Sprintf(`"$%s"`, envname) {
		t.Errorf("MarshalJSON should preserve source form, got: %s", string(mb))
	}
}

func TestStringRelativePath(t *testing.T) {
	tmpdata := "relative path content"
	tmpfile, err := os.CreateTemp(".", "tmpdata.*.txt")
	if err != nil {
		t.Errorf("Error creating temp file: %v", err)
		return
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	if err := os.WriteFile(tmpfile.Name(), []byte(tmpdata), 0600); err != nil {
		t.Errorf("Error writing temp file: %v", err)
		return
	}

	s := &String{}
	if err := s.Set("./" + tmpfile.Name()); err != nil {
		t.Errorf("Error loading string from relative path: %v", err)
	} else if s.String() != tmpdata {
		t.Errorf("Unexpected string value from file: %s", s.String())
	}
}

func TestTLSClientAuthInvalidMarshal(t *testing.T) {
	ta := TLSClientAuth{tls.ClientAuthType(99)}
	if _, err := ta.MarshalJSON(); err == nil {
		t.Errorf("Expected MarshalJSON of invalid TLSClientAuth to produce error")
	}
	if _, err := ta.MarshalYAML(); err == nil {
		t.Errorf("Expected MarshalYAML of invalid TLSClientAuth to produce error")
	}
}

func TestFileMalformed(t *testing.T) {
	var dummy interface{}

	if tmpj := getTempFile(t, "json", "NOT { valid JSON"); tmpj != "" {
		defer os.Remove(tmpj)
		if err := LoadJSON(&dummy, tmpj, true); err == nil {
			t.Errorf("Expected malformed JSON file to produce error")
		}
	} else {
		t.Errorf("Failed to get JSON tmpfile")
	}

	if tmpy := getTempFile(t, "yml", "key: [unclosed"); tmpy != "" {
		defer os.Remove(tmpy)
		if err := LoadYAML(&dummy, tmpy, true); err == nil {
			t.Errorf("Expected malformed YAML file to produce error")
		}
	} else {
		t.Errorf("Failed to get YAML tmpfile")
	}
}
