ROADMAP:
========

v0.1
----
- [x] simple forwarder (fiber handler)
- [x] simple tcp port reserve

v1.0
----
- [ ] control REST API
- [ ] configuration from yaml
- [ ] configuration from db
- [ ] Round Robin balancer(+errors counter)

v2.0
----
- [ ] ssh port forwarder(reverse/local)
- [ ] openvpn tunnel manager

v2.5
----
- [ ] proxy agents support
- [ ] lua scripting support


Commands:
=========

```bash
$ go run cmd/main.go --help

Usage: main <command>

Flags:
  -h, --help     Show context-sensitive help.
      --debug    Enable debug mode.

Commands:
  forward <target> <listen_port> <source_path_prefix> [<headers> ...]
    forward connection

  reserve <base_host> <base_port> <reserve_host> <reserve_port> [<listen_port> [<connect_timeout> [<max_idle_seconds>]]]
    reserve service

Run "main <command> --help" for more information on a command.
```

forward
--------
```bash
$ go run cmd/main forward --help
Usage: main forward <target> <listen_port> <source_path_prefix> [<headers> ...]

forward connection

Arguments:
  <target>                target
  <listen_port>           listen_port
  <source_path_prefix>    source path prefix
  [<headers> ...]         extra headers

Flags:
  -h, --help                   Show context-sensitive help.
      --debug                  Enable debug mode.

      --headers=HEADERS,...

$ go run cmd/main.go forward https://trx.getblock.io/28750cee-9025-42cd-9a9b-1f1c1423252a/mainnet/ 3005 / header1:value header2:value2
```

reserve
-------
```bash
$ go run cmd/main.go reserve --help
Usage: main reserve <base_host> <base_port> <reserve_host> <reserve_port> [<listen_port> [<connect_timeout> [<max_idle_seconds>]]]

reserve service

Arguments:
  <base_host>             main service host
  <base_port>             main service port
  <reserve_host>          reserve service host
  <reserve_port>          reserve service port
  [<listen_port>]         port to listen
  [<connect_timeout>]     connect_timeout
  [<max_idle_seconds>]    max_idle_seconds

Flags:
  -h, --help     Show context-sensitive help.
      --debug    Enable debug mode.

$ go run cmd/main.go reserve 127.0.0.1 8000 127.0.0.1 8080 4000
```
