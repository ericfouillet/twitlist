[![Build Status](https://travis-ci.org/ericfouillet/twitlist.svg?branch=master)](https://travis-ci.org/ericfouillet/twitlist)

# twitlist

A simple Twitter list manager.

## Introduction

The Twitter web interface doesn't allow adding/removing multiple users from a list in one go.

Lists are a good way to separate accounts by theme, so I created this small application to manage lists more easily.

The server is in Go (exposing a REST interface), the UI is in Elm. This is a work in progress.

## Installation:

- Install [Go](https://golang.org) and [Elm](http://elm-lang.org/)
- Clone the repository
- In _twitlistserver_: `go install .`
- In _cmd_: `go run local.go`

## Dependencies

This depends on a slightly modified version of [Anaconda](https://github.com/ChimeraCoder/anaconda), to add missing API calls. It is forked under my account.
