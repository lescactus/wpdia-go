# wpdia-go

This repository contains a simple cli written in go used to get the description of a given text in Wikipedia.

It takes in argument a given text and will retrieve the extract of page content using the TextExtracts API (https://www.mediawiki.org/wiki/Extension:TextExtracts#API).


`wpdia-go` allow to either return the content from Wikipedia before the first section (typically the text block before the table of contents): `exintro` or a given number of sentences between 1 and 10: `exsentences`.

Note that the [`TextExtracts` API](https://www.mediawiki.org/wiki/Extension:TextExtracts#API) recommends not to use `exsentences` as it does not work for HTML extracts and there are many edge cases for which it doesn't exist. For example "Arm. gen. Ing. John Smith was a soldier." will be treated as 4 sentences.

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
  -i, --exintro              Return only content before the first section. Mutually exclusive with 'exsentences'. (default true)
  -s, --exsentences string   How many sentences to return from Wikipedia. Must be between 1 and 10. If > 10, then default to 10. Mutually exclusive with 'exintro'. (default "10")
  -f, --full                 Also print the page Namespace and page ID.
  -h, --help                 help for wpdia-go
  -l, --lang string          Language. This will set the API endpoint used to retrieve data. (default "en")
  -o, --output string        Output type. Valid choices are [plain pretty json yaml]. (default "plain")
  -t, --timeout duration     Timeout value of the http client to the Wikipedia API. Examples values: '10s', '500ms' (default 15s)
  -v, --version              version for wpdia-go
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

## Examples

### Basic usage:
```
./wpdia-go golang
Title:
  Go (programming language)

Extract:
  Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency. It is often referred to as Golang because of its former domain name, golang.org, but its proper name is Go.There are two major implementations:

Google's self-hosting "gc" compiler toolchain, targeting multiple operating systems and WebAssembly.
gofrontend, a frontend to other compilers, with the libgo library. With GCC the combination is gccgo; with LLVM the combination is gollvm.A third-party source-to-source compiler, GopherJS, compiles Go to JavaScript for front-end web development.
```

### Change language

```
./wpdia-go --lang fr golang
Title:
  Go (langage)

Extract:
  Go est un langage de programmation compilé et concurrent inspiré de C et Pascal. Ce langage a été développé par Google à partir d’un concept initial de Robert Griesemer, Rob Pike et Ken Thompson. Go possède deux implémentations : la première utilise gc, le compilateur Go ; la seconde utilise gccgo, « frontend » GCC écrit en C++. Go est écrit en C en utilisant yacc et GNU Bison pour l’analyse syntaxique jusqu’à la version 1.4, et en Go lui-même pour les versions suivantes (1.5).

Un objectif de Go est donné par Rob Pike, l’un de ses trois créateurs, qui dit à propos des développeurs inexpérimentés :

« Ils ne sont pas capables de comprendre un langage brillant, mais nous voulons les amener à réaliser de bons programmes. Ainsi, le langage que nous leur donnons doit être facile à comprendre et facile à adopter »

Go veut faciliter et accélérer la programmation à grande échelle : en raison de sa simplicité, il est donc concevable de l’utiliser aussi bien pour écrire des applications, des scripts ou de grands systèmes. Cette simplicité est nécessaire aussi pour assurer la maintenance et l’évolution des programmes sur plusieurs générations de développeurs.
S’il vise aussi la rapidité d’exécution, indispensable à la programmation système, il considère le multithreading comme le moyen le plus robuste d’assurer sur les processeurs actuels cette rapidité tout en rendant la maintenance facile par séparation de tâches simples exécutées indépendamment afin d’éviter de créer des « usines à gaz ». Cette conception permet également le fonctionnement sans réécriture sur des architectures multi-cœurs en exploitant immédiatement l’augmentation de puissance correspondante.
```

```
./wpdia-go --lang it golang
Title:
  Go (linguaggio di programmazione)

Extract:
  Go è un linguaggio di programmazione open source sviluppato da Google.
Il lavoro su Go nacque nel settembre 2007 da Robert Griesemer, Rob Pike e Ken Thompson basandosi su un precedente lavoro correlato con il sistema operativo Inferno.
Secondo gli autori, l'esigenza di creare un nuovo linguaggio di programmazione nasce dal fatto che non esiste un linguaggio di programmazione che soddisfi le esigenze di una compilazione efficiente, di un'esecuzione veloce e di una facilità di programmazione.
Go viene annunciato ufficialmente nel novembre 2009.
```

```
./wpdia-go --lang es golang
Title:
  Go (lenguaje de programación)

Extract:
  Go es un lenguaje de programación concurrente y compilado inspirado en la sintaxis de C, que intenta ser dinámico como Python y con el rendimiento de C o C++. Ha sido desarrollado por Google[9]​ y sus diseñadores iniciales fueron Robert Griesemer, Rob Pike y Ken Thompson. [10]​ Actualmente está disponible en formato binario para los sistemas operativos Windows, GNU/Linux, FreeBSD  y Mac OS X, pudiendo también ser instalado en estos y en otros sistemas mediante el código fuente.[11]​[12]​ Go es un lenguaje de programación compilado, concurrente, imperativo, estructurado, orientado a objetos y con recolector de basura que de momento es soportado en diferentes tipos de sistemas UNIX, incluidos Linux, FreeBSD, Mac OS X y Plan 9 (puesto que parte del compilador está basado en un trabajo previo sobre el sistema operativo Inferno). Las arquitecturas soportadas son i386, amd64 y ARM.
```

### Return only the first 2 sentences

```
./wpdia-go -s 2 golang
Title:
  Go (programming language)

Extract:
  Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency.
```

### Pretty output

```
./wpdia-go --output pretty golang

  ## Go (programming language)                                                                    


  Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, 
  Rob Pike, and Ken Thompson. It is syntactically similar to C, but with memory safety, garbage   
  collection, structural typing, and CSP-style concurrency. It is often referred to as Golang     
  because of its former domain name, golang.org, but its proper name is Go.There are two major    
  implementations:                                                                                
                                                                                                  
  Google's self-hosting "gc" compiler toolchain, targeting multiple operating systems and         
  WebAssembly. gofrontend, a frontend to other compilers, with the libgo library. With GCC the    
  combination is gccgo; with LLVM the combination is gollvm.A third-party source-to-source compiler,
  GopherJS, compiles Go to JavaScript for front-end web development.
```

### Json output

```
./wpdia-go --output json golang  
{
    "title": "Go (programming language)",
    "extract": "Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency. It is often referred to as Golang because of its former domain name, golang.org, but its proper name is Go.There are two major implementations:\n\nGoogle's self-hosting \"gc\" compiler toolchain, targeting multiple operating systems and WebAssembly.\ngofrontend, a frontend to other compilers, with the libgo library. With GCC the combination is gccgo; with LLVM the combination is gollvm.A third-party source-to-source compiler, GopherJS, compiles Go to JavaScript for front-end web development."
}
```

### Yaml output

```
title: Go (programming language)
extract: |-
  Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency. It is often referred to as Golang because of its former domain name, golang.org, but its proper name is Go.There are two major implementations:

  Google's self-hosting "gc" compiler toolchain, targeting multiple operating systems and WebAssembly.
  gofrontend, a frontend to other compilers, with the libgo library. With GCC the combination is gccgo; with LLVM the combination is gollvm.A third-party source-to-source compiler, GopherJS, compiles Go to JavaScript for front-end web development.
```

### HTTP client timeout set to 3 seconds 

```
./wpdia-go --timeout 3s golang    
Title:
  Go (programming language)

Extract:
  Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency. It is often referred to as Golang because of its former domain name, golang.org, but its proper name is Go.There are two major implementations:

Google's self-hosting "gc" compiler toolchain, targeting multiple operating systems and WebAssembly.
gofrontend, a frontend to other compilers, with the libgo library. With GCC the combination is gccgo; with LLVM the combination is gollvm.A third-party source-to-source compiler, GopherJS, compiles Go to JavaScript for front-end web development.
```

### Output the page namespace and page id
```
Title:
  Go (programming language)

Ns:
  0

Pageid:
  25039021

Extract:
  Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency. It is often referred to as Golang because of its former domain name, golang.org, but its proper name is Go.There are two major implementations:

Google's self-hosting "gc" compiler toolchain, targeting multiple operating systems and WebAssembly.
gofrontend, a frontend to other compilers, with the libgo library. With GCC the combination is gccgo; with LLVM the combination is gollvm.A third-party source-to-source compiler, GopherJS, compiles Go to JavaScript for front-end web development.
```

### HTTP client timeout set to 500ms + json output + only 3 sentences + French language + full output

```
./wpdia-go --timeout 500ms --output json --exsentences 3 --lang fr --full golang
{
    "pageid": 4230561,
    "ns": 0,
    "title": "Go (langage)",
    "extract": "Go est un langage de programmation compilé et concurrent inspiré de C et Pascal. Ce langage a été développé par Google à partir d’un concept initial de Robert Griesemer, Rob Pike et Ken Thompson. Go possède deux implémentations : la première utilise gc, le compilateur Go ; la seconde utilise gccgo, « frontend » GCC écrit en C++."
}
```

---
**TODO:**

- [x] Improve display

- [x] Improve documentation

- [x] Improve http user agent 

- [x] Avoid code duplicate in http request builder

- [x] Parametrize `exsentences`, http timeout, etc ... (flag & env variable)

- [x] Output flag: `table`, `json`, etc...

- [x] Language support

- [ ] Implement "random article" 

- [ ] Fix 'may refer to:'

- [x] Improve base url

- [x] Dockerize

- [x] Debug flag (show page id, ns, timestamps, etc...)

- [ ] Verbose logs