# wiki-api

[![Go](https://github.com/ENG618/wiki-api/actions/workflows/go.yml/badge.svg)](https://github.com/ENG618/wiki-api/actions/workflows/go.yml)
[![CodeQL](https://github.com/ENG618/wiki-api/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/ENG618/wiki-api/actions/workflows/codeql-analysis.yml)

## Overview

This is a go package that exposes an API server that exposes the information on
the Wikipedia page about Heidenheim, Germany
(<https://en.wikipedia.org/wiki/Heidenheim_an_der_Brenz>) in a useful and
intuitive way.

## Running the server

- clone the repo
- `go run main.go` from the root of the directory

## Routes

### /

The base route simply says hello

### /panic

The panic route is a test route on what happens when a panic is invoked, and how the program still recovers.

### /mayor

The mayor route tells you who the current mayor of Heidenheim, Germany is

### /mayor/{year}

The mayor year route allows you to pass in a year, and returns the mayor/s for that year. There is some error handling on this route. So if you pass in anything other than a valid year, the response will be a `400` with information on what went wrong.
