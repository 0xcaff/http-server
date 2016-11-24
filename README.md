http-server
==============

This is a simple, fast and portable HTTP server built in go.

Performance
-----------

This HTTP server is about 7x faster than [node-http-server](//github.com/nodeapps/http-server)
and about 10x faster than `python -m SimpleHTTPServer`.

Installation
------------

    go get github.com/caffinatedmonkey/http-server
    http-server

Usage
-----

Run `http-server` from the directory you would like to serve.

### Options:
 - `-listen`: the address to listen on for example `0.0.0.0` or the default
[`:8080`]
 - `-redirect` Respond to unknown requests with this file (relative to `-path`). Can be used to serve Single Page Applications.
 - `-path`: the path to serve [`./`]
 - `-header`: header sent with every response, for example `X-Authentication:
global`
 - `-help`: displays the help message

