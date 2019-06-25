# Aeridya

Single domain webserver/CMS written and extendable using Golang

## Description

Aeridya extends the built-in HTTP functionality of Golang to deliver Web Pages where the logic is written in Golang.  The final render of the webpage uses Golang's Templating System to deliver static pages.  This provides you with the flexibility of HTML/CSS/JavaScript in an easy package, and the speed to do server side logic in Go.  The final application is recommended to run via a reverse proxy; specifically NGINX.  More documentation on this will be written once the application is further along.

Aeridya applications should consist of a Theme which routes incomming connections to the functions that it containts.  A Theme is recommended to contain Pages, which a Page consists of a set of instructions for each of the main HTTP Requests, ie. "GET", "PUT", "POST", "DELETE", "OPTIONS", "HEAD".

See "basic" for an example application.  You may also want to see "theme.go" on how the basic theme is implemented.

## Using Aeridya

To use Aeridya in your application, you must have a configuration file setup.  "See basic/conf" for the necessary basic configuration.

The following must be set in the configuration file:
```
## Configuration for basic site using Aeridya
## A "#" denotes a comment in the configuration
## A ";" denotes the beginning of a section
## The Aeridya Section is required by the Aeridya Core
;Aeridya

# Log location 
## NOTE: use stdout to get term output
Log = stdout

# Port to listen on
Port = 8000

# Domain name
Domain = domain.com

# HTTPS Setting.
## NOTE:  Aeridya itself does not handle HTTPS, and HTTPS 
## is only availble when used via a reverse proxy
## This setting will set the internal URL for routing
HTTPS = false

# Development Mode
## NOTE:  In production, it is recommended to disable
Development = true

# Workers
## Sets allowed connections able to do work 
Workers = 10
```