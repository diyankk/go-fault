/*
Package fault provides standard http middleware for fault injection in go.

Basics

Use the fault package to inject faults into the http request path of your service. Faults work by
modifying and/or delaying your service's http responses. Place the Fault middleware high enough in
the chain that it can act quickly, but after any other middlewares that should complete before fault
injection (auth, redirects, etc...).

The type and severity of injected faults is controlled by a single Options struct passed to a new
Fault struct. The Options struct must contain a field Injector, which is an interface that holds the
actual fault injection code in Injector.Handler. The Fault struct wraps Injector.Handler in another
Fault.Handler that applies generic Fault logic (such as what % of requests to run the Injector on)
to the Injector.

Package provided Handlers will always default to a "do nothing, pass request on" state if the
provided options are invalid. Make sure you use the NewFault() and NewTypeInjector() constructors to
create validated Faults and Injectors. If you are not seeing faults injected like you expect you may
have passed an out of bounds value, invalid http status code, incorrect percent, or other wrong
parameter.

Injectors

There are three main Injectors provided by the fault package:

    fault.RejectInjector
    fault.ErrorInjector
    fault.SlowInjector

RejectInjector

Use fault.RejectInjector to immediately return an empty response. For example, a curl for a rejected
response will produce:

    $ curl https://github.com
    curl: (52) Empty reply from server

ErrorInjector

Use fault.ErrorInjector to immediately return an http status code of your choosing along with the
standard HTTP response body for that code. For example, you can return a 200, 301, 418, 500, or any
other valid status code to test how your clients respond to different statuses. If ErrorInjector has
an invalid status code the middleware will pass on the request without evaluating.

SlowInjector

Use fault.SlowInjector to wait a configured time.Duration before proceeding with the request as
normal. For example, you can use the SlowInjector to add a 10ms delay to your incoming requests.

Combining Faults

It is easy to combine any of the Injectors into a chained action. There are two ways you might want
to combine Injectors.

First, you can create separate Faults for each Injector that are sequential but independent of each
other. For example, you can chain Faults such that 1% of requests will return a 500 error and
another 1% of requests will be rejected.

Second, you might want to combine Faults such that 1% of requests will be slowed for 10ms and then
rejected. You want these Faults to depend on each other. For this use the special ChainInjector,
which consolidates any number of Injectors into a single Injector that runs each of the provided
Injectors sequentially. When you add the ChainInjector to a Fault the entire chain will always
execute together.

Blacklisting Paths

The fault.Options struct has an option PathBlacklist. Any path you include in this list will never
have faults run against it. The paths that you include must match exactly the path in req.URL.Path,
including leading and trailing slashes.

Whitelisting Paths

The fault.Options struct has an option PathWhitelist. If you pass a non-empty list here then faults
will only be evaluated on the paths provided. Path blacklists take priority over whitelists. The paths
that you include must match exactly the path in req.URL.Path, including leading and trailing slashes.

Custom Injectors

The package provides an Injector interface and you can satisfy that interface to provide your own
Injector. Use custom injectors to add additional logic (logging, stats) to the package-provided
injectors or to create your own completely new Injector that can still be managed by the Fault
struct.

Configuration

All configuration for the fault package is done through the Options struct. There is no other way to
manage configuration for the package. It is up to the user of the fault package to manage how the
Options struct is generated. Common options are feature flags, environment variables, or code
changes in deploys.

*/
package fault
