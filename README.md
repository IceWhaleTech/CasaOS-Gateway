# CasaOS-Gateway

[![Go Reference](https://pkg.go.dev/badge/github.com/IceWhaleTech/CasaOS-Gateway.svg)](https://pkg.go.dev/github.com/IceWhaleTech/CasaOS-Gateway) [![Go Report Card](https://goreportcard.com/badge/github.com/IceWhaleTech/CasaOS-Gateway)](https://goreportcard.com/report/github.com/IceWhaleTech/CasaOS-Gateway) [![goreleaser](https://github.com/IceWhaleTech/CasaOS-Gateway/actions/workflows/release.yml/badge.svg)](https://github.com/IceWhaleTech/CasaOS-Gateway/actions/workflows/release.yml) [![codecov](https://codecov.io/gh/IceWhaleTech/CasaOS-Gateway/branch/main/graph/badge.svg?token=5JIHXF1RJ4)](https://codecov.io/gh/IceWhaleTech/CasaOS-Gateway)

CasaOS Gateway is a dynamic API gateway service that can be used to expose APIs from different other HTTP based services.

This gateway service comes with a simple management API for other services to register their APIs by route paths. A HTTP request arrived at gateway port will be forwarded to the service that is registered for the route path.

> As a best practice, a service behind this gateway should bind to localhost (`127.0.0.1` for IPv4, `::1` for IPv6) ONLY, so no external network access is allowed.

## Configuration

Upon launching, it will search for `gateway.ini` file in the following order:

```bash
./gateway.ini
./conf/gateway.ini
$HOME/.casaos/gateway.ini
/etc/casaos/gateway.ini
```

See [gateway.ini.sample](./build/etc/casaos/gateway.ini.sample) for default configuration.

## Running

Once running, gateway address and management address will be available in the files under `RuntimePath`  specified in configuration.

```bash
$ cat /var/run/casaos/gateway.url 
[::]:8080 # port is specified in configuration

$ cat /var/run/casaos/management.url 
[::]:34703 # port is randomly assigned
```

## Example

Assuming that

- the management API is running on port `34703`
- the gateway is running on port `8080`
- some API running at `http://localhost:12345/ping` that simply returns `pong`.

Then register the API as follows:

- POST `http://localhost:34703/v1/gateway/routes`

  ```json
  {
          "path": "/ping",
          "target": "http://localhost:12345"
  }
  ```

  or in command line:

  ```bash
  $ curl 'localhost:34703/v1/gateway/routes' --data-raw '
      {"path": "/ping", "target": "http://localhost:12345"}
    '
  ```

Now run

```bash
$ curl localhost:8080/ping
{"message":"pong"}
```

... which is equivalent as

```bash
$ curl localhost:12345/ping
{"message":"pong"}
```
