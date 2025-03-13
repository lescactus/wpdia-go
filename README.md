# wpdia-go [![Go CI](https://github.com/lescactus/wpdia-go/actions/workflows/go.yml/badge.svg)](https://github.com/lescactus/wpdia-go/actions/workflows/go.yml) [![goreleaser](https://github.com/lescactus/wpdia-go/actions/workflows/release.yml/badge.svg)](https://github.com/lescactus/wpdia-go/actions/workflows/release.yml)

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
  -a, --logformat string     Log format. Accepted values are [text json]. (default "text")
  -e, --loglevel string      Log level verbosity. Accepted values are [debug info warn error]. (default "error")
  -o, --output string        Output type. Valid choices are [plain pretty json yaml]. (default "plain")
  -r, --random               Return a random article.
  -t, --timeout duration     Timeout value of the http client to the Wikipedia API. Examples values: '10s', '500ms' (default 15s)
  -v, --version              version for wpdia-go
```
## Installation

Prebuilt binaries can be downloaded from the GitHub Releases [section](https://github.com/lescactus/wpdia-go/releases), or using a Docker image from the Github Container Registry.

### Running with Docker

```bash
docker run --rm -it --name wpdia-go ghcr.io/lescactus/wpdia-go
```

## Building

<details>

### Requirements

* Golang 1.23 or higher

### From source with go

You need a working [go](https://golang.org/doc/install) toolchain (It has been developped and tested with go >= 1.23). Refer to the official documentation for more information (or from your Linux/Mac/Windows distribution documentation to install it from your favorite package manager).

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
docker run --rm -it -v "$PWD":/app -w /app golang:1.23 go build -buildvcs=false

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

</details>

## Disambiguation pages

Sometimes the resulted page coming from Wikipedia's search is a disambiguation.
To quote [Wikipedia](https://en.wikipedia.org/wiki/Wikipedia:Disambiguation):
> Disambiguation in Wikipedia is the process of resolving conflicts that arise when a potential article title is ambiguous, most often because it refers to more than one subject covered by Wikipedia, either as the main topic of an article, or as a subtopic covered by an article in addition to the article's main topic. For example, Mercury can refer to a chemical element, a planet, a Roman god, and many other things.

> Disambiguation is required whenever, for a given word or phrase on which a reader might search, there is more than one existing Wikipedia article to which that word or phrase might be expected to lead. In this situation there must be a way for the reader to navigate quickly from the page that first appears to any of the other possible desired articles.

In this case, `wpdia-go` will print an error message asking the user to refine the query. Example:

```
./wpdia-go nancy
Title:
  Nancy

Extract:
  /!\ The requested page is a disambiguation page /!\

A disambiguation page is Wikipedia's way of resolving conflicts that arise when a potential article title is ambiguous - most often because it refers to more than one subject covered by Wikipedia.
For example, Mercury can refer to a chemical element, a planet, a Roman god, and many other things.

Try to refine the search in a more precise manner. Example:
	'Nancy France' instead of 'Nancy' - or 'Go verb' instead of 'Go'
```

When this happens, refining the query by beoing more precise will help.
For example, when looking for the description of the French city of Nancy, look for `Nancy France` instead of simply `Nancy`:

```
$ ./wpdia-go Nancy
Title:
  Nancy

Extract:
  /!\ The requested page is a disambiguation page /!\

A disambiguation page is Wikipedia's way of resolving conflicts that arise when a potential article title is ambiguous - most often because it refers to more than one subject covered by Wikipedia.
For example, Mercury can refer to a chemical element, a planet, a Roman god, and many other things.

Try to refine the search in a more precise manner. Example:
	'Nancy France' instead of 'Nancy' - or 'Go verb' instead of 'Go'


$ ./wpdia-go "Nancy france"
Title:
  Nancy, France

