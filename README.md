# Minimalist Go Config Library

The `config` package contains types useful for validating, parsing, and loading values of some useful types in config files. The `config` types generally embed a corresponding standard library type, and provide Unmarshaler and Marshaler support for loading configs from and saving configs to JSON and YAML formats.

The types in this package also implement `flag.Value` for use in command line arguments.

Types include:
  * `net/url.URL`
  * `time.Duration`
  * `net.{UDP,TCP,Unix}Addr`
  * `crypto/tls.Config`

All of these Unmarshal from and Marshal to their natural string representations, with the exception of `tls.Config`, which is represented in the configuration as a dictionary of filenames and other settings.

An additional utility `String` type holds a `string` value which can optionally be read from the environment or from a named file.

## Example

```go
import (
        "encoding/json"
        "github.com/farsightsec/go-config"
)

type Config struct {
        Server   config.URL
        Timeout  config.Duration
        Username config.String
        Key      config.String
}

// cfg.Server.URL = net/url.URL of Server.
// cfg.Timeout.Duration = time.Duration of Timeout.
var cfg Config

var serverURL config.URL

func init() {
        configText := `
        {
                "Server": "https://example.com:8080/api/",
                "Timeout": "90s",
                "Username": "$USER",
                "Key": "/etc/app/apikey"
        }`
        if err := json.Unmarshal([]byte(configText)); err != nil {
                panic(err)
        }

        flag.Var(&serverURL, "server", "URL for server")
}
```
