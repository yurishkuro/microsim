# microsim - Microservices Simulator

Used in https://github.com/PacktPublishing/Mastering-Distributed-Tracing/tree/master/Chapter12.

Published on Docker Hub: https://hub.docker.com/r/yurishkuro/microsim/tags

## Usage

To see all command line options:

```shell
docker run yurishkuro/microsim -h
```

`microsim` uses OpenTelemetry SDK to export traces. By default it will try to send them to `https://localhost:4318/v1/traces` (with TLS enabled). This can be changed by setting environment variables supported by the SDK:
  * `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` to point to a different host/port/endpoint.
  * `OTEL_EXPORTER_OTLP_ENDPOINT` to point to a different host/port (will still use default `/v1/traces` endpoint).
  * `OTEL_EXPORTER_OTLP_INSECURE=true` if you just want to switch default URL to `http` instead of `https`.

Note that when we run `microsim` as a container, the `localhost` refers to the container's inner network namespace, so it will not be able to reach the collector even if it's running on the same host. You might see an error like this:

```
2024/07/02 18:06:07 traces export: Post "https://localhost:4318/v1/traces": dial tcp [::1]:4318: connect: connection refused
```

To work around that, refer to the IP address of the host instead:

```shell
docker run --env OTEL_EXPORTER_OTLP_ENDPOINT=http://{YOUR_IP_ADDRESS}:4318/ \
  yurishkuro/microsim -w=1 -r=1
```

## Design

The tool takes a configuration file that describes the "architecture" of your desired system, such as services,
number of instances for each service, their endpoints and their dependecies. The tool simulates the traces that
may be generated from such system. Some random mutations are introduced during execution: latency of requests
is drawn randomly from a normal distribution with specified mean and stdev (error rates can be specified, but
currently not implemented).

The tool starts a configured number of "workers" (goroutines), and each worker sequentially runs simulations
of requests through the architecture, starting from the root service. Workers can run either for a certain
time period, or until they generate a predefined number of simulations. There is a configurable sleep interval
between simulations in each worker.

## Configuration

The tool comes with a built-in configuration for Jaeger's HotROD application that can be printed with:

```shell
docker run yurishkuro/microsim -o | jq
```

* A custom configuration can be provided via `-c` option.
* The schema for the config is hardcoded in the data model: [model/config.go](./model/config.go).
* The UML diagram below is generated with https://www.dumels.com/.

![UML diagram of configuration](model/uml-diagram.png)

## License

MIT license.