Extract:
  Nancy is the prefecture of the northeastern French department of Meurthe-et-Moselle. It was the capital of the Duchy of Lorraine which was annexed by France under King Louis XV in 1766 and replaced by a province with Nancy maintained as capital. Following its rise to prominence in the Age of Enlightenment, it was nicknamed the "capital of Eastern France" in the late 19th century. The metropolitan area of Nancy had a population of 511,257 inhabitants at the 2018 census, making it the 16th-largest urban area in France and Lorraine's largest. The population of the city of Nancy proper is 104,885.
The motto of the city is Non inultus premor, Latin for '"I am not injured unavenged"'—a reference to the thistle, which is a symbol of Lorraine. Place Stanislas, a large square built between 1752 and 1756 by architect Emmanuel Héré under the direction of Stanislaus I of Poland to link the medieval old town of Nancy and the new city built under Charles III, Duke of Lorraine in the 17th century, is now a UNESCO World Heritage Site, the first square in France to be given this distinction. The city also has many buildings listed as historical monuments and is one of the European centres of Art Nouveau thanks to the École de Nancy. Nancy is also a large university city; with the Centre Hospitalier Régional Universitaire de Brabois, the conurbation is home to one of the main health centres in Europe, renowned for its innovations in surgical robotics.
```

In the future, suggestions may be implemented.

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

### Output the page namespace, page id and page properties
```
./wpedia-go golang --full
Title:
  Go (programming language)

Ns:
  0

Pageid:
  25039021

WikiBase Short Description:
  Programming language

WikiBase Item:
  Q37227

Extract:
  Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency. It is often referred to as Golang because of its former domain name, golang.org, but its proper name is Go.There are two major implementations:

Google's self-hosting "gc" compiler toolchain, targeting multiple operating systems and WebAssembly.
gofrontend, a frontend to other compilers, with the libgo library. With GCC the combination is gccgo; with LLVM the combination is gollvm.A third-party source-to-source compiler, GopherJS, compiles Go to JavaScript for front-end web development.
```

### Info level logging

```
time=2025-01-21T11:29:03.826+01:00 level=INFO msg="Creating new Wiki client..." url=https://en.wikipedia.org/w/api.php
time=2025-01-21T11:29:03.826+01:00 level=INFO msg="Getting text extract..." title=golang random=false
time=2025-01-21T11:29:03.826+01:00 level=INFO msg="Searching title..." title=golang
time=2025-01-21T11:29:04.096+01:00 level=INFO msg="Search found a Page ID" pageid=25039021
Title:
  Go (programming language)

Extract:
  Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency. It is often referred to as Golang because of its former domain name, golang.org, but its proper name is Go.There are two major implementations:

