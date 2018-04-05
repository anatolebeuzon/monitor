# About

Monitor is a website monitoring project proposed by Datadog.

It is a client-server application written in Go.
Two binaries are available: `monitord`, the daemon, and `monitorctl`, the client.

![preview](TODO)

# Table of Contents

TODO

# Usage

## Requirements

Go 1.10 recommended. Go 1.7+ might be supported but has not been tested.

The packages have been tested on macOS and Linux.

## Install

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

Go has many great features, amongst which:

* As it is a compiled, statically typed language, it is faster and requires less resources than dynamically typed languages such as Python or JavaScript. Still, its type system is more straightforward than those of C++ or Java
* By design, Go is a concurrent language. It is an especially interesting feature for this project, as the daemon has to deal with numerous tasks at once, such as polling a potentially large number of websites while aggregating metrics and responding to the client. Gorountines and channels provide an effective way of doing all those tasks while keeping a logical, structured program.

### Why a client-server architecture?

Using a client-server architecture provides numerous benefits, the most notable ones being:

* separation of concerns: how websites are polled should be separate from how users interact with the result
* the ability to leave the daemon running in the background: `monitord` could be running 24/7 and controlled by a service manager such as `systemd`
* the ability to poll the websites from one machine and present the user interface on another. Typically, the daemon could be running on a server and users could occasionally take a look at the results from their laptop (without needing an interrupted network connection)

### Why store metrics in memory?

This choice was made in order to keep the project as simple as it needs to be. It results in less code and a more straightforward installation process (no need to install and configure a database).

It could evolve, in a future iteration, to use a time-series database that store the poll results, thus making the daemon stateless and more scalable.

### Why RPC?

Remote Procedure Call provides a clean and lightweigt means of communication between processes.

As `net/rpc` provides a straightforward implementation, it results in more idiomatic code than other solutions such as REST API endpoints.

### Thoughts on process daemonization

Tests were made for `monitord` to be self-daemonizing, but the result was not convincing.

It would allow the user to launch the daemon without needing a separate window for the dashboard (or without needing to append the `monitord` command with the `&` job control character). Yet it comes with a set of challenges that would not be worth the effort.

Go does not provide support for daemonization out of the box. While [a](https://github.com/takama/daemon) [few](https://github.com/sevlyar/go-daemon) [libraries](https://github.com/VividCortex/godaemon) were available on Github, they were generally cumbersome to use and added an undesired level of complexity.

In order to keep the code straightword, a decision was made not to use such libraries, and leave the user to deal with his platform-specific tools (`launchctl` on macOS, `systemctl` on Ubuntu, etc.) would he come to need 24/7 daemonization.

### Metrics: effective monitoring

**A choice was made _not_ to follow redirects.**
Indeed, monitoring redirections can be insightful in itself: it is important to know how fast a page responds, even if it gives a 301 reponse code. And the response time of the redirecting page should not be mixed with the response time of the page it redirects to.

**Another decision was made not to show minimum response times to the user.**
In an effort not to overwhelm the user with low-value information, minimum response times are not shown on the dashboard. Indeed, it would provide little insight into how long a website takes to respond for an average user. Infrastructure mainteners should optimize max and average response times, rather than optimizing a min response time that very few users will experience.

## Folder and files structure

TODO

# Future improvements

## Daemon

**Notifications:** as looking at a dashboard all day might get tiresome, a notification system could be implemented. Website maintainers would therefore be notified (e.g. on Slack) when a website is down.

**Database backend:** as mentioned in _[Why RPC?](TODO)_, if the project was used in a context were scalability is a concern, then using a time-series database would be more appropriate. Amongst others, it would reduce memory usage (above a certain number of websites), allow for longer data retention, and prevent data loss if the daemon was restarted.

* multiple pollers to handle more websites

## Dashboard

**Dynamically set dashboard height:** currently, the library used for displaying the dashboard (`gizak/termui`) does not support adapting the UI components' height to the window height. Therefore, users with small terminal windows may not see the bottom of the dashboard. In a future iteration, the dashboard's height could be computed from the window's height.

**Search engine:** navigating through the dashboard using left/right arrows is fine for a few websites, but can quickly get irritating when the number grows. A basic text input allowing the user to choose which website to show would then be more appropriate.

**Resiliency to network interruptions:** currently, the dashboard exits when it fails to connect to the daemon. This behavior was considered acceptable as long as the daemon and client are running on the same machine. However, if the daemon were to be used on a server, and the client on a user's laptop, the network connection between these two components would be less reliable. In this case, the dashboard should try to recover from a network failure by making new connection attempts to the daemon.

## Both

**Configuration check:**

* More tests

* Handle errors and valid response differently, to avoid false statistics
* Policy errors
* collect request counts bucketed by latencies

https://landing.google.com/sre/book/chapters/monitoring-distributed-systems.html
