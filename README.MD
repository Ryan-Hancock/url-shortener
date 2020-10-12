# URL Shortener

## Table of contents

- [General info](#general-info)
- [Technologies](#technologies)
- [Setup](#setup)
  - [Tests](#tests)
  - [Usage](#usage)

## General info

This is an API for a URL shortening service written in Golang (Go).

## Technologies

These are the technologies used to create this project:

- Go 1.14.7
- [Mux](https://github.com/gorilla/mux)
- [SQLite](https://github.com/mattn/go-sqlite3)
- [sqlx](https://github.com/jmoiron/sqlx)
- [cli](github.com/urfave/cli)

## Setup

To run this project, build it locally using Go:

```
$ go build -o url-shortner
```

Once built you can start the server with the CLI command:

```
$ ./url-shortener s --address 127.0.0.1:8000
```

### Tests

In order to run tests:

```
$ go test -v
```

### Usage

The following endpoints can be called to generate a short URL.

```
POST /url HTTP/1.1
Host: 127.0.0.1:8000
Content-Type: application/json

{
    "url": "http://google.com"
}

Returns
{
    "message": href("/somecode")
}
```

```
GET /somecode HTTP/1.1
Host: 127.0.0.1:8000

Returns <a href="/www.google.com">>Moved Permanently</a>.
```