Google's self-hosting "gc" compiler toolchain, targeting multiple operating systems and WebAssembly.
gofrontend, a frontend to other compilers, with the libgo library. With GCC the combination is gccgo; with LLVM the combination is gollvm.A third-party source-to-source compiler, GopherJS, compiles Go to JavaScript for front-end web development.
```

### Debug level logging

```
time=2025-01-21T11:30:35.175+01:00 level=INFO source=/home/amaldeme/gitclone/wpdia-go/cmd/root.go:83 msg="Creating new Wiki client..." url=https://en.wikipedia.org/w/api.php
time=2025-01-21T11:30:35.175+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:36 msg="Parsing base URL..." url=https://en.wikipedia.org/w/api.php
time=2025-01-21T11:30:35.175+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:47 msg="Base URL parsed" url=https://en.wikipedia.org/w/api.php
time=2025-01-21T11:30:35.175+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:62 msg="User-Agent set" user-agent="wpdia-go/0.4.1 (github.com/lescactus/wpdia-go) WikiClient/0.4.1"
time=2025-01-21T11:30:35.175+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/root.go:100 msg="New Wiki client created" url=https://en.wikipedia.org/w/api.php
time=2025-01-21T11:30:35.175+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/root.go:105 msg="Disabling 'exintro'..."
time=2025-01-21T11:30:35.175+01:00 level=INFO source=/home/amaldeme/gitclone/wpdia-go/cmd/root.go:118 msg="Getting text extract..." title=golang random=false
time=2025-01-21T11:30:35.175+01:00 level=INFO source=/home/amaldeme/gitclone/wpdia-go/cmd/root.go:129 msg="Searching title..." title=golang
time=2025-01-21T11:30:35.175+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:218 msg="Http request parameters set" params="map[list:[search] srlimit:[1] srsearch:[golang] utf8:[1]]"
time=2025-01-21T11:30:35.175+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:226 msg="Building http request..." params="map[list:[search] srlimit:[1] srsearch:[golang] utf8:[1]]" url=https://en.wikipedia.org/w/api.php user-agent="wpdia-go/0.4.1 (github.com/lescactus/wpdia-go) WikiClient/0.4.1"
time=2025-01-21T11:30:35.175+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:240 msg="Http request built" params="map[action:[query] format:[json] list:[search] srlimit:[1] srsearch:[golang] utf8:[1]]" url=https://en.wikipedia.org/w/api.php user-agent="wpdia-go/0.4.1 (github.com/lescactus/wpdia-go) WikiClient/0.4.1"
time=2025-01-21T11:30:35.175+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:245 msg="Sending http request..."
time=2025-01-21T11:30:35.462+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:255 msg="Http request sent..."
time=2025-01-21T11:30:35.462+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:260 msg="Reading http response body and unmarshalling.."
time=2025-01-21T11:30:35.462+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:277 msg="Http response body read and unmarshalled"
time=2025-01-21T11:30:35.462+01:00 level=INFO source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:292 msg="Search found a Page ID" pageid=25039021
time=2025-01-21T11:30:35.462+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/root.go:149 msg="Title found"
time=2025-01-21T11:30:35.462+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:82 msg="Setting http request parameters..."
time=2025-01-21T11:30:35.462+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:91 msg="Http request parameters set" params="map[exintro:[1] explaintext:[1] exsectionformat:[plain] pageids:[25039021] prop:[extracts|pageprops]]"
time=2025-01-21T11:30:35.462+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:137 msg="Building http request..." params="map[exintro:[1] explaintext:[1] exsectionformat:[plain] pageids:[25039021] prop:[extracts|pageprops]]" url=https://en.wikipedia.org/w/api.php user-agent="wpdia-go/0.4.1 (github.com/lescactus/wpdia-go) WikiClient/0.4.1"
time=2025-01-21T11:30:35.462+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:151 msg="Http request built" params="map[action:[query] exintro:[1] explaintext:[1] exsectionformat:[plain] format:[json] pageids:[25039021] prop:[extracts|pageprops]]" url=https://en.wikipedia.org/w/api.php user-agent="wpdia-go/0.4.1 (github.com/lescactus/wpdia-go) WikiClient/0.4.1"
time=2025-01-21T11:30:35.463+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:156 msg="Sending http request..."
time=2025-01-21T11:30:35.769+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:167 msg="Http request sent"
time=2025-01-21T11:30:35.770+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:178 msg="Reading http response body and unmarshalling..."
time=2025-01-21T11:30:35.770+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go:189 msg="Http response body read and unmarshalled"
time=2025-01-21T11:30:35.770+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/root.go:182 msg="Text extract found" title=golang random=false
time=2025-01-21T11:30:35.770+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/root.go:223 msg="Setting formatter..."
time=2025-01-21T11:30:35.770+01:00 level=DEBUG source=/home/amaldeme/gitclone/wpdia-go/cmd/root.go:241 msg="Formatter set to plain"
Title:
  Go (programming language)

Extract:
  Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency. It is often referred to as Golang because of its former domain name, golang.org, but its proper name is Go.There are two major implementations:

