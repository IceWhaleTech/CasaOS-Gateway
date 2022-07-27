# CasaOS-Gateway

## Configuration

Upon launching, it will search for `gateway.ini` file in the following order:

```bash
./gateway.ini
./conf/gateway.ini
$HOME/.casaos/gateway.ini
/etc/casaos/gateway.ini
```

Default configurations are:

```ini
[gateway]
port=8080
run-path=/var/run/casaos
```

## Running

Once running, gateway address and management address will be available in the files under `run-path`  specified in configuration.

```bash
$ cat /var/run/casaos/gateway.address 
[::]:8080 # port is specified in configuration

$ cat /var/run/casaos/management.address 
[::]:34703 # port is randomly assigned
```