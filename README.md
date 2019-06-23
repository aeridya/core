# Aeridya

Single domain webserver/CMS written and extendable using Golang

## Description

Aeridya extends the built-in HTTP functionality of Golang to deliver Web Pages where the logic is written in Golang.  The final render of the webpage uses Golang's Templating System to deliver static pages.  This provides you with the flexibility of HTML/CSS/JavaScript in an easy package, and the speed to do server side logic in Go.  The final application is recommended to run via a reverse proxy; specifically NGINX.  More documentation on this will be written once the application is further along.

Each Aeridya application consists of a Theme, and Pages.  A Theme is called directly by Aeridya and decides how an application's Pages are stored/accessed.  A Page consists of a set of instructions for each of the main HTTP Requests, ie. "GET", "PUT", "POST", "DELETE", "OPTIONS", "HEAD".

See "example" for the example application.  You may also want to see "theme.go" on how the basic theme is implemented.

## Using Aeridya

To use Aeridya in your application, you must have a configuration file setup.  "See example/conf" for the necessary basic configuration.

The following must be set in the configuration file:
```
Port = Port to start on (example: 5000)
Domain = Application's Domain: (example: domain.com)
Workers = Max Workers [you can be pretty generous here] (example: 1000)
Development = Deveopment Mode (true|false)
Log = File Path to log (Use stdout to print to terminal)
Statics = File Path to Statics Directory
```

Aeridya is meant to be easy to use in your own applications.  To get started, the following code is required:

```go
package main

import (
	"fmt"
	"github.com/hlfstr/aeridya"
	"os"
)

func main() {
	if e := aeridya.Create("/path/to/config"); e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	aeridya.Run()
}
```

### Customizing Aeridya

TODO