Google's self-hosting "gc" compiler toolchain, targeting multiple operating systems and WebAssembly.
gofrontend, a frontend to other compilers, with the libgo library. With GCC the combination is gccgo; with LLVM the combination is gollvm.A third-party source-to-source compiler, GopherJS, compiles Go to JavaScript for front-end web development. 
```

### Json logging

```
{"time":"2025-01-21T11:31:07.033969798+01:00","level":"INFO","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":83},"msg":"Creating new Wiki client...","url":"https://en.wikipedia.org/w/api.php"}
{"time":"2025-01-21T11:31:07.034029565+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.NewWikiClient","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":36},"msg":"Parsing base URL...","url":"https://en.wikipedia.org/w/api.php"}
{"time":"2025-01-21T11:31:07.034040644+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.NewWikiClient","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":47},"msg":"Base URL parsed","url":"https://en.wikipedia.org/w/api.php"}
{"time":"2025-01-21T11:31:07.034044089+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.NewWikiClient","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":62},"msg":"User-Agent set","user-agent":"wpdia-go/0.4.1 (github.com/lescactus/wpdia-go) WikiClient/0.4.1"}
{"time":"2025-01-21T11:31:07.03404882+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":100},"msg":"New Wiki client created","url":"https://en.wikipedia.org/w/api.php"}
{"time":"2025-01-21T11:31:07.034053185+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":105},"msg":"Disabling 'exintro'..."}
{"time":"2025-01-21T11:31:07.034056696+01:00","level":"INFO","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":118},"msg":"Getting text extract...","title":"golang","random":false}
{"time":"2025-01-21T11:31:07.034060131+01:00","level":"INFO","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":129},"msg":"Searching title...","title":"golang"}
{"time":"2025-01-21T11:31:07.03406576+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).SearchTitle","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":218},"msg":"Http request parameters set","params":{"list":["search"],"srlimit":["1"],"srsearch":["golang"],"utf8":["1"]}}
{"time":"2025-01-21T11:31:07.034088862+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).SearchTitle","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":226},"msg":"Building http request...","params":{"list":["search"],"srlimit":["1"],"srsearch":["golang"],"utf8":["1"]},"url":"https://en.wikipedia.org/w/api.php","user-agent":"wpdia-go/0.4.1 (github.com/lescactus/wpdia-go) WikiClient/0.4.1"}
{"time":"2025-01-21T11:31:07.03410538+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).SearchTitle","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":240},"msg":"Http request built","params":{"action":["query"],"format":["json"],"list":["search"],"srlimit":["1"],"srsearch":["golang"],"utf8":["1"]},"url":"https://en.wikipedia.org/w/api.php","user-agent":"wpdia-go/0.4.1 (github.com/lescactus/wpdia-go) WikiClient/0.4.1"}
{"time":"2025-01-21T11:31:07.034116426+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).SearchTitle","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":245},"msg":"Sending http request..."}
{"time":"2025-01-21T11:31:07.322802898+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).SearchTitle","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":255},"msg":"Http request sent..."}
{"time":"2025-01-21T11:31:07.322837498+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).SearchTitle","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":260},"msg":"Reading http response body and unmarshalling.."}
{"time":"2025-01-21T11:31:07.322963092+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).SearchTitle","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":277},"msg":"Http response body read and unmarshalled"}
{"time":"2025-01-21T11:31:07.322971978+01:00","level":"INFO","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).SearchTitle","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":292},"msg":"Search found a Page ID","pageid":25039021}
{"time":"2025-01-21T11:31:07.32298193+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":149},"msg":"Title found"}
{"time":"2025-01-21T11:31:07.322987044+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).GetExtract","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":82},"msg":"Setting http request parameters..."}
{"time":"2025-01-21T11:31:07.322995601+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).GetExtract","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":91},"msg":"Http request parameters set","params":{"exintro":["1"],"explaintext":["1"],"exsectionformat":["plain"],"pageids":["25039021"],"prop":["extracts|pageprops"]}}
{"time":"2025-01-21T11:31:07.323010539+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).do","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":137},"msg":"Building http request...","params":{"exintro":["1"],"explaintext":["1"],"exsectionformat":["plain"],"pageids":["25039021"],"prop":["extracts|pageprops"]},"url":"https://en.wikipedia.org/w/api.php","user-agent":"wpdia-go/0.4.1 (github.com/lescactus/wpdia-go) WikiClient/0.4.1"}
{"time":"2025-01-21T11:31:07.323024219+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).do","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":151},"msg":"Http request built","params":{"action":["query"],"exintro":["1"],"explaintext":["1"],"exsectionformat":["plain"],"format":["json"],"pageids":["25039021"],"prop":["extracts|pageprops"]},"url":"https://en.wikipedia.org/w/api.php","user-agent":"wpdia-go/0.4.1 (github.com/lescactus/wpdia-go) WikiClient/0.4.1"}
{"time":"2025-01-21T11:31:07.323031811+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).do","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":156},"msg":"Sending http request..."}
{"time":"2025-01-21T11:31:07.469715517+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).do","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":167},"msg":"Http request sent"}
{"time":"2025-01-21T11:31:07.469801848+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).do","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":178},"msg":"Reading http response body and unmarshalling..."}
{"time":"2025-01-21T11:31:07.469898822+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).do","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":189},"msg":"Http response body read and unmarshalled"}
{"time":"2025-01-21T11:31:07.469917014+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":182},"msg":"Text extract found","title":"golang","random":false}
{"time":"2025-01-21T11:31:07.469931546+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":223},"msg":"Setting formatter..."}
{"time":"2025-01-21T11:31:07.469946477+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":241},"msg":"Formatter set to plain"}
Title:
  Go (programming language)

