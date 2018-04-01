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

## Why RPC?

## About process daemonization

## Metrics: effective monitoring

# Usage

## Requirements

## Quick start

### Building

### Running

### Testing

### Documentation

Use godoc.
