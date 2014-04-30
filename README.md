http-server.go
==============

This is a simple, fast and portable HTTP server built in go.

Performance
-----------

This HTTP server is about 7x faster than [node-http-server](//github.com/nodeapps/http-server)
and about 10x faster than `python -m SimpleHTTPServer`.

Installation
------------

    go get github.com/caffinatedmonkey/http-server.go/http-server

Usage
-----

Run `http-server` from the directory you would like to serve.

### Options:
 - `-port`: the port to serve on [`8080`]
 - `-path`: the path to serve [`./`]
 - `-help`: displays the help message