Extract:
  Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency. It is often referred to as Golang because of its former domain name, golang.org, but its proper name is Go.There are two major implementations:

Google's self-hosting "gc" compiler toolchain, targeting multiple operating systems and WebAssembly.
gofrontend, a frontend to other compilers, with the libgo library. With GCC the combination is gccgo; with LLVM the combination is gollvm.A third-party source-to-source compiler, GopherJS, compiles Go to JavaScript for front-end web development.
```

### Random article

```
./wpdia-go --random
Title:
  John Matoian

Extract:
  John Matoian (born 1949) is a businessman and television industry executive. He was a vice-president of the CBS Entertainment division. He later became the president of Entertainment at Fox Broadcasting in September 1995. He was president at HBO from 1996 to 1999. He received both his B.A. and his J.D. from Duke University.
Matoian is a native of Fresno and is of Armenian descent.In the 2012 United States Presidential election, John Matoian had made $83,800 worth of contributions to Barack Obama's successful presidential campaign.In his book Springfield Confidential, Mike Reiss mentions Matoian by name as the Fox executive whose intense hatred of his and Al Jean's animated series The Critic led to its cancellation after a single season on the network.
```

### Random article + HTTP client timeout set to 500ms + json output + only 3 sentences + French language + full output + log level debug + log format json

```
./wpdia-go -t 500ms --output json --exsentences 3 --lang fr --full --loglevel debug --logformat json --random
{"time":"2025-01-21T11:32:09.017402147+01:00","level":"WARN","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":73},"msg":"The --random flag is set, the given arguments will be ignored: []"}
{"time":"2025-01-21T11:32:09.017458459+01:00","level":"INFO","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":83},"msg":"Creating new Wiki client...","url":"https://fr.wikipedia.org/w/api.php"}
{"time":"2025-01-21T11:32:09.017463063+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.NewWikiClient","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":36},"msg":"Parsing base URL...","url":"https://fr.wikipedia.org/w/api.php"}
{"time":"2025-01-21T11:32:09.017472394+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.NewWikiClient","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":47},"msg":"Base URL parsed","url":"https://fr.wikipedia.org/w/api.php"}
{"time":"2025-01-21T11:32:09.017475843+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.NewWikiClient","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":62},"msg":"User-Agent set","user-agent":"wpdia-go/0.4.1 (github.com/lescactus/wpdia-go) WikiClient/0.4.1"}
{"time":"2025-01-21T11:32:09.017480276+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":100},"msg":"New Wiki client created","url":"https://fr.wikipedia.org/w/api.php"}
{"time":"2025-01-21T11:32:09.017483799+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":105},"msg":"Disabling 'exintro'..."}
{"time":"2025-01-21T11:32:09.017487477+01:00","level":"INFO","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":118},"msg":"Getting text extract...","title":"","random":true}
{"time":"2025-01-21T11:32:09.017494719+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).GetExtractRandom","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":102},"msg":"Setting http request parameters..."}
{"time":"2025-01-21T11:32:09.017499661+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).GetExtractRandom","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":119},"msg":"Http request parameters set","params":{"explaintext":["1"],"exsectionformat":["plain"],"exsentences":["3"],"generator":["random"],"grnlimit":["1"],"grnnamespace":["0"],"prop":["extracts|pageprops"]}}
{"time":"2025-01-21T11:32:09.01752387+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).do","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":137},"msg":"Building http request...","params":{"explaintext":["1"],"exsectionformat":["plain"],"exsentences":["3"],"generator":["random"],"grnlimit":["1"],"grnnamespace":["0"],"prop":["extracts|pageprops"]},"url":"https://fr.wikipedia.org/w/api.php","user-agent":"wpdia-go/0.4.1 (github.com/lescactus/wpdia-go) WikiClient/0.4.1"}
{"time":"2025-01-21T11:32:09.017541852+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).do","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":151},"msg":"Http request built","params":{"action":["query"],"explaintext":["1"],"exsectionformat":["plain"],"exsentences":["3"],"format":["json"],"generator":["random"],"grnlimit":["1"],"grnnamespace":["0"],"prop":["extracts|pageprops"]},"url":"https://fr.wikipedia.org/w/api.php","user-agent":"wpdia-go/0.4.1 (github.com/lescactus/wpdia-go) WikiClient/0.4.1"}
{"time":"2025-01-21T11:32:09.017549683+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).do","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":156},"msg":"Sending http request..."}
{"time":"2025-01-21T11:32:09.224353509+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).do","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":167},"msg":"Http request sent"}
{"time":"2025-01-21T11:32:09.224639393+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).do","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":178},"msg":"Reading http response body and unmarshalling..."}
{"time":"2025-01-21T11:32:09.224900634+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.(*WikiClient).do","file":"/home/amaldeme/gitclone/wpdia-go/cmd/wpdia.go","line":189},"msg":"Http response body read and unmarshalled"}
{"time":"2025-01-21T11:32:09.224931725+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":182},"msg":"Text extract found","title":"","random":true}
{"time":"2025-01-21T11:32:09.224949915+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":223},"msg":"Setting formatter..."}
{"time":"2025-01-21T11:32:09.224972098+01:00","level":"DEBUG","source":{"function":"github.com/lescactus/wpdia-go/cmd.init.func1","file":"/home/amaldeme/gitclone/wpdia-go/cmd/root.go","line":241},"msg":"Formatter set to json"}
{
    "pageid": 4872215,
    "ns": 0,
    "title": "Château de Padern",
    "extract": "Le château de Padern surplombe le village du même nom, sur la route départementale 14 qui relie Cucugnan à Tuchan, dans le département français de l'Aude, en région Occitanie.\n\n\nSituation\nLe château est construit sur les hauteurs du village de Padern sur un piton de roches calcaires qui domine le Verdouble coulant en contrebas du village. Les chemins pour y accéder sont très escarpés, ce qui en faisait un ouvrage pratiquement imprenable.",
    "pageprops": {
        "wikibase_item": "Q2970165"
    }
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

- [x] Implement "random article" 

- [x] Fix 'may refer to:'

- [ ] Add suggestions for disambiguation pages

- [x] Improve base url

- [x] Dockerize

- [x] Debug flag (show page id, ns, timestamps, etc...)

- [x] Verbose logs