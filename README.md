# Architecture

Monitor is a client-server application written in Go.

The daemon, `monitord`, does all the heavy-lifting:

* reading the list of websites from a config file
* polling websites on a regular basis
* storing metrics in memory
* listening for `monitorctl` client requests
* aggregating metrics on-the-fly
* generating alerts when appropriate

The client, `monitorctl`:

* regularly polls the daemon for the latest aggregated metrics
* regularly polls the daemon for the latest alerts
* presents these results on a console dashboard

`monitord` and `monitorctl` communicate with each other using Remote Procedure Call.

## Why a client-server architecture?

## Why Go?

Go has many great features, amongst which:

* As it is a compiled, statically typed language, it is faster and requires less resources than dynamically typed languages such as Python or JavaScript. Still, its type system is more straightforward than those of C++ or Java
* By design, Go is a concurrent language. It is an especially interesting feature for this project, as the daemon has to deal with numerous tasks at once, such as polling thousands of websites while aggregating metrics and responding to the client. Gorountines and channels provide an effective way of doing all those things while keeping a logical, structured program.

## Why RPC?

## About process daemonization

## Metrics: effective monitoring

Redirects not followed, by choice.

## Folder and files structure

# Usage

## Requirements

## Quick start

### Building

### Running

### Testing

### Documentation

Use godoc.

# Future improvements

* stateless daemon
* multiple pollers
* dashboard with search engine to navigate through thousands of websites
* Alert web hooks, to be notified through Slack / Telegram / whatever
* Config validation
* Granular website config:
  * different polling interval
  * different availability threshold, depending on SLOs
* Handle errors and valid response differently, to avoid false statistics
