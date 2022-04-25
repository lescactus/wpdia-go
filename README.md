# wpdia-go

This repository contains a simple cli written in go used to get the description of a given text in Wikipedia.

It takes in argument a given text and will retrieve the extract of page content using the 
TextExtracts API (https://www.mediawiki.org/wiki/Extension:TextExtracts#API).

## Usage

```
wpdia-go is a simple cli used to get the description of a given text in Wikipedia.
It takes in argument a given text and will retrieve the extract of page content using the 
TextExtracts API (https://www.mediawiki.org/wiki/Extension:TextExtracts#API).

For multi-word search, enclose them using double quotes: "<multi word search>".


The source code is available at https://github.com/lescactus/wpedia-go.

Usage:
  wpdia-go [flags]

Flags:
  -h, --help          help for wpdia-go
  -l, --lang string   Language. This will set the API endpoint used to retrieve data. (default "en")
```
## Installation

### From source with go

You need a working [go](https://golang.org/doc/install) toolchain (It has been developped and tested with go 1.16 and go 1.16 only, but should work with go >= 1.14). Refer to the official documentation for more information (or from your Linux/Mac/Windows distribution documentation to install it from your favorite package manager).

```sh
# Clone this repository
git clone https://github.com/lescactus/wpdia-go.git && cd wpdia-go/

# Build from sources. Use the '-o' flag to change the compiled binary name
go build

# Default compiled binary is wpdia-go
# You can optionnaly move it somewhere in your $PATH to access it shell wide
./wpdia-go -h
```

### From source with docker

If you don't have [go](https://golang.org/) installed but have docker, run the following command to build inside a docker container:

```sh
# Build from sources inside a docker container. Use the '-o' flag to change the compiled binary name
# Warning: the compiled binary belongs to root:root
docker run --rm -it -v "$PWD":/app -w /app golang:1.16 go build

# Default compiled binary is dict-go
# You can optionnaly move it somewhere in your $PATH to access it shell wide
./wpdia-go -h
```

### From source with docker but built inside a docker image

If you don't want to pollute your computer with another program, this cli comes with its own docker image:

```sh
docker build -t wpdia-go .

docker run --rm wpdia-go "Rammstein"
```

---
**TODO:**

- [ ] Improve display

- [x] Improve documentation

- [x] Improve http user agent 

- [x] Avoid code duplicate in http request builder

- [ ] Parametrize `exsentences`, http timeout, etc ... (flag & env variable)

- [ ] Output flag: `table`, `json`, etc...

- [x] Language support

- [ ] Implement "random article" 

- [ ] Fix 'may refer to:'

- [x] Improve base url

- [x] Dockerize