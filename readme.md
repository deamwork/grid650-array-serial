# Grid 650 Array Module Serial Tool

中文说明请移步：[grid650 Array模块串口通信工具 文档](docs/readme_cn.md)

## Features & Roadmap
- [x] Custom text modification (introduced in v1.0.0)
    - [ ] Validation of text input
- [x] Time sync (introduced in v1.0.1)
- [ ] System performance info sync

## Configuration
You can use this config to map your device and http server settings.

The example config looks like this:
```yaml
http:
  listen: "0.0.0.0:80"
  tls:
    enable: false
    listen: "0.0.0.0:443"
    certificate_chain: "/path/to/your/cert.pem"
    private_key: "/path/to/your/cert.pem"
device:
  name: "/dev/tty.usbmodem14233301"
  baud: 115200
```

- `http` defines http server info, tls supported (not impl yet).
- `device` defines which grid650 needs to connect with.

### HTTP Server

TBD

### Device

You can find your grid650 device with following steps:

#### MacOS

1. plug your grid650
2. open `Terminal.app` (or `iTerm2.app`)
3. locate your device by commanding `ls -l /dev/tty*`
4. you should see a bunch of entry in the terminal
5. locate something like `tty.usbmodem14233301` (number section may various depending on your device)
6. assemble your device name as `/dev/tty.usbmodem14233301` and fill to config
7. fill baud in device section, grid650 array serial bit rate should be 115200 according to the user manual

#### Windows

TBD

#### Linux & *Unix

TBD

### Usage

```text
$ ./grid650-array-serial
grid 650 array module serial-comm server

Usage:
  grid650-array-serial [flags]
  grid650-array-serial [command]

Available Commands:
  config      Print config parse result
  help        Help about any command
  send        Send one custom text to the grid650 array device
  version     Print version information

Flags:
      --baud int                       Specific device bit rate [env BAUD]
      --config string                  Specific config file path. (default "../config/config.yaml")
      --device string                  Specific device tty [env DEVICE] (default "/dev/tty.usbmodem14233301")
  -d, --environment-development-mode   change environment mode to development [env ENVIRONMENT_DEVELOPMENT_MODE]
  -j, --environment-json-log-style     change log style to JSON [env ENVIRONMENT_JSON_LOG_STYLE]
  -h, --help                           help for grid650-array-serial
      --http-listen string             Specific http listen address. (default "0.0.0.0:80")

Use "grid650-array-serial [command] --help" for more information about a command.
```

#### Quick guide

Send text as custom text:
```bash
./grid650-array-serial send "i am grid650"
```

## Contribute

Any contribution are welcome. Please submit a pull request for feature implementation.

If you have question, please don't hesitate to open an issue.

### Build

```bash
make
```

### Run

```bash
bin/grid650-array-serial help
```

## Thanks
 - https://github.com/tarm/serial
 - https://github.com/XSAM/go-hybrid

## Chat with others

https://discord.gg/QExzFcmp

## References

The English version of documents are not available at this time.