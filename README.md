# addressenrich

[![Project status](https://img.shields.io/github/release/tmc/addressenrich.svg?style=flat-square)](https://github.com/tmc/addressenrich/releases/latest)
[![Build Status](https://github.com/tmc/addressenrich/workflows/Test/badge.svg)](https://github.com/tmc/addressenrich/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/tmc/addressenrich?cache=0)](https://goreportcard.com/report/tmc/addressenrich)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/tmc/addressenrich)

Command addressenrich adds a column to csv files that includes normalized aaddresses.

## Installation

addressenrich is a [Go](https://golang.org/) program for MacOS systems.

Presming  you have a working Go insallation, you can install `addressenrich` via:

```console
go install github.com/tmc/addressenrich
```

## Usage

```console
$ addressenrich -h
Usage of addressenrich:
  -input-file string
    	input file (default "-")
  -skip-rows int
    	rows to skip (default 1)
  -start-column int
    	column start index
  -v	verbose mode
```

