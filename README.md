# About

Monitor is a website monitoring project proposed by Datadog.

It is a client-server application written in Go.
Two binaries are available: `monitord`, the daemon, and `monitorctl`, the client.

![preview](https://github.com/oxlay/monitor/blob/master/preview.png)

# Table of contents

* [Usage](#usage)
  * [Requirements](#requirements)
  * [Install](#install)
  * [About config files](#about-config-files)
  * [Testing](#testing)
  * [Documentation](#documentation)
  * [About dependencies](#about-dependencies)
* [Architecture](#architecture)
  * [Overview](#overview)
  * [Design choices](#design-choices)
  * [Folder and files structure](#folder-and-files-structure)
* [Possible improvements](#possible-improvements)
  * [Project-wide improvements](#project-wide-improvements)
  * [Daemon-specific improvements](#daemon-specific-improvements)
  * [Dashboard-specific improvements](#dashboard-specific-improvements)

# Usage

## Requirements

**Go 1.10 recommended.** Go 1.7+ might be supported but has not been tested.

The packages have been tested on **macOS and Linux**.

## Install

```
go get github.com/oxlay/monitor/cmd/monitord github.com/oxlay/monitor/cmd/monitorctl
```

Providing that `$GOPATH/bin` is in your `$PATH`, you should be able to:

* **start the daemon** by simply running `monitord`
* **start the dashboard** in a separate window by running `monitorctl`

On daemon startup, you may need to wait a few seconds for poll results to be available.

## Configuration files

By default, `monitorctl` and `monitord` respectively look for the following config files provided in the repo:

* `$GOPATH/src/github.com/oxlay/monitor/cmd/monitorctl/config.json`
* `$GOPATH/src/github.com/oxlay/monitor/cmd/monitord/config.json`

You can override those defaults and pass any config flag using the `-config` flag:

```
monitord -config path/to/config-monitord.json &
monitorctl -config path/to/config-monitorctl.json
```

Documentation about the content of config files is available [through GoDoc](https://godoc.org/github.com/oxlay/monitor).

## Testing

To run tests for the alert logic:

```
cd $GOPATH/src/github.com/oxlay/monitor/cmd/monitord/daemon
go test
```

These tests are written following [table-driven testing](https://github.com/golang/go/wiki/TableDrivenTests) principles.

## Documentation

The project documentation is available [here](https://godoc.org/github.com/oxlay/monitor).

As an effort to provide easy access to the project's documentation, a choice was made to export all methods, thus making them available through `godoc`.
I believe it is an acceptable trade-off, as the `client` and `daemon` packages will not be distributed as libraries (the folder structure prevents such use cases), and are only meant to be used through `monitorctl` and `monitord` commands.

## About dependencies

Dependencies are included in the `vendor/` folder to allow for one-line install with `go get`.

# Architecture

## Overview

Monitor is a client-server application written in Go.

The daemon, `monitord`, does most of the heavy-lifting:

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

## Design choices

### Why Go?

Among Go's great features, its concurrency model, modern syntax, and low resource consumption made it a natural choice for the project.

Besides, `go get` is a plus, as it makes the install process of `monitord` and `monitorctl` remarkably simple.

### Why a client-server architecture?

Using a client-server architecture provides numerous benefits, the most notable ones being:

* separation of concerns: how websites are polled should be separate from how users interact with the result
* ability to leave the daemon running in the background: `monitord` could be running 24/7 and controlled by a service manager such as `systemd`
* ability to poll the websites from one machine and present the user interface on another. Typically, the daemon could be running on a server and users could occasionally take a look at the results from their laptop (without needing an interrupted network connection)

### Why store metrics in memory?

This choice was made in order to keep the project as simple as it needs to be. It results in less code and a more straightforward installation process than if a database needed to be installed and configured.

It could evolve, in a future iteration, to use a time-series database that stores the poll results, thus making the daemon stateless and more scalable. See [possible improvements](#daemon-specific-improvements).

### Why RPC?

Remote Procedure Call provides a clean and lightweight means of communication between processes.

As `net/rpc` provides a straightforward implementation, it results in more idiomatic code than other solutions such as REST API endpoints.

### Thoughts on process daemonization

Tests were made for `monitord` to be self-daemonizing, but the results were not convincing.

It would allow the user to launch the daemon without needing a separate window for the dashboard (or without the need to append the `monitord` command with the `&` job control character). Yet it comes with a set of challenges that would not be worth the effort.

Go does not provide support for daemonization out of the box. While [a](https://github.com/takama/daemon) [few](https://github.com/sevlyar/go-daemon) [libraries](https://github.com/VividCortex/godaemon) are available on Github, they are generally cumbersome to use and added an undesired level of complexity.

In order to keep the code straightforward, a decision was made not to use such libraries, and leave the user to deal with his platform-specific tools (`launchctl` on macOS, `systemctl` on Ubuntu, etc.), should he need a daemon that runs 24/7.

### Choosing the right metrics for effective monitoring

**In addition to response times, the dashboard provides the duration of each phase of HTTP requests.**
This request breakdown allows website maintainers to understand which parts of the request are the most problematic and should be improved in priority.

**A choice was made _not_ to follow redirects.**
Indeed, monitoring redirections can be insightful in itself: it is important to know how fast a page responds, even if it gives a 301 response code. And the response time of the redirecting page should not be mixed with the response time of the page it redirects to.

**Another decision was made not to show minimum response times to the user.**
In an effort not to overwhelm the user with low-value information, minimum response times are not shown on the dashboard. Indeed, it would provide little insight into how long a website takes to respond for an average user. Infrastructure maintainers should focus on optimizing max and average response times, rather than optimizing a min response time that very few users will experience.

## Folder and files structure

```
.
├── cmd                    # Contains standalone packages (CLI commands)
│   ├── monitorctl         # Client command
│   │   └── client         # Library used by monitorctl
│   └── monitord           # Daemon command
│       └── daemon         # Library used by monitord
├── internal               # Packages used by both monitorctl and monitord
│   └── payload            # Types used to communicate between monitorctl and monitord
└── vendor                 # Dependencies
```

# Possible improvements

## Project-wide improvements

**Metrics analysis:** Google's SRE team offers [valuable insights](https://landing.google.com/sre/book/chapters/monitoring-distributed-systems.html) into making an effective monitoring system. This project would benefit from implementing some of their suggestions, such as:

* separating the timing calculations of valid responses and those of error responses: indeed, if a website randomly throws 500 errors very fast, it does not mean that the website is fast, so those results should be separated from the average response time of valid responses
* creating configurable "policy errors": for example, if a website responds with a 200 response code in more than 3 seconds, it could be logged as an error
* bucketing results: displaying the distribution of response times (e.g. the number of requests with a response time between 0 and 100 ms, between 100 ms and 300 ms, between 300 and 800 ms, etc.) would show if there is a tail of slow responses that negatively impact the average response time

**Configuration check:** currently, the configuration file validity is not checked on startup of `monitord` or `monitorctl`. Basic checks such as URL validity checks could be implemented.

**Unit testing:** currently, only the alerting logic is tested by `go test`. If the project development would continue, improving code coverage by writing more tests might be a good investment, as it could make the project more reliable and maintainable.

## Daemon-specific improvements

**Notifications:** as looking at a dashboard all day might get tiresome, a notification system could be implemented. Website maintainers would therefore be notified (e.g. on Slack) when a website is down.

**Database backend:** as mentioned in _[Why store metrics in memory?](#why-store-metrics-in-memory)_, if the project was used in a context where scalability is a concern, then using a time-series database would be more appropriate. Amongst others, it would reduce memory usage (above a certain number of websites), allow for longer data retention, and prevent data loss if the daemon is restarted.

**Poller architecture:** currently, for each website in the config file, a goroutine is created to regularly poll the website. While this straightforward approach works well for moderate loads, it might not scale well as the number of websites grows. In this case, refactoring the polling logic might be necessary, and the [dispatcher-worker architecture proposed by Marcio Castilho](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/) could be a good source of inspiration.

## Dashboard-specific improvements

**Search engine:** navigating through the dashboard using left/right arrows is fine for a few websites, but can quickly get irritating when the number grows. In this case, a basic text input allowing the user to choose which website to show may be more appropriate.

**Resiliency to network interruptions:** currently, the dashboard exits when it fails to connect to the daemon. This behavior is considered acceptable as long as the daemon and client are running on the same machine. However, if the daemon was used on a server, and the client on a user's laptop, the network connection between these two components would be less reliable. In this case, the dashboard should try to recover from a network failure by making new connection attempts to the daemon.

**Dynamically set dashboard's height:** currently, the library used for displaying the dashboard ([`termui`](https://github.com/gizak/termui)) does not support adapting the UI components' height to the window's height. Therefore, users with small terminal windows may not see the bottom of the dashboard. In a future iteration, the dashboard's height could be computed from the window's height.
