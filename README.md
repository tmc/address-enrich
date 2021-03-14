# address-enrich

[![Project status](https://img.shields.io/github/release/tmc/address-enrich.svg?style=flat-square)](https://github.com/tmc/address-enrich/releases/latest)
[![Build Status](https://github.com/tmc/address-enrich/workflows/Test/badge.svg)](https://github.com/tmc/address-enrich/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/tmc/address-enrich?cache=0)](https://goreportcard.com/report/tmc/address-enrich)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/tmc/address-enrich)

Command address-enrich adds a column to csv files that includes normalized aaddresses.

## Installation

address-enrich is a [Go](https://golang.org/) program for MacOS systems.

Presming  you have a working Go insallation, you can install `address-enrich` via:

```console
go install github.com/tmc/address-enrich
```

## Usage

```console
$ address-enrich -h
Usage of address-enrich:
  -input-file string
    	input file (default "-")
  -skip-rows int
    	rows to skip (default 1)
  -start-column int
    	column start index
  -v	verbose mode
```

