/*
Package daemon implements regular website polling and metrics aggregation.

It regularly polls the websites specified in the config file,
stores the poll results in memory,
listens for RPC client request,
aggregates metrics on-the-fly,
and generates alerts when appropriate.
*/
package daemon
