# gopotato

a Go rewrite of Potato, a minimalist journal manager I originally wrote in Deno.

## Installation

```sh
go install github.com/n3dst4/gopotato@latest
```

## Usage

```
gopotato
```

## Running tests

```sh
go test ./...
```

## Check formatting

```sh
gofmt -l .
```

And to fix files in-place:

```sh
gofmt -w .
```

## Running linter

First, you need to have golangci-lint installed locally. They recommend [binary installetion](https://golangci-lint.run/docs/welcome/install/local/#binaries). Check link for current instructions. Then run:

```sh
golangci-lint run
```

## Configuration

Configuration is stored in a TOML file.
