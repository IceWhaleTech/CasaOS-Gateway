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
runtime-data-path=/var/run/casaos # See https://refspecs.linuxfoundation.org/FHS_3.0/fhs/ch05s13.html
```

## Running

Once running, gateway address and management address will be available in the files under `runtime-data-path`  specified in configuration.

```bash
$ cat /var/run/casaos/gateway.address 
[::]:8080 # port is specified in configuration

$ cat /var/run/casaos/management.address 
[::]:34703 # port is randomly assigned
```