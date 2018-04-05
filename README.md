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
No min showed, by choice.

## Folder and files structure

# Usage

## Requirements

Go 1.10 recommended. Go 1.7+ may be supported but has not been tested.

The packages have been tested on macOS and Linux.

## Quick start

Run:

```
go get github.com/oxlay/monitor/cmd/monitord github.com/oxlay/monitor/cmd/monitorctl
```

Providing that `$GOPATH/bin` is in your `$PATH`, you should be able:

* to start the daemon by simply running `monitord`
* to start the dashboard in a separate window by running `monitorctl`

On daemon startup, you may need to wait a few seconds for poll results to be available.

## About config files

By default, `monitorctl` and `monitord` respectively look for the following config files provided in the repo:

* `$GOPATH/src/github.com/oxlay/monitor/cmd/monitorctl/config.json`
* `$GOPATH/src/github.com/oxlay/monitor/cmd/monitord/config.json`

You can override those defaults and pass any config flag using the `-config` flag:

```
monitord -config path/to/config-monitord.json & monitorctl -config path/to/config-monitorctl.json
```

## Testing

To run tests for the alert logic:

```
cd $GOPATH/src/github.com/oxlay/monitor/daemon
go test
```

## Documentation

The project documentation is available [here](https://godoc.org/github.com/oxlay/monitor).

As an effort to provide easy access to the project's documentation, I purposefully chose to export all methods, thus making them available through `godoc`.
I believe it is an acceptable trade-off, as the `client` and `daemon` packages will not be distributed as libraries, and are only meant to be used through `monitorctl` and `monitord` commands.

## About dependencies

Dependencies are included in the `vendor/` folder to allow for one-line install with `go get`.

# Future improvements

* stateless daemon (here, stored in memory to keep it simple)
* multiple pollers
* dashboard with search engine to navigate through thousands of websites
* Alert web hooks, to be notified through Slack / Telegram / whatever
* Config validation
* Granular website config:
  * different polling interval
  * different availability threshold, depending on SLOs
* Handle errors and valid response differently, to avoid false statistics
* Improve network resiliency
* Adaptive dashboard height
* dashboard: accept more than two timeframes

```

```
