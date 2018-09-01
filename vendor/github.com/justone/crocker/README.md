[![GoDoc](https://godoc.org/github.com/justone/crocker?status.svg)](http://godoc.org/github.com/justone/crocker)

# crocker - Credential Helper Helper

This is a small helper library for integrating with the [Docker Credential
Helpers](https://github.com/docker/docker-credential-helpers/).  It handles
detecting if the right helper is installed and errors out if helpers aren't
found.

Example code:

```
import "github.com/justone/crocker"

func main() {
    c, err := crocker.New()

    err := c.Store(...)
    creds, err := c.Get(...)
    list, err := c.list(...)
    err := c.Erase(...)
}
```
