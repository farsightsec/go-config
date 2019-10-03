package config

import (
	"testing"
	"fmt"
	"os"
	"io/ioutil"

	"encoding/json"
	yaml "gopkg.in/yaml.v2"
)


func getTempFile(t *testing.T, ext string, content string) string {

	tmpfile, err := ioutil.TempFile("", "tmpdata." + ext)
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


func testStuff(t *testing.T, foo interface{}, teststr string, success bool) {
	var mb []byte
	var mb2 interface{}
	var err error = nil
	var err2 error = nil

	switch foo.(type) {
		case *TCPAddr:
			a := foo.(*TCPAddr)
			mb, err = a.MarshalJSON()
			mb2, err2 = a.MarshalYAML()
		case *UDPAddr:
			a := foo.(*UDPAddr)
			mb, err = a.MarshalJSON()
			mb2, err2 = a.MarshalYAML()
		case *Addr:
			a := foo.(*Addr)
			mb, err = a.MarshalJSON()
			mb2, err2 = a.MarshalYAML()
		case *Duration:
			d := foo.(*Duration)
			mb, err = d.MarshalJSON()
			mb2, err2 = d.MarshalYAML()
		case *String:
			s := foo.(*String)
			mb, err = s.MarshalJSON()
			mb2, err2 = s.MarshalYAML()
		case *URL:
			u := foo.(*URL)
			mb, err = u.MarshalJSON()
			mb2, err2 = u.MarshalYAML()
		default:
			t.Errorf("Test function received unexpected type: %T", foo)
			return
	}

	jstr := fmt.Sprintf("\"%s\"", teststr)

	if !success {
		if (err == nil && string(mb) == jstr) {
			t.Errorf("Unexpected successful marshaling of bad %T in JSON form", foo)
		}
	} else if err != nil {
		t.Errorf("Error marshaling %T in JSON form: %v", foo, err)
	} else if string(mb) != jstr {
		t.Errorf("Marshaled JSON %T was not expected value: %v", foo, string(mb))
	}

	if !success {
		if err2 == nil {
			if ystr, ok := mb2.(string); (ok != true || teststr == ystr) {
				t.Errorf("Unexpected successful marshaling of bad %T in YAML form", foo)
			}
		}
	} else if err2 != nil {
		t.Errorf("Error marshaling %T in YAML form: %v", foo, err2)
	} else if ystr, ok := mb2.(string); (ok != true || teststr != ystr) {
		t.Errorf("Marshaled YAML %T was not expected value: %v", foo, ystr)
	}

	if success {
		var uerr error = nil
		var uerr2 error = nil

		badjson := []byte("BAD, JSON, STRING")

		switch foo.(type) {
			case *TCPAddr:
			case *UDPAddr:
			case *Addr:
				a := &Addr{}
				uerr = a.UnmarshalJSON(mb)
				uerr2 = a.UnmarshalJSON(badjson)
			case *Duration:
				d := &Duration{}
				uerr = d.UnmarshalJSON(mb)
				uerr2 = d.UnmarshalJSON(badjson)
			case *String:
				s := &String{}
				uerr = s.UnmarshalJSON(mb)
				uerr2 = s.UnmarshalJSON(badjson)
			case *URL:
				u := &URL{}
				uerr = u.UnmarshalJSON(mb)
				uerr2 = u.UnmarshalJSON(badjson)
		}

		if uerr != nil {
			t.Errorf("Unmarshaling of JSON data into type %T produced error: %v", foo, uerr)
		}

		if uerr2 == nil {
			t.Errorf("Unmarshaling of bad JSON data into type %T failed to produce error", foo)
		}

	}

}

func TestBadAddr(t *testing.T) {
	ustr := "udp:1.2.3.4"
	a := &Addr{}
	a.Set(ustr)
	testStuff(t, a, ustr, true)

	u := &UDPAddr{}
	ta := &TCPAddr{}
	ua := &UnixAddr{}

	for i, _ := range []int{0, 1, 2, 3} {
		var err, err2 error = nil, nil

		switch(i) {
			case 0:
				err = a.Set("BADADDR")
				err2 = err
			case 1:
				err = u.Set("tcp:1.2.3.4")
				err2 = u.Set("BADADDR")
			case 2:
				err = ta.Set("udp:1.2.3.4")
				err2 = ta.Set("BADADDR")
			case 3:
				err = ua.Set("udp:1.2.3.4")
				err2 = ua.Set("BADADDR")
		}

		if err == nil {
			t.Errorf("Unexpected error-free setting of address from bad data / %d", i)
		} else {
			fmt.Printf("Bad address produced error: %T / %v\n", err, err)
		}

		if err2 == nil {
			t.Errorf("Expected set of bad address to produce error but it did not")
		}

	}

}

func TestDuration(t *testing.T) {
	ustr := "100ms"
	d := &Duration{}
	d.Set(ustr)
	testStuff(t, d, ustr, true)

	ustr = "1m30duif"
	d.Set(ustr)
	testStuff(t, d, ustr, false)
}

func TestString(t *testing.T) {
	ustr := "bla bla bla"
	s := &String{}
	s.Set(ustr)

	if s.String() != ustr {
		t.Errorf("Unexpected string set error")
	}

	testStuff(t, s, ustr, true)

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

	testStuff(t, u, ustr, true)
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
	Str1		String
	EnvString	String
	Num1		uint64
	Url		URL
	Dur		Duration
	Addr		Addr
	UAddr		UnixAddr
	TAddr		TCPAddr
	UDAddr		UDPAddr
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

	b, e := json.Marshal(dummyconf1)
	fmt.Printf("Its: %v, %v\n", string(b), e)

	if tmpy := getTempFile(t, "yml", ydata); tmpy != "" {
		defer os.Remove(tmpy)

		if err := LoadYAML(&dummyconf2, tmpy, true); err != nil {
			t.Errorf("Error loading dummy YAML config: %v", err)
		}

	} else {
		t.Errorf("Failed to get YAML tmpfile")
	}

	b, e = yaml.Marshal(dummyconf2)
}

type TLSConf struct {
	Conf		TLS
	Cauth		TLSClientAuth
}

func TestTLS(t *testing.T) {
	ustr := "BADAUTH"
	ta := &TLSClientAuth{}

	if err := ta.Set(ustr); err == nil {
		t.Errorf("Expected Bad TLS client auth type to produce error")
	} else {
		fmt.Printf("Bad TLS client auth produced error: %v\n", err)
	}

	ustr = "require"
	if err := ta.Set(ustr); err != nil {
		t.Errorf("Error in setting TLSClientAuth: %v", err)
	}

	fmt.Printf("Client auth val: %s\n", ta.String())

	jdata := `{ "cauth": "require" }`
	ydata := "cauth: require"

	if tmpj := getTempFile(t, "json", jdata); tmpj != "" {
		var testconf TLSConf
		defer os.Remove(tmpj)

		if err := LoadJSON(&testconf, tmpj, true); err != nil {
			t.Errorf("Error loading dummy JSON config: %v", err)
		}

		b, e := json.Marshal(testconf)
		fmt.Printf("TLSClientAuth JSON marshaling: %v, %v\n", string(b), e)
	} else {
		t.Errorf("Failed to get JSON tmpfile")
	}

	if tmpy := getTempFile(t, "yml", ydata); tmpy != "" {
		var testconf TLSConf
		defer os.Remove(tmpy)

		if err := LoadYAML(&testconf, tmpy, true); err != nil {
			t.Errorf("Error loading dummy YAML config: %v", err)
		}

		b, e := yaml.Marshal(testconf)
		fmt.Printf("TLSClientAuth YAML marshaling: %v, %v\n", string(b), e)
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

		b, e := json.Marshal(testconf)
		fmt.Printf("TLS JSON marshaling: %v, %v\n", string(b), e)
	} else {
		t.Errorf("Failed to get JSON tmpfile")
	}

	if tmpy := getTempFile(t, "yml", ydata); tmpy != "" {
		var testconf TLSConf
		defer os.Remove(tmpy)

		if err := LoadYAML(&testconf, tmpy, true); err != nil {
			t.Errorf("Error loading dummy YAML config: %v", err)
		}

		b, e := yaml.Marshal(testconf)
		fmt.Printf("TLS YAML marshaling: %v, %v\n", string(b), e)
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
